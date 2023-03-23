package destination

import (
	"context"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

func (p *Plugin) writeUnmanaged(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan arrow.Record) error {
	return p.client.Write(ctx, tables, res)
}
