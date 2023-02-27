package schema

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var calculateUniqueValueTestCases = []struct {
	Name          string
	Resource      any
	ExpectedValue *UUID
	Table         *Table
}{
	{
		Name: "Nil Value",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name: "string_column",
					Type: TypeString,
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x86, 0x60, 0x5d, 0xfd, 0x6c, 0xc0, 0x5b, 0x2d, 0x88, 0xa6, 0xf4, 0x9b, 0x65, 0x13, 0x96, 0x22}, Status: 0x2},
	}, {
		Name: "Nil Values",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name: "string_column",
					Type: TypeString,
				},
				{
					Name: "string_column2",
					Type: TypeString,
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0xc4, 0xf7, 0xd1, 0xeb, 0x1b, 0xc, 0x54, 0x11, 0x9e, 0x5d, 0xe6, 0x66, 0x79, 0x7f, 0x85, 0xa9}, Status: 0x2},
	},
	{
		Name: "Singular Value",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
				},
				{
					Name: "string_column2",
					Type: TypeString,
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x56, 0xb8, 0xd, 0xe7, 0x7f, 0xd7, 0x54, 0x93, 0x86, 0x5, 0x64, 0xf, 0x48, 0xec, 0xfb, 0x6}, Status: 0x2},
	},
	{
		Name: "Change Column Order from Singular Value",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name: "string_column2",
					Type: TypeString,
				},
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x56, 0xb8, 0xd, 0xe7, 0x7f, 0xd7, 0x54, 0x93, 0x86, 0x5, 0x64, 0xf, 0x48, 0xec, 0xfb, 0x6}, Status: 0x2},
	},
	{
		Name: "Multiple Values",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
				},
				{
					Name:     "string_column2",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver2"),
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x8f, 0x46, 0x30, 0x99, 0x8d, 0xaf, 0x5a, 0xe6, 0xb7, 0xe6, 0xde, 0x28, 0xe0, 0x37, 0xfd, 0xe}, Status: 0x2},
	},
	{
		Name: "Change Order From Multiple Values",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column2",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver2"),
				},
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x8f, 0x46, 0x30, 0x99, 0x8d, 0xaf, 0x5a, 0xe6, 0xb7, 0xe6, 0xde, 0x28, 0xe0, 0x37, 0xfd, 0xe}, Status: 0x2},
	},
	{
		Name: "Singular Primary Key",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "string_column2",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver2"),
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x51, 0x1e, 0x28, 0xc2, 0x14, 0x6d, 0x5a, 0x58, 0x8e, 0xda, 0xd3, 0x12, 0x50, 0x97, 0xcc, 0xf0}, Status: 0x2},
	},
	{
		Name: "Multiple Primary Keys",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "string_column2",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver2"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x9b, 0x8e, 0x13, 0x41, 0xa7, 0x23, 0x5c, 0x35, 0x81, 0x8b, 0x58, 0x2e, 0xed, 0x43, 0xfc, 0xb}, Status: 0x2},
	},
	{
		Name: "Change Order of Primary Keys",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column2",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver2"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"PathResolver":  "test",
			"PathResolver2": "test2",
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x9b, 0x8e, 0x13, 0x41, 0xa7, 0x23, 0x5c, 0x35, 0x81, 0x8b, 0x58, 0x2e, 0xed, 0x43, 0xfc, 0xb}, Status: 0x2},
	},
}

func resolveColumns(t *testing.T, resource *Resource, table *Table) {
	for _, column := range table.Columns {
		if column.Resolver == nil {
			continue
		}
		err := column.Resolver(context.Background(), nil, resource, column)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	}
}

func TestCalculateUniqueValue(t *testing.T) {
	for _, tc := range calculateUniqueValueTestCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			resource := NewResourceData(tc.Table, nil, tc.Resource)
			resolveColumns(t, resource, tc.Table)
			err := resource.CalculateUniqueValue()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, tc.ExpectedValue, resource.Get(CqIDColumn.Name))
		})
	}
}
