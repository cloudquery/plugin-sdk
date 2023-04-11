package deprecated

import (
	"time"

	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/google/uuid"
)

func GenTestData(table *schema.Table) schema.CQTypes {
	data := make(schema.CQTypes, len(table.Columns))
	for i, c := range table.Columns {
		switch c.Type {
		case schema.TypeBool:
			data[i] = &schema.Bool{
				Bool:   true,
				Status: schema.Present,
			}
		case schema.TypeInt:
			data[i] = &schema.Int8{
				Int:    1,
				Status: schema.Present,
			}
		case schema.TypeFloat:
			data[i] = &schema.Float8{
				Float:  1.1,
				Status: schema.Present,
			}
		case schema.TypeUUID:
			uuidColumn := &schema.UUID{}
			if err := uuidColumn.Set(uuid.NewString()); err != nil {
				panic(err)
			}
			data[i] = uuidColumn
		case schema.TypeString:
			if c.Name == "text_with_null" {
				data[i] = &schema.Text{
					Str:    "AStringWith\x00NullBytes",
					Status: schema.Present,
				}
			} else {
				data[i] = &schema.Text{
					Str:    "test",
					Status: schema.Present,
				}
			}
		case schema.TypeByteArray:
			data[i] = &schema.Bytea{
				Bytes:  []byte{1, 2, 3},
				Status: schema.Present,
			}
		case schema.TypeStringArray:
			if c.Name == "text_array_with_null" {
				data[i] = &schema.TextArray{
					Elements: []schema.Text{{Str: "test1", Status: schema.Present}, {Str: "test2\x00WithNull", Status: schema.Present}},
					Status:   schema.Present,
				}
			} else {
				data[i] = &schema.TextArray{
					Elements: []schema.Text{{Str: "test1", Status: schema.Present}, {Str: "test2", Status: schema.Present}},
					Status:   schema.Present,
				}
			}

		case schema.TypeIntArray:
			data[i] = &schema.Int8Array{
				Elements: []schema.Int8{{Int: 1, Status: schema.Present}, {Int: 2, Status: schema.Present}},
				Status:   schema.Present,
			}
		case schema.TypeTimestamp:
			data[i] = &schema.Timestamptz{
				Time:   time.Now().UTC().Round(time.Second),
				Status: schema.Present,
			}
		case schema.TypeJSON:
			data[i] = &schema.JSON{
				Bytes:  []byte(`{"test": "test"}`),
				Status: schema.Present,
			}
		case schema.TypeUUIDArray:
			uuidArrayColumn := &schema.UUIDArray{}
			if err := uuidArrayColumn.Set([]string{"00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}); err != nil {
				panic(err)
			}
			data[i] = uuidArrayColumn
		case schema.TypeInet:
			inetColumn := &schema.Inet{}
			if err := inetColumn.Set("192.0.2.0/24"); err != nil {
				panic(err)
			}
			data[i] = inetColumn
		case schema.TypeInetArray:
			inetArrayColumn := &schema.InetArray{}
			if err := inetArrayColumn.Set([]string{"192.0.2.1/24", "192.0.2.1/24"}); err != nil {
				panic(err)
			}
			data[i] = inetArrayColumn
		case schema.TypeCIDR:
			cidrColumn := &schema.CIDR{}
			if err := cidrColumn.Set("192.0.2.1"); err != nil {
				panic(err)
			}
			data[i] = cidrColumn
		case schema.TypeCIDRArray:
			cidrArrayColumn := &schema.CIDRArray{}
			if err := cidrArrayColumn.Set([]string{"192.0.2.1", "192.0.2.1"}); err != nil {
				panic(err)
			}
			data[i] = cidrArrayColumn
		case schema.TypeMacAddr:
			macaddrColumn := &schema.Macaddr{}
			if err := macaddrColumn.Set("aa:bb:cc:dd:ee:ff"); err != nil {
				panic(err)
			}
			data[i] = macaddrColumn
		case schema.TypeMacAddrArray:
			macaddrArrayColumn := &schema.MacaddrArray{}
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
