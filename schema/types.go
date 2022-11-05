package schema

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/cqtypes"
	"github.com/cloudquery/plugin-sdk/helpers"
)

func CQTypeFromSchema(value ValueType) cqtypes.CQType {
	switch value {
	case TypeBool:
		return &cqtypes.Bool{}
	case TypeByteArray:
		return &cqtypes.Bytea{}
	case TypeCIDRArray:
		return &cqtypes.CIDRArray{}
	case TypeCIDR:
		return &cqtypes.CIDR{}
	case TypeFloat:
		return &cqtypes.Float8{}
	case TypeInetArray:
		return &cqtypes.InetArray{}
	case TypeInet:
		return &cqtypes.Inet{}
	case TypeIntArray:
		return &cqtypes.Int8Array{}
	case TypeInt:
		return &cqtypes.Int8{}
	case TypeJSON:
		return &cqtypes.JSON{}
	case TypeMacAddrArray:
		return &cqtypes.MacaddrArray{}
	case TypeMacAddr:
		return &cqtypes.Macaddr{}
	case TypeStringArray:
		return &cqtypes.TextArray{}
	case TypeString:
		return &cqtypes.Text{}
	case TypeTimestamp:
		return &cqtypes.Timestamptz{}
	case TypeUUIDArray:
		return &cqtypes.UUIDArray{}
	case TypeUUID:
		return &cqtypes.UUID{}
	default:
		panic(fmt.Sprintf("unsupported type %d", value))
	}
}

// DefaultReverseTransformer tries best effort to convert a slice of values to CQTypes
// based on the provided table columns.
func DefaultReverseTransformer(table *Table, values []interface{}) (cqtypes.CQTypes, error) {
	valuesSlice := helpers.InterfaceSlice(values)
	res := make(cqtypes.CQTypes, len(valuesSlice))

	for i, v := range valuesSlice {
		t := CQTypeFromSchema(table.Columns[i].Type)
		if err := t.Set(v); err != nil {
			return nil, fmt.Errorf("failed to convert value %v to type %s: %w", v, table.Columns[i].Type, err)
		}
		res[i] = t
	}
	return res, nil
}
