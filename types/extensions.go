package types

import "github.com/apache/arrow/go/v12/arrow"


var ExtensionTypes = struct {
	UUID arrow.ExtensionType
	Inet arrow.ExtensionType
	Mac arrow.ExtensionType
	JSON arrow.ExtensionType
}{
	UUID: NewUUIDType(),
	Inet: NewInetType(),
	Mac: NewMacType(),
	JSON: NewJSONType(),
}