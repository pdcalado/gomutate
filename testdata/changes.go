package main

import "fmt"

// Change represents a mutation applied to an object.
type Change struct {
	Prefix    string          `json:"prefix,omitempty"`
	FieldName string          `json:"field_name,omitempty"`
	Operation ChangeOperation `json:"operation,omitempty"`
	OldValue  string          `json:"old_value,omitempty"`
	NewValue  string          `json:"new_value,omitempty"`
}

// ChangeOperation defines a type used for change operation enums.
type ChangeOperation string

const (
	ChangeOperationAdded   ChangeOperation = "Added"
	ChangeOperationRemoved ChangeOperation = "Removed"
	ChangeOperationUpdated ChangeOperation = "Updated"
	ChangeOperationSet     ChangeOperation = "Set"
	ChangeOperationClear   ChangeOperation = "Clear"
)

// ChangeFormatter defines an interface for formatting changes to human readable string.
type ChangeFormatter interface {
	Format(c *Change) string
}

// DefaultChangeFormatter provides basic change formatting functionality.
type DefaultChangeFormatter struct{}

// Format formats a change to a human readable string.
func (f *DefaultChangeFormatter) Format(c *Change) string {
	switch c.Operation {
	case ChangeOperationAdded:
		return fmt.Sprintf("%s%s %s: %s", c.Prefix, c.FieldName, c.Operation, c.NewValue)
	case ChangeOperationRemoved:
		return fmt.Sprintf("%s%s %s: %s", c.Prefix, c.FieldName, c.Operation, c.OldValue)
	case ChangeOperationUpdated:
		return fmt.Sprintf("%s%s %s: %s -> %s", c.Prefix, c.FieldName, c.Operation, c.OldValue, c.NewValue)
	case ChangeOperationSet:
		return fmt.Sprintf("%s%s %s: %s", c.Prefix, c.FieldName, c.Operation, c.NewValue)
	case ChangeOperationClear:
		return fmt.Sprintf("%s%s %s", c.Prefix, c.FieldName, c.Operation)
	}
	return ""
}

type DefaultChangeLogger struct {
	prefix    string
	changes   []Change
	formatter ChangeFormatter
}

func NewDefaultChangeLogger(prefix string) *DefaultChangeLogger {
	return &DefaultChangeLogger{
		prefix:    prefix,
		formatter: &DefaultChangeFormatter{},
	}
}

func (c *DefaultChangeLogger) Append(change Change) {
	change.Prefix = c.prefix + change.Prefix
	c.changes = append(c.changes, change)
}

func (c *DefaultChangeLogger) ToString() (result []string) {
	for i := range c.changes {
		result = append(result, c.formatter.Format(&c.changes[i]))
	}
	return
}

type ChainedChangeLogger struct {
	prefix string
	inner  ChangeLogger
}

func NewChainedChangeLogger(prefix string, inner ChangeLogger) *ChainedChangeLogger {
	return &ChainedChangeLogger{
		prefix: prefix,
		inner:  inner,
	}
}

func (c *ChainedChangeLogger) Append(change Change) {
	change.Prefix = c.prefix + change.Prefix
	c.inner.Append(change)
}

func (c *ChainedChangeLogger) ToString() []string {
	return c.inner.ToString()
}

type ChangeLogger interface {
	Append(change Change)
	ToString() []string
}
