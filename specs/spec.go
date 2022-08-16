package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Kind int

const (
	KindSource Kind = iota
	KindDestination
)

func (k Kind) String() string {
	return [...]string{"source", "destination"}[k]
}

func (k Kind) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(k.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (k *Kind) UnmarshalJSON(data []byte) (err error) {
	var kind string
	if err := json.Unmarshal(data, &kind); err != nil {
		return err
	}
	if *k, err = KindFromString(kind); err != nil {
		return err
	}
	return nil
}

func KindFromString(s string) (Kind, error) {
	switch s {
	case "source":
		return KindSource, nil
	case "destination":
		return KindDestination, nil
	default:
		return KindSource, fmt.Errorf("unknown registry %s", s)
	}
}

type Spec struct {
	Kind string      `json:"kind"`
	Spec interface{} `json:"-"`
}

func (s *Spec) UnmarshalJSON(data []byte) error {
	type S Spec
	type T struct {
		*S `json:",inline"`
		// Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}

	switch s.Kind {
	case "source":
		s.Spec = new(Source)
	case "destination":
		s.Spec = new(Destination)
	default:
		return fmt.Errorf("unknown kind %s", s.Kind)
	}
	return json.Unmarshal(data, s.Spec)
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
