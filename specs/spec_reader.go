package specs

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
}

func NewSpecReader(directory string) (*SpecReader, error) {
	reader := SpecReader{
		sources:      make(map[string]SourceSpec),
		destinations: make(map[string]DestinationSpec),
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
			default:
				return nil, fmt.Errorf("unknown kind %s", s.Kind)
			}
		}
	}
	return &reader, nil
}

func (s *SpecReader) GetSources() []SourceSpec {
	sources := make([]SourceSpec, 0, len(s.sources))
	for _, spec := range s.sources {
		sources = append(sources, spec)
	}
	return sources
}

func (s *SpecReader) GetSourceByName(name string) *SourceSpec {
	for _, spec := range s.sources {
		if spec.Name == name {
			return &spec
		}
	}
	return nil
}

func (s *SpecReader) GetDestinatinoByName(name string) *DestinationSpec {
	for _, spec := range s.destinations {
		if spec.Name == name {
			return &spec
		}
	}
	return nil
}
