package schema

import (
	"time"

	"github.com/google/uuid"
)

func TestSourceTable(name string) *Table {
	return &Table{
		Name:        name,
		Description: "Test table",
		Columns: ColumnList{
			CqIDColumn,
			CqParentIDColumn,
			{
				Name: "bool",
				Type: TypeBool,
			},
			{
				Name: "int",
				Type: TypeInt,
			},
			{
				Name: "float",
				Type: TypeFloat,
			},
			{
				Name:            "uuid",
				Type:            TypeUUID,
				CreationOptions: ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name: "text",
				Type: TypeString,
			},
			{
				Name: "text_with_null",
				Type: TypeString,
			},
			{
				Name: "bytea",
				Type: TypeByteArray,
			},
			{
				Name: "text_array",
				Type: TypeStringArray,
			},
			{
				Name: "text_array_with_null",
				Type: TypeStringArray,
			},
			{
				Name: "int_array",
				Type: TypeIntArray,
			},
			{
				Name: "timestamp",
				Type: TypeTimestamp,
			},
			{
				Name: "json",
				Type: TypeJSON,
			},
			{
				Name: "uuid_array",
				Type: TypeUUIDArray,
			},
			{
				Name: "inet",
				Type: TypeInet,
			},
			{
				Name: "inet_array",
				Type: TypeInetArray,
			},
			{
				Name: "cidr",
				Type: TypeCIDR,
			},
			{
				Name: "cidr_array",
				Type: TypeCIDRArray,
			},
			{
				Name: "macaddr",
				Type: TypeMacAddr,
			},
			{
				Name: "macaddr_array",
				Type: TypeMacAddrArray,
			},
		},
	}
}

// TestTable returns a table with columns of all type. useful for destination testing purposes
func TestTable(name string) *Table {
	sourceTable := TestSourceTable(name)
	sourceTable.Columns = append(ColumnList{
		CqSourceNameColumn,
		CqSyncTimeColumn,
	}, sourceTable.Columns...)
	return sourceTable
}

func GenTestData(table *Table) CQTypes {
	data := make(CQTypes, len(table.Columns))
	for i, c := range table.Columns {
		switch c.Type {
		case TypeBool:
			data[i] = &Bool{
				Bool:   true,
				Status: Present,
			}
		case TypeInt:
			data[i] = &Int8{
				Int:    1,
				Status: Present,
			}
		case TypeFloat:
			data[i] = &Float8{
				Float:  1.1,
				Status: Present,
			}
		case TypeUUID:
			uuidColumn := &UUID{}
			if err := uuidColumn.Set(uuid.NewString()); err != nil {
				panic(err)
			}
			data[i] = uuidColumn
		case TypeString:
			if c.Name == "text_with_null" {
				data[i] = &Text{
					Str:    "AStringWith\x00NullBytes",
					Status: Present,
				}
			} else {
				data[i] = &Text{
					Str:    "test",
					Status: Present,
				}
			}
		case TypeByteArray:
			data[i] = &Bytea{
				Bytes:  []byte{1, 2, 3},
				Status: Present,
			}
		case TypeStringArray:
			if c.Name == "text_array_with_null" {
				data[i] = &TextArray{
					Elements: []Text{{Str: "test1", Status: Present}, {Str: "test2\x00WithNull", Status: Present}},
					Status:   Present,
				}
			} else {
				data[i] = &TextArray{
					Elements: []Text{{Str: "test1", Status: Present}, {Str: "test2", Status: Present}},
					Status:   Present,
				}
			}

		case TypeIntArray:
			data[i] = &Int8Array{
				Elements: []Int8{{Int: 1, Status: Present}, {Int: 2, Status: Present}},
				Status:   Present,
			}
		case TypeTimestamp:
			data[i] = &Timestamptz{
				Time:   time.Now().UTC().Round(time.Second),
				Status: Present,
			}
		case TypeJSON:
			data[i] = &JSON{
				Bytes:  []byte(`{"test": "test"}`),
				Status: Present,
			}
		case TypeUUIDArray:
			uuidArrayColumn := &UUIDArray{}
			if err := uuidArrayColumn.Set([]string{"00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}); err != nil {
				panic(err)
			}
			data[i] = uuidArrayColumn
		case TypeInet:
			inetColumn := &Inet{}
			if err := inetColumn.Set("192.0.2.1/24"); err != nil {
				panic(err)
			}
			data[i] = inetColumn
		case TypeInetArray:
			inetArrayColumn := &InetArray{}
			if err := inetArrayColumn.Set([]string{"192.0.2.1/24", "192.0.2.1/24"}); err != nil {
				panic(err)
			}
			data[i] = inetArrayColumn
		case TypeCIDR:
			cidrColumn := &CIDR{}
			if err := cidrColumn.Set("192.0.2.1"); err != nil {
				panic(err)
			}
			data[i] = cidrColumn
		case TypeCIDRArray:
			cidrArrayColumn := &CIDRArray{}
			if err := cidrArrayColumn.Set([]string{"192.0.2.1", "192.0.2.1"}); err != nil {
				panic(err)
			}
			data[i] = cidrArrayColumn
		case TypeMacAddr:
			macaddrColumn := &Macaddr{}
			if err := macaddrColumn.Set("aa:bb:cc:dd:ee:ff"); err != nil {
				panic(err)
			}
			data[i] = macaddrColumn
		case TypeMacAddrArray:
			macaddrArrayColumn := &MacaddrArray{}
			if err := macaddrArrayColumn.Set([]string{"aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"}); err != nil {
				panic(err)
			}
			data[i] = macaddrArrayColumn
		default:
			panic("unknown type" + c.Type.String())
		}
	}
	return data
}
