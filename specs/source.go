package specs

import (
	"encoding/json"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Source is the shared configuration for all source plugins
type Source struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	// Path is the path in the registry
	Path string `json:"path,omitempty"`
	// Registry can be github,local,grpc. Might support things like https in the future.
	Registry      Registry    `json:"registry,omitempty"`
	MaxGoRoutines uint64      `json:"max_goroutines,omitempty"`
	Tables        []string    `json:"tables,omitempty"`
	SkipTables    []string    `json:"skip_tables,omitempty"`
	Destinations  []string    `json:"destinations,omitempty"`
	Spec          interface{} `json:"spec,omitempty"`
}

func (s *Source) SetDefaults() {
	if s.Registry.String() == "" {
		s.Registry = RegistryGithub
	}
	if s.Path == "" {
		s.Path = s.Name
	}
	if s.Version == "" {
		s.Version = "latest"
	}
	if s.Registry == RegistryGithub && !strings.Contains(s.Path, "/") {
		s.Path = "cloudquery/" + s.Path
	}
}

func (s *Source) UnmarshalSpec(out interface{}) error {
	b, err := json.Marshal(s.Spec)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(nil)
	dec.UseNumber()
	dec.DisallowUnknownFields()
	return json.Unmarshal(b, out)
}

func (s *Source) Validate() (*gojsonschema.Result, error) {
	return nil, nil
}
