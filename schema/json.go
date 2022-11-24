package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

type JSONTransformer interface {
	TransformJSON(*JSON) interface{}
}

type JSON struct {
	Bytes  []byte
	Status Status
}

func (*JSON) Type() ValueType {
	return TypeJSON
}

func (dst *JSON) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*JSON)
	if !ok {
		return false
	}

	return dst.Status == s.Status && bytes.Equal(dst.Bytes, s.Bytes)
}

func (dst *JSON) String() string {
	return string(dst.Bytes)
}

func (dst *JSON) Set(src interface{}) error {
	if src == nil {
		*dst = JSON{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		*dst = JSON{Bytes: []byte(value), Status: Present}
	case *string:
		if value == nil {
			*dst = JSON{Status: Null}
		} else {
			*dst = JSON{Bytes: []byte(*value), Status: Present}
		}
	case []byte:
		if value == nil {
			*dst = JSON{Status: Null}
		} else {
			*dst = JSON{Bytes: value, Status: Present}
		}

	// Encode* methods are defined on *JSON. If JSON is passed directly then the
	// struct itself would be encoded instead of Bytes. This is clearly a footgun
	// so detect and return an error. See https://github.com/jackc/pgx/issues/350.
	case JSON:
		return errors.New("use pointer to JSON instead of value")

	default:
		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}

		// For map and slice jsons, it is easier for users to work with '[]' or '{}' instead of JSON's 'null'.
		if bytes.Equal(buf, []byte(`null`)) {
			if isEmptyStringMap(value) {
				*dst = JSON{Bytes: []byte("{}"), Status: Present}
				return nil
			}

			if isEmptySlice(value) {
				*dst = JSON{Bytes: []byte("[]"), Status: Present}
				return nil
			}
		}

		*dst = JSON{Bytes: buf, Status: Present}
	}

	return nil
}

func (dst JSON) Get() interface{} {
	switch dst.Status {
	case Present:
		var i interface{}
		err := json.Unmarshal(dst.Bytes, &i)
		if err != nil {
			return dst
		}
		return i
	case Null:
		return nil
	default:
		return dst.Status
	}
}

// isEmptyStringMap returns true if the value is a map from string to any (i.e. map[string]interface{}).
// We need to use reflection for this, because it impossible to type-assert a map[string]string into a
// map[string]interface{}. See https://go.dev/doc/faq#convert_slice_of_interface.
func isEmptyStringMap(value interface{}) bool {
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return false
	}

	if reflect.TypeOf(value).Key().Kind() != reflect.String {
		return false
	}

	return reflect.ValueOf(value).Len() == 0
}

// isEmptySlice returns true if the value is a slice (i.e. []interface{}).
// We need to use reflection for this, because it impossible to type-assert a map[string]string into a
// map[string]interface{}. See https://go.dev/doc/faq#convert_slice_of_interface.
func isEmptySlice(value interface{}) bool {
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return false
	}

	return reflect.ValueOf(value).Len() == 0
}
