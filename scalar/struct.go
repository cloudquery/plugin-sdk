package scalar

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/arrow/go/v16/arrow"
)

// NameTransformer is used to alter the names of the struct fields.
// string    | bool  | action
// ""        | true  | skip field
// non-empty | true  | use new value
// any       | false | use field name
type NameTransformer func(field reflect.StructField) (string, bool)

// jsonTagTransformer is the default NameTransformer to be used
func jsonTagTransformer(field reflect.StructField) (string, bool) {
	tag, ok := field.Tag.Lookup("json")
	if !ok {
		return "", false
	}

	return strings.SplitN(tag, ",", 2)[0], true
}

type Struct struct {
	Valid bool
	Value map[string]any

	Type *arrow.StructType

	// This is used when traversing structs
	nameTransformer NameTransformer
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
		switch rv.Kind() {
		case reflect.Map:
			// map[string]??? is OK
			t := rv.Type()
			if t.Key().Kind() != reflect.String {
				s.Valid = false
				return fmt.Errorf("failed to set Struct to the value of type %T", val)
			}
			m := make(map[string]any, s.Type.NumFields())
			for _, sF := range s.Type.Fields() {
				v := rv.MapIndex(reflect.ValueOf(sF.Name))
				if v.IsValid() && v.CanInterface() {
					m[sF.Name] = v.Interface()
				} else {
					m[sF.Name] = nil
				}
			}
		case reflect.Struct:
			transformer := s.nameTransformer
			if transformer == nil {
				transformer = jsonTagTransformer
			}
			t := rv.Type()
			m := make(map[string]any, s.Type.NumFields())
			// prefill
			for _, f := range s.Type.Fields() {
				m[f.Name] = nil
			}
			// fill
			for i := 0; i < t.NumField(); i++ {
				tF := t.Field(i)
				name, ok := transformer(tF)
				if !ok {
					name = tF.Name
				}
				if name == "" {
					continue
				}
				if _, ok := s.Type.FieldByName(name); !ok {
					continue
				}
				v := rv.FieldByIndex(tF.Index)
				if v.IsValid() && v.CanInterface() {
					m[name] = v.Interface()
				}
			}
			s.Value = m
		default:
			s.Valid = false
			return fmt.Errorf("failed to set Struct to the value of type %T", val)
		}
	}

	s.Valid = true
	return nil
}

func (s *Struct) DataType() arrow.DataType {
	return s.Type
}
