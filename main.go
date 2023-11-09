package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

var (
	headerTemplate = `// Code generated by gomutations; DO NOT EDIT.
package {{.PackageName}}

import (
	{{range .Imports}}"{{.}}"
	{{end}}
)
`
	mainMutatorTemplate = `
type Mutator{{.TypeName}} struct {
	inner   *{{.TypeName}}
	changes ChangeLogger
}

func NewMutator{{.TypeName}}(obj *{{.TypeName}}) *Mutator{{.TypeName}} {
	return &Mutator{{.TypeName}}{
		inner:   obj,
		changes: NewDefaultChangeLogger(""),
	}
}

// FormatChanges returns the changes that were made to the object as strings
func (m *Mutator{{.TypeName}}) FormatChanges() []string {
	return m.changes.ToString()
}
`

	subMutatorTemplate = `
type Mutator{{.TypeName}} struct {
	inner   *{{.TypeName}}
	changes ChangeLogger
}

func NewMutator{{.TypeName}}(obj *{{.TypeName}}, changes ChangeLogger) *Mutator{{.TypeName}} {
	return &Mutator{{.TypeName}}{
		inner:   obj,
		changes: changes,
	}
}
`

	mutateFieldTemplate = `
// Mutate{{.FieldName}} mutates the {{.FieldName}} of the {{.TypeName}} object
func (m *Mutator{{.TypeName}}) Mutate{{.FieldName}}(value {{.FieldTypeName}}) bool {
	if m.inner.{{.FieldName}} == value {
		return false
	}

	m.changes.Append(Change{
		FieldName: "{{.FieldName}}",
		Operation: "Updated",
		OldValue:  fmt.Sprintf("%v", m.inner.{{.FieldName}}),
		NewValue:  fmt.Sprintf("%v", value),
	})
	m.inner.{{.FieldName}} = value

	return true
}
`

	mutateSetTemplate = `
// Mutate{{.FieldName}} sets {{.FieldName}} of the {{.TypeName}} object
func (m *Mutator{{.TypeName}}) Set{{.FieldName}}(value {{.FieldTypeName}}) bool {

	if len(value) == 0 && len(m.inner.{{.FieldName}}) == 0 {
		return false
	}

	operation := "Set"
	if len(value) == 0 {
		operation = "Clear"
	}

	m.changes.Append(Change{
		FieldName: "{{.FieldName}}",
		Operation: operation,
		OldValue:  fmt.Sprintf("%v", m.inner.{{.FieldName}}),
		NewValue:  fmt.Sprintf("%v", value),
	})
	m.inner.{{.FieldName}} = value

	return true
}
`

	mutateSetObjTemplate = `
// Mutate{{.FieldName}} sets {{.FieldName}} of the {{.TypeName}} object
func (m *Mutator{{.TypeName}}) Set{{.FieldName}}(value *{{.FieldTypeName}}) bool {

	m.changes.Append(Change{
		FieldName: "{{.FieldName}}",
		Operation: "Set",
		OldValue:  fmt.Sprintf("%v", m.inner.{{.FieldName}}),
		NewValue:  fmt.Sprintf("%v", value),
	})
	m.inner.{{.FieldName}} = *value

	return true
}
`

	mutateSetPtrTemplate = `
// Set{{.FieldName}} sets {{.FieldName}} of the {{.TypeName}} object
func (m *Mutator{{.TypeName}}) Set{{.FieldName}}(value {{.FieldTypeName}}) bool {

	if value == nil && m.inner.{{.FieldName}} == nil {
		return false
	}

	if value == m.inner.{{.FieldName}} {
		return false
	}

	operation := "Set"
	if value == nil {
		operation = "Clear"
	}

	m.changes.Append(Change{
		FieldName: "{{.FieldName}}",
		Operation: operation,
		OldValue:  fmt.Sprintf("%v", m.inner.{{.FieldName}}),
		NewValue:  fmt.Sprintf("%v", value),
	})
	m.inner.{{.FieldName}} = value

	return true
}
`

	mutatePtrTemplate = `
// Mutate{{.FieldName}} returns a mutator for {{.FieldName}} of the {{.TypeName}} object.
// If the field is nil, it will be initialized to a new {{.FieldTypeName}} object.
func (m *Mutator{{.TypeName}}) Mutate{{.FieldName}}() *Mutator{{.FieldTypeName}} {

	if m.inner.{{.FieldName}} == nil {
		m.inner.{{.FieldName}} = &{{.FieldTypeName}}{}
	}

	return NewMutator{{.FieldTypeName}}(m.inner.{{.FieldName}}, NewChainedChangeLogger(fmt.Sprintf("{{.FieldName}} "), m.changes))
}
`

	mutateSliceTemplate = `
// Mutate{{.FieldName}}At returns a mutator for {{.FieldName}} element at index of the {{.TypeName}} object.
func (m *Mutator{{.TypeName}}) Mutate{{.FieldName}}At(index int) *Mutator{{.FieldTypeName}} {
	return NewMutator{{.FieldTypeName}}({{if .FieldTypeIsPointer}}{{else}}&{{end}}m.inner.{{.FieldName}}[index], NewChainedChangeLogger(fmt.Sprintf("{{.FieldName}} "), m.changes))
}

// Append{{.FieldName}} appends a {{.FieldName}} element of the {{.TypeName}} object.
func (m *Mutator{{.TypeName}}) Append{{.FieldName}}(value ...{{if .FieldTypeIsPointer}}*{{end}}{{.FieldTypeName}}) {
	m.changes.Append(Change{
		FieldName: "{{.FieldName}}",
		Operation: "Added",
		OldValue:  "",
		NewValue:  fmt.Sprintf("%v", value),
	})
	m.inner.{{.FieldName}} = append(m.inner.{{.FieldName}}, value...)
}

// Remove{{.FieldName}} removes a {{.FieldName}} element of the {{.TypeName}} object.
func (m *Mutator{{.TypeName}}) Remove{{.FieldName}}(index int) {
	m.changes.Append(Change{
		FieldName: "{{.FieldName}}",
		Operation: "Removed",
		OldValue:  fmt.Sprintf("%v", m.inner.{{.FieldName}}[index]),
		NewValue:  "",
	})
	m.inner.{{.FieldName}} = append(m.inner.{{.FieldName}}[:index], m.inner.{{.FieldName}}[index+1:]...)
}
`

	mutateObjTemplate = `
// Mutate{{.FieldName}} returns a mutator for {{.FieldName}} of the {{.TypeName}} object.
func (m *Mutator{{.TypeName}}) Mutate{{.FieldName}}() *Mutator{{.FieldTypeName}} {

	return NewMutator{{.FieldTypeName}}(&m.inner.{{.FieldName}}, NewChainedChangeLogger(fmt.Sprintf("{{.FieldName}} "), m.changes))
}
`
)

