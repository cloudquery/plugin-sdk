package plugin

import "errors"

const (
	GoOSLinux   = "linux"
	GoOSWindows = "windows"
	GoOSDarwin  = "darwin"

	GoArchAmd64 = "amd64"
	GoArchArm64 = "arm64"
)

type Kind string

const (
	KindSource      Kind = "source"
	KindDestination Kind = "destination"
)

func (k Kind) Validate() error {
	switch k {
	case KindSource, KindDestination:
		return nil
	default:
		return errors.New("invalid plugin kind: must be one of source, destination")
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
