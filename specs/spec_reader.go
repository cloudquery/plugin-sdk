package specs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SpecReader struct {
	Sources      map[string]*Source
	Destinations map[string]*Destination
}

func (r *SpecReader) loadSpecsFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}
	data = []byte(os.ExpandEnv(string(data)))

	// support multiple yamls in one file
	for _, doc := range strings.Split(string(data), "\n---\n") {
		var s Spec
		if err := SpecUnmarshalYamlStrict([]byte(doc), &s); err != nil {
			return fmt.Errorf("failed to unmarshal file %s: %w", path, err)
		}
		switch s.Kind {
		case KindSource:
			source := s.Spec.(*Source)
			if r.Sources[source.Name] != nil {
				return fmt.Errorf("duplicate source name %s", source.Name)
			}
			source.SetDefaults()
			if err := source.Validate(); err != nil {
				return fmt.Errorf("failed to validate source %s: %w", source.Name, err)
			}
			r.Sources[source.Name] = source
		case KindDestination:
			destination := s.Spec.(*Destination)
			if r.Destinations[destination.Name] != nil {
				return fmt.Errorf("duplicate destination name %s", destination.Name)
			}
			destination.SetDefaults()
			if err := destination.Validate(); err != nil {
				return fmt.Errorf("failed to validate destination %s: %w", destination.Name, err)
			}
			r.Destinations[destination.Name] = destination
		default:
			return fmt.Errorf("unknown kind %s", s.Kind)
		}
	}
	return nil
}

func (r *SpecReader) loadSpecsFromDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", path, err)
	}
	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") && strings.HasSuffix(file.Name(), ".yml") {
			if err := r.loadSpecsFromFile(filepath.Join(path, file.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewSpecReader(paths []string) (*SpecReader, error) {
	reader := &SpecReader{
		Sources:      make(map[string]*Source),
		Destinations: make(map[string]*Destination),
	}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open path %s: %w", path, err)
		}
		fileInfo, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to stat path %s: %w", path, err)
		}
		file.Close()
		if fileInfo.IsDir() {
			if err := reader.loadSpecsFromDir(path); err != nil {
				return nil, err
			}
		} else {
			if err := reader.loadSpecsFromFile(path); err != nil {
				return nil, err
			}
		}
	}

	if len(reader.Sources) == 0 {
		return nil, fmt.Errorf("expecting at least one source in: %v", paths)
	}
	if len(reader.Destinations) == 0 {
		return nil, fmt.Errorf("expecting at least one destination in: %v", paths)
	}
	return reader, nil
}
