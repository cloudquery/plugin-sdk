package faker

type Option func(*faker)

func WithMaxDepth(depth int) Option {
	return func(f *faker) {
		f.maxDepth = depth
	}
}

func WithSkipFields(fieldName ...string) Option {
	return func(f *faker) {
		for _, field := range fieldName {
			f.skipFields[field] = struct{}{}
		}
	}
}
