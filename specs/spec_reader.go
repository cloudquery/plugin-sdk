package specs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
)

type SpecReader struct {
	Sources      map[string]*Source
	Destinations map[string]*Destination
}

var fileRegex = regexp.MustCompile(`\$\{file:([^}]+)\}`)
var envRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

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

// expand environment variables in the format ${ENV_VAR}
func expandEnv(cfg []byte) ([]byte, error) {
	var expandErr error
	cfg = envRegex.ReplaceAllFunc(cfg, func(match []byte) []byte {
		envVar := envRegex.FindSubmatch(match)[1]
		content, ok := os.LookupEnv(string(envVar))
		if !ok {
			expandErr = fmt.Errorf("env variable %s not found", envVar)
			return nil
		}
		return []byte(content)
	})

	return cfg, expandErr
}

func (r *SpecReader) loadSpecsFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// support multiple yamls in one file
	// this should work both on Windows and Unix
	normalizedConfig := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

	sections := bytes.Split(normalizedConfig, []byte("\n---\n"))
	for i, doc := range sections {
		doc, err = stripYamlComments(doc)
		if err != nil {
			return fmt.Errorf("failed to strip yaml comments in file %s (section %d): %w", path, i+1, err)
		}
		doc, err = expandFileConfig(doc)
		if err != nil {
			return fmt.Errorf("failed to expand file variable in file %s (section %d): %w", path, i+1, err)
		}
		doc, err = expandEnv(doc)
		if err != nil {
			return fmt.Errorf("failed to expand environment variable in file %s (section %d): %w", path, i+1, err)
		}
		var s Spec
		if err := SpecUnmarshalYamlStrict(doc, &s); err != nil {
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
			// We set the default value to 0, so it can be overridden later by plugins' defaults
			destination.SetDefaults(0, 0)
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
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") &&
			(strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml")) {
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

// strip yaml comments from the given yaml document by converting to JSON and back :)
func stripYamlComments(b []byte) ([]byte, error) {
	const openPlaceholder = "$$$OPEN$$$"
	const closePlaceholder = "$$$CLOSE$$$"

	// return an error if the yaml already contains our temporary placeholder for env variables by some unlucky coincidence
	if bytes.Contains(b, []byte(openPlaceholder)) || bytes.Contains(b, []byte(closePlaceholder)) {
		return nil, fmt.Errorf("%s and %s are reserved words in CloudQuery config", openPlaceholder, closePlaceholder)
	}

	// replace placeholder variables with valid yaml, otherwise it cannot be parsed
	// in some cases. Short of writing our own yaml parser to remove comments,
	// this seems like the best we can do.
	b = envRegex.ReplaceAllFunc(b, func(match []byte) []byte {
		content := envRegex.FindSubmatch(match)[1]
		return []byte(openPlaceholder + string(content) + closePlaceholder)
	})
	j, err := yaml.YAMLToJSON(b)
	if err != nil {
		return nil, err
	}
	b, err = yaml.JSONToYAML(j)
	if err != nil {
		return nil, err
	}
	// place back placeholder variables
	b = bytes.ReplaceAll(b, []byte(openPlaceholder), []byte("${"))
	b = bytes.ReplaceAll(b, []byte(closePlaceholder), []byte("}"))
	return b, nil
}
