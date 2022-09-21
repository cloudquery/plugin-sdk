package codegen

import (
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type (
	ResourceDefinition struct {
		Name  string
		Table *TableDefinition
	}
	TableDefinition struct {
		Name        string
		Columns     ColumnDefinitions
		Description string
		Relations   []string

		Resolver             string
		IgnoreError          string
		Multiplex            string
		PostResourceResolver string
		PreResourceResolver  string

		nameTransformer NameTransformer
		typeTransformer TypeTransformer

		extraColumns ColumnDefinitions

		skipFields           []string
		structFieldsToUnwrap []string

		unwrapAllEmbeddedStructFields bool

		logger zerolog.Logger
	}
)

func (t *TableDefinition) ignoreField(field reflect.StructField) bool {
	switch {
	case len(field.Name) == 0,
		slices.Contains(t.skipFields, field.Name),
		!field.IsExported(),
		isTypeIgnored(field.Type):
		return true
	default:
		return false
	}
}

// NewTableFromStruct creates a new TableDefinition from a struct by inspecting its fields
func NewTableFromStruct(name string, obj interface{}, opts ...TableOption) (*TableDefinition, error) {
	t := &TableDefinition{
		Name:            name,
		nameTransformer: DefaultNameTransformer,
		typeTransformer: DefaultTypeTransformer,
	}
	for _, opt := range opts {
		opt(t)
	}

	e := reflect.ValueOf(obj)
	if e.Kind() == reflect.Pointer {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", e.Kind())
	}

	t.Columns = append(t.Columns, t.extraColumns...)

	eType := e.Type()
	for i := 0; i < e.NumField(); i++ {
		field := eType.Field(i)

		switch {
		case t.shouldUnwrapField(field):
			if err := t.unwrapField(field); err != nil {
				return nil, err
			}
		default:
			if err := t.addColumnFromField(field, nil); err != nil {
				return nil, fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
			}
		}
	}

	return t, nil
}
