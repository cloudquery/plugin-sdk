package testdata

import (
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/google/uuid"
)

// TestTable returns a table with columns of all type. useful for destination testing purposes
func TestTable(name string) *schema.Table {
	return &schema.Table{
		Name:        name,
		Description: "Test table",
		Columns: schema.ColumnList{
			schema.CqIDColumn,
			schema.CqParentIDColumn,
			{
				Name: "bool",
				Type: schema.TypeBool,
			},
			{
				Name: "int",
				Type: schema.TypeInt,
			},
			{
				Name: "float",
				Type: schema.TypeFloat,
			},
			{
				Name:            "uuid",
				Type:            schema.TypeUUID,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name: "text",
				Type: schema.TypeString,
			},
			{
				Name: "bytea",
				Type: schema.TypeByteArray,
			},
			{
				Name: "text_array",
				Type: schema.TypeStringArray,
			},
			{
				Name: "int_array",
				Type: schema.TypeIntArray,
			},
			{
				Name: "timestamp",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "json",
				Type: schema.TypeJSON,
			},
			{
				Name: "uuid_array",
				Type: schema.TypeUUIDArray,
			},
			{
				Name: "inet",
				Type: schema.TypeInet,
			},
			{
				Name: "inet_array",
				Type: schema.TypeInetArray,
			},
			{
				Name: "cidr",
				Type: schema.TypeCIDR,
			},
			{
				Name: "cidr_array",
				Type: schema.TypeCIDRArray,
			},
			{
				Name: "macaddr",
				Type: schema.TypeMacAddr,
			},
			{
				Name: "macaddr_array",
				Type: schema.TypeMacAddrArray,
			},
		},
	}
}

func TestData() schema.CQTypes {
	cqID := &schema.UUID{}
	if err := cqID.Set(uuid.New()); err != nil {
		panic(err)
	}
	cqParentID := &schema.UUID{}
	if err := cqParentID.Set("00000000-0000-0000-0000-000000000003"); err != nil {
		panic(err)
	}
	boolColumn := &schema.Bool{
		Bool:   true,
		Status: schema.Present,
	}
	intColumn := &schema.Int8{
		Int:    1,
		Status: schema.Present,
	}
	floatColumn := &schema.Float8{
		Float:  1.1,
		Status: schema.Present,
	}
	uuidColumn := &schema.UUID{}
	if err := uuidColumn.Set("00000000-0000-0000-0000-000000000001"); err != nil {
		panic(err)
	}
	textColumn := &schema.Text{
		Str:    "test",
		Status: schema.Present,
	}
	byteaColumn := &schema.Bytea{
		Bytes:  []byte{1, 2, 3},
		Status: schema.Present,
	}
	textArrayColumn := &schema.TextArray{}
	if err := textArrayColumn.Set([]string{"test1", "test2"}); err != nil {
		panic(err)
	}
	intArrayColumn := &schema.Int8Array{}
	if err := intArrayColumn.Set([]int8{1, 2}); err != nil {
		panic(err)
	}
	timestampColumn := &schema.Timestamptz{
		Time:   time.Now(),
		Status: schema.Present,
	}
	jsonColumn := &schema.JSON{
		Bytes:  []byte(`{"test":"test"}`),
		Status: schema.Present,
	}
	uuidArrayColumn := &schema.UUIDArray{}
	if err := uuidArrayColumn.Set([]string{"00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}); err != nil {
		panic(err)
	}
	inetColumn := &schema.Inet{}
	if err := inetColumn.Set("192.0.2.1/24"); err != nil {
		panic(err)
	}
	inetArrayColumn := &schema.InetArray{}
	if err := inetArrayColumn.Set([]string{"192.0.2.1/24", "192.0.2.1/24"}); err != nil {
		panic(err)
	}
	cidrColumn := &schema.CIDR{}
	if err := cidrColumn.Set("192.0.2.1"); err != nil {
		panic(err)
	}
	cidrArrayColumn := &schema.CIDRArray{}
	if err := cidrArrayColumn.Set([]string{"192.0.2.1", "192.0.2.1"}); err != nil {
		panic(err)
	}
	macaddrColumn := &schema.Macaddr{}
	if err := macaddrColumn.Set("aa:bb:cc:dd:ee:ff"); err != nil {
		panic(err)
	}
	macaddrArrayColumn := &schema.MacaddrArray{}
	if err := macaddrArrayColumn.Set([]string{"aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"}); err != nil {
		panic(err)
	}

	data := schema.CQTypes{
		cqID,
		cqParentID,
		boolColumn,
		intColumn,
		floatColumn,
		uuidColumn,
		textColumn,
		byteaColumn,
		textArrayColumn,
		intArrayColumn,
		timestampColumn,
		jsonColumn,
		uuidArrayColumn,
		inetColumn,
		inetArrayColumn,
		cidrColumn,
		cidrArrayColumn,
		macaddrColumn,
		macaddrArrayColumn,
	}

	return data
}
