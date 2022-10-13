package faker

type Option func(*faker)

func WithMaxDepth(depth int) Option {
	return func(f *faker) {
		f.maxDepth = depth
	}
}

func WithSilent() Option {
	return func(f *faker) {
		f.silent = true
	}
}

func WithSkipFields(fieldName ...string) Option {
	return func(f *faker) {
		for _, field := range fieldName {
			f.skipFields[field] = struct{}{}
		}
	}
}
