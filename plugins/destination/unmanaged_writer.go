package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"golang.org/x/sync/errgroup"
)

func (p *Plugin) writeUnmanaged(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan schema.DestinationResource) error {
	ch := make(chan *ClientResource)
	eg, gctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return p.client.Write(gctx, tables, ch)
	})
	sourceColumn := &schema.Text{}
	_ = sourceColumn.Set(sourceSpec.Name)
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
		return fmt.Errorf("failed waiting for destination client %w", err)
	}
	return nil
}
