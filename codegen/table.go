package codegen

import (
	"github.com/cloudquery/plugin-sdk/schema"
)

type ResourceDefinition struct {
	Name  string
	Table *TableDefinition
}

type TableDefinition struct {
	Name        string
	Description string
	Columns     ColumnDefinitions
	Relations   []string

	Resolver             string
	IgnoreError          string
	Multiplex            string
	PostResourceResolver string
	PreResourceResolver  string
	nameTransformer      func(string) string
	skipFields           []string
	extraColumns         ColumnDefinitions
}

type ColumnDefinitions []ColumnDefinition

type ColumnDefinition struct {
	// Name name of the column
	Name          string
	Type          schema.ValueType
	Resolver      string
	Description   string
	IgnoreInTests bool
	Options       schema.ColumnCreationOptions
}

func (c ColumnDefinitions) GetByName(name string) *ColumnDefinition {
	for _, col := range c {
		if col.Name == name {
			return &col
		}
	}
	return nil
}
