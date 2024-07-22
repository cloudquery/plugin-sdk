package plugin

import (
	"context"

	"github.com/apache/arrow/go/v17/arrow"
)

type TransformerClient interface {
	Transform(context.Context, <-chan arrow.Record, chan<- arrow.Record) error
}

func (p *Plugin) Transform(ctx context.Context, recvRecords <-chan arrow.Record, sendRecords chan<- arrow.Record) error {
	return p.client.Transform(ctx, recvRecords, sendRecords)
}
