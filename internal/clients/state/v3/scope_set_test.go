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

// countingScope counts how often a compaction visits it.
type countingScope struct {
	liveScope
	visits *int
}

func (s *countingScope) Collapsed() checkpoint.Scope {
	*s.visits++
	return s
}

type endedScope struct {
	liveScope
	into checkpoint.Scope
}

func (s *endedScope) Collapsed() checkpoint.Scope { return s.into }

func TestScopeSetCompactsOnlyOnceItDoubles(t *testing.T) {
	const scopes = 100_000
	s := newScopeSet()
	visits := 0
	for i := range scopes {
		s.add(&countingScope{liveScope: liveScope{id: i}, visits: &visits})
	}
	require.Len(t, s.scopes, scopes, "live scopes are kept")
	// compactions at 64, 128, ..., 65536 visit 131008 scopes in all: the cost per add stays constant
	require.Equal(t, 131008, visits)
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
