package faker

import "github.com/rs/zerolog"

type Option func(*faker)

func WithMaxDepth(depth int) Option {
	return func(f *faker) {
		f.maxDepth = depth
	}
}

func WithLogger(logger zerolog.Logger) Option {
	return func(f *faker) {
		f.logger = logger
	}
}
