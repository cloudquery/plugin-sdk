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
	numberRegexp = regexp.MustCompile(`^[0-9]+$`)

	baseDurationSegmentPattern = `[-+]?([0-9]*(\.[0-9]*)?[a-z]+)+` // copied from time.ParseDuration
	baseDurationPattern        = fmt.Sprintf(`^%s$`, baseDurationSegmentPattern)
	baseDurationRegexp         = regexp.MustCompile(baseDurationPattern)

	humanDurationSignsPattern = `ago|from\s+now`

	humanDurationUnitsPattern = `nanoseconds?|ns|microseconds?|us|µs|μs|milliseconds?|ms|seconds?|s|minutes?|m|hours?|h|days?|d|months?|M|years?|Y`
	humanDurationUnitsRegex   = regexp.MustCompile(fmt.Sprintf(`^%s$`, humanDurationUnitsPattern))

	humanDurationSegmentPattern = fmt.Sprintf(`(([0-9]+\s+(%[1]s)|%[2]s))`, humanDurationUnitsPattern, baseDurationSegmentPattern)

	humanDurationPattern = fmt.Sprintf(`^%[1]s(\s+%[1]s)*$`, humanDurationSegmentPattern)

	humanRelativeDurationPattern = fmt.Sprintf(`^%[1]s(\s+%[1]s)*\s+(%[2]s)$`, humanDurationSegmentPattern, humanDurationSignsPattern)
	humanRelativeDurationRegexp  = regexp.MustCompile(humanRelativeDurationPattern)

	whitespaceRegexp = regexp.MustCompile(`\s+`)
)

// Duration is a wrapper around time.Duration that should be used in config
// when a duration type is required. We wrap the time.Duration type so that
// the spec can be extended in the future to support other types of durations
// (e.g. a duration that is specified in days).
type Duration struct {
	input string

	relative bool
	sign     int
	duration time.Duration
	days     int
	months   int
	years    int
}

func NewDuration(d time.Duration) Duration {
	return Duration{
		input:    d.String(),
		sign:     1,
		duration: d,
	}
}

func ParseDuration(s string) (Duration, error) {
	var d Duration
	d.input = s

	var inValue bool
	var value int64

	var inSign bool

	parts := whitespaceRegexp.Split(s, -1)

	var err error

	for _, part := range parts {
		switch {
		case inSign:
			if part != "now" {
				return Duration{}, fmt.Errorf("invalid duration format: invalid sign specifier: %q", part)
			}

			d.sign = 1
			inSign = false
		case inValue:
			if !humanDurationUnitsRegex.MatchString(part) {
				return Duration{}, fmt.Errorf("invalid duration format: invalid unit specifier: %q", part)
			}

			err = d.addUnit(part, value)
			if err != nil {
				return Duration{}, fmt.Errorf("invalid duration format: %w", err)
			}

			value = 0
			inValue = false
		case part == "ago":
			if d.sign != 0 {
				return Duration{}, fmt.Errorf("invalid duration format: more than one sign specifier")
			}

			d.sign = -1
		case part == "from":
			if d.sign != 0 {
				return Duration{}, fmt.Errorf("invalid duration format: more than one sign specifier")
			}

			inSign = true
		case numberRegexp.MatchString(part):
			value, err = strconv.ParseInt(part, 10, 64)
			if err != nil {
				return Duration{}, fmt.Errorf("invalid duration format: invalid value specifier: %q", part)
			}

			inValue = true
		case baseDurationRegexp.MatchString(part):
			duration, err := time.ParseDuration(part)
			if err != nil {
				return Duration{}, fmt.Errorf("invalid duration format: invalid value specifier: %q", part)
			}

			d.duration += duration
		default:
			return Duration{}, fmt.Errorf("invalid duration format: invalid value: %q", part)
		}
	}

	d.relative = d.sign != 0

	if !d.relative {
		d.sign = 1
	}

	return d, nil
}

func (d *Duration) addUnit(unit string, number int64) error {
	switch unit {
	case "nanosecond", "nanoseconds", "ns":
		d.duration += time.Nanosecond * time.Duration(number)
	case "microsecond", "microseconds", "us", "μs", "µs":
		d.duration += time.Microsecond * time.Duration(number)
	case "millisecond", "milliseconds":
		d.duration += time.Millisecond * time.Duration(number)
	case "second", "seconds":
		d.duration += time.Second * time.Duration(number)
	case "minute", "minutes":
		d.duration += time.Minute * time.Duration(number)
	case "hour", "hours":
		d.duration += time.Hour * time.Duration(number)
	case "day", "days":
		d.days += int(number)
	case "month", "months":
		d.months += int(number)
	case "year", "years":
		d.years += int(number)
	default:
		return fmt.Errorf("invalid unit: %q", unit)
	}

	return nil
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
	return json.Marshal(d.String())
}

func (d Duration) Duration() time.Duration {
	duration := d.duration
	duration += time.Duration(d.days) * 24 * time.Hour
	duration += time.Duration(d.months) * 30 * 24 * time.Hour
	duration += time.Duration(d.years) * 365 * 24 * time.Hour
	duration *= time.Duration(d.sign)
	return duration
}

func (d Duration) Equal(other Duration) bool {
	return d == other
}

func (d Duration) String() string {
	return d.input
}
