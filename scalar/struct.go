package scalar

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/apache/arrow/go/v16/arrow"
)

type Struct struct {
	Valid bool
	Value map[string]any

	Type *arrow.StructType
}

func (s *Struct) IsValid() bool {
	return s.Valid
}

func (s *Struct) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Struct)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && arrow.TypeEqual(s.Type, r.Type) && reflect.DeepEqual(s.Value, r.Value)
}

func (s *Struct) String() string {
	if !s.Valid {
		return nullValueStr
	}
	b, _ := json.Marshal(s.Value)
	return string(b)
}

func (s *Struct) Get() any {
	return s.Value
}

func (s *Struct) Set(val any) error {
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

	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			s.Value = nil
			s.Valid = false
			return nil
		}
		rv = rv.Elem()
	}

	value := rv.Interface()
	switch value := value.(type) {
	case string:
		var x map[string]any
		if err := json.Unmarshal([]byte(value), &x); err != nil {
			return err
		}
		s.Value = x

	case []byte:
		var x map[string]any
		if err := json.Unmarshal(value, &x); err != nil {
			return err
		}
		s.Value = x

	case *string:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case map[string]any:
		if value == nil {
			s.Valid = false
			return nil
		}
		s.Value = value

	default:
		if rv.Kind() != reflect.Struct {
			s.Valid = false
			return fmt.Errorf("failed to set Struct to the value of type %T", val)
		}
		t := rv.Type()
		m := make(map[string]any, t.NumField())
		for _, sF := range s.Type.Fields() {
			tF, ok := t.FieldByName(sF.Name)
			if !ok {
				return fmt.Errorf("failed to set Struct to the value of type %T: missing field %q", val, sF.Name)
			}
			v := rv.FieldByIndex(tF.Index)
			if v.IsValid() && v.CanInterface() {
				m[sF.Name] = v.Interface()
			} else {
				m[sF.Name] = nil
			}
		}
	}

	s.Valid = true
	return nil
}

func (s *Struct) DataType() arrow.DataType {
	return s.Type
}
