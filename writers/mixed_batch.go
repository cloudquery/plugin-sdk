package writers

import (
	"context"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

// MixedBatchClient is a client that will receive batches of messages for a mixture of tables.
type MixedBatchClient interface {
	CreateTableBatch(ctx context.Context, resources []plugin.MessageCreateTable) error
	InsertBatch(ctx context.Context, resources []plugin.MessageInsert) error
	DeleteStaleBatch(ctx context.Context, resources []plugin.MessageDeleteStale) error
}

type MixedBatchWriter struct {
	tables      schema.Tables
	client      MixedBatchClient
	workers     map[string]*worker
	workersLock *sync.Mutex

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSize      int
	batchSizeBytes int
}

// Assert at compile-time that MixedBatchWriter implements the Writer interface
var _ Writer = (*MixedBatchWriter)(nil)

type MixedBatchWriterOption func(writer *MixedBatchWriter)

func WithMixedBatchWriterLogger(logger zerolog.Logger) MixedBatchWriterOption {
	return func(p *MixedBatchWriter) {
		p.logger = logger
	}
}

func WithMixedBatchWriterBatchTimeout(timeout time.Duration) MixedBatchWriterOption {
	return func(p *MixedBatchWriter) {
		p.batchTimeout = timeout
	}
}

func WithMixedBatchWriterBatchSize(size int) MixedBatchWriterOption {
	return func(p *MixedBatchWriter) {
		p.batchSize = size
	}
}

func WithMixedBatchWriterSizeBytes(size int) MixedBatchWriterOption {
	return func(p *MixedBatchWriter) {
		p.batchSizeBytes = size
	}
}

func NewMixedBatchWriter(tables schema.Tables, client MixedBatchClient, opts ...MixedBatchWriterOption) (*MixedBatchWriter, error) {
	c := &MixedBatchWriter{
		tables:         tables,
		client:         client,
		workers:        make(map[string]*worker),
		workersLock:    &sync.Mutex{},
		logger:         zerolog.Nop(),
		batchTimeout:   defaultBatchTimeoutSeconds * time.Second,
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (c *MixedBatchWriter) Write(ctx context.Context, res <-chan plugin.Message) error {
	return nil // TODO
}
