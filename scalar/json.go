package scalar

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type JSON struct {
	Valid bool
	Value []byte
}

func (s *JSON) IsValid() bool {
	return s.Valid
}

func (*JSON) DataType() arrow.DataType {
	return types.ExtensionTypes.JSON
}

func (s *JSON) Get() any {
	return s.Value
}

func (s *JSON) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*JSON)
	if !ok {
		return false
	}
	if !s.Valid && !r.Valid {
		return true
	}

	if s.Valid != r.Valid {
		return false
	}

	equal, err := jsonBytesEqual(s.Value, r.Value)
	if err != nil {
		return false
	}
	return equal
}

func (s *JSON) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return string(s.Value)
}

func (s *JSON) Set(val any) error {
	if val == nil {
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
		if value == "" {
			return nil
		}
		if !json.Valid([]byte(value)) {
			return &ValidationError{Type: types.ExtensionTypes.JSON, Msg: "invalid json string", Value: value}
		}
		s.Value = []byte(value)
	case *string:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case []byte:
		if value == nil {
			return nil
		}
		if string(value) == "" {
			return nil
		}

		if !json.Valid(value) {
			return &ValidationError{Type: types.ExtensionTypes.UUID, Msg: "invalid json byte array", Value: value}
		}
		s.Value = value
	// Encode* methods are defined on *JSON. If JSON is passed directly then the
	// struct itself would be encoded instead of Bytes. This is clearly a footgun
	// so detect and return an error. See https://github.com/jackc/pgx/issues/350.
	case JSON:
		return &ValidationError{Type: types.ExtensionTypes.JSON, Msg: "use pointer to JSON instead of value", Value: value}
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
				s.Value = []byte("{}")
				s.Valid = true
				return nil
			}

			if isEmptySlice(value) {
				s.Value = []byte("[]")
				s.Valid = true
				return nil
			}
		}
		s.Value = buf
	}
	s.Valid = true
	return nil
}

func (s *JSON) ByteSize() int64 { return int64(len(s.Value)) }

var (
	_ Scalar = (*JSON)(nil)
)

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
