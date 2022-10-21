package plugins

import (
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

func (p *SourcePlugin) validate() error {
	if err := p.tables.ValidateDuplicateColumns(); err != nil {
		return fmt.Errorf("found duplicate columns in source plugin: %s: %w", p.name, err)
	}

	if err := p.tables.ValidateDuplicateTables(); err != nil {
		return fmt.Errorf("found duplicate tables in source plugin: %s: %w", p.name, err)
	}

	return nil
}

func (p *SourcePlugin) listAndValidateTables(tables, skipTables []string) ([]string, error) {
	if len(tables) == 0 {
		return nil, fmt.Errorf("list of tables is empty")
	}

	// return an error if skip tables contains a wildcard or glob pattern
	for _, t := range skipTables {
		if strings.Contains(t, "*") {
			return nil, fmt.Errorf("glob matching in skipped table name %q is not supported", t)
		}
	}

	// handle wildcard entry
	if funk.Equal(tables, []string{"*"}) {
		allResources := make([]string, 0, len(p.tables))
		for _, k := range p.tables {
			if funk.ContainsString(skipTables, k.Name) {
				continue
			}
			allResources = append(allResources, k.Name)
		}
		return allResources, nil
	}

	// wildcard should not be combined with other tables
	if funk.ContainsString(tables, "*") {
		return nil, fmt.Errorf("wildcard \"*\" table not allowed with explicit tables")
	}

	// return an error if other kinds of glob-matching is detected
	for _, t := range tables {
		if strings.Contains(t, "*") {
			return nil, fmt.Errorf("glob matching in table name %q is not supported", t)
		}
	}

	// return an error if a table is both explicitly included and skipped
	for _, t := range tables {
		if funk.ContainsString(skipTables, t) {
			return nil, fmt.Errorf("table %s cannot be both included and skipped", t)
		}
	}

	// return an error if a given table name doesn't match any known tables
	for _, t := range tables {
		if !funk.ContainsString(p.tables.TableNames(), t) {
			return nil, fmt.Errorf("name %s does not match any known table names", t)
		}
	}

	// return an error if child table is included, but not its parent table
	selectedTables := map[string]bool{}
	for _, t := range tables {
		selectedTables[t] = true
	}
	for _, t := range tables {
		for _, tt := range p.tables {
			if tt.Name != t {
				continue
			}
			if tt.Parent != nil && !selectedTables[tt.Parent.Name] {
				return nil, fmt.Errorf("table %s is a child table, and requires its parent table %s to also be synced", t, tt.Parent.Name)
			}
		}
	}

	return tables, nil
}
