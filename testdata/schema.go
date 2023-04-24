package testdata

import (
	"reflect"
	"sort"
	"strings"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/types"
)

func PrimitiveFields() []arrow.Field {
	primitiveTypesValue := reflect.ValueOf(arrow.PrimitiveTypes)
	primitiveTypesType := reflect.TypeOf(arrow.PrimitiveTypes)
	fields := make([]arrow.Field, primitiveTypesType.NumField())
	for i := 0; i < primitiveTypesType.NumField(); i++ {
		fieldName := primitiveTypesType.Field(i).Name
		dataType := primitiveTypesValue.FieldByName(fieldName).Interface().(arrow.DataType)
		fields[i] = arrow.Field{Name: strings.ToLower(fieldName), Type: dataType, Nullable: true}
	}
	return fields
}

func BinaryFields() []arrow.Field {
	binaryTypesValue := reflect.ValueOf(arrow.BinaryTypes)
	binaryTypesType := reflect.TypeOf(arrow.BinaryTypes)
	fields := make([]arrow.Field, binaryTypesType.NumField())
	for i := 0; i < binaryTypesType.NumField(); i++ {
		fieldName := binaryTypesType.Field(i).Name
		dataType := binaryTypesValue.FieldByName(fieldName).Interface().(arrow.DataType)
		fields[i] = arrow.Field{Name: strings.ToLower(fieldName), Type: dataType, Nullable: true}
	}
	return fields
}

func FixedWidthFields() []arrow.Field {
	fixedWidthTypesValue := reflect.ValueOf(arrow.FixedWidthTypes)
	fixedWidthTypesType := reflect.TypeOf(arrow.FixedWidthTypes)
	fields := make([]arrow.Field, fixedWidthTypesType.NumField())
	for i := 0; i < fixedWidthTypesType.NumField(); i++ {
		fieldName := fixedWidthTypesType.Field(i).Name
		dataType := fixedWidthTypesValue.FieldByName(fieldName).Interface().(arrow.DataType)
		fields[i] = arrow.Field{Name: strings.ToLower(fieldName), Type: dataType, Nullable: true}
	}
	return fields
}

func SortAndRemoveDuplicates(fields []arrow.Field) []arrow.Field {
	newFields := fields
	sort.Slice(newFields, func(i, j int) bool {
		return newFields[i].Name < newFields[j].Name
	})
	for i := 0; i < len(newFields)-1; i++ {
		if newFields[i].Name == newFields[i+1].Name {
			newFields = append(newFields[:i], newFields[i+1:]...)
			i--
		}
	}
	return newFields
}

func ListOfFields(baseFields []arrow.Field) []arrow.Field {
	fields := make([]arrow.Field, len(baseFields))
	for i := 0; i < len(baseFields); i++ {
		fields[i] = arrow.Field{Name: baseFields[i].Name + "_list", Type: arrow.ListOf(baseFields[i].Type), Nullable: true}
	}
	return fields
}

func MapOfFields(baseFields []arrow.Field) []arrow.Field {
	fields := make([]arrow.Field, len(baseFields))
	for i := 0; i < len(baseFields); i++ {
		fields[i] = arrow.Field{Name: baseFields[i].Name + "_map", Type: arrow.MapOf(baseFields[i].Type, baseFields[i].Type), Nullable: true}
	}
	return fields
}

// TestSourceFields returns fields for all Arrow types and composites thereof
func TestSourceFields() []arrow.Field {
	// cq fields
	cqFields := make([]arrow.Field, 0)
	cqIDMetadata := arrow.NewMetadata([]string{schema.MetadataUnique}, []string{"true"})
	cqFields = append(cqFields, arrow.Field{Name: schema.CqIDColumn.Name, Type: types.NewUUIDType(), Nullable: false, Metadata: cqIDMetadata})
	cqFields = append(cqFields, arrow.Field{Name: schema.CqParentIDColumn.Name, Type: types.NewUUIDType(), Nullable: false})

	basicFields := make([]arrow.Field, 0)
	basicFields = append(basicFields, PrimitiveFields()...)
	basicFields = append(basicFields, BinaryFields()...)
	basicFields = append(basicFields, FixedWidthFields()...)

	// add extensions
	basicFields = append(basicFields, arrow.Field{Name: "uuid", Type: types.NewUUIDType(), Nullable: true})
	basicFields = append(basicFields, arrow.Field{Name: "inet", Type: types.NewInetType(), Nullable: true})
	basicFields = append(basicFields, arrow.Field{Name: "mac", Type: types.NewMacType(), Nullable: true})

	// sort and remove duplicates (e.g. date32 and date64 appear twice)
	basicFields = SortAndRemoveDuplicates(basicFields)

	compositeFields := make([]arrow.Field, 0)
	compositeFields = append(compositeFields, ListOfFields(basicFields)...)
	compositeFields = append(compositeFields, MapOfFields(basicFields)...)

	// add JSON later, we don't want to include it as a list or map right now
	basicFields = append(basicFields, arrow.Field{Name: "json", Type: types.NewJSONType(), Nullable: true})

	// struct with all the types
	compositeFields = append(compositeFields, arrow.Field{Name: "struct", Type: arrow.StructOf(basicFields...), Nullable: true})

	// struct with nested struct
	compositeFields = append(compositeFields, arrow.Field{Name: "nested_struct", Type: arrow.StructOf(arrow.Field{Name: "inner", Type: arrow.StructOf(basicFields...), Nullable: true}), Nullable: true})

	allFields := append(append(cqFields, basicFields...), compositeFields...)
	return allFields
}

func TestSourceSchemaWithMetadata(md *arrow.Metadata) *arrow.Schema {
	fields := TestSourceFields()
	fields = append(fields, arrow.Field{Name: schema.CqSourceNameColumn.Name, Type: arrow.BinaryTypes.String, Nullable: true})
	fields = append(fields, arrow.Field{Name: schema.CqSyncTimeColumn.Name, Type: arrow.FixedWidthTypes.Timestamp_us, Nullable: true})
	return arrow.NewSchema(fields, md)
}

func TestSourceSchema(name string) *arrow.Schema {
	keys := []string{
		schema.MetadataTableName,
	}
	values := []string{
		name,
	}
	metadata := arrow.NewMetadata(keys, values)
	return TestSourceSchemaWithMetadata(&metadata)
}
