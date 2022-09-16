package codegen

import (
	"embed"
	"fmt"
	"io"
	"reflect"
	"text/template"
	"unicode"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/iancoleman/strcase"
)

type TableOptions func(*TableDefinition)

//go:embed templates/*.go.tpl
var TemplatesFS embed.FS

func valueToSchemaType(v reflect.Type) (schema.ValueType, error) {
	k := v.Kind()
	switch k {
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
		t := v.PkgPath() + "." + v.Name()
		if t == "time.Time" {
			return schema.TypeTimestamp, nil
		}
		return schema.TypeJSON, nil
	case reflect.Pointer:
		return valueToSchemaType(v.Elem())
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

func WithNameTransformer(transformer func(string) string) TableOptions {
	return func(t *TableDefinition) {
		t.Name = transformer(t.Name)
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
func WithUnwrapFieldsWithParentName(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.fieldsToUnwrapWithParentName = fields
	}
}

func WithUnwrapFieldsNoParentName(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.fieldsToUnwrapWithoutParentName = fields
	}
}

// Unwrap all fields that are embedded structs (1 level deep only)
func WithUnwrapAllEmbeddedStructs() TableOptions {
	return func(t *TableDefinition) {
		t.unwrapAllEmbeddedStructFields = true
	}
}

func defaultTransformer(name string) string {
	return strcase.ToSnake(name)
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

func (t *TableDefinition) shouldUnwrapFieldWithParentName(field reflect.StructField) bool {
	return isFieldStruct(field.Type) && sliceContains(t.fieldsToUnwrapWithParentName, field.Name)
}

func (t *TableDefinition) shouldUnwrapFieldWithoutParentName(field reflect.StructField) bool {
	return isFieldStruct(field.Type) && (t.unwrapAllEmbeddedStructFields && field.Anonymous) || sliceContains(t.fieldsToUnwrapWithoutParentName, field.Name)
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

// Adds a column with PathResolver("<parentFieldName>.<field.Name>").
// addParentFieldNameToColumnName is only used if 'parentFieldName' != "".
// If 'addParentFieldNameToColumnName' is true, the column will be named "<parent_field_name>.<field_name>" (in lowercase).
// If 'addParentFieldNameToColumnName' is false, the column will just be named "<field_name>"
func (t *TableDefinition) addColumnFromField(field reflect.StructField, parentFieldName string, addParentFieldNameToColumnName bool) {
	if t.ignoreField(field) {
		return
	}

	columnType, err := valueToSchemaType(field.Type)
	if err != nil {
		fmt.Printf("skipping field %s, got err: %v\n", field.Name, err)
		return
	}

	var pathResolver string
	var columnName string

	// generate a PathResolver to use by default
	if parentFieldName == "" {
		pathResolver = fmt.Sprintf(`schema.PathResolver("%s")`, field.Name)
		columnName = t.nameTransformer(field.Name)
	} else {
		pathResolver = fmt.Sprintf(`schema.PathResolver("%s.%s")`, parentFieldName, field.Name)
		if addParentFieldNameToColumnName {
			columnName = t.nameTransformer(parentFieldName) + "_" + t.nameTransformer(field.Name)
		} else {
			columnName = t.nameTransformer(field.Name)
		}
	}

	column := ColumnDefinition{
		Name:     columnName,
		Type:     columnType,
		Resolver: pathResolver,
	}
	t.Columns = append(t.Columns, column)
}

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

		if t.shouldUnwrapFieldWithParentName(field) {
			unwrappedFields := t.getUnwrappedFields(field)
			for _, f := range unwrappedFields {
				t.addColumnFromField(
					f,
					field.Name, /*parentFieldName*/
					true,       /*addParentFieldNameToColumnName*/
				)

			}
		} else if t.shouldUnwrapFieldWithoutParentName(field) {
			unwrappedField := t.getUnwrappedFields(field)
			for _, f := range unwrappedField {
				t.addColumnFromField(
					f,
					field.Name, /*parentFieldName*/
					false,      /*addParentFieldNameToColumnName*/
				)
			}
		} else {
			t.addColumnFromField(
				field,
				"",    /*parentFieldName*/
				false, /*addParentFieldNameToColumnName*/
			)
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
