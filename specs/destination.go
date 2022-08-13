package specs

import "gopkg.in/yaml.v3"

type WriteMode int

const (
	ModeAppendOnly WriteMode = iota
	ModeOverwrite
)

func (m WriteMode) String() string {
	return [...]string{"append-only", "overwrite"}[m]
}

type DestinationSpec struct {
	Name      string    `yaml:"name"`
	Version   string    `yaml:"version"`
	Path      string    `yaml:"path"`
	Registry  Registry  `yaml:"registry"`
	WriteMode WriteMode `yaml:"write_mode"`
	Spec      yaml.Node `yaml:"spec"`
}
