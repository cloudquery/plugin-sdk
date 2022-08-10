package specs

type ConnectionSpec struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (*ConnectionSpec) Validate() error {
	// TODO post process ConnectionSpec after unmarshalling here
	return nil
}
