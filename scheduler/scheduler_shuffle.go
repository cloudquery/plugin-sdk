package scheduler

import (
	"context"
	"hash/fnv"
	"math/rand"
	"strings"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *syncClient) syncShuffle(ctx context.Context, resolvedResources chan<- *schema.Resource) {
	// We have this because plugins can return sometimes clients in a random way which will cause
	// differences between this run and the next one.
	preInitialisedClients := make([][]schema.ClientMeta, len(s.tables))
	tableNames := make([]string, len(s.tables))
	for i, table := range s.tables {
		tableNames[i] = table.Name
		clients := []schema.ClientMeta{s.client}
		if table.Multiplex != nil {
			clients = table.Multiplex(s.client)
		}
		// Detect duplicate clients while multiplexing
		seenClients := make(map[string]bool)
		for _, c := range clients {
			if _, ok := seenClients[c.ID()]; !ok {
				seenClients[c.ID()] = true
			} else {
				s.logger.Warn().Str("client", c.ID()).Str("table", table.Name).Msg("multiplex returned duplicate client")
			}
		}
		preInitialisedClients[i] = clients
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		s.metrics.initWithClients(table, clients)
	}

	// First interleave the tables like in round-robin
	tableClients := roundRobinInterleave(s.tables, preInitialisedClients)

	// Then shuffle the tableClients to randomize the order in which they are retrieved.
	// We use a fixed seed so that runs with the same tables and clients perform similarly across syncs
	// however, if the table order changes, the seed will change and the shuffle order will be different,
	// so users have a little bit of control over the randomization.
	seed := hashTableNames(tableNames)
	shuffle(tableClients, seed)

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
			// Not checking for error here as nothing much to do.
			// the error is logged and this happens when context is cancelled.
			// This currently uses the DFS algorithm to resolve the tables, but this
			// may change in the future.
			s.resolveTableDfs(ctx, table, cl, nil, resolvedResources, 1)
		}()
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()
}

func hashTableNames(tableNames []string) int64 {
	h := fnv.New32a()
	h.Write([]byte(strings.Join(tableNames, ",")))
	return int64(h.Sum32())
}

func shuffle(tableClients []tableClient, seed int64) {
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(tableClients), func(i, j int) {
		tableClients[i], tableClients[j] = tableClients[j], tableClients[i]
	})
}
