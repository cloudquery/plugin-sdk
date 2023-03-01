package schema

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var calculateCQIDTestCases = []struct {
	Name              string
	Resource          any
	ExpectedValue     *UUID
	Table             *Table
	DeterministicCQID bool
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0xdb, 0x8e, 0xfb, 0x63, 0x15, 0xa9, 0x5f, 0x82, 0x98, 0x1a, 0x74, 0x98, 0xbf, 0x72, 0x38, 0x67}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x33, 0x49, 0x3d, 0xab, 0xce, 0xd6, 0x51, 0xcf, 0x81, 0x7e, 0x22, 0x62, 0x1, 0x93, 0x80, 0x9c}, Status: 0x2},
		DeterministicCQID: true,
	},
	{
		Name: "Multiple Identical Values",
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
			"PathResolver2": "test",
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x8a, 0x46, 0xd3, 0x40, 0x6a, 0xe3, 0x58, 0xd0, 0xa5, 0x92, 0x2c, 0xf3, 0xce, 0x9a, 0x6b, 0x2e}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0xd6, 0x6c, 0x68, 0xa4, 0x96, 0x3a, 0x53, 0xd3, 0x84, 0x7e, 0xab, 0xf, 0xfb, 0x8f, 0x1, 0x43}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0xd6, 0x6c, 0x68, 0xa4, 0x96, 0x3a, 0x53, 0xd3, 0x84, 0x7e, 0xab, 0xf, 0xfb, 0x8f, 0x1, 0x43}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x77, 0x70, 0x8a, 0x6e, 0x9d, 0xe1, 0x50, 0xaf, 0xb4, 0x1, 0x96, 0xa, 0x29, 0xc6, 0x40, 0x2a}, Status: 0x2},
		DeterministicCQID: true,
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x77, 0x70, 0x8a, 0x6e, 0x9d, 0xe1, 0x50, 0xaf, 0xb4, 0x1, 0x96, 0xa, 0x29, 0xc6, 0x40, 0x2a}, Status: 0x2},
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": map[string]any{
				"test":       "test",
				"testValInt": 1,
			},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x9e, 0x21, 0xba, 0xc4, 0xb7, 0xa2, 0x52, 0x38, 0x88, 0x36, 0xc8, 0x80, 0xa9, 0x44, 0xf2, 0xa8}, Status: 0x2},
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": map[string]any{
				"testValInt": 1,
				"test":       "test",
			},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x9e, 0x21, 0xba, 0xc4, 0xb7, 0xa2, 0x52, 0x38, 0x88, 0x36, 0xc8, 0x80, 0xa9, 0x44, 0xf2, 0xa8}, Status: 0x2},
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": []string{"test", "test2", "test3"},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x36, 0xee, 0xc3, 0x12, 0x34, 0x3, 0x57, 0xed, 0xbc, 0xe7, 0xa0, 0xe7, 0x64, 0xb8, 0x37, 0x9c}, Status: 0x2},
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
				},
			},
		},
		Resource: map[string]any{
			"PathResolver": []string{"test3", "test2", "test"},
		},
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x41, 0x4d, 0x19, 0x5c, 0xc9, 0x4, 0x56, 0x71, 0xb5, 0x33, 0x62, 0xcd, 0xae, 0x41, 0x41, 0x83}, Status: 0x2},
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
		ExpectedValue:     &UUID{Bytes: [16]uint8{0x28, 0x22, 0xbd, 0xf8, 0xce, 0xc5, 0x51, 0x8, 0xab, 0xab, 0xae, 0x6b, 0x5c, 0xa1, 0x8b, 0x44}, Status: 0x2},
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

func TestCalculateCQID(t *testing.T) {
	for _, tc := range calculateCQIDTestCases {
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
