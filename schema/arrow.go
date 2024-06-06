package schema

import (
	"github.com/apache/arrow/go/v16/arrow"
)

const (
	MetadataUnique              = "cq:extension:unique"
	MetadataPrimaryKey          = "cq:extension:primary_key"
	MetadataPrimaryKeyComponent = "cq:extension:primary_key_component"
	MetadataConstraintName      = "cq:extension:constraint_name"
	MetadataIncremental         = "cq:extension:incremental"

	MetadataTrue             = "true"
	MetadataFalse            = "false"
	MetadataTableName        = "cq:table_name"
	MetadataTableDescription = "cq:table_description"
	MetadataTableTitle       = "cq:table_title"
	MetadataTableDependsOn   = "cq:table_depends_on"
	MetadataTableIsPaid      = "cq:table_paid"
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
