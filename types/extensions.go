package types

import "github.com/apache/arrow/go/v15/arrow"

var ExtensionTypes = struct {
	UUID arrow.ExtensionType
	Inet arrow.ExtensionType
	MAC  arrow.ExtensionType
	JSON arrow.ExtensionType
}{
	UUID: NewUUIDType(),
	Inet: NewInetType(),
	MAC:  NewMACType(),
	JSON: NewJSONType(),
}
