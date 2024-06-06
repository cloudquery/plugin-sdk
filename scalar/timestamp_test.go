package scalar

import (
	"strconv"
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestTimestampDoubleSet(t *testing.T) {
	var r Timestamp
	assert.NoError(t, r.Set("2105-07-23 22:23:37.75007611 +0000 UTC"))

	r2 := r
	assert.NoError(t, r.Set(""))
	if r.Equal(&r2) {
		t.Errorf("%v = %v, expected null", r, r2)
	}
}

func TestAppendToBuilderTimestamp(t *testing.T) {
	for idx, tc := range []struct {
		Unit     arrow.TimeUnit
		Input    string
		Expected string
	}{
		// Input format: arrowStringFormat
		{
			Unit:     arrow.Second,
			Input:    "1999-01-08 04:05:06.123456789",
			Expected: "1999-01-08 04:05:06Z",
		},
		{
			Unit:     arrow.Millisecond,
			Input:    "1999-01-08 04:05:06.123456789",
			Expected: "1999-01-08 04:05:06.123Z",
		},
		{
			Unit:     arrow.Microsecond,
			Input:    "1999-01-08 04:05:06.123456789",
			Expected: "1999-01-08 04:05:06.123456Z",
		},
		{
			Unit:     arrow.Nanosecond,
			Input:    "1999-01-08 04:05:06.123456789",
			Expected: "1999-01-08 04:05:06.123456789Z",
		},
		// Input format: arrowStringFormatNew
		{
			Unit:     arrow.Second,
			Input:    "1999-01-08 04:05:06.123456789Z",
			Expected: "1999-01-08 04:05:06Z",
		},
		{
			Unit:     arrow.Millisecond,
			Input:    "1999-01-08 04:05:06.123456789Z",
			Expected: "1999-01-08 04:05:06.123Z",
		},
		{
			Unit:     arrow.Microsecond,
			Input:    "1999-01-08 04:05:06.123456789Z",
			Expected: "1999-01-08 04:05:06.123456Z",
		},
		{
			Unit:     arrow.Nanosecond,
			Input:    "1999-01-08 04:05:06.123456789Z",
			Expected: "1999-01-08 04:05:06.123456789Z",
		},
	} {
		tc := tc
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			timestamp := Timestamp{
				Type: &arrow.TimestampType{
					Unit:     tc.Unit,
					TimeZone: "UTC",
				},
			}
			err := timestamp.Set(tc.Input)
			if err != nil {
				t.Fatal(err)
			}

			bldr := array.NewTimestampBuilder(memory.DefaultAllocator, timestamp.Type)
			AppendToBuilder(bldr, &timestamp)

			arr := bldr.NewArray().(*array.Timestamp)
			actual := arr.ValueStr(0)

			require.Equal(t, tc.Expected, actual)
		})
	}
}
