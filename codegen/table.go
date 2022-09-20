package codegen

import (
	"embed"
	"fmt"
	"io"
	"reflect"
	"text/template"
	"unicode"

	"github.com/cloudquery/plugin-sdk/schema"
)

type TableDefinition struct {
	Name        string
	Columns     ColumnDefinitions
	Description string
	Relations   []string

	Resolver             string
	IgnoreError          string
	Multiplex            string
	PostResourceResolver string
	PreResourceResolver  string

	nameTransformer               func(reflect.StructField) string
	skipFields                    []string
	extraColumns                  ColumnDefinitions
	structFieldsToUnwrap          []string
	unwrapAllEmbeddedStructFields bool
	customTypeResolvers           map[schema.ValueType]string // value = ResolverProvider
}

func (t *TableDefinition) shouldUnwrapField(field reflect.StructField) bool {
	return isFieldStruct(field.Type) && (t.unwrapAllEmbeddedStructFields && field.Anonymous || sliceContains(t.structFieldsToUnwrap, field.Name))
}

func (t *TableDefinition) getUnwrappedFields(field reflect.StructField) []reflect.StructField {
	reflectType := field.Type
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	fields := make([]reflect.StructField, 0)
	for i := 0; i < reflectType.NumField(); i++ {
		sf := reflectType.Field(i)
		if t.ignoreField(sf) {
			continue
		}

		fields = append(fields, sf)
	}
	return fields
}

func (t *TableDefinition) ignoreField(field reflect.StructField) bool {
	return len(field.Name) == 0 || unicode.IsLower(rune(field.Name[0])) || sliceContains(t.skipFields, field.Name)
}

func (t *TableDefinition) addColumnFromField(field reflect.StructField, parent *reflect.StructField) {
	if t.ignoreField(field) {
		return
	}

	columnType, err := valueToSchemaType(field.Type)
	if err != nil {
		fmt.Printf("skipping field %s on table %s, got err: %v\n", field.Name, t.Name, err)
		return
	}

	path := field.Name
	name := t.nameTransformer(field)
	// skip field if there is no name
	if name == "" {
		return
	}
	if parent != nil {
		path = parent.Name + `.` + path
		name = t.nameTransformer(*parent) + "_" + name
	}
	resolver := `schema.PathResolver("` + path + `")`
	resolverProvider, ok := t.customTypeResolvers[columnType]
	if ok {
		resolver = resolverProvider + `("` + path + `")`
	}

	column := ColumnDefinition{
		Name:     name,
		Type:     columnType,
		Resolver: resolver,
	}
	t.Columns = append(t.Columns, column)
}

//go:embed templates/*.go.tpl
var TemplatesFS embed.FS

func (t *TableDefinition) GenerateTemplate(wr io.Writer) error {
	tpl, err := template.New("table.go.tpl").ParseFS(TemplatesFS, "templates/*")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tpl.Execute(wr, t); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}

// NewTableFromStruct creates a new TableDefinition from a struct by inspecting its fields
func NewTableFromStruct(name string, obj interface{}, opts ...TableOptions) (*TableDefinition, error) {
	t := &TableDefinition{
		Name:            name,
		nameTransformer: defaultTransformer,
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

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)

		if t.shouldUnwrapField(field) {
			unwrappedFields := t.getUnwrappedFields(field)
			var parent *reflect.StructField
			// For non embedded structs we need to add the parent field name to the path
			if !field.Anonymous {
				parent = &field
			}
			for _, f := range unwrappedFields {
				t.addColumnFromField(f, parent)
			}
		} else {
			t.addColumnFromField(field, nil)
		}
	}

	return t, nil
}
