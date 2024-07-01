package scalar

type Option func(Scalar)

func WithNameTransformer(transformer NameTransformer) Option {
	return func(s Scalar) {
		st, ok := s.(*Struct)
		if !ok {
			return
		}
		st.nameTransformer = transformer
	}
}
