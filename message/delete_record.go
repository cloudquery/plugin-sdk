package message

import (
	"fmt"
	"slices"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// NewDeleteRecordByKeys builds the delete message for the single resource
// identified by keys, a one-row record whose columns are the table's primary
// keys (see schema.NewKeysRecord).
//
// The where clause is one AND group with an EQ predicate per primary key column,
// and the relations cover every descendant table so destinations can cascade the
// delete.
func NewDeleteRecordByKeys(table *schema.Table, keys arrow.RecordBatch) (*SyncDeleteRecord, error) {
	if keys == nil {
		return nil, fmt.Errorf("no keys record for table %q", table.Name)
	}
	if keys.NumRows() != 1 {
		return nil, fmt.Errorf("keys record for table %q must have exactly one row, got %d", table.Name, keys.NumRows())
	}

	pks := table.PrimaryKeys()
	fields := keys.Schema().Fields()
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}
	if !slices.Equal(names, pks) {
		return nil, fmt.Errorf("keys record columns %v do not match the primary keys %v of table %q", names, pks, table.Name)
	}

	predicates := make(Predicates, len(pks))
	for i, name := range pks {
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