type templateStep struct {
	template string
	data     interface{}
}

type headerData struct {
	PackageName string
	Imports     []string
}

type mutatorData struct {
	TypeName string
}

type mutateFunctionData struct {
	TypeName           string
	FieldName          string
	FieldTypeName      string
	FieldTypeIsPointer bool
}

func trimPackagePrefix(typeName string, pkg string) string {
	start := strings.LastIndexAny(typeName, "[]*")
	if start == -1 {
		start = 0
	} else {
		start++
	}

	dot := strings.LastIndex(typeName, ".")
	if dot == -1 {
		return typeName
	}

	if typeName[start:dot] == pkg {
		return typeName[:start] + typeName[dot+1:]
	}

	return typeName
}

func trimAllPrefixes(typeName string, pkg string) string {
	start := strings.LastIndexAny(typeName, "[]*")
	if start == -1 {
		start = 0
	} else {
		start++
	}

	dot := strings.LastIndex(typeName, ".")
	if dot == -1 {
		return typeName[start:]
	}

	if typeName[start:dot] == pkg {
		return typeName[dot+1:]
	}

	return typeName
}

func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "gomutate generates Go code to mutate a Go type.\n")
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "\tgomutate [flags] -type Type [<directory>]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\ndirectory: \".\" if unspecified\n\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

var (
	flagType = flag.String("type", "", "list of types separated by comma (required)")
)

