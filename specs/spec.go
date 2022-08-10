package specs

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"-"`
}

type validator interface {
	Validate() error
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
	if err := obj.Spec.Decode(s.Spec); err != nil {
		return err
	}
	if v, ok := s.Spec.(validator); !ok {
		return nil
	} else {
		return v.Validate()
	}
}
