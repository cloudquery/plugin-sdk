package batcher

import "time"

type Batcher[T ~[]E, E any] struct {

	// BatchInterval controls the max interval between batches are produced
	BatchInterval time.Duration

	// BatchLen controls the maximum amount of entries in the produced batch
	BatchLen int

	// BatchSize controls the maximum data size of the produced batch
	BatchSize int
}
