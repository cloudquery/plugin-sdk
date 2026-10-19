package state

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/stretchr/testify/require"
)

func TestNewClientWithOptionsRejectsCheckpoints(t *testing.T) {
	_, err := NewClientWithOptions(context.Background(), nil, "cq_state", ClientOptions{Checkpoints: make(chan message.SyncMessage)})
	require.ErrorContains(t, err, "use NewConnectedClientWithOptions")
}
