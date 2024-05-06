package scalar

import (
	"encoding/json"

	"github.com/apache/arrow/go/v16/arrow"
)

type MonthInterval struct {
	Int
}

type monthIntervalData struct {
	Months int32 `json:"months"`
}

func (*MonthInterval) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.MonthInterval
}

func (s *MonthInterval) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*MonthInterval)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *MonthInterval) String() string {
	if !s.Valid {
		return nullValueStr
	}

	b, _ := json.Marshal(monthIntervalData{Months: int32(s.Value)})
	return string(b)
}

func (s *MonthInterval) Set(value any) error {
	if mi, ok := value.(arrow.MonthInterval); ok {
		return s.Int.Set(int32(mi))
	}

	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			s.Valid = false
			return nil
		}
		return s.Int.Set(value)
	case []byte:
		if len(v) == 0 {
			s.Valid = false
			return nil
		}

		var mi monthIntervalData
		if err := json.Unmarshal(v, &mi); err != nil {
			return err
		}
		s.Valid = true
		s.Value = int64(mi.Months)
		return nil
	case *string:
		if v == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*v)
	case map[string]any:
		b, _ := json.Marshal(v)
		return s.Set(b)
	default:
		return s.Int.Set(value)
	}
}

type DayTimeInterval struct {
	Value arrow.DayTimeInterval
	Valid bool
}

func (*DayTimeInterval) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.DayTimeInterval
}

func (s *DayTimeInterval) IsValid() bool {
	return s.Valid
}

func (s *DayTimeInterval) String() string {
	if !s.Valid {
		return nullValueStr
	}
	b, _ := json.Marshal(s.Value)
	return string(b)
}

func (s *DayTimeInterval) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*DayTimeInterval)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *DayTimeInterval) Set(value any) error {
	if value == nil {
		s.Valid = false
		return nil
	}

	if dti, ok := value.(arrow.DayTimeInterval); ok {
		s.Valid = true
		s.Value = dti
		return nil
	}

	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			s.Valid = false
			return nil
		}

		var dti arrow.DayTimeInterval
		if err := json.Unmarshal([]byte(v), &dti); err != nil {
			return err
		}
		s.Valid = true
		s.Value = dti
		return nil
	case []byte:
		if len(v) == 0 {
			s.Valid = false
			return nil
		}

		var dti arrow.DayTimeInterval
		if err := json.Unmarshal(v, &dti); err != nil {
			return err
		}
		s.Valid = true
		s.Value = dti
		return nil
	case *string:
		if v == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*v)
	case map[string]any:
		b, _ := json.Marshal(v)
		return s.Set(b)
	default:
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
}

func (s *DayTimeInterval) Get() any {
	if !s.Valid {
		return nil
	}

	return s.Value
}

type MonthDayNanoInterval struct {
	Value arrow.MonthDayNanoInterval
	Valid bool
}

func (*MonthDayNanoInterval) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.MonthDayNanoInterval
}

func (s *MonthDayNanoInterval) IsValid() bool {
	return s.Valid
}

func (s *MonthDayNanoInterval) String() string {
	if !s.Valid {
		return nullValueStr
	}
	b, _ := json.Marshal(s.Value)
	return string(b)
}

func (s *MonthDayNanoInterval) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*MonthDayNanoInterval)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *MonthDayNanoInterval) Set(value any) error {
	if value == nil {
		s.Valid = false
		return nil
	}

	if dti, ok := value.(arrow.MonthDayNanoInterval); ok {
		s.Valid = true
		s.Value = dti
		return nil
	}

	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			s.Valid = false
			return nil
		}

		var dti arrow.MonthDayNanoInterval
		if err := json.Unmarshal([]byte(v), &dti); err != nil {
			return err
		}
		s.Valid = true
		s.Value = dti
		return nil
	case []byte:
		if len(v) == 0 {
			s.Valid = false
			return nil
		}

		var dti arrow.MonthDayNanoInterval
		if err := json.Unmarshal(v, &dti); err != nil {
			return err
		}
		s.Valid = true
		s.Value = dti
		return nil
	case *string:
		if v == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*v)
	case map[string]any:
		b, _ := json.Marshal(v)
		return s.Set(b)
	default:
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
}

func (s *MonthDayNanoInterval) Get() any {
	if !s.Valid {
		return nil
	}

	return s.Value
}