func main() {
	log.SetPrefix("gomutations: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*flagType) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	directory := "."
	if len(flag.Args()) > 0 {
		directory = flag.Args()[0]
	}

	mainFilename := "acme.go"

	loadAllSyntax := packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo

	cfg := &packages.Config{
		Mode:  loadAllSyntax,
		Tests: false,
		Dir:   directory,
	}
	pkgs, err := packages.Load(cfg, "")
	if err != nil {
		log.Fatal(err)
	}

	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found, expected 1", len(pkgs))
	}

	pkg := pkgs[0]
	packageName := pkg.Name

	fset := token.NewFileSet()

	targetTypeName := *flagType

	// gather all struct declarations
	var typeSpecs []ast.Node

	var info types.Info

	for _, filename := range pkg.GoFiles {
		if !strings.HasSuffix(filename, "/"+mainFilename) {
			continue
		}

		node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("failed to parse file %s: %s", filename, err)
		}

		// Run type checker
		info = types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}

		_, err = (&types.Config{
			Importer: importer.Default(),
		}).Check(packageName, fset, []*ast.File{node}, &info)
		if err != nil {
			log.Fatal(err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			tspec, isTypeSpec := n.(*ast.TypeSpec)
			if !isTypeSpec {
				return true
			}

			_, isStructDecl := tspec.Type.(*ast.StructType)
			if !isStructDecl {
				return true
			}

			typeSpecs = append(typeSpecs, n)
			return true
		})
	}

	var mainDecl *ast.TypeSpec
	for _, node := range typeSpecs {
		spec := node.(*ast.TypeSpec)
		if spec.Name.Name == targetTypeName {
			mainDecl = spec
			break
		}
	}

	mainMutator := mutatorData{
		TypeName: mainDecl.Name.Name,
	}

	var otherMutators []mutatorData

	for _, node := range typeSpecs {
		spec := node.(*ast.TypeSpec)
		if spec.Name.Name == targetTypeName {
			continue
		}

		otherMutators = append(otherMutators, mutatorData{
			TypeName: spec.Name.Name,
		})
	}

	header := headerData{
		PackageName: packageName,
		Imports:     []string{"fmt", "time"},
	}

	var templateSteps []templateStep

	templateSteps = append(templateSteps, []templateStep{
		{
			template: headerTemplate,
			data:     header,
		}, {
			template: mainMutatorTemplate,
			data:     mainMutator,
		}}...)

	for i := range otherMutators {
		templateSteps = append(templateSteps, templateStep{
			template: subMutatorTemplate,
			data:     otherMutators[i],
		})
	}

	templateSteps = handleStructType(templateSteps, packageName, mainDecl, info.Types, typeSpecs)

	for i, step := range templateSteps {
		tmpl, err := template.New(fmt.Sprintf("template%d", i)).Parse(step.template)
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(os.Stdout, step.data)
		if err != nil {
			panic(err)
		}
	}
}

func handleStructType(
	steps []templateStep,
	packageName string,
	structSpec *ast.TypeSpec,
	typesInfo map[ast.Expr]types.TypeAndValue,
	typeSpecs []ast.Node,
) []templateStep {

	structType := structSpec.Type.(*ast.StructType)

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		ty := typesInfo[field.Type].Type

		locallyDefined := false

		trimmedTypeStr := trimAllPrefixes(ty.String(), packageName)
		for _, spec := range typeSpecs {
			if spec.(*ast.TypeSpec).Name.Name == trimmedTypeStr {
				steps = handleStructType(steps, packageName, spec.(*ast.TypeSpec), typesInfo, typeSpecs)
				locallyDefined = true
			}
		}

		tmpl := mutateFieldTemplate
		isPointer := false
		isSlice := false
		isObject := false

		switch ty.(type) {
		case *types.Slice:
			tmpl = mutateSetTemplate
			isSlice = true
		case *types.Map:
			tmpl = mutateSetTemplate
		case *types.Pointer:
			tmpl = mutateSetPtrTemplate
			isPointer = true
			// may be a struct non-pointer type
		default:
			if locallyDefined {
				tmpl = mutateSetObjTemplate
				isObject = true
			}
		}

		steps = append(steps, templateStep{
			template: tmpl,
			data: mutateFunctionData{
				TypeName:      structSpec.Name.Name,
				FieldName:     field.Names[0].Name,
				FieldTypeName: trimPackagePrefix(ty.String(), packageName),
			},
		})

		if !locallyDefined {
			continue
		}

		if isPointer {
			steps = append(steps, templateStep{
				template: mutatePtrTemplate,
				data: mutateFunctionData{
					TypeName:      structSpec.Name.Name,
					FieldName:     field.Names[0].Name,
					FieldTypeName: trimAllPrefixes(ty.String(), packageName),
				},
			})
		}

		if isSlice {
			_, fieldTypeIsPointer := ty.(*types.Slice).Elem().Underlying().(*types.Pointer)

			steps = append(steps, templateStep{
				template: mutateSliceTemplate,
				data: mutateFunctionData{
					TypeName:           structSpec.Name.Name,
					FieldName:          field.Names[0].Name,
					FieldTypeName:      trimAllPrefixes(ty.String(), packageName),
					FieldTypeIsPointer: fieldTypeIsPointer,
				},
			})
		}

		if isObject {
			steps = append(steps, templateStep{
				template: mutateObjTemplate,
				data: mutateFunctionData{
					TypeName:      structSpec.Name.Name,
					FieldName:     field.Names[0].Name,
					FieldTypeName: trimAllPrefixes(ty.String(), packageName),
				},
			})
		}
	}

	return steps
}

func isBasic(t types.Type) bool {
	switch x := t.(type) {
	case *types.Basic:
		return true
	case *types.Slice:
		return true
	case *types.Map:
		return true
	case *types.Pointer:
		return isBasic(x.Elem())
	default:
		return false
	}
}
