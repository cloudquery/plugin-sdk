package configtype

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/invopop/jsonschema"
)

var (
	durationPattern = `^[-+]?([0-9]*(\.[0-9]*)?[a-z]+)+$` // copied from time.ParseDuration
	durationRegexp  = regexp.MustCompile(durationPattern)
)

// Duration is a wrapper around time.Duration that should be used in config
// when a duration type is required. We wrap the time.Duration type so that
// the spec can be extended in the future to support other types of durations
// (e.g. a duration that is specified in days).
type Duration struct {
	duration time.Duration
}

func NewDuration(d time.Duration) Duration {
	return Duration{
		duration: d,
	}
}

func (Duration) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Pattern: durationPattern, // copied from time.ParseDuration
		Title:   "CloudQuery configtype.Duration",
	}
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration{duration: duration}
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.duration.String())
}

func (d *Duration) Duration() time.Duration {
	return d.duration
}

func (d Duration) Equal(other Duration) bool {
	return d.duration == other.duration
}
