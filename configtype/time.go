package configtype

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/invopop/jsonschema"
)

// Time is a wrapper around time.Time that should be used in config
// when a time type is required. We wrap the time.Time type so that
// the spec can be extended in the future to support other types of times
type Time struct {
	input    string
	time     time.Time
	duration *Duration
}

func ParseTime(s string) (Time, error) {
	var t Time
	t.input = s

	var err error
	switch {
	case timeNowRegexp.MatchString(s):
		t.duration = new(Duration)
		*t.duration = NewDuration(0)
	case timeRFC3339Regexp.MatchString(s):
		t.time, err = time.Parse(time.RFC3339, s)
	case dateRegexp.MatchString(s):
		t.time, err = time.Parse(time.DateOnly, s)
	case baseDurationRegexp.MatchString(s), humanRelativeDurationRegexp.MatchString(s):
		t.duration = new(Duration)
		*t.duration, err = ParseDuration(s)
	default:
		return t, fmt.Errorf("invalid time format: %s", s)
	}

	return t, err
}

var (
	timeRFC3339Pattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.(\d{1,9}))?(Z|((-|\+)\d{2}:\d{2}))$`
	timeRFC3339Regexp  = regexp.MustCompile(timeRFC3339Pattern)

	timeNowPattern = `^now$`
	timeNowRegexp  = regexp.MustCompile(timeNowPattern)

	datePattern = `^\d{4}-\d{2}-\d{2}$`
	dateRegexp  = regexp.MustCompile(datePattern)

	timePattern = patternCases(
		timeNowPattern,
		timeRFC3339Pattern,
		datePattern,
		baseDurationPattern,
		humanRelativeDurationPattern,
	)
)

func (Time) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Pattern: timePattern,
		Title:   "CloudQuery configtype.Time",
	}
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tim, err := ParseTime(s)
	if err != nil {
		return err
	}

	*t = tim
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.input)
}

func (t Time) AsTime(now time.Time) time.Time {
	if t.duration != nil {
		sign := t.duration.sign
		return now.Add(
			t.duration.duration*time.Duration(sign),
		).AddDate(
			t.duration.years*sign,
			t.duration.months*sign,
			t.duration.days*sign,
		)
	}

	return t.time
}

func (t Time) IsZero() bool {
	return t.duration == nil && t.time.IsZero()
}

func (t Time) String() string {
	return t.input
}
