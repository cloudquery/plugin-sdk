package local

type Spec struct {
	// Path is the path to the local directory.
	Path string `json:"path"`
}

func (s *Spec) SetDefaults() {
	s.Path = ".cq"
}
