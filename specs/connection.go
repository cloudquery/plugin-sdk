package specs

type ConnectionSpec struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (_ *ConnectionSpec) Validate() error {
	// TODO post process ConnectionSpec after unmarshalling here
	return nil
}
