package plugin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/santhosh-tekuri/jsonschema/v5"
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
	if p.skipTableValidation {
		return nil
	}
	tables, err := p.client.Tables(ctx, TableOptions{Tables: []string{"*"}})
	// ErrNotImplemented means it's a destination only plugin
	if err != nil && !errors.Is(err, ErrNotImplemented) {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	return validateTables(tables)
}

func JSONSchemaValidator(jsonSchema string) (*jsonschema.Schema, error) {
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	c.AssertFormat = true
	if err := c.AddResource("schema.json", strings.NewReader(jsonSchema)); err != nil {
		return nil, err
	}
	return c.Compile("schema.json")
}
