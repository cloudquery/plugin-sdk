package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Registry int

const (
	RegistryGithub Registry = iota
	RegistryLocal
	RegistryGrpc
)

func (r Registry) String() string {
	return [...]string{"github", "local", "grpc"}[r]
}

func (r Registry) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// This to implement to implement flag interface so we can use it with cobra
// https://github.com/spf13/pflag/blob/master/flag.go#L187
func (r Registry) Type() string {
	return "string"
}

func (r *Registry) Set(value string) error {
	switch value {
	case "github":
		*r = RegistryGithub
	case "local":
		*r = RegistryLocal
	case "grpc":
		*r = RegistryGrpc
	default:
		return fmt.Errorf("invalid registry %s", value)
	}
	return nil
}

func (r *Registry) UnmarshalJSON(data []byte) (err error) {
	var registry string
	if err := json.Unmarshal(data, &registry); err != nil {
		return err
	}
	if *r, err = RegistryFromString(registry); err != nil {
		return err
	}
	return nil
}

func RegistryFromString(s string) (Registry, error) {
	switch s {
	case "github":
		return RegistryGithub, nil
	case "local":
		return RegistryLocal, nil
	case "grpc":
		return RegistryGrpc, nil
	default:
		return RegistryGithub, fmt.Errorf("unknown registry %s", s)
	}
}
