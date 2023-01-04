package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/cloudquery/plugin-sdk/specs"
)

type Local struct {
	sourceName string
	spec       Spec
	tables     map[string]entries // table -> key -> value
	tablesLock sync.RWMutex
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
	l.tablesLock.RLock()
	defer l.tablesLock.RUnlock()

	if _, ok := l.tables[table]; !ok {
		return "", nil
	}
	return l.tables[table][key], nil
}

func (l *Local) Set(table, key, value string) error {
	l.tablesLock.Lock()
	defer l.tablesLock.Unlock()

	if _, ok := l.tables[table]; !ok {
		l.tables[table] = map[string]string{}
	}
	prev := l.tables[table][key]
	l.tables[table][key] = value
	if prev != value {
		// only flush if the value changed
		return l.flushTable(table, l.tables[table])
	}
	return nil
}

func (l *Local) Close() error {
	l.tablesLock.RLock()
	defer l.tablesLock.RUnlock()

	return l.flush()
}

func (l *Local) flush() error {
	for table, kv := range l.tables {
		err := l.flushTable(table, kv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Local) flushTable(table string, entries entries) error {
	if len(entries) == 0 {
		return nil
	}

	err := os.MkdirAll(l.spec.Path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create state directory %v: %w", l.spec.Path, err)
	}

	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state for table %v: %w", table, err)
	}
	f := path.Join(l.spec.Path, l.sourceName+"-"+table+".json")
	err = os.WriteFile(f, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state for table %v: %w", table, err)
	}

	return nil
}
