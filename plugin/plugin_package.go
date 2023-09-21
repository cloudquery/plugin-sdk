package plugin

import "errors"

const (
	GoOSLinux   = "linux"
	GoOSWindows = "windows"
	GoOSDarwin  = "darwin"

	GoArchAmd64 = "amd64"
	GoArchArm64 = "arm64"
)

type Type string

const (
	TypeSource      Type = "source"
	TypeDestination Type = "destination"
)

func (t Type) Validate() error {
	switch t {
	case TypeSource, TypeDestination:
		return nil
	default:
		return errors.New("invalid type: must be one of source, destination")
	}
}

type PackageType string

const (
	PackageTypeNative PackageType = "native"
)

type BuildTarget struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

var DefaultBuildTargets = []BuildTarget{
	{GoOSLinux, GoArchAmd64},
	{GoOSWindows, GoArchAmd64},
	{GoOSDarwin, GoArchAmd64},
	{GoOSDarwin, GoArchArm64},
}
