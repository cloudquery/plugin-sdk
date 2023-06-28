package writers

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type Writer interface {
	Write(ctx context.Context, res <-chan message.WriteMessage) error
}
