package scalar

import (
	"reflect"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
)

type List struct {
	Valid bool
	Value Vector
	Type  arrow.DataType
}

func (s *List) IsValid() bool {
	return s.Valid
}

func (s *List) DataType() arrow.DataType {
	return s.Type
}

func (s *List) String() string {
	if !s.Valid {
		return "(null)"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range s.Value {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

func (s *List) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*List)
	if !ok {
		return false
	}
	if s.Valid != r.Valid {
		return false
	}
	if len(s.Value) != len(r.Value) {
		return false
	}
	for i := range s.Value {
		if !s.Value[i].Equal(r.Value[i]) {
			return false
		}
	}
	return true
}

func (s *List) Get() any {
	return s.Value
}

func (s *List) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}
	if s.Type == nil {
		panic("List type is nil")
	}

	if sc, ok := val.(Scalar); ok {
		if !sc.IsValid() {
			s.Valid = false
			return nil
		}
		return s.Set(sc.Get())
	}

	reflectedValue := reflect.ValueOf(val)
	if !reflectedValue.IsValid() || reflectedValue.IsZero() {
		return nil
	}

	switch reflectedValue.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		length := reflectedValue.Len()
		s.Value = make(Vector, length)
		for i := 0; i < length; i++ {
			s.Value[i] = NewScalar(s.Type.(*arrow.ListType).Elem())
			if err := s.Value[i].Set(reflectedValue.Index(i).Interface()); err != nil {
				return err
			}
		}
	}

	s.Valid = true
	return nil
}
