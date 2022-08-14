package specs

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type WriteMode int

const (
	WriteModeAppendOnly WriteMode = iota
	WriteModeOverwrite
)

type DestinationSpec struct {
	Name      string    `yaml:"name,omitempty" json:"name,omitempty"`
	Version   string    `yaml:"version,omitempty" json:"version,omitempty"`
	Path      string    `yaml:"path,omitempty" json:"path,omitempty"`
	Registry  Registry  `yaml:"registry,omitempty" json:"registry,omitempty"`
	WriteMode WriteMode `yaml:"write_mode,omitempty" json:"write_mode,omitempty"`
	Spec      yaml.Node `yaml:"spec,omitempty" json:"spec,omitempty"`
}

func (m WriteMode) String() string {
	return [...]string{"append-only", "overwrite"}[m]
}

func (m WriteMode) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(m.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (m *WriteMode) UnmarshalJSON(data []byte) (err error) {
	var writeMode string
	if err := json.Unmarshal(data, &writeMode); err != nil {
		return err
	}
	if *m, err = WriteModeFromString(writeMode); err != nil {
		return err
	}
	return nil
}

func (m WriteMode) MarshalYAML() (interface{}, error) {
	return m.String(), nil
}

func (r *WriteMode) UnmarshalYAML(n *yaml.Node) (err error) {
	*r, err = WriteModeFromString(n.Value)
	return err
}

func WriteModeFromString(s string) (WriteMode, error) {
	switch s {
	case "append-only":
		return WriteModeAppendOnly, nil
	case "overwrite":
		return WriteModeOverwrite, nil
	}
	return 0, fmt.Errorf("invalid write mode: %s", s)
}
