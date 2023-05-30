package scalar

import (
	"encoding/json"
	"reflect"

	"github.com/apache/arrow/go/v13/arrow"
)

type Map struct {
	Valid bool
	Value any

	Type *arrow.MapType
}

func (s *Map) IsValid() bool {
	return s.Valid
}

func (s *Map) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Map)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && arrow.TypeEqual(s.Type, r.Type) && reflect.DeepEqual(s.Value, r.Value)
}

func (s *Map) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.Type.String() + " value"
}

func (s *Map) Get() any {
	return s.Value
}

func (s *Map) Set(val any) error {
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

	if str, ok := val.(string); ok {
		var x map[string]any
		if err := json.Unmarshal([]byte(str), &x); err != nil {
			return err
		}
		s.Value = x
		s.Valid = true
		return nil
	}

	s.Value = val
	s.Valid = true
	return nil
}

func (s *Map) DataType() arrow.DataType {
	return s.Type
}
