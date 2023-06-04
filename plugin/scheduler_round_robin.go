package plugin

import (
	"context"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"golang.org/x/sync/semaphore"
)

type tableClient struct {
	table  *schema.Table
	client schema.ClientMeta
}

func (p *Plugin) syncRoundRobin(ctx context.Context, options SyncOptions, client ManagedSyncClient, tables schema.Tables, resolvedResources chan<- *schema.Resource) {
	tableConcurrency := max(uint64(options.Concurrency/minResourceConcurrency), minTableConcurrency)
	resourceConcurrency := tableConcurrency * minResourceConcurrency

	p.tableSems = make([]*semaphore.Weighted, p.maxDepth)
	for i := uint64(0); i < p.maxDepth; i++ {
		p.tableSems[i] = semaphore.NewWeighted(int64(tableConcurrency))
		// reduce table concurrency logarithmically for every depth level
		tableConcurrency = max(tableConcurrency/2, minTableConcurrency)
	}
	p.resourceSem = semaphore.NewWeighted(int64(resourceConcurrency))

	// we have this because plugins can return sometimes clients in a random way which will cause
	// differences between this run and the next one.
	preInitialisedClients := make([][]schema.ClientMeta, len(tables))
	for i, table := range tables {
		clients := []schema.ClientMeta{client.(schema.ClientMeta)}
		if table.Multiplex != nil {
			clients = table.Multiplex(client.(schema.ClientMeta))
		}
		preInitialisedClients[i] = clients
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		p.metrics.initWithClients(table, clients)
	}

	// We start a goroutine that logs the metrics periodically.
	// It needs its own waitgroup
	var logWg sync.WaitGroup
	logWg.Add(1)

	logCtx, logCancel := context.WithCancel(ctx)
	go p.periodicMetricLogger(logCtx, &logWg)

	tableClients := roundRobinInterleave(tables, preInitialisedClients)

	var wg sync.WaitGroup
	for _, tc := range tableClients {
		table := tc.table
		cl := tc.client
		if err := p.tableSems[0].Acquire(ctx, 1); err != nil {
			// This means context was cancelled
			wg.Wait()
			// gracefully shut down the logger goroutine
			logCancel()
			logWg.Wait()
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer p.tableSems[0].Release(1)
			// not checking for error here as nothing much to do.
			// the error is logged and this happens when context is cancelled
			// Round Robin currently uses the DFS algorithm to resolve the tables, but this
			// may change in the future.
			p.resolveTableDfs(ctx, table, cl, nil, resolvedResources, 1)
		}()
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()

	// gracefully shut down the logger goroutine
	logCancel()
	logWg.Wait()
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
