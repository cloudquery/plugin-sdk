package batchingwriter

import (
	"time"

	"github.com/rs/zerolog"
)

type Option func(*Batching)

func WithLogger(l zerolog.Logger) Option {
	return func(w *Batching) {
		w.logger = l
	}
}

func WithDedupPK(dedup bool) Option {
	return func(w *Batching) {
		w.dedupPK = dedup
	}
}

func WithDefaultBatchSize(batchSize, batchSizeBytes int64) Option {
	return func(w *Batching) {
		if w.batchSize == 0 {
			w.batchSize = batchSize
		}
		if w.batchSizeBytes == 0 {
			w.batchSizeBytes = batchSizeBytes
		}
	}
}

func WithBatchTimeout(timeout time.Duration) Option {
	return func(w *Batching) {
		w.batchTimeout = timeout
	}
}
