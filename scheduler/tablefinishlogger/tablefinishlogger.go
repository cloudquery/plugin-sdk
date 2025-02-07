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

	activeCounts map[string]int // key: table name}
}

// New creates a new TableFinishLogger
func New(logger zerolog.Logger) *TableFinishLogger {
	return &TableFinishLogger{
		logger:       logger,
		activeCounts: make(map[string]int),
	}
}

// TableStarted signals that a table instance has started syncing
func (t *TableFinishLogger) TableStarted(table *schema.Table, parent *schema.Resource) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.activeCounts[table.Name]++
}

// TableFinished signals that a table instance has finished syncing
func (t *TableFinishLogger) TableFinished(table *schema.Table, client schema.ClientMeta, parent *schema.Resource) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.activeCounts[table.Name] == 0 && (parent == nil || t.activeCounts[parent.Table.Name] == 0) {
		t.logger.Info().Str("table", table.Name).Str("client", client.ID()).Msg("table fully finished syncing")
	}
}
