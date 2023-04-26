package types

import "github.com/apache/arrow/go/v12/arrow"

var ExtensionTypes = struct {
	Inet arrow.ExtensionType
	JSON arrow.ExtensionType
	MAC  arrow.ExtensionType
	UUID arrow.ExtensionType
}{
	Inet: NewInetType(),
	JSON: NewJSONType(),
	MAC:  NewMACType(),
	UUID: NewUUIDType(),
}
