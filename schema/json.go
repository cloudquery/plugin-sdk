package schema

import (
	"encoding/json"
	"reflect"
)

type Json struct {
	Json  []byte
	Valid bool
}

func (*Json) Type() ValueType {
	return TypeJSON
}

func jsonBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func (dst *Json) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*Json); ok {
		if dst.Valid != other.Valid {
			return false
		}
		t, err := jsonBytesEqual(dst.Json, other.Json)
		// this should never happen because we validate on scan
		if err != nil {
			panic(err)
		}
		return t
	}
	return false
}

func (dst *Json) Scan(src interface{}) error {
	if src == nil {
		*dst = Json{}
		return nil
	}

	switch src := src.(type) {
	case []byte:
		// doing validation
		var res interface{}
		if err := json.Unmarshal(src, &res); err != nil {
			return err
		}
		*dst = Json{Json: src, Valid: true}
	case string:
		// doing validation
		var res interface{}
		if err := json.Unmarshal([]byte(src), &res); err != nil {
			return err
		}
		*dst = Json{Json: []byte(src), Valid: true}
	default:
		// check if type and/or struct implements json.Marshaler
		b, err := json.Marshal(src)
		if err != nil {
			return err
		}
		*dst = Json{Json: b, Valid: true}
	}
	return nil
}
