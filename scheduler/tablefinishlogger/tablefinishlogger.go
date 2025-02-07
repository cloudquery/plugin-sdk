package tablefinishlogger

import (
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

// TableFinishLogger tracks when tables completely finish syncing
type TableFinishLogger struct {
	mu     sync.Mutex
	logger zerolog.Logger

	// Track which parent tables have finished
	parentsDone map[string]struct{} // key: parent table name

	// Track active count for child tables
	childActiveCounts map[string]int // key: child table name
}

// New creates a new TableFinishLogger
func New(logger zerolog.Logger) *TableFinishLogger {
	return &TableFinishLogger{
		logger:            logger,
		parentsDone:       make(map[string]struct{}),
		childActiveCounts: make(map[string]int),
	}
}

// TableStarted signals that a table instance has started syncing
func (t *TableFinishLogger) TableStarted(table *schema.Table, parent *schema.Resource) {
	if parent == nil {
		return // No need to track starts for parent tables
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.childActiveCounts[table.Name]++
}

// TableFinished signals that a table instance has finished syncing
func (t *TableFinishLogger) TableFinished(table *schema.Table, client schema.ClientMeta, parent *schema.Resource) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if parent == nil {
		// This is a parent table finishing
		t.parentsDone[table.Name] = struct{}{}
		t.logger.Info().Str("table", table.Name).Str("client", client.ID()).Msg("table fully finished syncing")
		return
	}

	// Handle child table
	t.childActiveCounts[table.Name]--
	remaining := t.childActiveCounts[table.Name]

	// Log when a child table is completely done, which happens when:
	// 1. The parent table is done
	// 2. There are no more active instances of this child table
	if _, parentDone := t.parentsDone[parent.Table.Name]; parentDone && remaining == 0 {
		t.logger.Info().Str("table", table.Name).Str("client", client.ID()).Str("parent_table", parent.Table.Name).Msg("table fully finished syncing")
	}
}
