package state

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/stretchr/testify/require"
)

type liveScope struct{ id int }

func (*liveScope) Checkpoint(context.Context, message.SyncMessage) error { return nil }

type endedScope struct {
	liveScope
	into checkpoint.Scope
}

func (s *endedScope) Collapsed() checkpoint.Scope { return s.into }

func TestScopeSetCompactsOnlyOnceItDoubles(t *testing.T) {
	const scopes = 100_000
	s := newScopeSet()
	compactions := 0
	for i := range scopes {
		before := s.compactAt
		s.add(&liveScope{id: i})
		if s.compactAt != before {
			compactions++
		}
	}
	require.Len(t, s.scopes, scopes, "live scopes are kept")
	// each compaction visits the whole set, so doubling keeps the total cost linear
	require.Equal(t, 11, compactions, "compactions at 64, 128, ..., 65536")
	require.Equal(t, 131072, s.compactAt)
}

func TestScopeSetCollapsesEndedScopes(t *testing.T) {
	shared := &liveScope{}
	s := newScopeSet()
	for i := range 100_000 {
		s.add(&endedScope{liveScope: liveScope{id: i}, into: shared})
	}
	require.Less(t, len(s.scopes), minCompactSize, "ended scopes collapse into the one they share")
	require.Contains(t, s.scopes, checkpoint.Scope(shared))
}
