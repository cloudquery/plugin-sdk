package writers

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

const (
	msgTypeCreateTable = iota
	msgTypeInsert
	msgTypeDeleteStale
)

var allMsgTypes = []int{msgTypeCreateTable, msgTypeInsert, msgTypeDeleteStale}

// MixedBatchClient is a client that will receive batches of messages with a mixture of tables.
type MixedBatchClient interface {
	CreateTableBatch(ctx context.Context, messages []plugin.MessageCreateTable) error
	InsertBatch(ctx context.Context, messages []plugin.MessageInsert) error
	DeleteStaleBatch(ctx context.Context, messages []plugin.MessageDeleteStale) error
}

type MixedBatchWriter struct {
	tables         schema.Tables
	client         MixedBatchClient
	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSize      int
	batchSizeBytes int

	workerCreateTable *mixedBatchWorker[plugin.MessageCreateTable]
	workerInsert      *mixedBatchWorker[plugin.MessageInsert]
	workerDeleteStale *mixedBatchWorker[plugin.MessageDeleteStale]
	workersLock       *sync.Mutex
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

type mixedBatchWorker[T plugin.Message] struct {
	count     int
	wg        *sync.WaitGroup
	ch        chan T
	flush     chan chan bool
	messages  []T
	writeFunc func(ctx context.Context, messages []T) error
}

func newWorker[T plugin.Message](writeFunc func(ctx context.Context, messages []T) error) *mixedBatchWorker[T] {
	w := &mixedBatchWorker[T]{
		writeFunc: writeFunc,
		messages:  make([]T, 0, defaultBatchSize),
		count:     0,
		ch:        make(chan T),
		wg:        &sync.WaitGroup{},
	}
	return w
}

func (w *mixedBatchWorker[T]) listen(ctx context.Context, ch <-chan T) chan chan bool {
	flush := make(chan chan bool, 1)
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.start(ctx, ch, flush)
	}()
	return flush
}

func (w *mixedBatchWorker[T]) start(ctx context.Context, ch <-chan T, flush chan chan bool) {
	sizeBytes := int64(0)
	messages := make([]T, 0)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				if len(messages) > 0 {
					w.writeFunc(ctx, messages)
				}
				return
			}
			if uint64(len(messages)) == 1000 || sizeBytes+util.TotalRecordSize(r) > int64(1000) {
				w.writeFunc(ctx, messages)
				messages = make([]T, 0)
				sizeBytes = 0
			}
			messages = append(messages, msg)
			sizeBytes += util.TotalRecordSize(msg)
		case <-time.After(w.batchTimeout):
			if len(messages) > 0 {
				w.writeFunc(ctx, messages)
				messages = make([]T, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(messages) > 0 {
				w.writeFunc(ctx, messages)
				messages = make([]T, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func NewMixedBatchWriter(tables schema.Tables, client MixedBatchClient, opts ...MixedBatchWriterOption) (*MixedBatchWriter, error) {
	c := &MixedBatchWriter{
		tables:         tables,
		client:         client,
		workersLock:    &sync.Mutex{},
		logger:         zerolog.Nop(),
		batchTimeout:   defaultBatchTimeoutSeconds * time.Second,
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,

		workerCreateTable: newWorker[plugin.MessageCreateTable](client.CreateTableBatch),
		workerInsert:      newWorker[plugin.MessageInsert](client.InsertBatch),
		workerDeleteStale: newWorker[plugin.MessageDeleteStale](client.DeleteStaleBatch),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// Write starts listening for messages on the msgChan channel and writes them to the client in batches.
func (w *MixedBatchWriter) Write(ctx context.Context, msgChan <-chan plugin.Message) error {
	w.workersLock.Lock()
	flushCreateTable := w.workerCreateTable.listen(ctx, msgChan)
	flushInsert := w.workerInsert.listen(ctx, msgChan)
	flushDeleteStale := w.workerDeleteStale.listen(ctx, msgChan)
	w.workersLock.Unlock()

	done := make(chan bool)
	for msg := range msgChan {
		switch v := msg.(type) {
		case plugin.MessageCreateTable:
			w.workerCreateTable.ch <- v
		case plugin.MessageInsert:
			flushCreateTable <- done
			<-done
			flushDeleteStale <- done
			<-done
			w.workerInsert.ch <- v
		case plugin.MessageDeleteStale:
			flushCreateTable <- done
			<-done
			flushInsert <- done
			<-done
			w.workerDeleteStale.ch <- v
		}
	}

	flushCreateTable <- done
	<-done

	flushInsert <- done
	<-done

	flushDeleteStale <- done
	<-done

	w.workersLock.Lock()
	close(w.workerCreateTable.ch)
	close(w.workerInsert.ch)
	close(w.workerDeleteStale.ch)

	w.workersLock.Unlock()
	return nil
}

func (w *MixedBatchWriter) flush(ctx context.Context, messageID int, messages []plugin.Message) error {
	var err error
	switch messageID {
	case msgTypeCreateTable:
		msgs := make([]plugin.MessageCreateTable, len(messages))
		for i := range messages {
			msgs[i] = messages[i].(plugin.MessageCreateTable)
		}
		err = w.client.CreateTableBatch(ctx, msgs)
	case msgTypeInsert:
		// TODO: should we remove duplicates here?
		w.writeInsert(ctx, messages)
	case msgTypeDeleteStale:
		w.writeDeleteStale(ctx, messages)
	}
	if err != nil {

	}
	start := time.Now()
	batchSize := len(resources)
	if err := w.client.WriteTableBatch(ctx, table, resources); err != nil {
		w.logger.Err(err).Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
	} else {
		w.logger.Info().Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
	}
}

func messageID(msg plugin.Message) int {
	switch msg.(type) {
	case plugin.MessageCreateTable:
		return msgTypeCreateTable
	case plugin.MessageInsert:
		return msgTypeInsert
	case plugin.MessageDeleteStale:
		return msgTypeDeleteStale
	default:
		panic("unknown message type: " + reflect.TypeOf(msg).String())
	}
}
