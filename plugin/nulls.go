package plugin

import (
	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
)

func stripNullsFromLists(list array.ListLike) array.ListLike {
	// TODO: handle Arrow maps separately if required
	bldr := array.NewBuilder(memory.DefaultAllocator, list.DataType()).(array.ListLikeBuilder)
	for j := 0; j < list.Len(); j++ {
		if list.IsNull(j) {
			bldr.AppendNull()
			continue
		}
		bldr.Append(true)
		vBldr := bldr.ValueBuilder()
		from, to := list.ValueOffsets(j)
		slc := array.NewSlice(list.ListValues(), from, to)
		for k := 0; k < int(to-from); k++ {
			if slc.IsNull(k) {
				continue
			}
			err := vBldr.AppendValueFromString(slc.ValueStr(k))
			if err != nil {
				panic(err)
			}
		}
	}

	return bldr.NewArray().(array.ListLike)
}

type AllowNullFunc func(arrow.DataType) bool

func (s *WriterTestSuite) replaceNullsByEmpty(arr arrow.Array) arrow.Array {
	if s.allowNull == nil {
		return arr
	}

	if !s.allowNull(arr.DataType()) && arr.NullN() > 0 {
		builder := array.NewBuilder(memory.DefaultAllocator, arr.DataType())
		for j := 0; j < arr.Len(); j++ {
			if arr.IsNull(j) {
				builder.AppendEmptyValue()
				continue
			}

			if err := builder.AppendValueFromString(arr.ValueStr(j)); err != nil {
				panic(err)
			}
		}

		arr = builder.NewArray()
	}

	// we need to process the nested arrays, too
	return s.replaceNullsByEmptyNestedArray(arr)
}

func (s *WriterTestSuite) replaceNullsByEmptyNestedArray(arr arrow.Array) arrow.Array {
	if s.allowNull == nil {
		return arr
	}

	switch arr := arr.(type) {
	case array.ListLike: // TODO: handle Arrow maps separately if required
		values := s.handleNullsArray(arr.ListValues())
		return array.MakeFromData(
			array.NewData(arr.DataType(), arr.Len(),
				arr.Data().Buffers(),
				[]arrow.ArrayData{values.Data()},
				arr.NullN(), arr.Data().Offset(),
			),
		)
	case *array.Struct:
		children := make([]arrow.ArrayData, arr.NumField())
		for i := 0; i < arr.NumField(); i++ {
			children[i] = s.handleNullsArray(arr.Field(i)).Data()
		}
		return array.MakeFromData(
			array.NewData(arr.DataType(), arr.Len(),
				arr.Data().Buffers(),
				children,
				arr.NullN(), arr.Data().Offset(),
			),
		)
	default:
		return arr
	}
}

func (s *WriterTestSuite) handleNulls(record arrow.Record) arrow.Record {
	cols := record.Columns()
	for c, col := range cols {
		cols[c] = s.handleNullsArray(col)
	}
	return array.NewRecord(record.Schema(), cols, record.NumRows())
}

func (s *WriterTestSuite) handleNullsArray(arr arrow.Array) arrow.Array {
	if list, ok := arr.(array.ListLike); ok && s.ignoreNullsInLists {
		arr = stripNullsFromLists(list) // TODO: handle Arrow maps separately if required
	}

	return s.replaceNullsByEmpty(arr)
}
