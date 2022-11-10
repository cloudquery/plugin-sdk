package schema

import (
	"testing"
	"time"
)

type Timestamp struct {
	time.Time
}

func TestTimestamptzSet(t *testing.T) {
	type _time time.Time

	timeInstance := time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC)
	timeRFC3339NanoBytes, _ := timeInstance.MarshalText()

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
		{source: string(timeRFC3339NanoBytes), result: Timestamptz{Time: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Status: Present}},
		{source: "2150-10-15 07:25:09.75007611 +0000 UTC", result: Timestamptz{Time: time.Date(2150, 10, 15, 7, 25, 9, 750076110, time.UTC), Status: Present}},
		{source: timeInstance.String(), result: Timestamptz{Time: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Status: Present}},
		{source: Timestamp{timeInstance}, result: Timestamptz{Time: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Status: Present}},
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
