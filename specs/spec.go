package specs

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
	default:
		return fmt.Errorf("unknown kind %s", s.Kind)
	}
	return obj.Spec.Decode(s.Spec)
}

func (s Spec) MarshalYAML() (interface{}, error) {
	type T struct {
		Kind string      `yaml:"kind,omitempty"`
		Spec interface{} `yaml:"spec,omitempty"`
	}
	tmp := T{
		Kind: s.Kind,
		Spec: s.Spec,
	}
	return tmp, nil
}
