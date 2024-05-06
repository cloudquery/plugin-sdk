package scalar

import (
	"net"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

func TestListSet(t *testing.T) {
	ipOne := net.IP{192, 168, 1, 1}
	ipNet := net.IPNet{IP: ipOne, Mask: net.IPMask{255, 255, 255, 255}}
	typedNil := (*net.IP)(nil)
	successfulTests := []struct {
		source any
		result List
	}{
		{source: []int{1, 2}, result: List{Value: []Scalar{
			&Int{Value: 1, Valid: true},
			&Int{Value: 2, Valid: true},
		}, Valid: true, Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)}},
		{source: &List{Value: []Scalar{
			&Int{Value: 1, Valid: true},
			&Int{Value: 2, Valid: true},
		}, Valid: true, Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
			result: List{Value: []Scalar{
				&Int{Value: 1, Valid: true},
				&Int{Value: 2, Valid: true},
			}, Valid: true, Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)}},
		{source: []*net.IPNet{&ipNet, nil}, result: List{Value: []Scalar{
			&Inet{Value: &ipNet, Valid: true},
			&Inet{Valid: false},
		}, Valid: true, Type: arrow.ListOf(types.ExtensionTypes.Inet)}},
		{source: []*net.IP{&ipOne, typedNil, nil}, result: List{Value: []Scalar{
			&Inet{Value: &ipNet, Valid: true},
			&Inet{Valid: false},
			&Inet{Valid: false},
		}, Valid: true, Type: arrow.ListOf(types.ExtensionTypes.Inet)}},
		{source: `[1, 2]`, result: List{Value: []Scalar{
			&Int{Value: 1, Valid: true},
			&Int{Value: 2, Valid: true},
		}, Valid: true, Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)}},
		{source: `[1, null, 2]`, result: List{Value: []Scalar{
			&Int{Value: 1, Valid: true},
			&Int{Valid: false},
			&Int{Value: 2, Valid: true},
		}, Valid: true, Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)}},
	}

	for i, tt := range successfulTests {
		r := List{
			Type: tt.result.Type,
		}
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
