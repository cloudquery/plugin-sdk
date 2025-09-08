package schema

import (
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/internal/sha1"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
)

const (
	MetadataUnique              = "cq:extension:unique"
	MetadataPrimaryKey          = "cq:extension:primary_key"
	MetadataPrimaryKeyComponent = "cq:extension:primary_key_component"
	MetadataConstraintName      = "cq:extension:constraint_name"
	MetadataIncremental         = "cq:extension:incremental"
	MetadataTypeSchema          = "cq:extension:type_schema"

	MetadataTrue                   = "true"
	MetadataFalse                  = "false"
	MetadataTableName              = "cq:table_name"
	MetadataTableDescription       = "cq:table_description"
	MetadataTableTitle             = "cq:table_title"
	MetadataTableDependsOn         = "cq:table_depends_on"
	MetadataTableIsPaid            = "cq:table_paid"
	MetadataTablePermissionsNeeded = "cq:table_permissions_needed"
	MetadataTableSensitiveColumns  = "cq:table_sensitive_columns"
)

type Schemas []*arrow.Schema

func (s Schemas) Len() int {
	return len(s)
}

func (s Schemas) SchemaByName(name string) *arrow.Schema {
	for _, sc := range s {
		tableName, ok := sc.Metadata().GetValue(MetadataTableName)
		if !ok {
			continue
		}
		if tableName == name {
			return sc
		}
	}
	return nil
}

func hashRecord(record arrow.RecordBatch) arrow.Array {
	numRows := int(record.NumRows())
	fields := record.Schema().Fields()
	hashArray := types.NewUUIDBuilder(memory.DefaultAllocator)
	hashArray.Reserve(numRows)
	for row := range numRows {
		rowHash := sha1.New()
		for col := 0; col < int(record.NumCols()); col++ {
			fieldName := fields[col].Name
			rowHash.Write([]byte(fieldName))
			value := record.Column(col).ValueStr(row)
			_, _ = rowHash.Write([]byte(value))
		}
		// This part ensures that we conform to the UUID spec
		hashArray.Append(newUUID(uuid.NameSpaceURL, rowHash.Sum(nil)))
	}
	return hashArray.NewArray()
}

func newUUID(space uuid.UUID, data []byte) uuid.UUID {
	return uuid.NewHash(sha1.New(), space, data, 5)
}

func nullUUIDsForRecord(numRows int) arrow.Array {
	uuidArray := types.NewUUIDBuilder(memory.DefaultAllocator)
	uuidArray.AppendNulls(numRows)
	return uuidArray.NewArray()
}

func StringArrayFromValue(value string, nRows int) arrow.Array {
	arrayBuilder := array.NewStringBuilder(memory.DefaultAllocator)
	arrayBuilder.Reserve(nRows)
	for range nRows {
		arrayBuilder.AppendString(value)
	}
	return arrayBuilder.NewArray()
}

func TimestampArrayFromTime(t time.Time, unit arrow.TimeUnit, timeZone string, nRows int) (arrow.Array, error) {
	ts, err := arrow.TimestampFromTime(t, unit)
	if err != nil {
		return nil, err
	}
	arrayBuilder := array.NewTimestampBuilder(memory.DefaultAllocator, &arrow.TimestampType{Unit: unit, TimeZone: timeZone})
	arrayBuilder.Reserve(nRows)
	for range nRows {
		arrayBuilder.Append(ts)
	}
	return arrayBuilder.NewArray(), nil
}

func ReplaceFieldInRecord(src arrow.RecordBatch, fieldName string, field arrow.Array) (record arrow.RecordBatch, err error) {
	fieldIndexes := src.Schema().FieldIndices(fieldName)
	for i := range fieldIndexes {
		record, err = src.SetColumn(fieldIndexes[i], field)
		if err != nil {
			return nil, err
		}
	}
	return record, nil
}

func AddInternalColumnsToRecord(record arrow.RecordBatch, cqClientIDValue string) (arrow.RecordBatch, error) {
	schema := record.Schema()
	nRows := int(record.NumRows())

	newFields := []arrow.Field{}
	newColumns := []arrow.Array{}

	var err error
	if !schema.HasField(CqIDColumn.Name) {
		cqID := hashRecord(record)
		newFields = append(newFields, CqIDColumn.ToArrowField())
		newColumns = append(newColumns, cqID)
	}
	if !schema.HasField(CqParentIDColumn.Name) {
		cqParentID := nullUUIDsForRecord(nRows)
		newFields = append(newFields, CqParentIDColumn.ToArrowField())
		newColumns = append(newColumns, cqParentID)
	}

	clientIDArray := StringArrayFromValue(cqClientIDValue, nRows)
	if !schema.HasField(CqClientIDColumn.Name) {
		newFields = append(newFields, CqClientIDColumn.ToArrowField())
		newColumns = append(newColumns, clientIDArray)
	} else {
		record, err = ReplaceFieldInRecord(record, CqClientIDColumn.Name, clientIDArray)
		if err != nil {
			return nil, err
		}
	}

	allFields := append(schema.Fields(), newFields...)
	allColumns := append(record.Columns(), newColumns...)
	metadata := schema.Metadata()
	newSchema := arrow.NewSchema(allFields, &metadata)
	return array.NewRecordBatch(newSchema, allColumns, int64(nRows)), nil
}
