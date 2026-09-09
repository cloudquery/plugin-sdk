package plugin

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/require"
)

func TestAppendFromStringPreservesInt64Precision(t *testing.T) {
	for _, tc := range []struct {
		name string
		dt   arrow.DataType
		str  string
	}{
		{name: "int64", dt: arrow.PrimitiveTypes.Int64, str: "-8717895732742165505"},
		{name: "list_of_int64", dt: arrow.ListOf(arrow.PrimitiveTypes.Int64), str: "[-8717895732742165505]"},
		{name: "large_list_of_int64", dt: arrow.LargeListOf(arrow.PrimitiveTypes.Int64), str: "[-8717895732742165505]"},
		{name: "struct_with_int64", dt: arrow.StructOf(arrow.Field{Name: "v", Type: arrow.PrimitiveTypes.Int64, Nullable: true}), str: `{"v":-8717895732742165505}`},
		{name: "uint64", dt: arrow.PrimitiveTypes.Uint64, str: "18428615660272232523"},
		{name: "list_of_uint64", dt: arrow.ListOf(arrow.PrimitiveTypes.Uint64), str: "[18428615660272232523]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bldr := array.NewBuilder(memory.DefaultAllocator, tc.dt)
			defer bldr.Release()

			require.NoError(t, appendFromString(bldr, tc.str))

			arr := bldr.NewArray()
			defer arr.Release()

			require.Equal(t, tc.str, arr.ValueStr(0))
		})
	}
}
