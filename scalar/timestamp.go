package scalar

import (
	"encoding"
	"fmt"
	"math"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
)

const (
	// this is the default format used by time.Time.String()
	defaultStringFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

	// these are used by Arrow string format (time is in UTC)
	arrowStringFormat    = "2006-01-02 15:04:05.999999999"
	arrowStringFormatNew = "2006-01-02 15:04:05.999999999Z"
)

type Timestamp struct {
	Valid bool
	Value time.Time
	Type  *arrow.TimestampType
}

func (s *Timestamp) IsValid() bool {
	return s.Valid
}

func (s *Timestamp) DataType() arrow.DataType {
	return s.Type
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
		return nullValueStr
	}
	return s.Value.Format(time.RFC3339)
}

func (s *Timestamp) Get() any {
	return s.Value
}

func (s *Timestamp) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	if sc, ok := val.(Scalar); ok {
		if !sc.IsValid() {
			s.Valid = false
			return nil
		}
		return s.Set(sc.Get())
	}

	switch value := val.(type) {
	case int:
		s.Value = time.Unix(int64(value), 0).UTC()
	case int64:
		s.Value = time.Unix(value, 0).UTC()
	case uint64:
		if value > math.MaxInt64 {
			return &ValidationError{Type: s.DataType(), Msg: "uint64 greater than MaxInt64", Value: value}
		}
		s.Value = time.Unix(int64(value), 0).UTC()
	case time.Time:
		s.Value = value.UTC()
	case *time.Time:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case string:
		return s.DecodeText([]byte(value))
	case *string:
		if value == nil {
			s.Valid = false
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
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Timestamp) DecodeText(src []byte) error {
	if len(src) == 0 {
		s.Valid = false
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

		// there is no good way of detecting format, so we just try few of them
		for _, format := range []string{
			time.RFC3339,
			defaultStringFormat,
			arrowStringFormat,
			arrowStringFormatNew,
		} {
			tim, err = time.Parse(format, sbuf)
			if err == nil {
				s.Value = tim.UTC()
				s.Valid = true
				return nil
			}
		}
		return &ValidationError{Type: s.DataType(), Msg: "cannot parse timestamp", Value: sbuf, Err: err}
	}
}
