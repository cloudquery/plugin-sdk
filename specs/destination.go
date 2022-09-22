package specs

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"
)

//go:embed templates/destination.go.tpl
var destinationExampleTemplate string

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
	if d.Version == "" {
		d.Version = "latest"
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
	dec := json.NewDecoder(nil)
	dec.UseNumber()
	dec.DisallowUnknownFields()
	return json.Unmarshal(b, out)
}

func (s *Destination) WriteExample(w io.Writer) error {
	tmpSource := *s
	if tmpSource.Registry == RegistryGithub && strings.HasPrefix(tmpSource.Path, "cloudquery/") {
		tmpSource.Spec = fmt.Sprintf("Check documentation here: https://github.com/cloudquery/cloudquery/tree/main/cli/internal/destinations/%s", tmpSource.Name)
	} else if tmpSource.Registry == RegistryGithub {
		splitPath := strings.Split(tmpSource.Path, "/")
		if len(splitPath) != 2 {
			return fmt.Errorf("invalid path: %s", tmpSource.Path)
		}
		tmpSource.Spec = fmt.Sprintf("Check documentation here: https://github.com/%s/cq-destination-%s", splitPath[0], splitPath[1])
	}
	tpl, err := template.New("sourceTemplate").Parse(sourceExampleTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse source template: %w", err)
	}
	if err := tpl.Execute(w, tmpSource); err != nil {
		return fmt.Errorf("failed to execute source template: %w", err)
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


