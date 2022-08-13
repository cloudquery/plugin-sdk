package specs

type Registry int

const (
	RegistryGithub Registry = iota
	RegistryLocal
	RegistryGrpc
)

func (m Registry) String() string {
	return [...]string{"github", "local", "grpc"}[m]
}
