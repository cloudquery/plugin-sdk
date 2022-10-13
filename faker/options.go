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

func WithSkipFields(fieldName ...string) Option {
	return func(f *faker) {
		for _, field := range fieldName {
			f.skipFields[field] = struct{}{}
		}
	}
}
