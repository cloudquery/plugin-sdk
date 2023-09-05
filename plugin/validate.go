package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func validateTables(tables schema.Tables) error {
	if err := tables.ValidateDuplicateTables(); err != nil {
		return fmt.Errorf("found duplicate tables in plugin: %w", err)
	}

	if err := tables.ValidateDuplicateColumns(); err != nil {
		return fmt.Errorf("found duplicate columns in plugin: %w", err)
	}

	return nil
}

func (p *Plugin) validate(ctx context.Context) error {
	tables, err := p.client.Tables(ctx, TableOptions{Tables: []string{"*"}})
	// ErrNotImplemented means it's a destination only plugin
	if err != nil && !errors.Is(err, ErrNotImplemented) {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	return validateTables(tables)
}
