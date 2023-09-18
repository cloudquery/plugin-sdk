package scheduler

import (
	"context"
	"math/rand"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *syncClient) syncShuffle(ctx context.Context, resolvedResources chan<- *schema.Resource) {
	// we have this because plugins can return sometimes clients in a random way which will cause
	// differences between this run and the next one.
	preInitialisedClients := make([][]schema.ClientMeta, len(s.tables))
	for i, table := range s.tables {
		clients := []schema.ClientMeta{s.client}
		if table.Multiplex != nil {
			clients = table.Multiplex(s.client)
		}
		preInitialisedClients[i] = clients
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		s.metrics.initWithClients(table, clients)
	}

	// first interleave the tables like in round-robin
	tableClients := roundRobinInterleave(s.tables, preInitialisedClients)

	// then shuffle the tableClients to randomize the order in which they are retrieved
	shuffle(tableClients)

	var wg sync.WaitGroup
	for _, tc := range tableClients {
		table := tc.table
		cl := tc.client
		if err := s.scheduler.tableSems[0].Acquire(ctx, 1); err != nil {
			// This means context was cancelled
			wg.Wait()
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer s.scheduler.tableSems[0].Release(1)
			// not checking for error here as nothing much to do.
			// the error is logged and this happens when context is cancelled
			// This currently uses the DFS algorithm to resolve the tables, but this
			// may change in the future.
			s.resolveTableDfs(ctx, table, cl, nil, resolvedResources, 1)
		}()
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()
}

func shuffle(tableClients []tableClient) {
	// use a fixed seed so that runs with the same tables and clients perform similarly across syncs
	r := rand.New(rand.NewSource(99))
	r.Shuffle(len(tableClients), func(i, j int) {
		tableClients[i], tableClients[j] = tableClients[j], tableClients[i]
	})
}
