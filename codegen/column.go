package codegen

import (
	"github.com/cloudquery/plugin-sdk/schema"
)

type (
	ColumnDefinitions []ColumnDefinition
	ColumnDefinition  struct {
		// Name of the column
		Name          string
		Type          schema.ValueType
		Resolver      string
		Description   string
		IgnoreInTests bool
		Options       schema.ColumnCreationOptions
	}
)

func (c ColumnDefinitions) GetByName(name string) *ColumnDefinition {
	for _, col := range c {
		if col.Name == name {
			return &col
		}
	}
	return nil
}
