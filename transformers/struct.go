package transformers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/thoas/go-funk"
)

const maxJSONTypeSchemaDepth = 5

type structTransformer struct {
	table                         *schema.Table
	skipFields                    []string
	nameTransformer               NameTransformer
	typeTransformer               TypeTransformer
	resolverTransformer           ResolverTransformer
	ignoreInTestsTransformer      IgnoreInTestsTransformer
	unwrapAllEmbeddedStructFields bool
	structFieldsToUnwrap          []string
	pkFields                      []string
	pkFieldsFound                 []string
	pkComponentFields             []string
	pkComponentFieldsFound        []string
}

func isFieldStruct(reflectType reflect.Type) bool {
	switch reflectType.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		return reflectType.Elem().Kind() == reflect.Struct
	default:
		return false
	}
}

func isTypeIgnored(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Func,
		reflect.Chan,
		reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

func (t *structTransformer) getUnwrappedFields(field reflect.StructField) []reflect.StructField {
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

func (t *structTransformer) unwrapField(field reflect.StructField) error {
	unwrappedFields := t.getUnwrappedFields(field)
	var parent *reflect.StructField
	// For non embedded structs we need to add the parent field name to the path
	if !field.Anonymous {
		parent = &field
	}
	for _, f := range unwrappedFields {
		if err := t.addColumnFromField(f, parent); err != nil {
			return fmt.Errorf("failed to add column from field %s: %w", f.Name, err)
		}
	}
	return nil
}

func (t *structTransformer) shouldUnwrapField(field reflect.StructField) bool {
	switch {
	case !isFieldStruct(field.Type):
		return false
	case slices.Contains(t.structFieldsToUnwrap, field.Name):
		return true
	case !field.Anonymous:
		return false
	case t.unwrapAllEmbeddedStructFields:
		return true
	default:
		return false
	}
}

func (t *structTransformer) ignoreField(field reflect.StructField) bool {
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

func (t *structTransformer) addColumnFromField(field reflect.StructField, parent *reflect.StructField) error {
	if t.ignoreField(field) {
		return nil
	}

	columnType, err := t.getColumnType(field)
	if err != nil {
		return err
	}

	if columnType == nil {
		return nil // ignored
	}

	path := field.Name
	name, err := t.nameTransformer(field)
	if err != nil {
		return fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
	}
	// skip field if there is no name
	if name == "" {
		return nil
	}
	if parent != nil {
		parentName, err := t.nameTransformer(*parent)
		if err != nil {
			return fmt.Errorf("failed to transform field name for parent field %s: %w", parent.Name, err)
		}
		name = parentName + "_" + name
		path = parent.Name + `.` + path
	}
	if t.table.Columns.Get(name) != nil {
		return nil
	}

	resolver := t.resolverTransformer(field, path)
	if resolver == nil {
		resolver = DefaultResolverTransformer(field, path)
	}

	column := schema.Column{
		Name:          name,
		Type:          columnType,
		Resolver:      resolver,
		IgnoreInTests: t.ignoreInTestsTransformer(field),
	}

	// Enrich JSON column with detailed schema
	if columnType == types.ExtensionTypes.JSON {
		column.TypeSchema = structSchemaToJSON(t.fieldToJSONSchema(field, 0))
	}

	for _, pk := range t.pkFields {
		if pk == path {
			// use path to allow the following
			// 1. Don't duplicate the PK fields if the unwrapped struct contains a fields with the same name
			// 2. Allow specifying the nested unwrapped field as part of the PK.
			column.PrimaryKey = true
			t.pkFieldsFound = append(t.pkFieldsFound, pk)
		}
	}

	for _, pk := range t.pkComponentFields {
		if pk == path {
			// use path to allow the following
			// 1. Don't duplicate the PK fields if the unwrapped struct contains a fields with the same name
			// 2. Allow specifying the nested unwrapped field as part of the PK.
			column.PrimaryKeyComponent = true
			t.pkComponentFieldsFound = append(t.pkComponentFieldsFound, pk)
		}
	}

	t.table.Columns = append(t.table.Columns, column)

	return nil
}

func TransformWithStruct(st any, opts ...StructTransformerOption) schema.Transform {
	t := &structTransformer{
		nameTransformer:          DefaultNameTransformer,
		typeTransformer:          DefaultTypeTransformer,
		resolverTransformer:      DefaultResolverTransformer,
		ignoreInTestsTransformer: DefaultIgnoreInTestsTransformer,
	}
	for _, opt := range opts {
		opt(t)
	}

	return func(table *schema.Table) error {
		t.table = table
		e := reflect.ValueOf(st)
		if e.Kind() == reflect.Pointer {
			e = e.Elem()
		}
		if e.Kind() == reflect.Slice {
			e = reflect.MakeSlice(e.Type(), 1, 1).Index(0)
		}
		if e.Kind() != reflect.Struct {
			return fmt.Errorf("expected struct, got %s", e.Kind())
		}
		eType := e.Type()
		for i := 0; i < e.NumField(); i++ {
			field := eType.Field(i)

			switch {
			case t.shouldUnwrapField(field):
				if err := t.unwrapField(field); err != nil {
					return err
				}
			default:
				if err := t.addColumnFromField(field, nil); err != nil {
					return fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
				}
			}
		}
		// Validate that all expected PK fields were found
		if diff := funk.SubtractString(t.pkFields, t.pkFieldsFound); len(diff) > 0 {
			return fmt.Errorf("failed to create all of the desired primary keys: %v", diff)
		}

		if diff := funk.SubtractString(t.pkComponentFields, t.pkComponentFieldsFound); len(diff) > 0 {
			return fmt.Errorf("failed to find all of the desired primary key components: %v", diff)
		}
		return nil
	}
}

func (t *structTransformer) getColumnType(field reflect.StructField) (arrow.DataType, error) {
	columnType, err := t.typeTransformer(field)
	if err != nil {
		return nil, fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
	}

	if columnType == nil {
		columnType, err = DefaultTypeTransformer(field)
		if err != nil {
			return nil, fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
	}
	return columnType, nil
}

func structSchemaToJSON(s any) string {
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(s)
	return strings.TrimSpace(b.String())
}

func normalizePointer(field reflect.StructField) reflect.Value {
	if field.Type.Kind() == reflect.Ptr {
		return reflect.New(field.Type.Elem())
	}
	return reflect.New(field.Type)
}

func (t *structTransformer) fieldToJSONSchema(field reflect.StructField, depth int) any {
	transformInput := normalizePointer(field)
	switch transformInput.Elem().Kind() {
	case reflect.Struct:
		fieldsMap := make(map[string]any)
		fieldType := transformInput.Elem().Type()
		for i := 0; i < fieldType.NumField(); i++ {
			name, err := t.nameTransformer(fieldType.Field(i))
			if err != nil {
				continue
			}
			columnType, err := t.getColumnType(fieldType.Field(i))
			if err != nil {
				continue
			}
			if columnType == nil {
				fieldsMap[name] = "any"
				continue
			}
			if columnType == types.ExtensionTypes.JSON && depth < maxJSONTypeSchemaDepth {
				fieldsMap[name] = t.fieldToJSONSchema(fieldType.Field(i), depth+1)
				continue
			}
			if arrow.IsListLike(columnType.ID()) {
				fieldsMap[name] = []any{columnType.(*arrow.ListType).Elem().String()}
				continue
			}
			fieldsMap[name] = columnType.String()
		}
		return fieldsMap
	case reflect.Map:
		keySchema, ok := t.fieldToJSONSchema(reflect.StructField{
			Type: field.Type.Key(),
		}, depth+1).(string)
		if keySchema == "" || !ok {
			return ""
		}
		valueSchema := t.fieldToJSONSchema(reflect.StructField{
			Type: field.Type.Elem(),
		}, depth+1)
		if valueSchema == "" {
			return ""
		}
		return map[string]any{
			keySchema: valueSchema,
		}
	case reflect.Slice:
		valueSchema := t.fieldToJSONSchema(reflect.StructField{
			Type: field.Type.Elem(),
		}, depth+1)
		if valueSchema == "" {
			return ""
		}
		return []any{valueSchema}
	}

	columnType, err := t.getColumnType(field)
	if err != nil {
		return ""
	}
	if columnType == nil {
		return "any"
	}
	return columnType.String()
}
