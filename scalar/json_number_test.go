package scalar

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/require"
)

func TestNestedIntPrecision(t *testing.T) {
	for _, tc := range []struct {
		name string
		dt   arrow.DataType
		str  string
	}{
		{name: "struct_int64", dt: arrow.StructOf(arrow.Field{Name: "v", Type: arrow.PrimitiveTypes.Int64, Nullable: true}), str: `{"v":-8717895732742165505}`},
		{name: "struct_uint64", dt: arrow.StructOf(arrow.Field{Name: "v", Type: arrow.PrimitiveTypes.Uint64, Nullable: true}), str: `{"v":18428615660272232523}`},
		{name: "list_int64", dt: arrow.ListOf(arrow.PrimitiveTypes.Int64), str: `[-8717895732742165505]`},
		{name: "list_uint64", dt: arrow.ListOf(arrow.PrimitiveTypes.Uint64), str: `[18428615660272232523]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScalar(tc.dt)
			require.NoError(t, s.Set(tc.str))
			require.Equal(t, tc.str, s.String())
		})
	}
}
