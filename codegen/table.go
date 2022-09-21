package codegen

import (
	"reflect"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
)

type ResourceDefinition struct {
	Name  string
	Table *TableDefinition
}

type NameTransformer func(reflect.StructField) (string, error)
type TypeTransformer func(reflect.StructField) (schema.ValueType, error)

type TableDefinition struct {
	Name        string
	Columns     ColumnDefinitions
	Description string
	Relations   []string

	Resolver                      string
	IgnoreError                   string
	Multiplex                     string
	PostResourceResolver          string
	PreResourceResolver           string
	nameTransformer               NameTransformer
	typeTransformer               TypeTransformer
	skipFields                    []string
	extraColumns                  ColumnDefinitions
	structFieldsToUnwrap          []string
	unwrapAllEmbeddedStructFields bool
	logger                        zerolog.Logger
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
