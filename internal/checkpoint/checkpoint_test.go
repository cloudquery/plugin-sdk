package checkpoint

import (
	"context"
	"sync"
	"sync/atomic"
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

// collapsingScope collapses into into once ended.
type collapsingScope struct {
	heldScope
	ended bool
	into  Scope
}

func (s *collapsingScope) Collapsed() Scope {
	if s.ended {
		return s.into
	}
	return s
}

func TestJoinCollapsesEndedScopes(t *testing.T) {
	shared := &heldScope{}
	a := &collapsingScope{ended: true, into: shared}
	b := &collapsingScope{ended: true, into: shared}
	running := &collapsingScope{into: shared}

	require.Same(t, shared, Join(a, b))
	joined, ok := Join(a, b, running).(*joined)
	require.True(t, ok)
	require.Len(t, joined.scopes, 2)
}

func TestAbortReleasesLaterCheckpoints(t *testing.T) {
	ctx := context.Background()
	res := make(chan message.SyncMessage, 10)
	a, b := &heldScope{}, &heldScope{}

	first, err := Send(ctx, res, a, nil, marker("first"))
	require.NoError(t, err)
	second, err := Send(ctx, res, b, []*Sent{first}, marker("second"))
	require.NoError(t, err)

	for _, msg := range a.msgs {
		Drop(msg)
	}
	require.True(t, first.Aborted())
	require.True(t, second.Aborted(), "a checkpoint waiting for an aborted one is aborted too")
	b.deliver(t, res)
	require.Empty(t, received(res))

	_, err = Send(ctx, res, b, []*Sent{second}, marker("third"))
	require.ErrorIs(t, err, ErrSyncAborted)
}

func TestNotifyAbort(t *testing.T) {
	var flag atomic.Bool
	sent := &Sent{}
	sent.NotifyAbort(&flag)
	require.False(t, flag.Load())
	sent.Abort()
	require.True(t, flag.Load())

	var late atomic.Bool
	sent.NotifyAbort(&late)
	require.True(t, late.Load(), "a checkpoint aborted already sets the flag at once")

	var finished atomic.Bool
	done := &Sent{}
	done.NotifyAbort(&finished)
	require.NoError(t, done.Finish(context.Background(), make(chan message.SyncMessage)))
	done.Abort()
	require.False(t, finished.Load(), "a sent checkpoint is never aborted")
}

func TestNotifyAbortFollowsMerges(t *testing.T) {
	var flag atomic.Bool
	held := NewHeld(&heldScope{})
	sent := &Sent{}
	sent.NotifyAbort(&flag)
	sent.MergeInto(held)
	held.Abort()
	require.True(t, flag.Load())
	require.True(t, sent.Aborted())
}

func TestConcurrentMergeKeepsAborts(t *testing.T) {
	// with the root read before its lock is taken, a flag added while the checkpoint merges can be lost
	for range 1000 {
		var flag atomic.Bool
		held := NewHeld(&heldScope{})
		sent := &Sent{}
		var wg sync.WaitGroup
		wg.Go(func() { sent.MergeInto(held) })
		wg.Go(func() { sent.NotifyAbort(&flag) })
		wg.Wait()
		held.Abort()
		require.True(t, flag.Load())
	}
}
