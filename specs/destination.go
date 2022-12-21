package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

type WriteMode int

type Destination struct {
	Name      string      `json:"name,omitempty"`
	Version   string      `json:"version,omitempty"`
	Path      string      `json:"path,omitempty"`
	Registry  Registry    `json:"registry,omitempty"`
	WriteMode WriteMode   `json:"write_mode,omitempty"`
	BatchTimeout int     `json:"batch_timeout,omitempty"`
	BatchSize int         `json:"batch_size,omitempty"`
	Workers   int         `json:"workers,omitempty"`
	Spec      interface{} `json:"spec,omitempty"`
}

const (
	WriteModeOverwriteDeleteStale WriteMode = iota
	WriteModeOverwrite
	WriteModeAppend
)

const defaultBatchSize = 10000
const defaultWorkers = 1
// batchTimeout is the timeout for a batch to be sent to the destination if no resources are received
const defaultBatchTimeoutSeconds = 20

var (
	writeModeStrings = []string{"overwrite-delete-stale", "overwrite", "append"}
)

func (d *Destination) SetDefaults() {
	if d.Registry.String() == "" {
		d.Registry = RegistryGithub
	}
	if d.BatchSize == 0 {
		d.BatchSize = defaultBatchSize
	}
	if d.Workers == 0 {
		d.Workers = defaultWorkers
	}
	if d.BatchTimeout == 0 {
		d.BatchTimeout = defaultBatchTimeoutSeconds
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
	if d.Path == "" {
		msg := "path is required"
		// give a small hint to help users transition from the old config format that didn't require path
		officialPlugins := []string{"postgresql", "csv"}
		if funk.ContainsString(officialPlugins, d.Name) {
			msg += fmt.Sprintf(". Hint: try setting path to cloudquery/%s in your config", d.Name)
		}
		return fmt.Errorf(msg)
	}

	if d.Registry == RegistryGithub {
		if d.Version == "" {
			return fmt.Errorf("version is required")
		}
		if !strings.HasPrefix(d.Version, "v") {
			return fmt.Errorf("version must start with v")
		}
	}
	if d.BatchSize < 0 {
		return fmt.Errorf("batch_size must be greater than 0")
	}
	if d.Workers < 0 {
		return fmt.Errorf("workers must be greater than 0")
	}
	return nil
}

func (m WriteMode) String() string {
	return writeModeStrings[m]
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
	case "overwrite-delete-stale":
		return WriteModeOverwriteDeleteStale, nil
	}
	return 0, fmt.Errorf("invalid write mode: %s", s)
}
