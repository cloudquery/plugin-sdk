package schema

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v3/types"
)

var (
	TypeBool   = arrow.FixedWidthTypes.Boolean
	TypeInt    = arrow.PrimitiveTypes.Int64
	TypeFloat  = arrow.PrimitiveTypes.Float64
	TypeString = arrow.BinaryTypes.String
	TypeText   = arrow.BinaryTypes.String

	TypeUUID    = types.ExtensionTypes.UUID
	TypeJSON    = types.ExtensionTypes.JSON
	TypeInet    = types.ExtensionTypes.Inet
	TypeCIDR    = TypeInet
	TypeMacAddr = types.ExtensionTypes.MAC

	TypeTimestamp = arrow.FixedWidthTypes.Timestamp_us

	TypeByteArray   = arrow.BinaryTypes.Binary
	TypeStringArray = arrow.ListOf(arrow.BinaryTypes.String)
	TypeTextArray   = arrow.ListOf(arrow.BinaryTypes.String)
	TypeIntArray    = arrow.ListOf(arrow.PrimitiveTypes.Int64)

	TypeUUIDArray    = arrow.ListOf(types.ExtensionTypes.UUID)
	TypeInetArray    = arrow.ListOf(types.ExtensionTypes.Inet)
	TypeCIDRArray    = TypeInetArray
	TypeMacAddrArray = arrow.ListOf(types.ExtensionTypes.MAC)
)
