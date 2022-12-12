package plugins

import (
	"fmt"
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
