package schema

import (
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
)

// NewKeysRecord builds the one-row record that identifies a single resource of
// the given table. Its schema is the table's identity columns - the primary key
// components when the table declares any, otherwise the primary keys - in
// Table.IdentityColumns() order, so that two records for the same resource are
// always byte-identical regardless of the iteration order of values.
//
// It is an error for values to omit an identity column or to carry any column
// that is not one. In particular a table whose only primary key is _cq_id is
// keyed on its components, and passing _cq_id is rejected.
func NewKeysRecord(table *Table, values map[string]any) (arrow.RecordBatch, error) {
	identity := table.IdentityColumns()
	if len(identity) == 0 {
		return nil, fmt.Errorf("table %q has no identity columns", table.Name)
	}
	wanted := make(map[string]struct{}, len(identity))
	for _, name := range identity {
		wanted[name] = struct{}{}
	}

	for name := range values {
		if table.Column(name) == nil {
			return nil, fmt.Errorf("column %q does not exist on table %q", name, table.Name)
		}
		if _, ok := wanted[name]; !ok {
			return nil, fmt.Errorf("column %q is not an identity column of table %q (expected %v)", name, table.Name, identity)
		}
	}

	fields := make([]arrow.Field, len(identity))
	scalars := make(scalar.Vector, len(identity))
	for i, name := range identity {
		value, ok := values[name]
		if !ok {
			return nil, fmt.Errorf("missing identity column %q for table %q", name, table.Name)
		}
		col := table.Column(name)
		fields[i] = arrow.Field{Name: col.Name, Type: col.Type, Nullable: true}
		s := scalar.NewScalar(col.Type)
		if err := s.Set(value); err != nil {
			return nil, fmt.Errorf("failed to set identity column %q for table %q: %w", name, table.Name, err)
		}
		scalars[i] = s
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, arrow.NewSchema(fields, nil))
	defer bldr.Release()
	for i, s := range scalars {
		scalar.AppendToBuilder(bldr.Field(i), s)
	}
	return bldr.NewRecordBatch(), nil
}

// KeysString renders a keys record as a canonical "col=val,col=val" string in
// column order. It is used for logging and for coalescing changes that refer to
// the same resource.
func KeysString(keys arrow.RecordBatch) string {
	if keys == nil || keys.NumRows() == 0 {
		return ""
	}

	var sb strings.Builder
	sc := keys.Schema()
	for i, field := range sc.Fields() {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(field.Name)
		sb.WriteByte('=')
		sb.WriteString(columnValueString(field.Type, keys.Column(i)))
	}
	return sb.String()
}

// columnValueString renders the first row of arr through the scalar that matches
// the column type, so that the output matches what NewKeysRecord was given.
func columnValueString(dt arrow.DataType, arr arrow.Array) string {
	s := scalar.NewScalar(dt)
	if err := s.Set(arr.GetOneForMarshal(0)); err != nil {
		// Not every arrow type round-trips through its marshalled form; the array's
		// own rendering is still stable for a given value.
		return arr.ValueStr(0)
	}
	return s.String()
}
