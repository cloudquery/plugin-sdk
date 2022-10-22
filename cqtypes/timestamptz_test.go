package cqtypes

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
		{source: "2020-04-05 06:07:08Z", result: Timestamptz{Time: time.Date(2020, 4, 5, 6, 7, 8, 0, time.UTC), Status: Present}},
	}

	for i, tt := range successfulTests {
		var r Timestamptz
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if diff := cmp.Diff(r, tt.result); diff != "" {
			t.Errorf("%d: got diff:\n%s", i, diff)
		}
	}
}