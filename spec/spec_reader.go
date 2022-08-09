package spec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type SpecReader struct {
	sources      map[string]SourceSpec
	destinations map[string]DestinationSpec
	connections  map[string]ConnectionSpec
}

func NewSpecReader(directory string) (*SpecReader, error) {
	reader := SpecReader{
		sources:      make(map[string]SourceSpec),
		destinations: make(map[string]DestinationSpec),
		connections:  make(map[string]ConnectionSpec),
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".cq.yml") {
			data, err := os.ReadFile(filepath.Join(directory, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}
			var s Spec
			if err := yaml.Unmarshal(data, &s); err != nil {
				return nil, fmt.Errorf("failed to unmarshal file %s: %w", file.Name(), err)
			}
			switch s.Kind {
			case "source":
				reader.sources[file.Name()] = *s.Spec.(*SourceSpec)
			case "destination":
				reader.destinations[file.Name()] = *s.Spec.(*DestinationSpec)
			case "connection":
				reader.connections[file.Name()] = *s.Spec.(*ConnectionSpec)
			default:
				return nil, fmt.Errorf("unknown kind %s", s.Kind)
			}
		}
	}
	return &reader, nil
}

func (s *SpecReader) GetSourceByName(name string) SourceSpec {
	for _, spec := range s.sources {
		if spec.Name == name {
			return spec
		}
	}
	return SourceSpec{}
}

func (s *SpecReader) GetDestinatinoByName(name string) DestinationSpec {
	for _, spec := range s.destinations {
		if spec.Name == name {
			return spec
		}
	}
	return DestinationSpec{}
}

func (s *SpecReader) GetConnectionByName(name string) ConnectionSpec {
	for _, spec := range s.connections {
		if spec.Source == name {
			return spec
		}
	}
	return ConnectionSpec{}
}

func (s *SpecReader) Connections() []ConnectionSpec {
	var connections []ConnectionSpec
	for _, spec := range s.connections {
		connections = append(connections, spec)
	}
	return connections
}
