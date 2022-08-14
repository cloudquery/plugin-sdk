package specs

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
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

func (r Registry) MarshalYAML() (interface{}, error) {
	return r.String(), nil
}

func (r *Registry) UnmarshalYAML(n *yaml.Node) (err error) {
	*r, err = RegistryFromString(n.Value)
	return err
}

// func (r *Registry) UnmarshalYaml(data []byte) (err error) {
// 	var registry string
// 	if err := yaml.Unmarshal(data, &registry); err != nil {
// 		return err
// 	}
// 	if *r, err = RegistryFromString(registry); err != nil {
// 		return err
// 	}
// 	return nil
// }

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
