package types

import "github.com/apache/arrow/go/v16/arrow"

func RegisterAllExtensions() error {
	if err := arrow.RegisterExtensionType(&UUIDType{}); err != nil {
		return err
	}
	if err := arrow.RegisterExtensionType(&JSONType{}); err != nil {
		return err
	}
	if err := arrow.RegisterExtensionType(&InetType{}); err != nil {
		return err
	}
	return arrow.RegisterExtensionType(&MACType{})
}

func UnregisterAllExtensions() error {
	if err := arrow.UnregisterExtensionType(ExtensionTypes.MAC.ExtensionName()); err != nil {
		return err
	}
	if err := arrow.UnregisterExtensionType(ExtensionTypes.Inet.ExtensionName()); err != nil {
		return err
	}
	if err := arrow.UnregisterExtensionType(ExtensionTypes.JSON.ExtensionName()); err != nil {
		return err
	}
	return arrow.UnregisterExtensionType(ExtensionTypes.UUID.ExtensionName())
}
