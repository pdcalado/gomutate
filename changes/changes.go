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
}

// Prefix defines a type used for change prefix enums.
// Prefixes are used to identify the object's field that was changed.
type Prefix string

const (
	// PrefixEmpty is the empty prefix for root level field changes.
	PrefixEmpty Prefix = ""
)

// Operation defines a type used for change operation enums.
type Operation string

const (
	// OperationAdded is used for slice appends and map inserts.
	OperationAdded Operation = "Added"
	// OperationRemoved is used for slice removals and map deletes.
	OperationRemoved Operation = "Removed"
	// OperationUpdated is used for field updates.
	OperationUpdated Operation = "Updated"
	// OperationSet is used for field sets.
	// Setting a whole slice/map, or setting a field which was previously unset.
	OperationSet Operation = "Set"
	// OperationClear is used for setting a field's value to its zero value.
	OperationClear Operation = "Clear"
)

// Formatter defines an interface for formatting changes to human readable string.
type Formatter interface {
	Format(c *Change) string
}

// DefaultFormatter provides basic change formatting functionality.
type DefaultFormatter struct{}

func joinPrefixes(prefixes []Prefix) string {
	result := ""
	for i := range prefixes {
		if prefixes[i] == PrefixEmpty {
			continue
		}

		if result == "" {
			result = string(prefixes[i])
		} else {
			result = fmt.Sprintf("%s %s", result, prefixes[i])
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
	switch c.Operation {
	case OperationAdded:
		return fmt.Sprintf("%s%s %s: %s", prefix, c.FieldName, c.Operation, c.NewValue)
	case OperationRemoved:
		return fmt.Sprintf("%s%s %s: %s", prefix, c.FieldName, c.Operation, c.OldValue)
	case OperationUpdated:
		return fmt.Sprintf("%s%s %s: %s -> %s", prefix, c.FieldName, c.Operation, c.OldValue, c.NewValue)
	case OperationSet:
		return fmt.Sprintf("%s%s %s: %s", prefix, c.FieldName, c.Operation, c.NewValue)
	case OperationClear:
		return fmt.Sprintf("%s%s %s", prefix, c.FieldName, c.Operation)
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
