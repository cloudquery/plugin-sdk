package scalar

import (
	"encoding/base64"
	"encoding/json"
	"reflect"

	"github.com/apache/arrow/go/v15/arrow"
)

type Struct struct {
	Valid bool
	Value any

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
	// this will check for typed nils as well, so no need to check below
	if IsNil(val) {
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
	case string:
		var x map[string]any
		if err := json.Unmarshal([]byte(value), &x); err != nil {
			return err
		}
		for name := range x {
			if f, ok := s.Type.FieldByName(name); ok {
				xs, ok := x[name].(string)
				if !ok {
					continue
				}
				switch {
				case arrow.TypeEqual(f.Type, arrow.BinaryTypes.Binary):
					v, err := base64.StdEncoding.DecodeString(xs)
					if err != nil {
						return err
					}
					x[name] = v
				case arrow.TypeEqual(f.Type, arrow.BinaryTypes.LargeBinary):
					v, err := base64.StdEncoding.DecodeString(xs)
					if err != nil {
						return err
					}
					x[name] = v
				}
			}
		}
		s.Value = x

	case []byte:
		var x map[string]any
		if err := json.Unmarshal(value, &x); err != nil {
			return err
		}
		s.Value = x

	case *string:
		return s.Set(*value)

	default:
		s.Value = val
	}

	s.Valid = true
	return nil
}

func (s *Struct) DataType() arrow.DataType {
	return s.Type
}
