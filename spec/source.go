package spec

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// SourceSpec is the shared configuration for all source plugins
type SourceSpec struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	// Path is the path in the registry
	Path string `json:"path" yaml:"path"`
	// Registry can be github,local,grpc. Might support things like https in the future.
	Registry      string    `json:"registry" yaml:"registry"`
	MaxGoRoutines uint64    `json:"max_goroutines" yaml:"max_goroutines"`
	Tables        []string  `json:"tables" yaml:"tables"`
	SkipTables    []string  `json:"skip_tables" yaml:"skip_tables"`
	Spec          yaml.Node `json:"spec" yaml:"spec"`
}

func (s *SourceSpec) UnmarshalYAML(n *yaml.Node) error {
	if err := n.Decode(s); err != nil {
		return err
	}
	// set default
	if s.Registry == "" {
		s.Registry = "github"
	}
	if s.Path == "" {
		s.Path = s.Name
	}
	if s.Version == "" {
		s.Version = "latest"
	}
	if !strings.Contains(s.Path, "/") {
		s.Path = "cloudquery/" + s.Path
	}
	return nil
}
