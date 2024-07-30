package plugin

import (
	"context"

	"github.com/apache/arrow/go/v17/arrow"
)

type TransformerClient interface {
	Transform(context.Context, <-chan arrow.Record, chan<- arrow.Record) error
	TransformSchema(context.Context, *arrow.Schema) (*arrow.Schema, error)
}

func (p *Plugin) Transform(ctx context.Context, recvRecords <-chan arrow.Record, sendRecords chan<- arrow.Record) error {
	return p.client.Transform(ctx, recvRecords, sendRecords)
}
func (p *Plugin) TransformSchema(ctx context.Context, old *arrow.Schema) (*arrow.Schema, error) {
	return p.client.TransformSchema(ctx, old)
}
