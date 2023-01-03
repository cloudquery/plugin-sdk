package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/cloudquery/plugin-sdk/specs"
)

type Local struct {
	sourceName string
	spec       Spec
	tables     map[string]entries // table -> key -> value
}

type entries map[string]string

func New(sourceSpec specs.Source) (*Local, error) {
	spec := Spec{}
	err := sourceSpec.UnmarshalBackendSpec(&spec)
	if err != nil {
		return nil, err
	}
	spec.SetDefaults()

	l := &Local{
		sourceName: sourceSpec.Name,
		spec:       spec,
	}
	tables, err := l.loadPreviousState()
	if err != nil {
		return nil, err
	}
	if tables == nil {
		tables = map[string]entries{}
	}
	l.tables = tables
	return l, nil
}

func (l *Local) loadPreviousState() (map[string]entries, error) {
	files, err := os.ReadDir(l.spec.Path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	var tables = map[string]entries{}
	for _, f := range files {
		if f.IsDir() || !f.Type().IsRegular() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(name, ".json") || !strings.HasPrefix(name, l.sourceName+"-") {
			continue
		}
		p := path.Join(l.spec.Path, name)
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("failed to open state file %v: %w", name, err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read state file %v: %w", name, err)
		}
		var kv entries
		err = json.Unmarshal(b, &kv)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal state file %v: %w", name, err)
		}
		table := strings.TrimPrefix(strings.TrimSuffix(name, ".json"), l.sourceName+"-")
		tables[table] = kv
	}
	return tables, nil
}

func (l *Local) Get(table, key string) (string, error) {
	if _, ok := l.tables[table]; !ok {
		return "", nil
	}
	return l.tables[table][key], nil
}

func (l *Local) Set(table, key, value string) error {
	if _, ok := l.tables[table]; !ok {
		l.tables[table] = map[string]string{}
	}
	l.tables[table][key] = value
	return l.flush()
}

func (l *Local) flush() error {
	for table, kv := range l.tables {
		b, err := json.MarshalIndent(kv, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal state for table %v: %w", table, err)
		}
		f := path.Join(l.spec.Path, l.sourceName+"-"+table+".json")
		err = os.WriteFile(f, b, 0644)
		if err != nil {
			return fmt.Errorf("failed to write state for table %v: %w", table, err)
		}
	}
	return nil
}

func (l *Local) Close() error {
	return l.flush()
}
