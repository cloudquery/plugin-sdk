package specs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SpecReader struct {
	sources      map[string]Source
	destinations map[string]Destination
}

func NewSpecReader(directory string) (*SpecReader, error) {
	reader := SpecReader{
		sources:      make(map[string]Source),
		destinations: make(map[string]Destination),
	}
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") && strings.HasSuffix(file.Name(), ".yml") {
			data, err := os.ReadFile(filepath.Join(directory, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}
			var s Spec
			if err := SpecUnmarshalYamlStrict(data, &s); err != nil {
				return nil, fmt.Errorf("failed to unmarshal file %s: %w", file.Name(), err)
			}
			switch s.Kind {
			case KindSource:
				reader.sources[file.Name()] = *s.Spec.(*Source)
			case KindDestination:
				reader.destinations[file.Name()] = *s.Spec.(*Destination)
			default:
				return nil, fmt.Errorf("unknown kind %s", s.Kind)
			}
		}
	}
	return &reader, nil
}

func (s *SpecReader) GetSources() []Source {
	sources := make([]Source, 0, len(s.sources))
	for _, spec := range s.sources {
		sources = append(sources, spec)
	}
	return sources
}

func (s *SpecReader) GetSourceByName(name string) *Source {
	for _, spec := range s.sources {
		if spec.Name == name {
			return &spec
		}
	}
	return nil
}

func (s *SpecReader) GetDestinations() []Destination {
	destinations := make([]Destination, 0, len(s.destinations))
	for _, spec := range s.destinations {
		destinations = append(destinations, spec)
	}
	return destinations
}

func (s *SpecReader) GetDestinationByName(name string) *Destination {
	for _, spec := range s.destinations {
		if spec.Name == name {
			return &spec
		}
	}
	return nil
}
