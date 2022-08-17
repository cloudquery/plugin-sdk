package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type WriteMode int

const (
	WriteModeAppendOnly WriteMode = iota
	WriteModeOverwrite
)

type Destination struct {
	Name      string      `json:"name,omitempty"`
	Version   string      `json:"version,omitempty"`
	Path      string      `json:"path,omitempty"`
	Registry  Registry    `json:"registry,omitempty"`
	WriteMode WriteMode   `json:"write_mode,omitempty"`
	Spec      interface{} `json:"spec,omitempty"`
}

func (d *Destination) SetDefaults() {
	if d.Registry.String() == "" {
		d.Registry = RegistryGithub
	}
	if d.Path == "" {
		d.Path = d.Name
	}
	if d.Version == "" {
		d.Version = "latest"
	}
	if d.Registry == RegistryGithub && !strings.Contains(d.Path, "/") {
		d.Path = "cloudquery/" + d.Path
	}
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

func WriteModeFromString(s string) (WriteMode, error) {
	switch s {
	case "append-only":
		return WriteModeAppendOnly, nil
	case "overwrite":
		return WriteModeOverwrite, nil
	}
	return 0, fmt.Errorf("invalid write mode: %s", s)
}
