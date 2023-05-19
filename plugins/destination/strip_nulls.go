package destination

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
)

func stripNullsFromLists(records []arrow.Record) {
	for i := range records {
		cols := make([]arrow.Array, records[i].NumCols())
		for c := range records[i].Columns() {
			if records[i].Column(c).DataType().ID() == arrow.LIST {
				list := records[i].Column(c).(*array.List)
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
				continue
			}
			cols[c] = records[i].Column(c)
		}
		records[i] = array.NewRecord(records[i].Schema(), cols, records[i].NumRows())
	}
}
