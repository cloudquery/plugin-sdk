package plugins

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
)

func (p *SourcePlugin) validate() error {
	if err := p.tables.ValidateDuplicateColumns(); err != nil {
		return fmt.Errorf("found duplicate columns in source plugin: %s: %w", p.name, err)
	}

	if err := p.tables.ValidateDuplicateTables(); err != nil {
		return fmt.Errorf("found duplicate tables in source plugin: %s: %w", p.name, err)
	}

	if err := p.tables.ValidateTableNames(); err != nil {
		return fmt.Errorf("found table with invalid name in source plugin: %s: %w", p.name, err)
	}

	if err := p.tables.ValidateColumnNames(); err != nil {
		return fmt.Errorf("found column with invalid name in source plugin: %s: %w", p.name, err)
	}

	return nil
}

func (p *SourcePlugin) listAndValidateTables(tables, skipTables []string) (schema.Tables, error) {
	if len(tables) == 0 {
		return nil, fmt.Errorf("list of tables is empty")
	}

	// return an error if a table pattern doesn't match any known tables
	var includedTables schema.Tables
	for _, t := range tables {
		tt := p.tables.GlobMatch(t)
		if len(tt) == 0 {
			return nil, fmt.Errorf("tables entry matches no known tables: %q", t)
		}
		for _, ttt := range tt {
			if includedTables.Get(ttt.Name) != nil {
				continue
			}
			includedTables = append(includedTables, ttt)
		}
	}

	// return an error if skip tables doesn't match any known tables
	var skippedTables schema.Tables
	skippedTableMap := map[string]bool{}
	for _, t := range skipTables {
		tt := p.tables.GlobMatch(t)
		if len(tt) == 0 {
			return nil, fmt.Errorf("skip_tables entry matches no known tables: %q", t)
		}
		for _, ttt := range tt {
			if skippedTables.Get(ttt.Name) != nil {
				continue
			}
			skippedTables = append(skippedTables, ttt)
		}
		for _, st := range tt {
			skippedTableMap[st.Name] = true
		}
	}

	// return an error if a table is both explicitly included and skipped
	var remainingTables schema.Tables
	for _, included := range includedTables {
		if skippedTableMap[included.Name] {
			continue
		}
		remainingTables = append(remainingTables, included)
	}

	// return an error if child table is included without its parent
	for _, t := range remainingTables {
		if t.Parent != nil {
			pt := t.Parent
			for pt.Parent != nil {
				pt = pt.Parent
			}

			return nil, fmt.Errorf("table %s is a descendant table and cannot be included without its top-level parent table %s", t.Name, pt.Name)
		}
	}

	return remainingTables, nil
}
