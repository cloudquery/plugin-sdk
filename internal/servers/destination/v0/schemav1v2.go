package destination

import (
	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/schemav2"
	"github.com/cloudquery/plugin-sdk/v2/types"
)

func TablesV1ToV2(tables []*schema.Table) schemav2.Tables {
	res := make(schemav2.Tables, len(tables))
	for i, t := range tables {
		res[i] = TableV1ToV2(t)
	}
	return res
}

func TableV1ToV2(table *schema.Table) *schemav2.Table {
	return &schemav2.Table{
		Name:          table.Name,
		Description:   table.Description,
		Columns:       ColumnsV1ToV2(table.Columns),
		IgnoreInTests: table.IgnoreInTests,
		IsIncremental: table.IsIncremental,
	}
}

func ColumnsV1ToV2(columns []schema.Column) []schemav2.Column {
	res := make([]schemav2.Column, len(columns))
	for i, c := range columns {
		res[i] = ColumnV1ToV2(c)
	}
	return res
}

func ColumnV1ToV2(column schema.Column) schemav2.Column {
	return schemav2.Column{
		Name:        column.Name,
		Description: column.Description,
		Type:        TypeV1ToV2(column.Type),
		CreationOptions: schemav2.ColumnCreationOptions{
			NotNull:    column.CreationOptions.NotNull,
			Unique:     column.CreationOptions.Unique,
			PrimaryKey: column.CreationOptions.PrimaryKey,
		},
		IgnoreInTests: column.IgnoreInTests,
	}
}

func TypeV1ToV2(dataType schema.ValueType) arrow.DataType {
	var typ arrow.DataType

	switch dataType {
	case schema.TypeBool:
		return arrow.FixedWidthTypes.Boolean
	case schema.TypeInt:
		return arrow.PrimitiveTypes.Int64
	case schema.TypeFloat:
		return arrow.PrimitiveTypes.Float64
	case schema.TypeUUID:
		return types.ExtensionTypes.UUID
	case schema.TypeString:
		return arrow.BinaryTypes.String
	case schema.TypeByteArray:
		return arrow.BinaryTypes.Binary
	case schema.TypeStringArray:
		return arrow.ListOf(arrow.BinaryTypes.String)
	case schema.TypeIntArray:
		return arrow.ListOf(arrow.PrimitiveTypes.Int64)
	case schema.TypeTimestamp:
		return arrow.FixedWidthTypes.Timestamp_us
	case schema.TypeJSON:
		return types.ExtensionTypes.JSON
	case schema.TypeUUIDArray:
		return arrow.ListOf(types.ExtensionTypes.UUID)
	case schema.TypeInet:
		return types.ExtensionTypes.Inet
	case schema.TypeInetArray:
		return arrow.ListOf(types.ExtensionTypes.Inet)
	case schema.TypeCIDR:
		return types.ExtensionTypes.Inet
	case schema.TypeCIDRArray:
		return arrow.ListOf(types.ExtensionTypes.Inet)
	case schema.TypeMacAddr:
		return types.ExtensionTypes.Mac
	case schema.TypeMacAddrArray:
		return arrow.ListOf(types.ExtensionTypes.Mac)
	default:
		panic("unknown type " + typ.Name())
	}
}
