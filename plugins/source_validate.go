package plugins

import (
	"fmt"
	"strings"

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

// listAndValidateTables returns all the tables matched by the `tables` and `skip_tables` config settings.
// It will return ALL tables, including descendent tables. Callers should take care to only use the top-level
// tables if that is what they need.
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
				// prevent duplicates
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
				// prevent duplicates
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
		var missingParents []string
		pt := t
		for pt.Parent != nil {
			if includedTables.Get(pt.Parent.Name) == nil {
				missingParents = append(missingParents, pt.Parent.Name)
			}
			pt = pt.Parent
		}
		if len(missingParents) > 0 {
			return nil, fmt.Errorf("table %s is a descendant table and cannot be included without %s", t.Name, strings.Join(missingParents, ", "))
		}
	}

	return remainingTables, nil
}
