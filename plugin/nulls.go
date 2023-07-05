package plugin

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
)

func stripNullsFromLists(record arrow.Record) arrow.Record {
	cols := record.Columns()
	for c, col := range cols {
		list, ok := col.(array.ListLike)
		if !ok {
			continue
		}
		if _, ok := list.(*array.Map); ok {
			// maps also correspond to array.ListLike
			continue
		}

		bldr := array.NewListBuilder(memory.DefaultAllocator, list.DataType().(arrow.ListLikeType).Elem())
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
		cols[c] = bldr.NewArray()
	}
	return array.NewRecord(record.Schema(), cols, record.NumRows())
}

type AllowNullFunc func(arrow.DataType) bool

func (f AllowNullFunc) replaceNullsByEmpty(record arrow.Record) arrow.Record {
	if f == nil {
		return record
	}

	cols := record.Columns()
	for c, col := range cols {
		if col.NullN() == 0 || f(col.DataType()) {
			continue
		}

		builder := array.NewBuilder(memory.DefaultAllocator, col.DataType())
		for j := 0; j < col.Len(); j++ {
			if col.IsNull(j) {
				builder.AppendEmptyValue()
				continue
			}

			if err := builder.AppendValueFromString(col.ValueStr(j)); err != nil {
				panic(err)
			}
		}
		cols[c] = builder.NewArray()
	}
	return array.NewRecord(record.Schema(), cols, record.NumRows())
}

func (s *WriterTestSuite) handleNulls(record arrow.Record) arrow.Record {
	if s.ignoreNullsInLists {
		record = stripNullsFromLists(record)
	}
	return s.allowNull.replaceNullsByEmpty(record)
}
