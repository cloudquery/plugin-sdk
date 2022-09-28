package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type WriteMode int

type Destination struct {
	Name      string      `json:"name,omitempty"`
	Version   string      `json:"version,omitempty"`
	Path      string      `json:"path,omitempty"`
	Registry  Registry    `json:"registry,omitempty"`
	WriteMode WriteMode   `json:"write_mode,omitempty"`
	Spec      interface{} `json:"spec,omitempty"`
}

const (
	WriteModeAppend WriteMode = iota
	WriteModeOverwrite
)

func (d *Destination) SetDefaults() {
	if d.Registry.String() == "" {
		d.Registry = RegistryGithub
	}
	if d.Path == "" {
		d.Path = d.Name
	}
	if d.Registry == RegistryGithub && !strings.Contains(d.Path, "/") {
		d.Path = "cloudquery/" + d.Path
	}
}

func (d *Destination) UnmarshalSpec(out interface{}) error {
	b, err := json.Marshal(d.Spec)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func (d *Destination) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}

	if d.Registry == RegistryGithub {
		if d.Version == "" {
			return fmt.Errorf("version is required")
		}
		if !strings.HasPrefix(d.Version, "v") {
			return fmt.Errorf("version must start with v")
		}
	}

	return nil
}

func (m WriteMode) String() string {
	return [...]string{"append", "overwrite"}[m]
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

func WriteModeFromString(s string) (WriteMode, error) {
	switch s {
	case "append":
		return WriteModeAppend, nil
	case "overwrite":
		return WriteModeOverwrite, nil
	}
	return 0, fmt.Errorf("invalid write mode: %s", s)
}
