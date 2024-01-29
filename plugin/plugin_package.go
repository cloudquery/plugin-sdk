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
	Env  []string `json:"env"`
}

func (t BuildTarget) GetEnvVariables() []string {
	return append([]string{
		"GOOS=" + t.OS,
		"GOARCH=" + t.Arch,
		"CGO_ENABLED=0", // default is this, but adding `CGO_ENABLED=1` to BuildTarget.Env solves the issue
	}, t.Env...)
}

var DefaultBuildTargets = []BuildTarget{
	{OS: GoOSLinux, Arch: GoArchAmd64},
	{OS: GoOSLinux, Arch: GoArchArm64},
	{OS: GoOSWindows, Arch: GoArchAmd64},
	{OS: GoOSDarwin, Arch: GoArchAmd64},
	{OS: GoOSDarwin, Arch: GoArchArm64},
}
