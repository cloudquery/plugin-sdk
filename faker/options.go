package faker

type Option func(*faker)

func WithMaxDepth(depth int) Option {
	return func(f *faker) {
		f.maxDepth = depth
	}
}

func WithVerbose() Option {
	return func(f *faker) {
		f.verbose = true
	}
}
