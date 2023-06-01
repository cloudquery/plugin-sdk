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
		{source: time.Date(1900, 1, 1, 0, 0, 1, 0, time.UTC), result: Date64{Value: -2208988800000, Valid: true}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC), result: Date64{Value: 946598400000, Valid: true}},
		{source: "2150-10-15", result: Date64{Value: 5705078400000, Valid: true}},
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
