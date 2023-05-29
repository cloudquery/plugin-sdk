package scalar

import (
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v3/types"
)

// Scalar represents a single value of a specific DataType as opposed to
// an array.
//
// Scalars are useful for passing single value inputs to compute functions
// (not yet implemented) or for representing individual array elements,
// (with a non-trivial cost though).
type Scalar interface {
	fmt.Stringer
	// IsValid returns true if the value is non-null, otherwise false.
	IsValid() bool
	// The datatype of the value in this scalar
	DataType() arrow.DataType
	// Performs cheap validation checks, returns nil if successful
	// Validate() error
	// tries to set the value of the scalar to the given value
	Set(val any) error

	Get() any
	Equal(other Scalar) bool
}

type Vector []Scalar

func (v Vector) Equal(r Vector) bool {
	if len(v) != len(r) {
		return false
	}
	for i := range v {
		if !v[i].Equal(r[i]) {
			return false
		}
	}
	return true
}

func NewScalar(dt arrow.DataType) Scalar {
	switch dt.ID() {
	case arrow.TIMESTAMP:
		return &Timestamp{}
	case arrow.BINARY:
		return &Binary{}
	case arrow.STRING:
		return &String{}
	case arrow.INT64:
		return &Int{Type: arrow.PrimitiveTypes.Int64}
	case arrow.INT32:
		return &Int{Type: arrow.PrimitiveTypes.Int32}
	case arrow.INT16:
		return &Int{Type: arrow.PrimitiveTypes.Int16}
	case arrow.INT8:
		return &Int{Type: arrow.PrimitiveTypes.Int8}
	case arrow.UINT64:
		return &Uint{Type: arrow.PrimitiveTypes.Uint64}
	case arrow.UINT32:
		return &Uint{Type: arrow.PrimitiveTypes.Uint32}
	case arrow.UINT16:
		return &Uint{Type: arrow.PrimitiveTypes.Uint16}
	case arrow.UINT8:
		return &Uint{Type: arrow.PrimitiveTypes.Uint8}
	case arrow.FLOAT32:
		return &Float{Type: arrow.PrimitiveTypes.Float32}
	case arrow.FLOAT64:
		return &Float{Type: arrow.PrimitiveTypes.Float64}
	case arrow.BOOL:
		return &Bool{}
	case arrow.EXTENSION:
		switch {
		case arrow.TypeEqual(dt, types.ExtensionTypes.UUID):
			return &UUID{}
		case arrow.TypeEqual(dt, types.ExtensionTypes.JSON):
			return &JSON{}
		case arrow.TypeEqual(dt, types.ExtensionTypes.MAC):
			return &Mac{}
		case arrow.TypeEqual(dt, types.ExtensionTypes.Inet):
			return &Inet{}
		default:
			panic("not implemented extension: " + dt.Name())
		}
	case arrow.LIST:
		return &List{
			Type: dt,
		}
	default:
		panic("not implemented: " + dt.Name())
	}
}

func AppendToBuilder(bldr array.Builder, s Scalar) {
	if !s.IsValid() {
		bldr.AppendNull()
		return
	}
	switch s.DataType().ID() {
	case arrow.BINARY:
		bldr.(*array.BinaryBuilder).Append(s.(*Binary).Value)
	case arrow.LARGE_BINARY:
		bldr.(*array.BinaryBuilder).Append(s.(*LargeBinary).Value)
	case arrow.STRING:
		bldr.(*array.StringBuilder).Append(s.(*String).Value)
	case arrow.INT64:
		bldr.(*array.Int64Builder).Append(s.(*Int).Value)
	case arrow.INT32:
		bldr.(*array.Int32Builder).Append(int32(s.(*Int).Value))
	case arrow.INT16:
		bldr.(*array.Int16Builder).Append(int16(s.(*Int).Value))
	case arrow.INT8:
		bldr.(*array.Int8Builder).Append(int8(s.(*Int).Value))
	case arrow.UINT64:
		bldr.(*array.Uint64Builder).Append(s.(*Uint).Value)
	case arrow.UINT32:
		bldr.(*array.Uint32Builder).Append(uint32(s.(*Uint).Value))
	case arrow.UINT16:
		bldr.(*array.Uint16Builder).Append(uint16(s.(*Uint).Value))
	case arrow.UINT8:
		bldr.(*array.Uint8Builder).Append(uint8(s.(*Uint).Value))
	case arrow.FLOAT32:
		bldr.(*array.Float32Builder).Append(float32(s.(*Float).Value))
	case arrow.FLOAT64:
		bldr.(*array.Float64Builder).Append(s.(*Float).Value)
	case arrow.BOOL:
		bldr.(*array.BooleanBuilder).Append(s.(*Bool).Value)
	case arrow.TIMESTAMP:
		bldr.(*array.TimestampBuilder).Append(arrow.Timestamp(s.(*Timestamp).Value.UnixMicro()))
	case arrow.LIST:
		lb := bldr.(*array.ListBuilder)
		if s.IsValid() {
			lb.Append(true)
			for _, v := range s.(*List).Value {
				AppendToBuilder(lb.ValueBuilder(), v)
			}
		} else {
			lb.AppendNull()
		}
	case arrow.EXTENSION:
		switch {
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.UUID):
			bldr.(*types.UUIDBuilder).Append(s.(*UUID).Value)
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.JSON):
			bldr.(*types.JSONBuilder).AppendBytes(s.(*JSON).Value)
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.MAC):
			bldr.(*types.MACBuilder).Append(s.(*Mac).Value)
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.Inet):
			bldr.(*types.InetBuilder).Append(s.(*Inet).Value)
		default:
			panic("not implemented extension: " + s.DataType().Name())
		}
	default:
		panic("not implemented: " + s.DataType().String())
	}
}

func AppendToRecordBuilder(bldr *array.RecordBuilder, vector Vector) {
	for i, scalar := range vector {
		AppendToBuilder(bldr.Field(i), scalar)
	}
}
