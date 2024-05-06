package mixedbatchwriter

import (
	"context"
	"time"

	"github.com/apache/arrow/go/v16/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"github.com/rs/zerolog"
)

// Client is a client that will receive batches of messages with a mixture of tables.
type Client interface {
	MigrateTableBatch(ctx context.Context, messages message.WriteMigrateTables) error
	InsertBatch(ctx context.Context, messages message.WriteInserts) error
	DeleteStaleBatch(ctx context.Context, messages message.WriteDeleteStales) error
	DeleteRecordsBatch(ctx context.Context, messages message.WriteDeleteRecords) error
}

type MixedBatchWriter struct {
	client         Client
	logger         zerolog.Logger
	batchSize      int
	batchSizeBytes int
	batchTimeout   time.Duration
	tickerFn       writers.TickerFunc
}

// Assert at compile-time that MixedBatchWriter implements the Writer interface
var _ writers.Writer = (*MixedBatchWriter)(nil)

type Option func(writer *MixedBatchWriter)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *MixedBatchWriter) {
		p.logger = logger
	}
}

func WithBatchSize(size int) Option {
	return func(p *MixedBatchWriter) {
		p.batchSize = size
	}
}

func WithBatchSizeBytes(size int) Option {
	return func(p *MixedBatchWriter) {
		p.batchSizeBytes = size
	}
}

func WithBatchTimeout(timeout time.Duration) Option {
	return func(p *MixedBatchWriter) {
		p.batchTimeout = timeout
	}
}

func withTickerFn(tickerFn writers.TickerFunc) Option {
	return func(p *MixedBatchWriter) {
		p.tickerFn = tickerFn
	}
}

const (
	defaultBatchTimeout   = 20 * time.Second
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 5 * 1024 * 1024 // 5 MiB
)

func New(client Client, opts ...Option) (*MixedBatchWriter, error) {
	c := &MixedBatchWriter{
		client:         client,
		logger:         zerolog.Nop(),
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
		batchTimeout:   defaultBatchTimeout,
		tickerFn:       writers.NewTicker,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// Write starts listening for messages on the msgChan channel and writes them to the client in batches.
func (w *MixedBatchWriter) Write(ctx context.Context, msgChan <-chan message.WriteMessage) error {
	migrateTable := &batchManager[message.WriteMigrateTables, *message.WriteMigrateTable]{
		batch:     make([]*message.WriteMigrateTable, 0, w.batchSize),
		writeFunc: w.client.MigrateTableBatch,
	}
	insert := &insertBatchManager{
		batch:             make([]*message.WriteInsert, 0, w.batchSize),
		writeFunc:         w.client.InsertBatch,
		maxBatchSizeBytes: int64(w.batchSizeBytes),
		logger:            w.logger,
	}
	deleteStale := &batchManager[message.WriteDeleteStales, *message.WriteDeleteStale]{
		batch:     make([]*message.WriteDeleteStale, 0, w.batchSize),
		writeFunc: w.client.DeleteStaleBatch,
	}

	deleteRecord := &batchManager[message.WriteDeleteRecords, *message.WriteDeleteRecord]{
		batch:     make([]*message.WriteDeleteRecord, 0, w.batchSize),
		writeFunc: w.client.DeleteRecordsBatch,
	}

	flush := func(msgType writers.MsgType) error {
		if msgType == writers.MsgTypeUnset {
			return nil
		}
		switch msgType {
		case writers.MsgTypeMigrateTable:
			return migrateTable.flush(ctx)
		case writers.MsgTypeInsert:
			return insert.flush(ctx)
		case writers.MsgTypeDeleteStale:
			return deleteStale.flush(ctx)
		case writers.MsgTypeDeleteRecord:
			return deleteRecord.flush(ctx)
		default:
			panic("unknown message type")
		}
	}
	prevMsgType := writers.MsgTypeUnset
	var err error
	ticker := w.tickerFn(w.batchTimeout)
	defer ticker.Stop()
loop:
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				break loop
			}
			msgType := writers.MsgID(msg)
			if prevMsgType != msgType {
				if err := flush(prevMsgType); err != nil {
					return err
				}
				ticker.Reset(w.batchTimeout)
			}
			prevMsgType = msgType
			switch v := msg.(type) {
			case *message.WriteMigrateTable:
				err = migrateTable.append(ctx, v)
			case *message.WriteInsert:
				err = insert.append(ctx, v)
			case *message.WriteDeleteStale:
				err = deleteStale.append(ctx, v)
			case *message.WriteDeleteRecord:
				err = deleteRecord.append(ctx, v)
			default:
				panic("unknown message type")
			}
			if err != nil {
				return err
			}
		case <-ticker.Chan():
			if err := flush(prevMsgType); err != nil {
				return err
			}
			ticker.Reset(w.batchTimeout)
			prevMsgType = writers.MsgTypeUnset
		}
	}
	return flush(prevMsgType)
}

// generic batch manager for most message types
type batchManager[A ~[]T, T message.WriteMessage] struct {
	batch     []T
	writeFunc func(ctx context.Context, messages A) error
}

func (m *batchManager[A, T]) append(ctx context.Context, msg T) error {
	if len(m.batch) == cap(m.batch) {
		if err := m.flush(ctx); err != nil {
			return err
		}
	}
	m.batch = append(m.batch, msg)
	return nil
}

func (m *batchManager[A, T]) flush(ctx context.Context) error {
	if len(m.batch) == 0 {
		return nil
	}

	err := m.writeFunc(ctx, m.batch)
	if err != nil {
		return err
	}
	clear(m.batch) // GC can work
	m.batch = m.batch[:0]
	return nil
}

// special batch manager for insert messages that also keeps track of the total size of the batch
type insertBatchManager struct {
	batch             []*message.WriteInsert
	writeFunc         func(ctx context.Context, messages message.WriteInserts) error
	curBatchSizeBytes int64
	maxBatchSizeBytes int64
	logger            zerolog.Logger
}

func (m *insertBatchManager) append(ctx context.Context, msg *message.WriteInsert) error {
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
	start := time.Now()
	batchSize := len(m.batch)
	err := m.writeFunc(ctx, m.batch)
	if err != nil {
		m.logger.Err(err).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
		return err
	}
	m.logger.Debug().Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")

	clear(m.batch) // GC can work
	m.batch = m.batch[:0]
	m.curBatchSizeBytes = 0
	return nil
}
