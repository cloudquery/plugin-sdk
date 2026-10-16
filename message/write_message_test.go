package message

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFailedTables(t *testing.T) {
	failed := func(tables ...string) error {
		return &FailedTablesError{Tables: tables, Err: errors.New("write failed")}
	}
	for _, tc := range []struct {
		name string
		err  error
		want []string
	}{
		{name: "nil", err: nil},
		{name: "unattributed", err: errors.New("write failed")},
		{name: "tables", err: failed("b", "a", "b"), want: []string{"a", "b"}},
		{name: "wrapped", err: fmt.Errorf("flush: %w", failed("a")), want: []string{"a"}},
		{name: "joined", err: errors.Join(failed("a"), fmt.Errorf("flush: %w", failed("b"))), want: []string{"a", "b"}},
		{name: "joined with unattributed", err: errors.Join(failed("a"), errors.New("migrate failed"))},
		{name: "no tables", err: failed()},
		{name: "empty table name", err: failed("a", "")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, FailedTables(tc.err))
		})
	}
}

func TestFailedTablesErrorWithoutCause(t *testing.T) {
	err := &FailedTablesError{Tables: []string{"a", "b"}}
	require.EqualError(t, err, "failed to write tables a, b")
	require.NoError(t, errors.Unwrap(err))
	require.Equal(t, []string{"a", "b"}, FailedTables(err))
}
