package schema

import (
	"bytes"
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
)

const (
	MetadataUnique         = "cq:extension:unique"
	MetadataPrimaryKey     = "cq:extension:primary_key"
	MetadataConstraintName = "cq:extension:constraint_name"
	MetadataIncremental    = "cq:extension:incremental"

	MetadataTrue             = "true"
	MetadataFalse            = "false"
	MetadataTableName        = "cq:table_name"
	MetadataTableDescription = "cq:table_description"
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

func (s Schemas) Encode() ([][]byte, error) {
	ret := make([][]byte, len(s))
	for i, sc := range s {
		var buf bytes.Buffer
		wr := ipc.NewWriter(&buf, ipc.WithSchema(sc))
		if err := wr.Close(); err != nil {
			return nil, err
		}
		ret[i] = buf.Bytes()
	}
	return ret, nil
}

func NewSchemasFromBytes(b [][]byte) (Schemas, error) {
	ret := make([]*arrow.Schema, len(b))
	for i, buf := range b {
		rdr, err := ipc.NewReader(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret[i] = rdr.Schema()
	}
	return ret, nil
}

func NewTablesFromBytes(b [][]byte) (Tables, error) {
	schemas, err := NewSchemasFromBytes(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode schemas: %w", err)
	}
	return NewTablesFromArrowSchemas(schemas)
}
