package schema

import "encoding/json"

type Json struct {
	Json []byte
	Valid bool
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