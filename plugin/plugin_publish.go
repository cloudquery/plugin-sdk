package plugin

const (
	GoOslinux   = "linux"
	GoOswindows = "windows"
	GoOsDarwin  = "darwin"

	GoArchAmd64 = "amd64"
	GoArchArm64 = "arm64"
)

type PackageType string

const (
	PackageTypeNative PackageType = "native"
	PackageTypeDocker PackageType = "docker"
)

type BuildTarget struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

var buildTargets = []BuildTarget{
	{GoOslinux, GoArchAmd64},
	{GoOswindows, GoArchAmd64},
	{GoOsDarwin, GoArchAmd64},
	{GoOsDarwin, GoArchArm64},
}
