// Changes package provides functionality for tracking changes to objects.
package changes

import (
	"fmt"
)

// Change represents a mutation applied to an object.
type Change struct {
	Prefix    []Prefix  `json:"prefix,omitempty"`
	FieldName string    `json:"field_name,omitempty"`
	Operation Operation `json:"operation,omitempty"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
	Key       string    `json:"key,omitempty"`
}

// FieldName defines a type for the field name enums used in prefixes.
type FieldName string

var (
	// FieldNameEmpty is the empty prefix name for root level field changes.
	FieldNameEmpty FieldName = ""
)

// Prefix defines a type used for change prefix enums.
// Prefixes are used to identify the object's field that was changed.
// If a field is a map or slice, the key is also included in the prefix.
//
// For example, if a FieldName "Foo" was added to a map with key "bar",
// the prefix would be {Name: "Foo", Key: "bar"}, and the default change logger
// would print "Foo[bar] added with value 'value'".
type Prefix struct {
	Name FieldName
	Key  string
}

// NewPrefix creates a new instance of Prefix using name only.
func NewPrefix(name FieldName) Prefix {
	return Prefix{
		Name: name,
	}
}

// NewPrefixWithKey creates a new instance of Prefix using name and key.
func NewPrefixWithKey(name FieldName, key string) Prefix {
	return Prefix{
		Name: name,
		Key:  key,
	}
}

var (
	// PrefixEmpty is the empty prefix for root level field changes.
	PrefixEmpty = NewPrefix(FieldNameEmpty)
)

// Key interface can be implemented for types so that the resulting string is
// used as the key in the prefix.
type Key interface {
	KeyForChanges() string
}

// IntoKey converts an interface to a string key, uses KeyForChanges
// if the interface implements Key.
func IntoKey(i interface{}) string {
	switch v := i.(type) {
	case Key:
		return v.KeyForChanges()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Operation defines a type used for change operation enums.
type Operation string

const (
	// OperationAdded is used for slice appends and map inserts.
	OperationAdded Operation = "added"
	// OperationRemoved is used for slice removals and map deletes.
	OperationRemoved Operation = "removed"
	// OperationUpdated is used for field updates.
	OperationUpdated Operation = "updated"
	// OperationSet is used for field sets.
	// Setting a whole slice/map, or setting a field which was previously unset.
	OperationSet Operation = "set"
	// OperationCleared is used for setting a field's value to its zero value.
	OperationCleared Operation = "cleared"
)

// Formatter defines an interface for formatting changes to human readable string.
type Formatter interface {
	Format(c *Change) string
}

// DefaultFormatter provides basic change formatting functionality.
type DefaultFormatter struct{}

// default method for printing a prefix
func printPrefix(prefix *Prefix) string {
	if prefix.Key == "" {
		return string(prefix.Name)
	}
	return fmt.Sprintf("%s[%s]", prefix.Name, prefix.Key)
}

func joinPrefixes(prefixes []Prefix) string {
	result := ""
	for i := range prefixes {
		if prefixes[i].Name == FieldNameEmpty {
			continue
		}

		if result == "" {
			result = printPrefix(&prefixes[i])
		} else {
			result = fmt.Sprintf("%s %s", result, printPrefix(&prefixes[i]))
		}
	}

	if result == "" {
		return result
	}
	return result + " "
}

// Format formats a change to a human readable string.
func (f *DefaultFormatter) Format(c *Change) string {
	prefix := joinPrefixes(c.Prefix)

	fieldName := c.FieldName
	if c.Key != "" {
		fieldName = fmt.Sprintf("%s[%s]", fieldName, c.Key)
	}

	switch c.Operation {
	case OperationAdded:
		return fmt.Sprintf("%s%s %s with value '%s'", prefix, fieldName, c.Operation, c.NewValue)
	case OperationRemoved:
		return fmt.Sprintf("%s%s %s, value was '%s'", prefix, fieldName, c.Operation, c.OldValue)
	case OperationUpdated:
		return fmt.Sprintf("%s%s %s from '%s' to '%s'", prefix, fieldName, c.Operation, c.OldValue, c.NewValue)
	case OperationSet:
		return fmt.Sprintf("%s%s %s to '%s'", prefix, fieldName, c.Operation, c.NewValue)
	case OperationCleared:
		return fmt.Sprintf("%s%s %s, value was '%s'", prefix, fieldName, c.Operation, c.OldValue)
	}
	return ""
}

// DefaultLogger provides basic change logging functionality.
type DefaultLogger struct {
	prefix    Prefix
	changes   []Change
	formatter Formatter
}

// NewDefaultLogger creates a new instance of DefaultChangeLogger.
func NewDefaultLogger(prefix Prefix) *DefaultLogger {
	return &DefaultLogger{
		prefix:    prefix,
		formatter: &DefaultFormatter{},
	}
}

// Append appends a change to the change logger.
func (c *DefaultLogger) Append(change Change) {
	change.Prefix = append([]Prefix{c.prefix}, change.Prefix...)
	c.changes = append(c.changes, change)
}

// ToString converts the change logger to a slice of human readable strings.
func (c *DefaultLogger) ToString() (result []string) {
	for i := range c.changes {
		result = append(result, c.formatter.Format(&c.changes[i]))
	}
	return
}

// ChainedLogger implements Logger interface using an inner change logger.
// Multiple change loggers are chained together by prepending prefixes.
type ChainedLogger struct {
	prefix Prefix
	inner  Logger
}

// NewChainedLogger creates a new instance of ChainedLogger.
func NewChainedLogger(prefix Prefix, inner Logger) *ChainedLogger {
	return &ChainedLogger{
		prefix: prefix,
		inner:  inner,
	}
}

// Append appends a change to the change logger.
func (c *ChainedLogger) Append(change Change) {
	change.Prefix = append([]Prefix{c.prefix}, change.Prefix...)
	c.inner.Append(change)
}

// ToString converts the change logger to a slice of human readable strings.
func (c *ChainedLogger) ToString() []string {
	return c.inner.ToString()
}

// Logger defines an interface for logging changes.
type Logger interface {
	Append(change Change)
	ToString() []string
}
