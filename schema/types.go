package schema

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/cqtypes"
	"github.com/cloudquery/plugin-sdk/helpers"
)

type Status byte

const (
	Undefined Status = iota
	Null
	Present
)

// CQTypesFromValues tries best effort to convert a slice of values to CQTypes
// based on the provided table columns.
func CQTypesFromValues(table *Table, values []interface{}) (cqtypes.CQTypes, error) {
	valuesSlice := helpers.InterfaceSlice(values)
	res := make(cqtypes.CQTypes, len(valuesSlice))

	for i, v := range valuesSlice {
		if v == nil {
			valuesSlice[i] = nil
		}
		var t cqtypes.CQType
		var err error
		switch table.Columns[i].Type {
		case TypeBool:
			t := &cqtypes.Bool{}
			err = t.Set(v)
		case TypeInt:
			t := &cqtypes.Int8{}
			err = t.Set(v)
		case TypeFloat:
			t := &cqtypes.Float8{}
			err = t.Set(v)
		case TypeUUID:
			t := &cqtypes.UUID{}
			err = t.Set(v)
		case TypeString:
			t := &cqtypes.Text{}
			err = t.Set(v)
		case TypeByteArray:
			t := &cqtypes.Bytea{}
			err = t.Set(v)
		case TypeStringArray:
			t := &cqtypes.TextArray{}
			err = t.Set(v)
		case TypeIntArray:
			t := &cqtypes.Int8Array{}
			err = t.Set(v)
		case TypeTimestamp:
			t := &cqtypes.Timestamptz{}
			err = t.Set(v)
		case TypeJSON:
			t := &cqtypes.JSON{}
			err = t.Set(v)
		case TypeUUIDArray:
			t := &cqtypes.UUIDArray{}
			err = t.Set(v)
		case TypeInet:
			t := &cqtypes.Inet{}
			err = t.Set(v)
		case TypeInetArray:
			t := &cqtypes.InetArray{}
			err = t.Set(v)
		case TypeCIDR:
			t := &cqtypes.CIDR{}
			err = t.Set(v)
		case TypeCIDRArray:
			t := &cqtypes.CIDRArray{}
			err = t.Set(v)
		case TypeMacAddr:
			t := &cqtypes.Macaddr{}
			err = t.Set(v)
		case TypeMacAddrArray:
			t := &cqtypes.MacaddrArray{}
			err = t.Set(v)
		default:
			return nil, fmt.Errorf("unsupported type %s", table.Columns[i].Type)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to convert value %v to type %s: %w", v, table.Columns[i].Type, err)
		}
		res[i] = t
	}
	return res, nil
}
