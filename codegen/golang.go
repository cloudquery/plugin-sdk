package codegen

import (
	"embed"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/iancoleman/strcase"
)

type TableOptions func(*TableDefinition)

//go:embed templates/*.go.tpl
var TemplatesFS embed.FS

func defaultTypeTransformer(v reflect.Type) (schema.ValueType, error) {
	k := v.Kind()
	switch k {
	case reflect.Pointer:
		return defaultTypeTransformer(v.Elem())
	case reflect.String:
		return schema.TypeString, nil
	case reflect.Bool:
		return schema.TypeBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return schema.TypeInt, nil
	case reflect.Float32, reflect.Float64:
		return schema.TypeFloat, nil
	case reflect.Map:
		return schema.TypeJSON, nil
	case reflect.Struct:
		timeValue := time.Time{}
		if v == reflect.TypeOf(timeValue) {
			return schema.TypeTimestamp, nil
		}
		return schema.TypeJSON, nil
	case reflect.Slice:
		switch v.Elem().Kind() {
		case reflect.String:
			return schema.TypeStringArray, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return schema.TypeIntArray, nil
		default:
			return schema.TypeJSON, nil
		}
	default:
		return schema.TypeInvalid, fmt.Errorf("unsupported type: %s", k)
	}
}

func WithNameTransformer(transformer NameTransformer) TableOptions {
	return func(t *TableDefinition) {
		t.nameTransformer = transformer
	}
}

func WithSkipFields(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.skipFields = fields
	}
}

func WithExtraColumns(columns []ColumnDefinition) TableOptions {
	return func(t *TableDefinition) {
		t.extraColumns = columns
	}
}

// Unwrap specific struct fields (1 level deep only)
func WithUnwrapFieldsStructs(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.structFieldsToUnwrap = fields
	}
}

// WithUnwrapAllEmbeddedStructs unwraps all fields that are embedded structs (1 level deep only)
func WithUnwrapAllEmbeddedStructs() TableOptions {
	return func(t *TableDefinition) {
		t.unwrapAllEmbeddedStructFields = true
	}
}

// WithLogger
func WithLogger(logger zerolog.Logger) TableOptions {
	return func(t *TableDefinition) {
		t.logger = logger
	}
}

// WithValueTypeTransformer sets a function that can override the schema type for specific fields. Return `schema.TypeInvalid` to fall back to default behavior.
func WithTypeTransformer(transformer TypeTransformer) TableOptions {
	return func(t *TableDefinition) {
		t.typeTransformer = transformer
	}
}

func DefaultTransformer(field reflect.StructField) (string, error) {
	name := field.Name
	if jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]; len(jsonTag) > 0 {
		// return empty string if the field is not related api response
		if jsonTag == "-" {
			return "", nil
		}
		name = jsonTag
	}
	return strcase.ToSnake(name), nil
}

func sliceContains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func isFieldStruct(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Struct || (reflectType.Kind() == reflect.Ptr && reflectType.Elem().Kind() == reflect.Struct)
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

func (t *TableDefinition) addColumnFromField(field reflect.StructField, parent *reflect.StructField) error {
	if t.ignoreField(field) {
		return nil
	}

	var columnType schema.ValueType = schema.TypeInvalid
	var err error
	if t.typeTransformer != nil {
		columnType, err = t.typeTransformer(field.Type)
		if err != nil {
			return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
	}

	if columnType == schema.TypeInvalid {
		var err error
		columnType, err = defaultTypeTransformer(field.Type)
		if err != nil {
			return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
	}

	// generate a PathResolver to use by default
	pathResolver := fmt.Sprintf(`schema.PathResolver("%s")`, field.Name)
	name, err := t.nameTransformer(field)
	if err != nil {
		return fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
	}
	// skip field if there is no name
	if name == "" {
		return nil
	}
	if parent != nil {
		pathResolver = fmt.Sprintf(`schema.PathResolver("%s.%s")`, parent.Name, field.Name)
		parentName, err := t.nameTransformer(*parent)
		if err != nil {
			return fmt.Errorf("failed to transform field name for parent field %s: %w", parent.Name, err)
		}
		name = fmt.Sprintf("%s_%s", parentName, name)
	}

	column := ColumnDefinition{
		Name:     name,
		Type:     columnType,
		Resolver: pathResolver,
	}
	t.Columns = append(t.Columns, column)
	return nil
}

// NewTableFromStruct creates a new TableDefinition from a struct by inspecting its fields
func NewTableFromStruct(name string, obj interface{}, opts ...TableOptions) (*TableDefinition, error) {
	t := &TableDefinition{
		Name:            name,
		nameTransformer: DefaultTransformer,
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
				if err := t.addColumnFromField(f, parent); err != nil {
					return nil, fmt.Errorf("failed to add column from field %s: %w", f.Name, err)
				}
			}
		} else {
			if err := t.addColumnFromField(field, nil); err != nil {
				return nil, fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
			}
		}
	}

	return t, nil
}

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
