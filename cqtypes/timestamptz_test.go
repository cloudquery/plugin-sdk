package cqtypes

import (
	"testing"
	"time"
)

func TestTimestamptzSet(t *testing.T) {
	type _time time.Time

	successfulTests := []struct {
		source interface{}
		result Timestamptz
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamptz{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), Status: Present}},
		{source: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamptz{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Status: Present}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), result: Timestamptz{Time: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), Status: Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Status: Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local), result: Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local), Status: Present}},
		{source: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamptz{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), Status: Present}},
		{source: _time(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)), result: Timestamptz{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Status: Present}},
		{source: Infinity, result: Timestamptz{InfinityModifier: Infinity, Status: Present}},
		{source: NegativeInfinity, result: Timestamptz{InfinityModifier: NegativeInfinity, Status: Present}},
		{source: "2105-07-23 22:23:37.133334968 +0300 UTC m=+2610529275.484148261", result: Timestamptz{Time: time.Date(2105, 7, 23, 22, 23, 37, 133334968, time.UTC), Status: Present}},
		// {source: "2020-04-05 06:07:08Z", result: Timestamptz{Time: time.Date(2020, 4, 5, 6, 7, 8, 0, time.UTC), Status: Present}},
	}

	for i, tt := range successfulTests {
		var r Timestamptz
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
