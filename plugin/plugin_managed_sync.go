package plugin

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (p *Plugin) managedSync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error {
	if len(p.sessionTables) == 0 {
		return fmt.Errorf("no tables to sync - please check your spec 'tables' and 'skip_tables' settings")
	}

	managedClient, err := p.client.NewManagedSyncClient(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to create managed sync client: %w", err)
	}

	resources := make(chan *schema.Resource)
	go func() {
		defer close(resources)
		switch options.Scheduler {
		case SchedulerDFS:
			p.syncDfs(ctx, options, managedClient, p.sessionTables, resources)
		case SchedulerRoundRobin:
			p.syncRoundRobin(ctx, options, managedClient, p.sessionTables, resources)
		default:
			panic(fmt.Errorf("unknown scheduler %s", options.Scheduler))
		}
	}()
	for resource := range resources {
		vector := resource.GetValues()
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, resource.Table.ToArrowSchema())
		scalar.AppendToRecordBuilder(bldr, vector)
		rec := bldr.NewRecord()
		res <- rec
	}
	return nil
}
