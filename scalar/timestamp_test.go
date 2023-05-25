package scalar

import (
	"testing"
	"time"
)

type TimestampSt struct {
	time.Time
}

func TestTimestampSet(t *testing.T) {
	type _time time.Time

	timeInstance := time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC)
	timeRFC3339NanoBytes, _ := timeInstance.MarshalText()

	successfulTests := []struct {
		source any
		result Timestamp
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamp{Value: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamp{Value: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), result: Timestamp{Value: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), Valid: true}},
		{source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamp{Value: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local), result: Timestamp{Value: time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local), Valid: true}},
		{source: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), result: Timestamp{Value: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: int(time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local).Unix()), result: Timestamp{Value: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: uint64(time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local).Unix()), result: Timestamp{Value: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local).Unix(), result: Timestamp{Value: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: _time(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)), result: Timestamp{Value: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Valid: true}},
		{source: string(timeRFC3339NanoBytes), result: Timestamp{Value: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Valid: true}},
		{source: "2150-10-15 07:25:09.75007611 +0000 UTC", result: Timestamp{Value: time.Date(2150, 10, 15, 7, 25, 9, 750076110, time.UTC), Valid: true}},
		{source: timeInstance.String(), result: Timestamp{Value: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Valid: true}},
		{source: TimestampSt{timeInstance}, result: Timestamp{Value: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Valid: true}},
		{source: "", result: Timestamp{}},
		{source: &Timestamp{Value: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Valid: true}, result: Timestamp{Value: time.Date(2105, 7, 23, 22, 23, 37, 750076110, time.UTC), Valid: true}},
	}

	for i, tt := range successfulTests {
		var r Timestamp
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
