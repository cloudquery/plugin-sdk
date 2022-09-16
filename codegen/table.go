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
	Columns     ColumnDefinitions
	Description string
	Relations   []string

	Resolver                        string
	IgnoreError                     string
	Multiplex                       string
	PostResourceResolver            string
	PreResourceResolver             string
	nameTransformer                 func(string) string
	skipFields                      []string
	extraColumns                    ColumnDefinitions
	fieldsToUnwrapWithParentName    []string // These structs are unwrapped into columns named "<parent>_<field>"
	fieldsToUnwrapWithoutParentName []string // These structs are unwrapped into columns named "<field>"
	unwrapAllEmbeddedStructFields   bool     // Embedded (anonymous) structs are always unwrapped without the parent name
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
