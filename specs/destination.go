package specs

import "gopkg.in/yaml.v3"

type DestinationSpec struct {
	Name     string    `yaml:"name"`
	Version  string    `yaml:"version"`
	Path     string    `yaml:"path"`
	Registry string    `yaml:"registry"`
	Spec     yaml.Node `yaml:"spec"`
}

func (_ *DestinationSpec) Validate() error {
	// TODO post process DestinationSpec after unmarshalling here
	return nil
}
