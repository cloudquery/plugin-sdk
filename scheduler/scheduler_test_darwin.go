//go:build darwin
// +build darwin

package scheduler

import (
	"testing"

	"github.com/rs/zerolog"
)

func getTestLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	// return zerolog.New(zerolog.NewTestWriter(t))
	return zerolog.Nop()
}
