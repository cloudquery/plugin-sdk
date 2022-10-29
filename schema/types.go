package schema

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/cqtypes"
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
	res := make(cqtypes.CQTypes, len(values))

	for i, v := range values {
		if v == nil {
			values[i] = nil
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
			return nil, err
		}
		res[i] = t
	}
	return res, nil
}

// func (c CQTypes) Equal(other CQTypes) bool {
// 	if other == nil {
// 		return false
// 	}
// 	if len(c) != len(other) {
// 		return false
// 	}
// 	for i := range c {
// 		if c[i] == nil {
// 			if other[i] != nil {
// 				return false
// 			}
// 		} else {
// 			if !c[i].Equal(other[i]) {
// 				return false
// 			}
// 		}

// 	}
// 	return true
// }
