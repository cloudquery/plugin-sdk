package message

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// NewDeleteRecordByKeys builds the delete message for the resource identified by
// keys, a one-row record whose columns are a non-empty subset of the table's
// identity columns - the primary key components when the table declares any,
// otherwise the primary keys (see schema.NewKeysRecord).
//
// The where clause is one AND group with an EQ predicate for exactly the columns
// the record carries, so a record holding part of a resource's identity deletes
// everything matching that part. The relations cover every descendant table so
// destinations can cascade the delete.
func NewDeleteRecordByKeys(table *schema.Table, keys arrow.RecordBatch) (*SyncDeleteRecord, error) {
	if keys == nil {
		return nil, fmt.Errorf("no keys record for table %q", table.Name)
	}
	if keys.NumRows() != 1 {
		return nil, fmt.Errorf("keys record for table %q must have exactly one row, got %d", table.Name, keys.NumRows())
	}

	identity := table.IdentityColumns()
	fields := keys.Schema().Fields()
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("keys record for table %q has no columns, expected a subset of %v", table.Name, identity)
	}
	// The record must be a subset of the identity columns, in the table's
	// declaration order, which is what schema.NewKeysRecord produces.
	if !isOrderedSubset(names, identity) {
		return nil, fmt.Errorf("keys record columns %v are not a subset of the identity columns %v of table %q", names, identity, table.Name)
	}

	predicates := make(Predicates, len(names))
	for i, name := range names {
		column := keys.Column(i)
		sc := arrow.NewSchema([]arrow.Field{fields[i]}, nil)
		predicates[i] = Predicate{
			Operator: pb.Predicate_EQ.String(),
			Column:   name,
			Record:   array.NewRecordBatch(sc, []arrow.Array{column}, 1),
		}
	}

	return &SyncDeleteRecord{
		DeleteRecord: DeleteRecord{
			TableName:      table.Name,
			TableRelations: descendantRelations(table),
			WhereClause: PredicateGroups{{
				GroupingType: pb.PredicatesGroup_AND.String(),
				Predicates:   predicates,
			}},
		},
	}, nil
}

// isOrderedSubset reports whether every name appears in whole, in the same order
// and without repeats.
func isOrderedSubset(names, whole []string) bool {
	i := 0
	for _, name := range names {
		for i < len(whole) && whole[i] != name {
			i++
		}
		if i == len(whole) {
			return false
		}
		i++
	}
	return true
}

// descendantRelations lists every table below the given one, each paired with its
// immediate parent.
func descendantRelations(table *schema.Table) TableRelations {
	relations := make(TableRelations, 0, len(table.Relations))
	for _, child := range table.Relations {
		relations = append(relations, TableRelation{
			TableName:   child.Name,
			ParentTable: table.Name,
		})
		relations = append(relations, descendantRelations(child)...)
	}
	if len(relations) == 0 {
		return nil
	}
	return relations
}
