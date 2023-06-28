package writers

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type Writer interface {
	Write(ctx context.Context, res <-chan message.WriteMessage) error
}

const (
	DefaultBatchTimeoutSeconds = 20
	DefaultBatchSize           = 10000
	DefaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)
