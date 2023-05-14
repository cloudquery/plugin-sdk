package types

import "github.com/apache/arrow/go/v13/arrow"

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
	if err := arrow.RegisterExtensionType(&MacType{}); err != nil {
		return err
	}
	return nil
}
