package plugins

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/glob"
	"github.com/cloudquery/plugin-sdk/specs"
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

func (p *SourcePlugin) validateGlobSpec(spec specs.Source) error {
	flattenedTables := p.tables.FlattenTables()
	for _, includePattern := range spec.Tables {
		matched := false
		for _, table := range flattenedTables {
			if glob.Glob(includePattern, table.Name) {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("no table that matches %s exists", includePattern)
		}
	}
	for _, excludePattern := range spec.SkipTables {
		matched := false
		for _, table := range flattenedTables {
			if glob.Glob(excludePattern, table.Name) {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("no table that matches %s exists", excludePattern)
		}
	}
	return nil
}