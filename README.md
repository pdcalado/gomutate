# gomutate

![ci workflow](https://github.com/pdcalado/gomutate/actions/workflows/ci.yml/badge.svg)

Generate code to mutate your Go types.

For a type like this:

```go
type Acme struct {
 Name        string
 YearOfBirth int
 Employees   []*Employee
}

type Employee struct {
 Name     string
 Position string
 JoinedAt time.Time
}
```

Generates code allowing you to mutate it like this:

```go
acme := Acme{
    Name:        "Acme Inc.",
    YearOfBirth: 2000,
    Employees: []*Employee{
        {
            Name:     "John Doe",
            Position: "CEO",
            JoinedAt: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
        },
        {
            Name:     "Jane Doe",
            Position: "CTO",
            JoinedAt: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
        },
    },
}

mutator := NewMutatorAcme(&acme)
mutator.SetYearOfBirth(2019)
mutator.EmployeesAt(0).SetName("John Smith")
```

and reporting changes like this:

```go
fmt.Println(mutator.FormatChanges())
```

which outputs:

```
YearOfBirth Updated: 2018 -> 2019
Employees Name Updated: John Doe -> John Smith
```

## Usage

```console
go get github.com/pdcalado/gomutate
```

```console
go run github.com/pdcalado/gomutate -type <type-name> -w <path-to-output-file> <path-to-input-file>
```

Input and output files must be in the same package. Omit the `-w` flag to print to stdout.

## Features

See our [tests](./testdata/main.go) for examples of other possibly unlisted supported operations.

- operations with pointers and basic types are idempotent (the same mutation performed twice must only report one change)
- a custom formatter and custom change logger can be provided
- mutate a field with basic type
- mutate a field with a slice or map of basic types
- mutate a field with a slice of structs or struct pointers defined in the same package
- mutate a field with a map of structs or struct pointers defined in the same package, including struct pointers as keys
- append and delete from a slice
- insert and delete from a map

## Limitations

- Only supports structs
- Only supports chaining mutators for types defined in the same package
- Only supports types with exported fields
- Does not support nested slices or maps, like `[][]string` or `map[string]map[int]string`
