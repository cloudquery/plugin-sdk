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

func sortAndRemoveDuplicates(fields []arrow.Field) []arrow.Field {
	newFields := make([]arrow.Field, len(fields))
	copy(newFields, fields)
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

func removeFieldsByType(fields []arrow.Field, dt ...arrow.Type) []arrow.Field {
	newFields := make([]arrow.Field, len(fields))
	copy(newFields, fields)
	for _, d := range dt {
		for i := 0; i < len(newFields); i++ {
			if newFields[i].Type.ID() == d {
				newFields = append(newFields[:i], newFields[i+1:]...)
				i--
			}
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

type TestSourceOptions struct {
	IncludeDates      bool
	IncludeMaps       bool
	IncludeStructs    bool
	IncludeIntervals  bool
	IncludeDurations  bool
	IncludeTimes      bool // time of day types
	IncludeLargeTypes bool // e.g. large binary, large string
}

// TestSourceFields returns fields for all Arrow types and composites thereof. TestSourceOptions controls
// which types are included.
func TestSourceFields(opts TestSourceOptions) []arrow.Field {
	// cq fields
	var cqFields []arrow.Field
	cqIDMetadata := arrow.NewMetadata([]string{schema.MetadataUnique}, []string{"true"})
	cqFields = append(cqFields, arrow.Field{Name: schema.CqIDColumn.Name, Type: types.NewUUIDType(), Nullable: false, Metadata: cqIDMetadata})
	cqFields = append(cqFields, arrow.Field{Name: schema.CqParentIDColumn.Name, Type: types.NewUUIDType(), Nullable: false})

	var basicFields []arrow.Field
	basicFields = append(basicFields, PrimitiveFields()...)
	basicFields = append(basicFields, BinaryFields()...)
	basicFields = append(basicFields, FixedWidthFields()...)

	// add extensions
	basicFields = append(basicFields, arrow.Field{Name: "uuid", Type: types.NewUUIDType(), Nullable: true})
	basicFields = append(basicFields, arrow.Field{Name: "inet", Type: types.NewInetType(), Nullable: true})
	basicFields = append(basicFields, arrow.Field{Name: "mac", Type: types.NewMacType(), Nullable: true})

	// sort and remove duplicates (e.g. date32 and date64 appear twice)
	basicFields = sortAndRemoveDuplicates(basicFields)

	// we don't support float16 right now
	basicFields = removeFieldsByType(basicFields, arrow.FLOAT16)

	if !opts.IncludeDates {
		basicFields = removeFieldsByType(basicFields, arrow.DATE32, arrow.DATE64)
	}
	if !opts.IncludeTimes {
		basicFields = removeFieldsByType(basicFields, arrow.TIME32, arrow.TIME64)
	}
	if !opts.IncludeIntervals {
		basicFields = removeFieldsByType(basicFields, arrow.INTERVAL_DAY_TIME, arrow.INTERVAL_MONTHS, arrow.INTERVAL_MONTH_DAY_NANO)
	}
	if !opts.IncludeDurations {
		basicFields = removeFieldsByType(basicFields, arrow.DURATION)
	}
	if !opts.IncludeLargeTypes {
		basicFields = removeFieldsByType(basicFields, arrow.LARGE_BINARY, arrow.LARGE_STRING)
	}

	var compositeFields []arrow.Field
	compositeFields = append(compositeFields, ListOfFields(basicFields)...)

	if opts.IncludeMaps {
		compositeFields = append(compositeFields, MapOfFields(basicFields)...)
	}

	// add JSON later, we don't want to include it as a list or map right now (it causes complications with JSON unmarshalling)
	basicFields = append(basicFields, arrow.Field{Name: "json", Type: types.NewJSONType(), Nullable: true})

	if opts.IncludeStructs {
		// struct with all the types
		compositeFields = append(compositeFields, arrow.Field{Name: "struct", Type: arrow.StructOf(basicFields...), Nullable: true})

		// struct with nested struct
		compositeFields = append(compositeFields, arrow.Field{Name: "nested_struct", Type: arrow.StructOf(arrow.Field{Name: "inner", Type: arrow.StructOf(basicFields...), Nullable: true}), Nullable: true})
	}

	allFields := append(append(cqFields, basicFields...), compositeFields...)
	return allFields
}

var PKColumnNames = []string{"uuid_pk", "string_pk"}

// TestSourceSchemaWithMetadata returns a schema for all Arrow types and composites thereof.
func TestSourceSchemaWithMetadata(md *arrow.Metadata, opts TestSourceOptions) *arrow.Schema {
	var fields []arrow.Field
	pkMetadata := map[string]string{
		schema.MetadataPrimaryKey: "true",
		schema.MetadataUnique:     "true",
	}
	fields = append(fields, arrow.Field{Name: "uuid_pk", Type: types.NewUUIDType(), Nullable: false, Metadata: arrow.MetadataFrom(pkMetadata)})
	fields = append(fields, arrow.Field{Name: "string_pk", Type: arrow.BinaryTypes.String, Nullable: false, Metadata: arrow.MetadataFrom(pkMetadata)})
	fields = append(fields, arrow.Field{Name: schema.CqSourceNameColumn.Name, Type: arrow.BinaryTypes.String, Nullable: true})
	fields = append(fields, arrow.Field{Name: schema.CqSyncTimeColumn.Name, Type: arrow.FixedWidthTypes.Timestamp_us, Nullable: true})
	fields = append(fields, TestSourceFields(opts)...)
	return arrow.NewSchema(fields, md)
}

// TestSourceSchema returns a schema for all Arrow types and composites thereof.
func TestSourceSchema(name string, opts TestSourceOptions) *arrow.Schema {
	metadata := arrow.MetadataFrom(map[string]string{
		schema.MetadataTableName: name,
	})
	return TestSourceSchemaWithMetadata(&metadata, opts)
}
