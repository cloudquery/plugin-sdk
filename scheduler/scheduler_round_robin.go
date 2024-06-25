package scheduler

import (
	"context"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type tableClient struct {
	table  *schema.Table
	client schema.ClientMeta
}

func (s *syncClient) syncRoundRobin(ctx context.Context, resolvedResources chan<- *schema.Resource) {
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
		s.metrics.initWithClients(table, clients, s.invocationID)
	}

	tableClients := roundRobinInterleave(s.tables, preInitialisedClients)

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
			// Round Robin currently uses the DFS algorithm to resolve the tables, but this
			// may change in the future.
			s.resolveTableDfs(ctx, table, cl, nil, resolvedResources, 1)
		}()
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()
}

// interleave table-clients so that we get:
// table1-client1, table2-client1, table3-client1, table1-client2, table2-client2, table3-client2, ...
func roundRobinInterleave(tables schema.Tables, preInitialisedClients [][]schema.ClientMeta) []tableClient {
	tableClients := make([]tableClient, 0)
	c := 0
	for {
		addedNew := false
		for i, table := range tables {
			if c < len(preInitialisedClients[i]) {
				tableClients = append(tableClients, tableClient{table: table, client: preInitialisedClients[i][c]})
				addedNew = true
			}
		}
		c++
		if !addedNew {
			break
		}
	}
	return tableClients
}
