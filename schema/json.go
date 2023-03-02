package schema

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type JSONTransformer interface {
	TransformJSON(*JSON) any
}

type JSON struct {
	Bytes  []byte
	Status Status
}

func (dst *JSON) GetStatus() Status {
	return dst.Status
}

func (*JSON) Type() ValueType {
	return TypeJSON
}

func (dst *JSON) Size() int {
	return len(dst.Bytes)
}

// JSONBytesEqual compares the JSON in two byte slices.
func jsonBytesEqual(a, b []byte) (bool, error) {
	var j, j2 any
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func (dst *JSON) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*JSON)
	if !ok {
		return false
	}
	if dst.Status == s.Status && dst.Status != Present {
		return true
	}

	if dst.Status != s.Status {
		return false
	}

	equal, err := jsonBytesEqual(dst.Bytes, s.Bytes)
	if err != nil {
		return false
	}
	return equal
}

func (dst *JSON) String() string {
	return string(dst.Bytes)
}

func (dst *JSON) Set(src any) error {
	if src == nil {
		*dst = JSON{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() any }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		if value == "" {
			*dst = JSON{Bytes: []byte(""), Status: Null}
			return nil
		}
		if !json.Valid([]byte(value)) {
			return &ValidationError{Type: TypeJSON, Msg: "invalid json string", Value: value}
		}
		*dst = JSON{Bytes: []byte(value), Status: Present}
	case *string:
		if value == nil {
			*dst = JSON{Status: Null}
		} else {
			if *value == "" {
				*dst = JSON{Bytes: []byte(""), Status: Null}
				return nil
			}
			if !json.Valid([]byte(*value)) {
				return &ValidationError{Type: TypeJSON, Msg: "invalid json string pointer", Value: value}
			}
			*dst = JSON{Bytes: []byte(*value), Status: Present}
		}
	case []byte:
		if value == nil {
			*dst = JSON{Status: Null}
		} else {
			if string(value) == "" {
				*dst = JSON{Bytes: []byte(""), Status: Null}
				return nil
			}

			if !json.Valid(value) {
				return &ValidationError{Type: TypeJSON, Msg: "invalid json byte array", Value: value}
			}
			*dst = JSON{Bytes: value, Status: Present}
		}

	// Encode* methods are defined on *JSON. If JSON is passed directly then the
	// struct itself would be encoded instead of Bytes. This is clearly a footgun
	// so detect and return an error. See https://github.com/jackc/pgx/issues/350.
	case JSON:
		return &ValidationError{Type: TypeJSON, Msg: "use pointer to JSON instead of value", Value: value}

	default:
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(value)
		if err != nil {
			return err
		}

		// JSON encoder adds a newline to the end of the output that we don't want.
		buf := bytes.TrimSuffix(buffer.Bytes(), []byte("\n"))
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

func (dst JSON) Get() any {
	switch dst.Status {
	case Present:
		var i any
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

// isEmptyStringMap returns true if the value is a map from string to any (i.e. map[string]any).
// We need to use reflection for this, because it impossible to type-assert a map[string]string into a
// map[string]any. See https://go.dev/doc/faq#convert_slice_of_interface.
func isEmptyStringMap(value any) bool {
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return false
	}

	if reflect.TypeOf(value).Key().Kind() != reflect.String {
		return false
	}

	return reflect.ValueOf(value).Len() == 0
}

// isEmptySlice returns true if the value is a slice (i.e. []any).
// We need to use reflection for this, because it impossible to type-assert a map[string]string into a
// map[string]any. See https://go.dev/doc/faq#convert_slice_of_interface.
func isEmptySlice(value any) bool {
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return false
	}

	return reflect.ValueOf(value).Len() == 0
}
