package scheduler

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/queue"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *syncClient) syncPriorityQueue(ctx context.Context, resolvedResources chan<- *schema.Resource) {
	// we have this because plugins can return sometimes clients in a random way which will cause
	// differences between this run and the next one.
	preInitialisedClients := make([][]schema.ClientMeta, len(s.tables))
	tableNames := make([]string, len(s.tables))
	for i, table := range s.tables {
		tableNames[i] = table.Name
		clients := []schema.ClientMeta{s.client}
		if table.Multiplex != nil {
			clients = table.Multiplex(s.client)
		}
		preInitialisedClients[i] = clients
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		s.metrics.InitWithClients(table, clients)
	}

	tableClients := roundRobinInterleave(s.tables, preInitialisedClients)
	seed := hashTableNames(tableNames)
	shuffle(tableClients, seed)

	dispatcher := queue.NewDispatcher(
		s.logger,
		s.metrics,
		queue.WithWorkerCount(s.scheduler.concurrency),
		queue.WithCaser(s.scheduler.caser),
		queue.WithDeterministicCQID(s.deterministicCQID),
		queue.WithInvocationID(s.invocationID),
	)
	queueClients := make([]queue.TableClientPair, 0, len(tableClients))
	for _, tc := range tableClients {
		queueClients = append(queueClients, queue.TableClientPair{
			Table:  tc.table,
			Client: tc.client,
		})
	}
	dispatcher.Dispatch(ctx, queueClients, resolvedResources)
}
