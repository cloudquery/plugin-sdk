package scalar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNil(t *testing.T) {
	s := make([]*string, 1)
	s[0] = nil
	require.True(t, IsNil(s[0]))
	require.False(t, isNilTrivial(s[0]))
}

func isNilTrivial(a any) bool {
	return a == nil
}
