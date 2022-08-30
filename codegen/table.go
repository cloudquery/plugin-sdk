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
	Columns     []ColumnDefinition
	Relations   []*TableDefinition

	Resolver             string
	IgnoreError          string
	Multiplex            string
	PostResourceResolver string
	Options              schema.TableCreationOptions
	nameTransformer      func(string) string
	skipFields           []string
	overrideColumns      ColumnDefinitions
	extraColumns         ColumnDefinitions
	descriptionsEnabled  bool
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
