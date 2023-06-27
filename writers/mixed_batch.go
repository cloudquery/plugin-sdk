package writers

import (
	"context"
	"time"

	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/rs/zerolog"
)

// MixedBatchClient is a client that will receive batches of messages with a mixture of tables.
type MixedBatchClient interface {
	MigrateTableBatch(ctx context.Context, messages []*message.MigrateTable, options plugin.WriteOptions) error
	InsertBatch(ctx context.Context, messages []*message.Insert, options plugin.WriteOptions) error
	DeleteStaleBatch(ctx context.Context, messages []*message.DeleteStale, options plugin.WriteOptions) error
}

type MixedBatchWriter struct {
	client         MixedBatchClient
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

func NewMixedBatchWriter(client MixedBatchClient, opts ...MixedBatchWriterOption) (*MixedBatchWriter, error) {
	c := &MixedBatchWriter{
		client:         client,
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

// Write starts listening for messages on the msgChan channel and writes them to the client in batches.
func (w *MixedBatchWriter) Write(ctx context.Context, options plugin.WriteOptions, msgChan <-chan message.Message) error {
	migrateTable := &batchManager[*message.MigrateTable]{
		batch:        make([]*message.MigrateTable, 0, w.batchSize),
		writeFunc:    w.client.MigrateTableBatch,
		writeOptions: options,
	}
	insert := &insertBatchManager{
		batch:             make([]*message.Insert, 0, w.batchSize),
		writeFunc:         w.client.InsertBatch,
		maxBatchSizeBytes: int64(w.batchSizeBytes),
		writeOptions:      options,
	}
	deleteStale := &batchManager[*message.DeleteStale]{
		batch:        make([]*message.DeleteStale, 0, w.batchSize),
		writeFunc:    w.client.DeleteStaleBatch,
		writeOptions: options,
	}
	flush := func(msgType msgType) error {
		switch msgType {
		case msgTypeMigrateTable:
			return migrateTable.flush(ctx)
		case msgTypeInsert:
			return insert.flush(ctx)
		case msgTypeDeleteStale:
			return deleteStale.flush(ctx)
		default:
			panic("unknown message type")
		}
	}
	prevMsgType := msgTypeUnset
	var err error
	for msg := range msgChan {
		msgType := msgID(msg)
		if prevMsgType != msgTypeUnset && prevMsgType != msgType {
			if err := flush(prevMsgType); err != nil {
				return err
			}
		}
		prevMsgType = msgType
		switch v := msg.(type) {
		case *message.MigrateTable:
			err = migrateTable.append(ctx, v)
		case *message.Insert:
			err = insert.append(ctx, v)
		case *message.DeleteStale:
			err = deleteStale.append(ctx, v)
		default:
			panic("unknown message type")
		}
		if err != nil {
			return err
		}
	}
	if prevMsgType == msgTypeUnset {
		return nil
	}
	return flush(prevMsgType)
}

// generic batch manager for most message types
type batchManager[T message.Message] struct {
	batch        []T
	writeFunc    func(ctx context.Context, messages []T, options plugin.WriteOptions) error
	writeOptions plugin.WriteOptions
}

func (m *batchManager[T]) append(ctx context.Context, msg T) error {
	if len(m.batch) == cap(m.batch) {
		if err := m.flush(ctx); err != nil {
			return err
		}
	}
	m.batch = append(m.batch, msg)
	return nil
}

func (m *batchManager[T]) flush(ctx context.Context) error {
	if len(m.batch) == 0 {
		return nil
	}

	err := m.writeFunc(ctx, m.batch, m.writeOptions)
	if err != nil {
		return err
	}
	m.batch = m.batch[:0]
	return nil
}

// special batch manager for insert messages that also keeps track of the total size of the batch
type insertBatchManager struct {
	batch             []*message.Insert
	writeFunc         func(ctx context.Context, messages []*message.Insert, writeOptions plugin.WriteOptions) error
	curBatchSizeBytes int64
	maxBatchSizeBytes int64
	writeOptions      plugin.WriteOptions
}

func (m *insertBatchManager) append(ctx context.Context, msg *message.Insert) error {
	if len(m.batch) == cap(m.batch) || m.curBatchSizeBytes+util.TotalRecordSize(msg.Record) > m.maxBatchSizeBytes {
		if err := m.flush(ctx); err != nil {
			return err
		}
	}
	m.batch = append(m.batch, msg)
	m.curBatchSizeBytes += util.TotalRecordSize(msg.Record)
	return nil
}

func (m *insertBatchManager) flush(ctx context.Context) error {
	if len(m.batch) == 0 {
		return nil
	}

	err := m.writeFunc(ctx, m.batch, m.writeOptions)
	if err != nil {
		return err
	}
	m.batch = m.batch[:0]
	return nil
}
