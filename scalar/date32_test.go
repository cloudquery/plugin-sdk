package scalar

import (
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/stretchr/testify/assert"
)

func TestDate32Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Date32
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			result: Date32{Value: arrow.Date32FromTime(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), Valid: true}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC),
			result: Date32{Value: arrow.Date32FromTime(time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)), Valid: true}},
		{source: "2150-10-15", result: Date32{Value: arrow.Date32FromTime(time.Date(2150, 10, 15, 0, 0, 0, 0, time.UTC)), Valid: true}},
		{source: "", result: Date32{}},
	}

	for i, tt := range successfulTests {
		t.Run(tt.result.String(), func(t *testing.T) {
			var r Date32
			assert.NoError(t, r.Set(tt.source))
			assert.Truef(t, r.Equal(&tt.result), "%d: %v != %v", i, r, tt.result)
		})
	}
}
