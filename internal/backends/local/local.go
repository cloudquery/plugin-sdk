package local

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/cloudquery/plugin-sdk/v2/specs"
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
		table, kv, err := l.readFile(name)
		if err != nil {
			return nil, err
		}
		tables[table] = kv
	}
	return tables, nil
}

func (l *Local) readFile(name string) (table string, kv entries, err error) {
	p := path.Join(l.spec.Path, name)
	f, err := os.Open(p)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open state file: %w", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read state file: %w", err)
	}
	err = f.Close()
	if err != nil {
		return "", nil, fmt.Errorf("failed to close state file: %w", err)
	}
	err = json.Unmarshal(b, &kv)
	if err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal state file: %w", err)
	}
	table = strings.TrimPrefix(strings.TrimSuffix(name, ".json"), l.sourceName+"-")
	return table, kv, nil
}

func (l *Local) Get(_ context.Context, table, clientID string) (string, error) {
	l.tablesLock.RLock()
	defer l.tablesLock.RUnlock()

	if _, ok := l.tables[table]; !ok {
		return "", nil
	}
	return l.tables[table][clientID], nil
}

func (l *Local) Set(_ context.Context, table, clientID, value string) error {
	l.tablesLock.Lock()
	defer l.tablesLock.Unlock()

	if _, ok := l.tables[table]; !ok {
		l.tables[table] = map[string]string{}
	}
	prev := l.tables[table][clientID]
	l.tables[table][clientID] = value
	if prev != value {
		// only flush if the value changed
		return l.flushTable(table, l.tables[table])
	}
	return nil
}

func (l *Local) Close(_ context.Context) error {
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
