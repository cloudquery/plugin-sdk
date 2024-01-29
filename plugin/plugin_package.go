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
	OS   string   `json:"os"`
	Arch string   `json:"arch"`
	CGO  bool     `json:"cgo"`
	Env  []string `json:"env"`
}

func (t BuildTarget) GetEnvVariables() []string {
	cgo := "CGO_ENABLED="
	if t.CGO {
		cgo += "1"
	} else {
		cgo += "0"
	}
	return append([]string{
		"GOOS=" + t.OS,
		"GOARCH=" + t.Arch,
		cgo, // default is to tool at the param. Can be overridden by adding `CGO_ENABLED=1` to BuildTarget.Env
	}, t.Env...)
}

var DefaultBuildTargets = []BuildTarget{
	{OS: GoOSLinux, Arch: GoArchAmd64},
	{OS: GoOSLinux, Arch: GoArchArm64},
	{OS: GoOSWindows, Arch: GoArchAmd64},
	{OS: GoOSDarwin, Arch: GoArchAmd64},
	{OS: GoOSDarwin, Arch: GoArchArm64},
}
