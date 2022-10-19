package schema

import "fmt"

type String struct {
	String string
	Valid  bool
}

func (dst *String) Scan(src interface{}) error {
	if src == nil {
		*dst = String{}
		return nil
	}

	switch src := src.(type) {
	case string:
		*dst = String{String: src, Valid: true}
		return nil
	case []byte:
		*dst = String{String: string(src), Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}