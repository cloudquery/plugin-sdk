package types

import "github.com/apache/arrow/go/v17/arrow"

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
