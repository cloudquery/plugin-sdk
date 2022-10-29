//nolint:revive,gocritic
package cqtypes

import (
	"fmt"
	"time"
)

const pgTimestamptzHourFormat = "2006-01-02 15:04:05.999999999Z07"
const pgTimestamptzMinuteFormat = "2006-01-02 15:04:05.999999999Z07:00"
const pgTimestamptzSecondFormat = "2006-01-02 15:04:05.999999999Z07:00:00"

// const microsecFromUnixEpochToY2K = 946684800 * 1000000

const (
// negativeInfinityMicrosecondOffset = -9223372036854775808
// infinityMicrosecondOffset         = 9223372036854775807
)

type Timestamptz struct {
	Time             time.Time
	Status           Status
	InfinityModifier InfinityModifier
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
		return dst.Time.String()
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
		return fmt.Errorf("cannot convert %v to Timestamptz", value)
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
		var format string
		if len(sbuf) >= 9 && (sbuf[len(sbuf)-9] == '-' || sbuf[len(sbuf)-9] == '+') {
			format = pgTimestamptzSecondFormat
		} else if len(sbuf) >= 6 && (sbuf[len(sbuf)-6] == '-' || sbuf[len(sbuf)-6] == '+') {
			format = pgTimestamptzMinuteFormat
		} else {
			format = pgTimestamptzHourFormat
		}

		tim, err := time.Parse(format, sbuf)
		if err != nil {
			return err
		}

		*dst = Timestamptz{Time: normalizePotentialUTC(tim), Status: Present}
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
