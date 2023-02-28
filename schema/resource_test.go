package schema

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var calculateUniqueValueTestCases = []struct {
	Name          string
	Resource      any
	ExpectedValue *UUID
	Table         *Table
	ConsistentID  bool
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
		ConsistentID:  true,
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": map[string]any{
				"test":       "test",
				"testValInt": 1,
			},
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x3f, 0xb7, 0xc3, 0xac, 0x6b, 0x9a, 0x54, 0xa2, 0xb4, 0xd5, 0xd, 0x73, 0x11, 0xc0, 0x76, 0x6c}, Status: 0x2},
		ConsistentID:  true,
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": map[string]any{
				"testValInt": 1,
				"test":       "test",
			},
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x3f, 0xb7, 0xc3, 0xac, 0x6b, 0x9a, 0x54, 0xa2, 0xb4, 0xd5, 0xd, 0x73, 0x11, 0xc0, 0x76, 0x6c}, Status: 0x2},
		ConsistentID:  true,
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": []string{"test", "test2", "test3"},
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x28, 0xac, 0xa0, 0xec, 0xb, 0x34, 0x5b, 0xe5, 0xb7, 0x9d, 0xc, 0xae, 0xcc, 0x19, 0xa4, 0xeb}, Status: 0x2},
		ConsistentID:  true,
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": []string{"test3", "test2", "test"},
		},
		ExpectedValue: &UUID{Bytes: [16]uint8{0x30, 0x55, 0x9c, 0x58, 0xfe, 0x8, 0x5f, 0x2a, 0xb9, 0x2f, 0xb3, 0x2d, 0x35, 0x4d, 0x31, 0xca}, Status: 0x2},
		ConsistentID:  true,
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
				},
				{
					Name:     "int_column",
					Type:     TypeInt,
					Resolver: PathResolver("IntValue"),
				},
				{
					Name:     "float_column",
					Type:     TypeFloat,
					Resolver: PathResolver("FloatValue"),
				},
				{
					Name:     "uuid_column",
					Type:     TypeUUID,
					Resolver: PathResolver("UUIDValue"),
				},
				{
					Name:     "string_column",
					Type:     TypeString,
					Resolver: PathResolver("StringValue"),
				},
				{
					Name:     "byte_array_column",
					Type:     TypeByteArray,
					Resolver: PathResolver("ByteArrayValue"),
				},
				{
					Name:     "string_array_column",
					Type:     TypeStringArray,
					Resolver: PathResolver("StringArrayValue"),
				},
				{
					Name:     "int_array_column",
					Type:     TypeIntArray,
					Resolver: PathResolver("IntArrayValue"),
				},
				{
					Name:     "timestamp_column",
					Type:     TypeTimestamp,
					Resolver: PathResolver("TimestampValue"),
				},
				{
					Name:     "json_map_column",
					Type:     TypeJSON,
					Resolver: PathResolver("JSONMapValue"),
				},
				{
					Name:     "json_array_column",
					Type:     TypeJSON,
					Resolver: PathResolver("JSONArrayValue"),
				},
				{
					Name:     "uuid_array_column",
					Type:     TypeUUIDArray,
					Resolver: PathResolver("UUIDArrayValue"),
				},
				{
					Name:     "inet_column",
					Type:     TypeInet,
					Resolver: PathResolver("InetValue"),
				},
				{
					Name:     "inet_array_column",
					Type:     TypeInetArray,
					Resolver: PathResolver("InetArrayValue"),
				},
				{
					Name:     "cidr_column",
					Type:     TypeCIDR,
					Resolver: PathResolver("CidrValue"),
				},
				{
					Name:     "cidr_array_column",
					Type:     TypeCIDRArray,
					Resolver: PathResolver("CidrArrayValue"),
				},
				{
					Name:     "mac_address_column",
					Type:     TypeMacAddr,
					Resolver: PathResolver("MacAddressValue"),
				},
				{
					Name:     "mac_address_array_column",
					Type:     TypeMacAddrArray,
					Resolver: PathResolver("MacAddressArrayValue"),
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
		ExpectedValue: &UUID{Bytes: [16]uint8{0xa3, 0x99, 0x29, 0x81, 0xb3, 0xf3, 0x5f, 0x8c, 0xaf, 0x31, 0x1d, 0x7c, 0xee, 0x7b, 0x1c, 0x58}, Status: 0x2},
		ConsistentID:  true,
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
			err := resource.CalculateUniqueValue(tc.ConsistentID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, tc.ExpectedValue, resource.Get(CqIDColumn.Name))
		})
	}
}
