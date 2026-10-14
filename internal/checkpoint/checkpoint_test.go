package checkpoint

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/stretchr/testify/require"
)

// heldScope holds messages until the test delivers them, as a scope does until its resources are sent.
type heldScope struct {
	msgs []message.SyncMessage
}

func (s *heldScope) Checkpoint(_ context.Context, msg message.SyncMessage) error {
	s.msgs = append(s.msgs, msg)
	return nil
}

func (s *heldScope) deliver(t *testing.T, res chan<- message.SyncMessage) {
	t.Helper()
	for _, msg := range s.msgs {
		require.NoError(t, Deliver(context.Background(), res, msg))
	}
	s.msgs = nil
}

func marker(name string) message.SyncMessage {
	return &message.SyncError{TableName: name}
}

func received(res chan message.SyncMessage) []string {
	var names []string
	for len(res) > 0 {
		names = append(names, (<-res).GetTable().Name)
	}
	return names
}

func TestSendWaitsForScopesAndPreviousCheckpoints(t *testing.T) {
	ctx := context.Background()
	res := make(chan message.SyncMessage, 10)
	a, b := &heldScope{}, &heldScope{}

	first, err := Send(ctx, res, a, nil, marker("first"))
	require.NoError(t, err)
	second, err := Send(ctx, res, b, []*Sent{first}, marker("second"))
	require.NoError(t, err)

	b.deliver(t, res)
	require.Empty(t, received(res), "second must wait for first")
	require.False(t, second.Done())
	a.deliver(t, res)
	require.Equal(t, []string{"first", "second"}, received(res))
	require.True(t, first.Done())
	require.True(t, second.Done())

	// without a scope, a checkpoint waits only for previous ones
	third, err := Send(ctx, res, nil, []*Sent{second}, marker("third"))
	require.NoError(t, err)
	require.True(t, third.Done())
	require.Equal(t, []string{"third"}, received(res))
}

func TestJoin(t *testing.T) {
	ctx := context.Background()
	res := make(chan message.SyncMessage, 10)
	a, b := &heldScope{}, &heldScope{}

	require.Nil(t, Join(nil, nil))
	require.Same(t, a, Join(nil, a, a))
	joined := Join(Join(a, nil), b, a)
	require.NoError(t, joined.Checkpoint(ctx, marker("joined")))

	a.deliver(t, res)
	require.Empty(t, received(res))
	b.deliver(t, res)
	require.Equal(t, []string{"joined"}, received(res))
}

func TestMergeInto(t *testing.T) {
	ctx := context.Background()
	res := make(chan message.SyncMessage, 10)
	a := &heldScope{}

	merged := &Sent{}
	into := &Sent{}
	merged.MergeInto(into)
	_, err := Send(ctx, res, a, []*Sent{merged}, marker("after merged"))
	require.NoError(t, err)

	a.deliver(t, res)
	require.Empty(t, received(res))
	require.False(t, merged.Done())
	require.NoError(t, into.Finish(ctx, res))
	require.True(t, merged.Done())
	require.Equal(t, []string{"after merged"}, received(res))
}
