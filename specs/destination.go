package specs

import "gopkg.in/yaml.v3"

type DestinationSpec struct {
	Name     string    `yaml:"name"`
	Version  string    `yaml:"version"`
	Path     string    `yaml:"path"`
	Registry string    `yaml:"registry"`
	Spec     yaml.Node `yaml:"spec"`
}

func (s *DestinationSpec) Validate() error {
	return nil
}
