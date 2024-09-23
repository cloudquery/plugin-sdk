package configtype

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/invopop/jsonschema"
)

type NowFunc func() time.Time

// Time is a wrapper around time.Time that should be used in config
// when a time type is required. We wrap the time.Time type so that
// the spec can be extended in the future to support other types of times
type Time struct {
	time     time.Time
	duration time.Duration
}

func NewTime(t time.Time) Time {
	return Time{
		time: t,
	}
}

func NewRelativeTime(d time.Duration) Time {
	return Time{
		duration: d,
	}
}

var (
	timeRFC3339Pattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.(\d{3}|\d{6}|\d{9}))?(Z|((-|\+)\d{2}:\d{2}))$`
	timeRFC3339Regexp  = regexp.MustCompile(timeRFC3339Pattern)

	datePattern = `^\d{4}-\d{2}-\d{2}$`
	dateRegexp  = regexp.MustCompile(datePattern)

	timePattern = patternCases(timeRFC3339Pattern, datePattern, durationPattern)
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

	var err error
	switch {
	case timeRFC3339Regexp.MatchString(s):
		t.time, err = time.Parse(time.RFC3339, s)
	case dateRegexp.MatchString(s):
		t.time, err = time.Parse(time.DateOnly, s)
	case durationRegexp.MatchString(s):
		t.duration, err = time.ParseDuration(s)
	default:
		return fmt.Errorf("invalid time format: %s", s)
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *Time) MarshalJSON() ([]byte, error) {
	if !d.time.IsZero() {
		return json.Marshal(d.time)
	}

	return json.Marshal(d.duration.String())
}

func (t *Time) Time(nowFunc NowFunc) time.Time {
	if !t.time.IsZero() {
		return t.time
	}

	return nowFunc().Add(t.duration)
}

func (t Time) IsRelative() bool {
	return t.time.IsZero()
}

// Equal compares two Time structs. Note that relative and fixed times are never equal
func (t Time) Equal(other Time) bool {
	return t.time.Equal(other.time) && t.duration == other.duration
}

func (t Time) String() string {
	if !t.time.IsZero() {
		return t.time.String()
	}
	return t.duration.String()
}
