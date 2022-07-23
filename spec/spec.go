package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"-"`
}

func (s *Spec) UnmarshalYAML(n *yaml.Node) error {
	type S Spec
	type T struct {
		*S   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := n.Decode(obj); err != nil {
		return err
	}
	switch s.Kind {
	case "source":
		s.Spec = new(SourceSpec)
	case "destination":
		s.Spec = new(DestinationSpec)
	case "connection":
		s.Spec = new(ConnectionSpec)
	default:
		return fmt.Errorf("unknown kind %s", s.Kind)
	}
	return obj.Spec.Decode(s.Spec)
}
