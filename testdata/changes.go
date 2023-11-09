package main

import "fmt"

type Change struct {
	Prefix    string
	FieldName string
	Operation string
	OldValue  string
	NewValue  string
}

func (c *Change) String() string {
	switch c.Operation {
	case "Added":
		return fmt.Sprintf("%s%s %s: %s", c.Prefix, c.FieldName, c.Operation, c.NewValue)
	case "Removed":
		return fmt.Sprintf("%s%s %s: %s", c.Prefix, c.FieldName, c.Operation, c.OldValue)
	case "Updated":
		return fmt.Sprintf("%s%s %s: %s -> %s", c.Prefix, c.FieldName, c.Operation, c.OldValue, c.NewValue)
	case "Set":
		return fmt.Sprintf("%s%s %s: %s", c.Prefix, c.FieldName, c.Operation, c.NewValue)
	case "Clear":
		return fmt.Sprintf("%s%s %s", c.Prefix, c.FieldName, c.Operation)
	}
	return ""
}

type DefaultChangeLogger struct {
	prefix  string
	changes []Change
}

func NewDefaultChangeLogger(prefix string) *DefaultChangeLogger {
	return &DefaultChangeLogger{
		prefix: prefix,
	}
}

func (c *DefaultChangeLogger) Append(change Change) {
	change.Prefix = c.prefix + change.Prefix
	c.changes = append(c.changes, change)
}

func (c *DefaultChangeLogger) ToString() (result []string) {
	for _, change := range c.changes {
		result = append(result, change.String())
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
