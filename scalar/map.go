package scalar

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/apache/arrow/go/v16/arrow"
)

type Map struct {
	Valid bool
	Value any

	Type *arrow.MapType
}

func (m *Map) IsValid() bool {
	return m.Valid
}

func (m *Map) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Map)
	if !ok {
		return false
	}
	return m.Valid == r.Valid && arrow.TypeEqual(m.Type, r.Type) && reflect.DeepEqual(m.Value, r.Value)
}

func (m *Map) String() string {
	if !m.Valid {
		return nullValueStr
	}
	b, _ := json.Marshal(m.Value)
	return string(b)
}

func (m *Map) Get() any {
	return m.Value
}

func (m *Map) Set(val any) error {
	if val == nil {
		m.Valid = false
		return nil
	}

	if sc, ok := val.(Scalar); ok {
		if !sc.IsValid() {
			m.Valid = false
			return nil
		}
		return m.Set(sc.Get())
	}

	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			m.Value = nil
			m.Valid = false
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
		m.Value = x

	case []byte:
		var x map[string]any
		if err := json.Unmarshal(value, &x); err != nil {
			return err
		}
		m.Value = x

	case *string:
		if value == nil {
			m.Valid = false
			return nil
		}
		return m.Set(*value)

	default:
		if rv.Kind() != reflect.Map {
			m.Valid = false
			return fmt.Errorf("failed to set Map to the value of type %T", val)
		}

		if rv.IsNil() {
			m.Valid = false
			return nil
		}
		m.Value = value
	}

	m.Valid = true
	return nil
}

func (m *Map) DataType() arrow.DataType {
	return m.Type
}

var _ Scalar = (*Map)(nil)
