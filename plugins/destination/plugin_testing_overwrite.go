package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow/array"
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

	resources := createTestResources(schema.CQSchemaToArrow(table), sourceName, syncTime, 2)
	// st := array.RecordToStructArray(resources)
	// tbl := array.NewTableFromRecords(schema.CQSchemaToArrow(table), []arrow.Record{resources})
	// array.NewRecord()
	if err := p.writeAll(ctx, sourceSpec, tables, syncTime, resources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	// sortResources(table, resources)

	resourcesRead, err := p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	// sortCQTypes(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}
	
	if array.RecordEqual(resources[0], resourcesRead[0]) {
		return fmt.Errorf("expected first resource to be equal")
	}
	// if diff := resources[0].Data.Diff(resourcesRead[0]); diff != "" {
	// 	return fmt.Errorf("expected first resource diff: %s", diff)
	// }
	
	if array.RecordEqual(resources[1], resourcesRead[1]) {
		return fmt.Errorf("expected second resource to be equal")
	}
	// if diff := resources[1].Data.Diff(resourcesRead[1]); diff != "" {
	// 	return fmt.Errorf("expected second resource diff: %s", diff)
	// }

	secondSyncTime := syncTime.Add(time.Second).UTC()

	// copy first resource but update the sync time
	updatedResource := createTestResources(schema.CQSchemaToArrow(table), sourceName, secondSyncTime, 1)[0]
	// write second time
	if err := p.writeOne(ctx, sourceSpec, tables, secondSyncTime, updatedResource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	// sortCQTypes(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("after overwrite expected 2 resources, got %d", len(resourcesRead))
	}

	if array.RecordEqual(resources[1], resourcesRead[0]) {
		return fmt.Errorf("after overwrite expected first resource to be equal")
	}
	// if diff := resources[1].Data.Diff(resourcesRead[0]); diff != "" {
	// 	return fmt.Errorf("after overwrite expected first resource diff: %s", diff)
	// }
	if array.RecordEqual(updatedResource, resourcesRead[1]) {
		return fmt.Errorf("after overwrite expected second resource to be equal")
	}
	// if diff := updatedResource.Data.Diff(resourcesRead[1]); diff != "" {
	// 	return fmt.Errorf("after overwrite expected second resource diff: %s", diff)
	// }

	return nil
}
