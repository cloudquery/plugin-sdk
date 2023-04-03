package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/testdata"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (s *PluginTestSuite) destinationPluginTestWriteAppend(ctx context.Context, mem memory.Allocator, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeAppend
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := spec.Name
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*schema.Table{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testAppendSource" + uuid.NewString()
	specSource := specs.Source{
		Name: sourceName,
	}
	record1 := testdata.GenTestData(mem, table.ToArrowSchema(), sourceName, syncTime, uuid.Nil, 1)[0]
	defer record1.Release()
	if err := p.writeOne(ctx, specSource, tables, syncTime, record1); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	secondSyncTime := syncTime.Add(10 * time.Second).UTC()
	record2 := testdata.GenTestData(mem, table.ToArrowSchema(), sourceName, secondSyncTime, uuid.Nil, 1)[0]
	defer record2.Release()

	if !s.tests.SkipSecondAppend {
		// write second time
		if err := p.writeOne(ctx, specSource, tables, secondSyncTime, record2); err != nil {
			return fmt.Errorf("failed to write one second time: %w", err)
		}
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)

	expectedResource := 2
	if s.tests.SkipSecondAppend {
		expectedResource = 1
	}

	if len(resourcesRead) != expectedResource {
		return fmt.Errorf("expected %d resources, got %d", expectedResource, len(resourcesRead))
	}

	if !array.RecordEqual(record1, resourcesRead[0]) {
		diff := RecordDiff(record1, resourcesRead[0])
		return fmt.Errorf("first expected resource diff: %s", diff)
	}

	if !s.tests.SkipSecondAppend {
		if !array.RecordEqual(record2, resourcesRead[1]) {
			diff := RecordDiff(record2, resourcesRead[1])
			return fmt.Errorf("second expected resource diff: %s", diff)
		}
	}

	return nil
}
