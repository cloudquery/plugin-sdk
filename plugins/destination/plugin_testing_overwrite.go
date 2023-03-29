package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/testdata"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeOverwrite
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_%s_%d", spec.Name, time.Now().Unix())
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*schema.Table{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testOverwriteSource" + uuid.NewString()
	sourceSpec := specs.Source{
		Name: sourceName,
	}

	resources := createTestResources(table, sourceName, syncTime, 2)
	if err := p.writeAll(ctx, sourceSpec, tables, syncTime, resources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	sortResources(table, resources)

	resourcesRead, err := p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortCQTypes(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}

	if diff := resources[0].Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("expected first resource diff: %s", diff)
	}

	if diff := resources[1].Data.Diff(resourcesRead[1]); diff != "" {
		return fmt.Errorf("expected second resource diff: %s", diff)
	}

	secondSyncTime := time.Now().UTC().Round(1 * time.Second).Add(time.Hour).UTC()

	updatedResource := createTestResources(table, sourceName, secondSyncTime, 1)[0]
	for _, colIndex := range []int{2, 3, 7} {
		old := resources[0].Data[colIndex].Get()
		return updatedResource.Data[colIndex].Set(old)
	}

	// write second time
	if err := p.writeOne(ctx, sourceSpec, tables, secondSyncTime, updatedResource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortCQTypes(table, resourcesRead)

	if p.spec.PartitionMinutes > 0 {
		if len(resourcesRead) != 3 {
			return fmt.Errorf("after overwrite with partitioning expected 3 resources, got %d", len(resourcesRead))
		}
		for i, val := range resources[0].Data {
			if i == 1 {
				if val.Equal(resourcesRead[2][i]) {
					return fmt.Errorf("after overwrite with partitioning expected timestamp diff: %s", val.Get())
				}
				continue
			}
			if !val.Equal(resourcesRead[2][i]) {
				return fmt.Errorf("after overwrite with partitioning expected resource diff: %s", val.Get())
			}
		}
	} else {
		if len(resourcesRead) != 2 {
			return fmt.Errorf("after overwrite expected 2 resources, got %d", len(resourcesRead))
		}
		if diff := resources[1].Data.Diff(resourcesRead[0]); diff != "" {
			return fmt.Errorf("after overwrite expected first resource diff: %s", diff)
		}

		if diff := updatedResource.Data.Diff(resourcesRead[1]); diff != "" {
			return fmt.Errorf("after overwrite expected second resource diff: %s", diff)
		}
	}
	return nil
}
