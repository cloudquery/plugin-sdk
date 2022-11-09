//nolint:revive,gocritic
package schema

import (
	"encoding"
	"fmt"
	"time"
)

// const pgTimestamptzHourFormat = "2006-01-02 15:04:05.999999999Z07"
// const pgTimestamptzMinuteFormat = "2006-01-02 15:04:05.999999999Z07:00"
// const pgTimestamptzSecondFormat = "2006-01-02 15:04:05.999999999Z07:00:00"

// this is the default format used by time.Time.String()
const defaultStringFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

// const microsecFromUnixEpochToY2K = 946684800 * 1000000

const (
// negativeInfinityMicrosecondOffset = -9223372036854775808
// infinityMicrosecondOffset         = 9223372036854775807
)

type TimestamptzTransformer interface {
	TransformTimestamptz(*Timestamptz) interface{}
}

type Timestamptz struct {
	Time             time.Time
	Status           Status
	InfinityModifier InfinityModifier
}

func (*Timestamptz) Type() ValueType {
	return TypeTimestamp
}

func (dst *Timestamptz) Equal(src CQType) bool {
	if src == nil {
		return false
	}

	if value, ok := src.(*Timestamptz); ok {
		return dst.Status == value.Status && dst.Time.Equal(value.Time) && dst.InfinityModifier == value.InfinityModifier
	}

	return false
}

func (dst *Timestamptz) String() string {
	if dst.Status == Present {
		return dst.Time.Format(time.RFC3339)
	} else {
		return ""
	}
}

func (dst *Timestamptz) Set(src interface{}) error {
	if src == nil {
		*dst = Timestamptz{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case time.Time:
		*dst = Timestamptz{Time: value, Status: Present}
	case *time.Time:
		if value == nil {
			*dst = Timestamptz{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case string:
		return dst.DecodeText([]byte(value))
	case *string:
		if value == nil {
			*dst = Timestamptz{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case InfinityModifier:
		*dst = Timestamptz{InfinityModifier: value, Status: Present}
	default:
		if originalSrc, ok := underlyingTimeType(src); ok {
			return dst.Set(originalSrc)
		}
		if value, ok := value.(encoding.TextMarshaler); ok {
			s, err := value.MarshalText()
			if err == nil {
				return dst.Set(string(s))
			}
			// fall through to String() method
		}
		if value, ok := value.(fmt.Stringer); ok {
			s := value.String()
			return dst.Set(s)
		}
		return fmt.Errorf("cannot convert %v of type %T to Timestamptz", value, value)
	}

	return nil
}

func (dst Timestamptz) Get() interface{} {
	switch dst.Status {
	case Present:
		if dst.InfinityModifier != None {
			return dst.InfinityModifier
		}
		return dst.Time
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (dst *Timestamptz) DecodeText(src []byte) error {
	if src == nil {
		*dst = Timestamptz{Status: Null}
		return nil
	}

	sbuf := string(src)
	switch sbuf {
	case "infinity":
		*dst = Timestamptz{Status: Present, InfinityModifier: Infinity}
	case "-infinity":
		*dst = Timestamptz{Status: Present, InfinityModifier: -Infinity}
	default:
		var tim time.Time
		var err error

		if len(sbuf) > len(defaultStringFormat)+1 && sbuf[len(defaultStringFormat)+1] == 'm' {
			sbuf = sbuf[:len(defaultStringFormat)]
		}

		// there is no good way of detecting format so we just try few of them
		tim, err = time.Parse(time.RFC3339, sbuf)
		if err == nil {
			*dst = Timestamptz{Time: normalizePotentialUTC(tim), Status: Present}
			return nil
		}
		tim, err = time.Parse(defaultStringFormat, sbuf)
		if err == nil {
			*dst = Timestamptz{Time: normalizePotentialUTC(tim), Status: Present}
			return nil
		}
		return fmt.Errorf("cannot parse %s as Timestamptz", sbuf)

	}

	return nil
}

// Normalize timestamps in UTC location to behave similarly to how the Golang
// standard library does it: UTC timestamps lack a .loc value.
//
// Reason for this: when comparing two timestamps with reflect.DeepEqual (generally
// speaking not a good idea, but several testing libraries (for example testify)
// does this), their location data needs to be equal for them to be considered
// equal.
func normalizePotentialUTC(timestamp time.Time) time.Time {
	if timestamp.Location().String() != time.UTC.String() {
		return timestamp
	}

	return timestamp.UTC()
}
