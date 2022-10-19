package schema

import "fmt"

func ScanType(v interface{}, t ValueType) error {
	switch t {
		case TypeUUID:
	}
	return fmt.Errorf("cannot scan type %T", v)
}