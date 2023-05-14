package scalar

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
)

func TestListSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result List
	}{
		{source: []int{1,2}, result: List{Value: []Scalar{
			&Int64{Value: 1, Valid: true},
			&Int64{Value: 2, Valid: true},
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