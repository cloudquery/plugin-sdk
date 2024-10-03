package scheduler

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

const (
	// This is an environment variable and not a spec option in each plugin to make it easier to enable it
	cqDebugSyncMultiplier = "CQ_DEBUG_SYNC_MULTIPLIER"
)

func getTestMultiplier() (int, error) {
	strValue, ok := os.LookupEnv(cqDebugSyncMultiplier)
	if ok {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return 0, fmt.Errorf("failed to parse %s=%s as integer: %w", cqDebugSyncMultiplier, strValue, err)
		}
		return intValue, nil
	}
	return 0, nil
}

func (s *syncClient) syncTest(ctx context.Context, syncMultiplier int, resolvedResources chan<- *schema.Resource) {
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
		// we do this here to avoid locks so we initialize the metrics structure once in the main goroutine
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		s.metrics.InitWithClients(table, clients)
	}

	// First interleave the tables like in round-robin
	tableClients := roundRobinInterleave(s.tables, preInitialisedClients)
	// Then shuffle the tableClients to randomize the order in which they are retrieved.
	// We use a fixed seed so that runs with the same tables and clients perform similarly across syncs
	// however, if the table order changes, the seed will change and the shuffle order will be different,
	// so users have a little bit of control over the randomization.
	seed := hashTableNames(tableNames)
	allClients := make([]tableClient, 0, len(tableClients)*syncMultiplier)
	for _, tc := range tableClients {
		for i := 0; i < syncMultiplier; i++ {
			allClients = append(allClients, tc)
		}
	}
	shuffle(allClients, seed)
	allClients = shardTableClients(allClients, s.shard)

	var wg sync.WaitGroup
	for _, tc := range allClients {
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
