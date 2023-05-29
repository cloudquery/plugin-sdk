package scalar

import (
	"testing"
	"time"
)

func TestDate64Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Date64
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 1, 0, time.Local), result: Date64{Value: time.Date(1900, 1, 1, 0, 0, 1, 0, time.Local).UnixMilli(), Valid: true}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), result: Date64{Value: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local).UnixMilli(), Valid: true}},
		{source: "2150-10-15 12:13:14.123456789", result: Date64{Value: time.Date(2150, 10, 15, 12, 13, 14, 123456789, time.UTC).UnixMilli(), Valid: true}},
		{source: "", result: Date64{}},
	}

	for i, tt := range successfulTests {
		var r Date64
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
