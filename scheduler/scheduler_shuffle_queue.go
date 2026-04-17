package scheduler

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/queue"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage/inmemory"
)

func (s *syncClient) syncShuffleQueue(ctx context.Context, resolvedResources chan<- *schema.Resource) {
	preInitialisedClients := make([][]schema.ClientMeta, len(s.tables))
	tableNames := make([]string, len(s.tables))
	for i, table := range s.tables {
		tableNames[i] = table.Name
		clients := []schema.ClientMeta{s.client}
		if table.Multiplex != nil {
			clients = table.Multiplex(s.client)
		}
		preInitialisedClients[i] = clients
		s.metrics.InitWithClients(table, clients)
	}

	tableClients := roundRobinInterleave(s.tables, preInitialisedClients)
	tableClients = shardTableClients(tableClients, s.shard)
	seed := hashTableNames(tableNames)
	shuffle(tableClients, seed)

	// Storage: use the scheduler-provided storage, or construct an in-memory
	// default if none was configured. This preserves backward compatibility —
	// users not setting spec.queue get the same random-pop in-memory queue as
	// before.
	store := s.scheduler.storage
	if store == nil {
		store = inmemory.New(seed)
		defer func() {
			if err := store.Close(ctx); err != nil {
				s.logger.Warn().Err(err).Msg("failed to close in-memory storage")
			}
		}()
	}

	codec := queue.NewCodec(s.tables.FlattenTables())

	scheduler := queue.NewShuffleQueueScheduler(
		s.logger,
		s.metrics,
		seed,
		queue.WithWorkerCount(s.scheduler.concurrency),
		queue.WithCaser(s.scheduler.caser),
		queue.WithDeterministicCQID(s.deterministicCQID),
		queue.WithInvocationID(s.invocationID),
		queue.WithStorage(store),
		queue.WithCodec(codec),
	)
	queueClients := make([]queue.WorkUnit, 0, len(tableClients))
	for _, tc := range tableClients {
		queueClients = append(queueClients, queue.WorkUnit{
			Table:  tc.table,
			Client: tc.client,
		})
	}
	scheduler.Sync(ctx, queueClients, resolvedResources, s.msgChan)
}
