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
// the given table. Its columns are a non-empty subset of the table's identity
// columns - the primary key components when the table declares any, otherwise
// the primary keys - always in Table.IdentityColumns() order, so that two
// records naming the same columns of the same resource are byte-identical
// regardless of the iteration order of values.
//
// A subset is allowed because an event often knows only part of a resource's
// identity: a CloudTrail event carries the ARN, while the table's components may
// also include derived columns such as an input hash. Whatever is supplied has
// to identify the resource well enough for the caller's purpose; this function
// only enforces that it is non-empty and that every column is an identity
// column. In particular a table whose only primary key is _cq_id is keyed on its
// components, and passing _cq_id is rejected.
func NewKeysRecord(table *Table, values map[string]any) (arrow.RecordBatch, error) {
	identity := table.IdentityColumns()
	if len(identity) == 0 {
		return nil, fmt.Errorf("table %q has no identity columns", table.Name)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no values given for table %q, expected at least one of %v", table.Name, identity)
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
			return nil, fmt.Errorf("column %q is not an identity column of table %q (expected a subset of %v)", name, table.Name, identity)
		}
	}

	// Follow the table's declaration order, skipping what was not supplied.
	fields := make([]arrow.Field, 0, len(values))
	scalars := make(scalar.Vector, 0, len(values))
	for _, name := range identity {
		value, ok := values[name]
		if !ok {
			continue
		}
		col := table.Column(name)
		fields = append(fields, arrow.Field{Name: col.Name, Type: col.Type, Nullable: true})
		s := scalar.NewScalar(col.Type)
		if err := s.Set(value); err != nil {
			return nil, fmt.Errorf("failed to set identity column %q for table %q: %w", name, table.Name, err)
		}
		scalars = append(scalars, s)
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
