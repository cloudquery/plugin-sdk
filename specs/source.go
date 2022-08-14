package specs

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// SourceSpec is the shared configuration for all source plugins
type SourceSpec struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// Path is the path in the registry
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Registry can be github,local,grpc. Might support things like https in the future.
	Registry      Registry  `json:"registry,omitempty" yaml:"registry,omitempty"`
	MaxGoRoutines uint64    `json:"max_goroutines,omitempty" yaml:"max_goroutines,omitempty"`
	Tables        []string  `json:"tables,omitempty" yaml:"tables,omitempty"`
	SkipTables    []string  `json:"skip_tables,omitempty" yaml:"skip_tables,omitempty"`
	Destinations  []string  `json:"destinations,omitempty" yaml:"destinations,omitempty"`
	Spec          yaml.Node `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (s *SourceSpec) UnmarshalYAML(n *yaml.Node) error {
	type S SourceSpec
	type T struct {
		*S `yaml:",inline"`
	}
	// This is a neat trick to avoid recursion and use unmarshal as a one stop shop for default setting
	obj := &T{S: (*S)(s)}
	if err := n.Decode(&obj); err != nil {
		return err
	}

	// set default
	if s.Registry.String() == "" {
		s.Registry = RegistryGithub
	}
	if s.Path == "" {
		s.Path = s.Name
	}
	if s.Version == "" {
		s.Version = "latest"
	}
	if s.Registry == RegistryGithub && !strings.Contains(s.Path, "/") {
		s.Path = "cloudquery/" + s.Path
	}
	return nil
}
