package faker

type Option func(*faker)

func WithMaxDepth(depth int) Option {
	return func(f *faker) {
		f.maxDepth = depth
	}
}

func WithSkipEFace() Option {
	return func(f *faker) {
		f.ignoreEFace = true
	}
}

func WithSkipFields(fieldName ...string) Option {
	return func(f *faker) {
		for _, field := range fieldName {
			f.skipFields[field] = struct{}{}
		}
	}
}
