package schema

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var calculateCQIDPrimaryKeyTestCases = []struct {
	Name              string
	Resource          any
	ExpectedValue     *UUID
	Table             *Table
	DeterministicCQID bool
}{
	{
		Name: "Multiple Identical PK Values",
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
			"PathResolver2": "test",
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x50, 0x3c, 0x31, 0xc6, 0x9, 0x71, 0x5c, 0x89, 0x83, 0x1b, 0x17, 0x74, 0x9e, 0xf, 0xf5, 0xc7}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "Multiple Identical PK Values- In Different Columns",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "string_column2",
					Type:     TypeString,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "string_column3",
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
			"PathResolver2": "test",
		},
		// This should be a different value than the previous test case ("Multiple Identical Values") because the column names are different
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x42, 0x1d, 0x8a, 0x42, 0x6a, 0xf0, 0x51, 0x2d, 0xb8, 0x49, 0xc7, 0xaf, 0xb1, 0xaf, 0x56, 0xec}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0xa4, 0xe7, 0x20, 0x33, 0xb7, 0x52, 0x58, 0x70, 0x8d, 0x10, 0xe8, 0xa0, 0x54, 0x60, 0x43, 0xec}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x6e, 0xdd, 0xa8, 0x8a, 0x80, 0xab, 0x50, 0x3a, 0x8c, 0xdc, 0xac, 0x91, 0xbb, 0x2e, 0xa4, 0x37}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x6e, 0xdd, 0xa8, 0x8a, 0x80, 0xab, 0x50, 0x3a, 0x8c, 0xdc, 0xac, 0x91, 0xbb, 0x2e, 0xa4, 0x37}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "Singular JSON Map",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "json_column",
					Type:     TypeJSON,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": map[string]any{
				"test":       "test",
				"testValInt": 1,
			},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x59, 0xa7, 0x9a, 0x5c, 0x66, 0xc2, 0x5f, 0xb3, 0xb9, 0x8a, 0x9e, 0x81, 0xb8, 0x9e, 0x3f, 0xc7}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "Singular JSON Map- Values Change order",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "json_column",
					Type:     TypeJSON,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": map[string]any{
				"testValInt": 1,
				"test":       "test",
			},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x59, 0xa7, 0x9a, 0x5c, 0x66, 0xc2, 0x5f, 0xb3, 0xb9, 0x8a, 0x9e, 0x81, 0xb8, 0x9e, 0x3f, 0xc7}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "Singular JSON Array",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "json_column",
					Type:     TypeJSON,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": []string{"test", "test2", "test3"},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x3, 0xf, 0x3e, 0xb6, 0x65, 0x24, 0x51, 0x1c, 0xac, 0xa, 0x91, 0x4c, 0x7, 0xa4, 0x1f, 0x6c}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "Singular JSON Array- Values Changes order- And CQ_ID",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "json_column",
					Type:     TypeJSON,
					Resolver: PathResolver("PathResolver"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": []string{"test3", "test2", "test"},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x5, 0x44, 0x97, 0x89, 0x98, 0x8c, 0x5b, 0xb, 0x81, 0x12, 0xdd, 0x7e, 0x74, 0xae, 0x70, 0x56}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "All CQ Types",
		Table: &Table{
			Name: "test_table",
			Columns: []Column{
				CqIDColumn,
				{
					Name:     "bool_column",
					Type:     TypeBool,
					Resolver: PathResolver("BooleanValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "int_column",
					Type:     TypeInt,
					Resolver: PathResolver("IntValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "float_column",
					Type:     TypeFloat,
					Resolver: PathResolver("FloatValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "uuid_column",
					Type:     TypeUUID,
					Resolver: PathResolver("UUIDValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("StringValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "byte_array_column",
					Type:     TypeByteArray,
					Resolver: PathResolver("ByteArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "string_array_column",
					Type:     TypeStringArray,
					Resolver: PathResolver("StringArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "int_array_column",
					Type:     TypeIntArray,
					Resolver: PathResolver("IntArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "timestamp_column",
					Type:     TypeTimestamp,
					Resolver: PathResolver("TimestampValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "json_map_column",
					Type:     TypeJSON,
					Resolver: PathResolver("JSONMapValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "json_array_column",
					Type:     TypeJSON,
					Resolver: PathResolver("JSONArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "uuid_array_column",
					Type:     TypeUUIDArray,
					Resolver: PathResolver("UUIDArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "inet_column",
					Type:     TypeInet,
					Resolver: PathResolver("InetValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "inet_array_column",
					Type:     TypeInetArray,
					Resolver: PathResolver("InetArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "cidr_column",
					Type:     TypeCIDR,
					Resolver: PathResolver("CidrValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "cidr_array_column",
					Type:     TypeCIDRArray,
					Resolver: PathResolver("CidrArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "mac_address_column",
					Type:     TypeMacAddr,
					Resolver: PathResolver("MacAddressValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
				{
					Name:     "mac_address_array_column",
					Type:     TypeMacAddrArray,
					Resolver: PathResolver("MacAddressArrayValue"),
					CreationOptions: ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		},
		Resource: map[string]any{
			"BooleanValue":     true,
			"IntValue":         1456,
			"FloatValue":       1456.12,
			"UUIDValue":        "14625e33-4c0a-44e6-909c-0f9865c1b0f9",
			"StringValue":      "test",
			"ByteArrayValue":   []byte{'G', 'O', 'L', 'A', 'N', 'G'},
			"StringArrayValue": []string{"test", "test2", "test3"},
			"IntArrayValue":    []int{1, 2, 3},
			"TimestampValue":   time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			"JSONArrayValue":   []string{"test3", "test2", "test"},
			"JSONMapValue": map[string]any{
				"testValInt": 1,
				"test":       "test",
			},
			"UUIDArrayValue":       []string{"14625e33-4c0a-44e6-909c-0f9865c1b0f9", "14525e33-4c0a-44e6-909c-0f9865c1b0f0"},
			"InetValue":            netip.MustParseAddr("192.0.2.1"),
			"InetArrayValue":       []netip.Addr{netip.MustParseAddr("192.0.2.1"), netip.MustParseAddr("192.0.2.1")},
			"CidrValue":            "192.0.2.1/24",
			"CidrArrayValue":       []string{"192.0.2.1/24", "192.0.2.1/16"},
			"MacAddressValue":      "aa:bb:cc:dd:ee:ff",
			"MacAddressArrayValue": []string{"aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x82, 0xaa, 0xa2, 0xc7, 0x75, 0x10, 0x5c, 0x6d, 0xb8, 0xef, 0x18, 0x49, 0x33, 0x99, 0x4a, 0x9d}, Status: 0x2},
		DeterministicCQID: true,
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

func TestCalculateCQIDWithPrimaryKeys(t *testing.T) {
	for _, tc := range calculateCQIDPrimaryKeyTestCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			resource := NewResourceData(tc.Table, nil, tc.Resource)
			resolveColumns(t, resource, tc.Table)
			err := resource.CalculateCQID(tc.DeterministicCQID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, tc.ExpectedValue, resource.Get(CqIDColumn.Name))
		})
	}
}

var calculateCQIDNoPrimaryKeyTestCases = []struct {
	Name              string
	Resource          any
	ExpectedValue     *UUID
	Table             *Table
	DeterministicCQID bool
}{
	{
		Name: "No Primary Keys",
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0xa5, 0x55, 0x12, 0x57, 0x65, 0x41, 0x4e, 0x2a, 0xa4, 0x78, 0x4c, 0x88, 0x54, 0x4d, 0x38, 0x34}, Status: 0x2},
		DeterministicCQID: true,
	},
}

// This test is to ensure that the CQID is not deterministic when there are no primary keys
func TestCalculateCQIDNoPrimaryKeys(t *testing.T) {
	for _, tc := range calculateCQIDNoPrimaryKeyTestCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			resource := NewResourceData(tc.Table, nil, tc.Resource)
			resolveColumns(t, resource, tc.Table)
			err := resource.CalculateCQID(tc.DeterministicCQID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			initialCQID := resource.Get(CqIDColumn.Name).String()

			err = resource.CalculateCQID(tc.DeterministicCQID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.NotEqual(t, initialCQID, resource.Get(CqIDColumn.Name).String())
		})
	}
}
