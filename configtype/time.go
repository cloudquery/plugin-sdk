package configtype

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/invopop/jsonschema"
)

// Time is a wrapper around time.Time that should be used in config
// when a time type is required. We wrap the time.Time type so that
// the spec can be extended in the future to support other types of times
type Time struct {
	input       string
	time        time.Time
	duration    *timeDuration
	hashNowFunc func() time.Time
}

func ParseTime(s string) (Time, error) {
	var t Time
	t.input = s

	var err error
	switch {
	case timeNowRegexp.MatchString(s):
		t.duration = new(timeDuration)
		*t.duration = newTimeDuration(0)
	case timeRFC3339Regexp.MatchString(s):
		t.time, err = time.Parse(time.RFC3339, s)
	case dateRegexp.MatchString(s):
		t.time, err = time.Parse(time.DateOnly, s)
	case baseDurationRegexp.MatchString(s), humanRelativeDurationRegexp.MatchString(s):
		t.duration = new(timeDuration)
		*t.duration, err = parseTimeDuration(s)
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

	numberRegexp = regexp.MustCompile(`^[0-9]+$`)

	baseDurationSegmentPattern = `[-+]?([0-9]*(\.[0-9]*)?[a-z]+)+` // copied from time.ParseDuration
	baseDurationPattern        = fmt.Sprintf(`^%s$`, baseDurationSegmentPattern)
	baseDurationRegexp         = regexp.MustCompile(baseDurationPattern)

	humanDurationSignsPattern = `ago|from\s+now`

	humanDurationUnitsPattern = `nanoseconds?|ns|microseconds?|us|µs|μs|milliseconds?|ms|seconds?|s|minutes?|m|hours?|h|days?|d|months?|M|years?|Y`
	humanDurationUnitsRegex   = regexp.MustCompile(fmt.Sprintf(`^%s$`, humanDurationUnitsPattern))

	humanDurationSegmentPattern = fmt.Sprintf(`(([0-9]+\s+(%[1]s)|%[2]s))`, humanDurationUnitsPattern, baseDurationSegmentPattern)

	humanRelativeDurationPattern = fmt.Sprintf(`^%[1]s(\s+%[1]s)*\s+(%[2]s)$`, humanDurationSegmentPattern, humanDurationSignsPattern)
	humanRelativeDurationRegexp  = regexp.MustCompile(humanRelativeDurationPattern)

	whitespaceRegexp = regexp.MustCompile(`\s+`)

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
		Description: "Allows for defining timestamps in both absolute(RFC3339) and relative formats. " +
			"Absolute timestamp example: `2024-01-01T12:00:00+00:00`.\n" +
			"Relative timestamps can take this format:\n" +
			"- `now`\n" +
			"- `x seconds [ago|from now]`\n" +
			"- `x minutes [ago|from now]`\n" +
			"- `x hours [ago|from now]`\n" +
			"- `x days [ago|from now]`\n" +
			"`until` field usage:\n" +
			"- `until: now`\n" +
			"- `until: 2 days ago`\n" +
			"- `until: 10 months 3 days 4h20m from now`\n" +
			"- `until: 2024-01-01T12:00:00+00:00`",
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

func (t Time) Hash() (uint64, error) {
	nowFunc := t.hashNowFunc
	if nowFunc == nil {
		nowFunc = time.Now
	}
	at := t.AsTime(nowFunc())
	return uint64(at.UnixNano()), nil
}

func (t *Time) SetHashNowFunc(f func() time.Time) {
	t.hashNowFunc = f
}

type timeDuration struct {
	input string

	relative bool
	sign     int
	duration time.Duration
	days     int
	months   int
	years    int
}

func newTimeDuration(d time.Duration) timeDuration {
	return timeDuration{
		input:    d.String(),
		sign:     1,
		duration: d,
	}
}

func parseTimeDuration(s string) (timeDuration, error) {
	var d timeDuration
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
				return timeDuration{}, fmt.Errorf("invalid duration format: invalid sign specifier: %q", part)
			}

			d.sign = 1
			inSign = false
		case inValue:
			if !humanDurationUnitsRegex.MatchString(part) {
				return timeDuration{}, fmt.Errorf("invalid duration format: invalid unit specifier: %q", part)
			}

			err = d.addUnit(part, value)
			if err != nil {
				return timeDuration{}, fmt.Errorf("invalid duration format: %w", err)
			}

			value = 0
			inValue = false
		case part == "ago":
			if d.sign != 0 {
				return timeDuration{}, fmt.Errorf("invalid duration format: more than one sign specifier")
			}

			d.sign = -1
		case part == "from":
			if d.sign != 0 {
				return timeDuration{}, fmt.Errorf("invalid duration format: more than one sign specifier")
			}

			inSign = true
		case numberRegexp.MatchString(part):
			value, err = strconv.ParseInt(part, 10, 64)
			if err != nil {
				return timeDuration{}, fmt.Errorf("invalid duration format: invalid value specifier: %q", part)
			}

			inValue = true
		case baseDurationRegexp.MatchString(part):
			duration, err := time.ParseDuration(part)
			if err != nil {
				return timeDuration{}, fmt.Errorf("invalid duration format: invalid value specifier: %q", part)
			}

			d.duration += duration
		default:
			return timeDuration{}, fmt.Errorf("invalid duration format: invalid value: %q", part)
		}
	}

	d.relative = d.sign != 0

	if !d.relative {
		d.sign = 1
	}

	return d, nil
}

func (d *timeDuration) addUnit(unit string, number int64) error {
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
		if number < math.MinInt || number > math.MaxInt {
			return fmt.Errorf("invalid %s value: %d. Out of bounds", unit, number)
		}
		d.days += int(number)
	case "month", "months":
		if number < math.MinInt || number > math.MaxInt {
			return fmt.Errorf("invalid %s value: %d. Out of bounds", unit, number)
		}
		d.months += int(number)
	case "year", "years":
		if number < math.MinInt || number > math.MaxInt {
			return fmt.Errorf("invalid %s value: %d. Out of bounds", unit, number)
		}
		d.years += int(number)
	default:
		return fmt.Errorf("invalid unit: %q", unit)
	}

	return nil
}

func (d timeDuration) Duration() time.Duration {
	duration := d.duration
	duration += time.Duration(d.days) * 24 * time.Hour
	duration += time.Duration(d.months) * 30 * 24 * time.Hour
	duration += time.Duration(d.years) * 365 * 24 * time.Hour
	duration *= time.Duration(d.sign)
	return duration
}
