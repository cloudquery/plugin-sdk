package plugin

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
)

type TransformerClient interface {
	Transform(context.Context, <-chan arrow.RecordBatch, chan<- arrow.RecordBatch) error
	TransformSchema(context.Context, *arrow.Schema) (*arrow.Schema, error)
}

func (p *Plugin) Transform(ctx context.Context, recvRecords <-chan arrow.RecordBatch, sendRecords chan<- arrow.RecordBatch) error {
	err := p.client.Transform(ctx, recvRecords, sendRecords)
	close(sendRecords)
	return err
}
func (p *Plugin) TransformSchema(ctx context.Context, old *arrow.Schema) (*arrow.Schema, error) {
	return p.client.TransformSchema(ctx, old)
}
