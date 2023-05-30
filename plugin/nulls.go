package plugin

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
)

func stripNullsFromLists(records []arrow.Record) {
	for i := range records {
		cols := records[i].Columns()
		for c, col := range cols {
			if col.DataType().ID() != arrow.LIST {
				continue
			}

			list := col.(*array.List)
			bldr := array.NewListBuilder(memory.DefaultAllocator, list.DataType().(*arrow.ListType).Elem())
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
		records[i] = array.NewRecord(records[i].Schema(), cols, records[i].NumRows())
	}
}

type AllowNullFunc func(arrow.DataType) bool

func (f AllowNullFunc) replaceNullsByEmpty(records []arrow.Record) {
	if f == nil {
		return
	}
	for i := range records {
		cols := records[i].Columns()
		for c, col := range records[i].Columns() {
			if col.NullN() == 0 || f(col.DataType()) {
				continue
			}

			builder := array.NewBuilder(memory.DefaultAllocator, records[i].Column(c).DataType())
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
		records[i] = array.NewRecord(records[i].Schema(), cols, records[i].NumRows())
	}
<<<<<<< HEAD:plugins/destination/nulls.go
}
=======
}
>>>>>>> 5ba1713 (wip):plugin/nulls.go
