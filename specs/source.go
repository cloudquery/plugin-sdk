package specs

import (
	"encoding/json"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Source is the spec for a source plugin
type Source struct {
	// Name of the source plugin to use
	Name string `json:"name,omitempty"`
	// Version of the source plugin to use
	Version string `json:"version,omitempty"`
	// Path is the canonical path to the source plugin in a given registry
	// For example:
	// in github the path will be: org/repo
	// For the local registry the path will be the path to the binary: ./path/to/binary
	// For the gRPC registry the path will be the address of the gRPC server: host:port
	Path string `json:"path,omitempty"`
	// Registry can be github,local,grpc.
	Registry      Registry `json:"registry,omitempty"`
	MaxGoRoutines uint64   `json:"max_goroutines,omitempty"`
	// Tables to sync from the source plugin
	Tables []string `json:"tables,omitempty"`
	// SkipTables defines tables to skip when syncing data. Useful if a glob pattern is used in Tables
	SkipTables []string `json:"skip_tables,omitempty"`
	// Destinations are the names of destination plugins to send sync data to
	Destinations []string `json:"destinations,omitempty"`
	// Spec defines plugin specific configuration
	// This is different in every source plugin.
	Spec interface{} `json:"spec,omitempty"`
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

// UnmarshalSpec unmarshals the internal spec into the given interface
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

func (*Source) Validate() (*gojsonschema.Result, error) {
	return nil, nil
}
