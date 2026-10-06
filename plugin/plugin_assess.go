package plugin

import (
	"context"
	"errors"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type AssessCategory int

const (
	AssessCategoryUnknown AssessCategory = iota
	AssessCategoryNoChange
	AssessCategoryAutomaticallyMigratable
	AssessCategoryManualMigrationRequired
	AssessCategoryTableRemoved
	AssessCategoryFileSchemaChanged
)

const AssessNotSupportedReason = "destination does not support assessment"

// TablePair holds a table before and after a schema change. Old is nil for an added table, New is nil for a removed one.
type TablePair struct {
	Old *schema.Table
	New *schema.Table
}

func (p TablePair) TableName() string {
	if p.New != nil {
		return p.New.Name
	}
	if p.Old != nil {
		return p.Old.Name
	}
	return ""
}

type AssessOptions struct {
	MigrateForce bool
}

type Evidence struct {
	SyntheticValue string
	Before         string
	After          string
}

type ColumnFinding struct {
	ColumnName         string
	Category           AssessCategory
	OldType            string
	NewType            string
	SafeModeBehavior   string
	ForcedModeBehavior string
	Evidence           []Evidence
}

type TableFinding struct {
	TableName                string
	Category                 AssessCategory
	SafeModeBehavior         string
	ForcedModeBehavior       string
	Columns                  []ColumnFinding
	Evidence                 []Evidence
	IncompleteCoverageReason string
}

// Assessor is an optional DestinationClient interface that reports how schema changes would be applied, without writing anything.
// It is called after Init with NoConnection set, so implementations must not open database or cloud connections.
type Assessor interface {
	AssessTables(ctx context.Context, tables []TablePair, options AssessOptions) ([]TableFinding, error)
}

// AssessTables returns one Unknown finding per table when the client does not implement Assessor.
func (p *Plugin) AssessTables(ctx context.Context, tables []TablePair, options AssessOptions) ([]TableFinding, error) {
	if p.client == nil {
		return nil, errors.New("plugin not initialized. call Init() first")
	}
	if assessor, ok := p.client.(Assessor); ok {
		return assessor.AssessTables(ctx, tables, options)
	}
	findings := make([]TableFinding, len(tables))
	for i, table := range tables {
		findings[i] = TableFinding{
			TableName:                table.TableName(),
			Category:                 AssessCategoryUnknown,
			IncompleteCoverageReason: AssessNotSupportedReason,
		}
	}
	return findings, nil
}
