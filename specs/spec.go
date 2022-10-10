package specs

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

type Kind int

type Spec struct {
	Kind   Kind        `json:"kind"`
	Plugin interface{} `json:"plugin"`
}

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
		return KindSource, fmt.Errorf("unknown kind %s", s)
	}
}

func (s *Spec) UnmarshalJSON(data []byte) error {
	var t struct {
		Kind   Kind        `json:"kind"`
		Plugin interface{} `json:"plugin"`

		// Deprecated
		Spec interface{} `json:"spec"`
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	dec.UseNumber()
	if err := dec.Decode(&t); err != nil {
		return err
	}
	s.Kind = t.Kind

	if t.Plugin != nil && t.Spec != nil {
		return fmt.Errorf("spec must have either plugin or spec, not both")
	}

	switch s.Kind {
	case KindSource:
		s.Plugin = new(Source)
	case KindDestination:
		s.Plugin = new(Destination)
	default:
		return fmt.Errorf("unknown kind %s", s.Kind)
	}

	var (
		b   []byte
		err error
	)
	if t.Plugin == nil {
		b, err = json.Marshal(t.Spec) // fallback to deprecated field if Plugin is not set
	} else {
		b, err = json.Marshal(t.Plugin)
	}
	if err != nil {
		return err
	}

	dec = json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	return dec.Decode(s.Plugin)
}

func SpecUnmarshalYamlStrict(b []byte, spec *Spec) error {
	jb, err := yaml.YAMLToJSON(b)
	if err != nil {
		return fmt.Errorf("failed to convert yaml to json: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(jb))
	dec.DisallowUnknownFields()
	dec.UseNumber()
	if err := dec.Decode(spec); err != nil {
		return fmt.Errorf("failed to decode spec: %w", err)
	}
	return nil
}
