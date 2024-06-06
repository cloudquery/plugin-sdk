package scalar

import (
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
)

func TestDate32Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Date32
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), result: Date32{Value: arrow.Date32FromTime(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), Valid: true}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC), result: Date32{Value: arrow.Date32FromTime(time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)), Valid: true}},
		{source: "2150-10-15", result: Date32{Value: arrow.Date32FromTime(time.Date(2150, 10, 15, 0, 0, 0, 0, time.UTC)), Valid: true}},
		{source: "", result: Date32{}},
	}

	for i, tt := range successfulTests {
		var r Date32
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
