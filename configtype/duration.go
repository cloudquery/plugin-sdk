package configtype

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/invopop/jsonschema"
)

var (
	baseDurationPattern = `^[-+]?([0-9]*(\.[0-9]*)?[a-z]+)+$` // copied from time.ParseDuration
	baseDurationRegexp  = regexp.MustCompile(baseDurationPattern)

	humanDurationUnits = "seconds?|minutes?|hours?|days?|months?|years?"

	humanDurationPattern = fmt.Sprintf(`^[0-9]+\s+(%s)$`, humanDurationUnits)
	humanDurationRegexp  = regexp.MustCompile(humanDurationPattern)

	humanRelativeDurationPattern = fmt.Sprintf(`^[0-9]+\s+(%s)\s+(ago|from\s+now)$`, humanDurationUnits)
	humanRelativeDurationRegexp  = regexp.MustCompile(humanRelativeDurationPattern)

	whitespaceRegexp = regexp.MustCompile(`\s+`)

	fromNowRegexp = regexp.MustCompile(`from\s+now`)
)

// Duration is a wrapper around time.Duration that should be used in config
// when a duration type is required. We wrap the time.Duration type so that
// the spec can be extended in the future to support other types of durations
// (e.g. a duration that is specified in days).
type Duration struct {
	duration time.Duration
	days     int
	months   int
	years    int
}

func NewDuration(d time.Duration) Duration {
	return Duration{
		duration: d,
	}
}

func ParseDuration(s string) (Duration, error) {
	var d Duration
	var err error
	switch {
	case humanDurationRegexp.MatchString(s):
		d, err = parseHumanDuration(s)
	case humanRelativeDurationRegexp.MatchString(s):
		d, err = parseHumanRelativeDuration(s)
	case baseDurationRegexp.MatchString(s):
		d.duration, err = time.ParseDuration(s)
	default:
		return d, fmt.Errorf("invalid duration format: %q", s)
	}

	return d, err
}

func parseHumanDuration(s string) (Duration, error) {
	parts := whitespaceRegexp.Split(s, 2)
	if len(parts) != 2 {
		return Duration{}, fmt.Errorf("invalid duration format: %q", s)
	}

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return Duration{}, fmt.Errorf("invalid duration format: invalid number: %q", s)
	}

	d, err := parseHumanDurationUnit(parts[1], 1, number)
	if err != nil {
		return Duration{}, fmt.Errorf("invalid duration format: %w", err)
	}

	return d, nil
}

func parseHumanRelativeDuration(s string) (Duration, error) {
	parts := whitespaceRegexp.Split(s, 3)
	if len(parts) != 3 {
		return Duration{}, fmt.Errorf("invalid duration format: %q", s)
	}

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return Duration{}, fmt.Errorf("invalid duration format: invalid number: %q", s)
	}

	sign, err := parseHumanDurationSign(parts[2])
	if err != nil {
		return Duration{}, fmt.Errorf("invalid duration format: %w", err)
	}

	d, err := parseHumanDurationUnit(parts[1], sign, number)
	if err != nil {
		return Duration{}, fmt.Errorf("invalid duration format: %w", err)
	}

	return d, nil
}

func parseHumanDurationUnit(unit string, sign, number int) (Duration, error) {
	var d Duration
	switch unit {
	case "second", "seconds":
		d.duration = time.Second * time.Duration(sign) * time.Duration(number)
	case "minute", "minutes":
		d.duration = time.Minute * time.Duration(sign) * time.Duration(number)
	case "hour", "hours":
		d.duration = time.Hour * time.Duration(sign) * time.Duration(number)
	case "day", "days":
		d.days = sign * number
	case "month", "months":
		d.months = sign * number
	case "year", "years":
		d.years = sign * number
	default:
		return Duration{}, fmt.Errorf("invalid unit: %q", unit)
	}

	return d, nil
}

func parseHumanDurationSign(sign string) (int, error) {
	switch {
	case sign == "ago":
		return -1, nil
	case fromNowRegexp.MatchString(sign):
		return 1, nil
	default:
		return 0, fmt.Errorf("invalid duration format: invalid sign specifier: %q", sign)
	}
}

func (Duration) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Pattern: patternCases(baseDurationPattern, humanDurationPattern, humanRelativeDurationPattern),
		Title:   "CloudQuery configtype.Duration",
	}
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	duration, err := ParseDuration(s)
	if err != nil {
		return err
	}

	*d = duration
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration().String())
}

func (d *Duration) Duration() time.Duration {
	duration := d.duration
	duration += time.Duration(d.days) * 24 * time.Hour
	duration += time.Duration(d.months) * 30 * 24 * time.Hour
	duration += time.Duration(d.years) * 365 * 24 * time.Hour
	return duration
}

func (d Duration) Equal(other Duration) bool {
	return d == other
}

func (d Duration) String() string {
	return d.Duration().String()
}
