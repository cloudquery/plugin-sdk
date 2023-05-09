package schema

import (
	"fmt"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/stretchr/testify/require"
)

func TestFieldChange_String(t *testing.T) {
	type testCase struct {
		change   FieldChange
		expected string
	}

	for _, tc := range []testCase{
		{
			change: FieldChange{
				Type:       TableColumnChangeTypeUnknown,
				ColumnName: "name",
				Current: arrow.Field{
					Name: "name",
					Type: new(arrow.BooleanType),
					Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
						PrimaryKey: true,
						Unique:     true,
					}),
				},
				Previous: arrow.Field{
					Name:     "name",
					Type:     new(arrow.BooleanType),
					Nullable: true,
				},
			},
			expected: `? name: nullable(bool) -> name: bool, metadata: ["cq:extension:primary_key": "true", "cq:extension:unique": "true"]`,
		},
		{
			change: FieldChange{
				Type:       TableColumnChangeTypeAdd,
				ColumnName: "name",
				Current: arrow.Field{
					Name: "name",
					Type: new(arrow.BooleanType),
					Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
						PrimaryKey: true,
						Unique:     true,
					}),
				},
			},
			expected: `+ name: bool, metadata: ["cq:extension:primary_key": "true", "cq:extension:unique": "true"]`,
		},
		{
			change: FieldChange{
				Type:       TableColumnChangeTypeUpdate,
				ColumnName: "name",
				Current: arrow.Field{
					Name: "name",
					Type: new(arrow.BooleanType),
					Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
						PrimaryKey: true,
						Unique:     true,
					}),
				},
				Previous: arrow.Field{
					Name:     "name",
					Type:     new(arrow.BooleanType),
					Nullable: true,
				},
			},
			expected: `~ name: nullable(bool) -> name: bool, metadata: ["cq:extension:primary_key": "true", "cq:extension:unique": "true"]`,
		},
		{
			change: FieldChange{
				Type:       TableColumnChangeTypeRemove,
				ColumnName: "name",
				Previous: arrow.Field{
					Name:     "name",
					Type:     new(arrow.BooleanType),
					Nullable: true,
				},
			},
			expected: `- name: nullable(bool)`,
		},
	} {
		require.Equal(t, tc.expected, tc.change.String())
	}
}

func TestFieldChanges_String(t *testing.T) {
	changes := FieldChanges{
		{
			Type:       TableColumnChangeTypeUnknown,
			ColumnName: "unknown",
			Current: arrow.Field{
				Name: "unknown",
				Type: new(arrow.BooleanType),
				Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
					PrimaryKey: true,
					Unique:     true,
				}),
			},
			Previous: arrow.Field{
				Name:     "unknown",
				Type:     new(arrow.BooleanType),
				Nullable: true,
			},
		},
		{
			Type:       TableColumnChangeTypeAdd,
			ColumnName: "add",
			Current: arrow.Field{
				Name: "add",
				Type: new(arrow.BooleanType),
				Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
					PrimaryKey: true,
					Unique:     true,
				}),
			},
		},
		{
			Type:       TableColumnChangeTypeUpdate,
			ColumnName: "update",
			Current: arrow.Field{
				Name: "update",
				Type: new(arrow.BooleanType),
				Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
					PrimaryKey: true,
					Unique:     true,
				}),
			},
			Previous: arrow.Field{
				Name:     "update",
				Type:     new(arrow.BooleanType),
				Nullable: true,
			},
		},
		{
			Type:       TableColumnChangeTypeRemove,
			ColumnName: "remove",
			Current: arrow.Field{
				Name: "remove",
				Type: new(arrow.BooleanType),
				Metadata: NewFieldMetadataFromOptions(MetadataFieldOptions{
					PrimaryKey: true,
					Unique:     true,
				}),
			},
			Previous: arrow.Field{
				Name:     "remove",
				Type:     new(arrow.BooleanType),
				Nullable: true,
			},
		},
	}

	const expected = `? unknown: nullable(bool) -> unknown: bool, metadata: ["cq:extension:primary_key": "true", "cq:extension:unique": "true"]
+ add: bool, metadata: ["cq:extension:primary_key": "true", "cq:extension:unique": "true"]
~ update: nullable(bool) -> update: bool, metadata: ["cq:extension:primary_key": "true", "cq:extension:unique": "true"]
- remove: nullable(bool)`
	require.Equal(t, expected, changes.String())
	require.Equal(t, expected, fmt.Sprintf("%v", changes))
}
