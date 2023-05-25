package scalar

import (
	"encoding"
	"fmt"
	"math"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
)

// const pgTimestamptzHourFormat = "2006-01-02 15:04:05.999999999Z07"
// const pgTimestamptzMinuteFormat = "2006-01-02 15:04:05.999999999Z07:00"
// const pgTimestamptzSecondFormat = "2006-01-02 15:04:05.999999999Z07:00:00"

// this is the default format used by time.Time.String()
const defaultStringFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

// this is used by arrow string format (time is in UTC)
const arrowStringFormat = "2006-01-02 15:04:05.999999999"

// const microsecFromUnixEpochToY2K = 946684800 * 1000000

const (
// negativeInfinityMicrosecondOffset = -9223372036854775808
// infinityMicrosecondOffset         = 9223372036854775807
)

type Timestamp struct {
	Valid bool
	Value time.Time
}

func (s *Timestamp) IsValid() bool {
	return s.Valid
}

func (*Timestamp) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.Timestamp_us
}

func (s *Timestamp) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Timestamp)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value.Equal(r.Value)
}

func (s *Timestamp) String() string {
	if !s.Valid {
		return "(null)"
	}
	return s.Value.Format(time.RFC3339)
}

func (s *Timestamp) Set(val any) error {
	if val == nil {
		return nil
	}

	switch value := val.(type) {
	case int:
		if value < 0 {
			return &ValidationError{Type: arrow.FixedWidthTypes.Timestamp_us, Msg: "negative timestamp"}
		}
		s.Value = time.Unix(int64(value), 0).UTC()
	case int64:
		if value < 0 {
			return &ValidationError{Type: arrow.FixedWidthTypes.Timestamp_us, Msg: "negative timestamp"}
		}
		s.Value = time.Unix(value, 0).UTC()
	case uint64:
		if value > math.MaxInt64 {
			return &ValidationError{Type: arrow.FixedWidthTypes.Timestamp_us, Msg: "uint64 bigger than MaxInt64", Value: value}
		}
		s.Value = time.Unix(int64(value), 0).UTC()
	case time.Time:
		s.Value = value.UTC()
	case *time.Time:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case string:
		return s.DecodeText([]byte(value))
	case *string:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingTimeType(val); ok {
			return s.Set(originalSrc)
		}
		if value, ok := value.(encoding.TextMarshaler); ok {
			text, err := value.MarshalText()
			if err == nil {
				return s.Set(string(text))
			}
			// fall through to String() method
		}
		if value, ok := value.(fmt.Stringer); ok {
			str := value.String()
			return s.Set(str)
		}
		return &ValidationError{Type: arrow.FixedWidthTypes.Timestamp_us, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Timestamp) DecodeText(src []byte) error {
	if len(src) == 0 {
		return nil
	}

	sbuf := string(src)
	// nolint:gocritic,revive
	switch sbuf {
	default:
		var tim time.Time
		var err error

		if len(sbuf) > len(defaultStringFormat)+1 && sbuf[len(defaultStringFormat)+1] == 'm' {
			sbuf = sbuf[:len(defaultStringFormat)]
		}

		// there is no good way of detecting format so we just try few of them
		tim, err = time.Parse(time.RFC3339, sbuf)
		if err == nil {
			s.Value = tim.UTC()
			s.Valid = true
			return nil
		}
		tim, err = time.Parse(defaultStringFormat, sbuf)
		if err == nil {
			s.Value = tim.UTC()
			s.Valid = true
			return nil
		}
		tim, err = time.Parse(arrowStringFormat, sbuf)
		if err == nil {
			s.Value = tim.UTC()
			s.Valid = true
			return nil
		}
		return &ValidationError{Type: arrow.FixedWidthTypes.Timestamp_us, Msg: "cannot parse timestamp", Value: sbuf, Err: err}
	}
}
