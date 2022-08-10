package specs

type ConnectionSpec struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (s *ConnectionSpec) Validate() error {
	return nil
}
