package destination

import (
	"context"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"golang.org/x/sync/errgroup"
)

func (p *Plugin) writeUnmanaged(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, res <-chan schema.DestinationResource) error {
	ch := make(chan *ClientResource)
	eg, gctx := errgroup.WithContext(ctx)
	// given most destination plugins writing in batch we are using a worker pool to write in parallel
	// it might not generalize well and we might need to move it to each destination plugin implementation.
	for i := 0; i < writeWorkers; i++ {
		eg.Go(func() error {
			return p.client.Write(gctx, tables, ch)
		})
	}
	sourceColumn := &schema.Text{}
	_ = sourceColumn.Set(sourceName)
	syncTimeColumn := &schema.Timestamptz{}
	_ = syncTimeColumn.Set(syncTime)
	for r := range res {
		if len(r.Data) < len(tables.Get(r.TableName).Columns) {
			r.Data = append([]schema.CQType{sourceColumn, syncTimeColumn}, r.Data...)
		}
		clientResource := &ClientResource{
			TableName: r.TableName,
			Data:      schema.TransformWithTransformer(p.client, r.Data),
		}
		select {
		case <-gctx.Done():
			close(ch)
			return eg.Wait()
		case ch <- clientResource:
		}
	}

	close(ch)
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}