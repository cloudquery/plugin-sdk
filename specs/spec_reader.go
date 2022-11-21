package specs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SpecReader struct {
	Sources      map[string]*Source
	Destinations map[string]*Destination
}

var fileRegex = regexp.MustCompile(`\$\{file:([^}]+)\}`)

func expandFileConfig(cfg []byte) ([]byte, error) {
	var expandErr error
	cfg = fileRegex.ReplaceAllFunc(cfg, func(match []byte) []byte {
		filename := fileRegex.FindSubmatch(match)[1]
		content, err := os.ReadFile(string(filename))
		if err != nil {
			expandErr = err
			return nil
		}
		return content
	})
	return cfg, expandErr
}

func (r *SpecReader) loadSpecsFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}
	data, err = expandFileConfig(data)
	if err != nil {
		return fmt.Errorf("failed to expand file variable in file %s: %w", path, err)
	}
	data = []byte(os.ExpandEnv(string(data)))

	// support multiple yamls in one file
	// this should work both on Windows and Unix
	normalizedConfig := strings.ReplaceAll(string(data), "\r\n", "\n")
	for _, doc := range strings.Split(normalizedConfig, "\n---\n") {
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

func (r *SpecReader) validate() error {
	if len(r.Sources) == 0 {
		return fmt.Errorf("expecting at least one source")
	}
	if len(r.Destinations) == 0 {
		return fmt.Errorf("expecting at least one destination")
	}

	// here we check if source with different versions use the same destination and error out if yes
	var destinationSourceMap = make(map[string]string)
	for _, source := range r.Sources {
		for _, destination := range source.Destinations {
			if r.Destinations[destination] == nil {
				return fmt.Errorf("source %s references unknown destination %s", source.Name, destination)
			}
			destinationToSourceKey := fmt.Sprintf("%s-%s", destination, source.Path)
			if destinationSourceMap[destinationToSourceKey] == "" {
				destinationSourceMap[destinationToSourceKey] = source.Path + "@" + source.Version
			} else if destinationSourceMap[destinationToSourceKey] != source.Path+"@"+source.Version {
				return fmt.Errorf("destination %s is used by multiple sources %s with different versions", destination, source.Path)
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
			return nil, err
		}
		fileInfo, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, err
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

	if err := reader.validate(); err != nil {
		return nil, err
	}

	return reader, nil
}
