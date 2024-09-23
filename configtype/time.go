package configtype

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/invopop/jsonschema"
)

type timeType int

const (
	timeTypeZero timeType = iota
	timeTypeFixed
	timeTypeRelative
)

// Time is a wrapper around time.Time that should be used in config
// when a time type is required. We wrap the time.Time type so that
// the spec can be extended in the future to support other types of times
type Time struct {
	typ      timeType
	time     time.Time
	duration time.Duration
}

func NewTime(t time.Time) Time {
	return Time{
		typ:  timeTypeFixed,
		time: t,
	}
}

func NewRelativeTime(d time.Duration) Time {
	return Time{
		typ:      timeTypeRelative,
		duration: d,
	}
}

var (
	timeRFC3339Pattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.(\d{1,9}))?(Z|((-|\+)\d{2}:\d{2}))$`
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
		if t.time.IsZero() {
			t.typ = timeTypeZero
		} else {
			t.typ = timeTypeFixed
		}
	case dateRegexp.MatchString(s):
		t.typ = timeTypeFixed
		t.time, err = time.Parse(time.DateOnly, s)
	case durationRegexp.MatchString(s):
		t.typ = timeTypeRelative
		t.duration, err = time.ParseDuration(s)
	default:
		return fmt.Errorf("invalid time format: %s", s)
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	switch t.typ {
	case timeTypeFixed:
		return json.Marshal(t.time)
	case timeTypeRelative:
		return json.Marshal(t.duration.String())
	default:
		return json.Marshal(time.Time{})
	}
}

func (t Time) Time(nowFunc func() time.Time) time.Time {
	switch t.typ {
	case timeTypeFixed:
		return t.time
	case timeTypeRelative:
		return nowFunc().Add(t.duration)
	default:
		return time.Time{}
	}
}

func (t Time) IsRelative() bool {
	return t.typ == timeTypeRelative
}

func (t Time) IsZero() bool {
	return t.typ == timeTypeZero
}

func (t Time) IsFixed() bool {
	return t.typ == timeTypeFixed
}

// Equal compares two Time structs. Note that relative and fixed times are never equal
func (t Time) Equal(other Time) bool {
	return t.typ == other.typ && t.time.Equal(other.time) && t.duration == other.duration
}

func (t Time) String() string {
	switch t.typ {
	case timeTypeFixed:
		return t.time.String()
	case timeTypeRelative:
		return t.duration.String()
	default:
		return time.Time{}.String()
	}
}
