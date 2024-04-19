package plugin

import (
	"errors"
)

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

func (t BuildTarget) EnvVariables() []string {
	variables := append(t.cgoEnvVariables(), "GOOS="+t.OS, "GOARCH="+t.Arch)
	return append(variables, t.Env...)
}

func (t BuildTarget) cgoEnvVariables() []string {
	// default is to tool at the param. Can be overridden by adding `CGO_ENABLED=1` to BuildTarget.Env
	if !t.CGO {
		return []string{"CGO_ENABLED=0"}
	}

	switch t.OS {
	case GoOSWindows:
		return []string{"CGO_ENABLED=1", "CC=x86_64-w64-mingw32-gcc", "CXX=x86_64-w64-mingw32-g++"}
	case GoOSDarwin:
		return []string{"CGO_ENABLED=1", "CC=o64-clang", "CXX=o64-clang++"}
	case GoOSLinux:
		// nop, see below
	default:
		return []string{"CGO_ENABLED=1"}
	}

	// linux
	switch t.Arch {
	case GoArchAmd64:
		return []string{"CGO_ENABLED=1", "CC=gcc", "CXX=g++"}
	case GoArchArm64:
		return []string{"CGO_ENABLED=1", "CC=aarch64-linux-gnu-gcc", "CXX=aarch64-linux-gnu-g++"}
	default:
		return []string{"CGO_ENABLED=1"}
	}
}

var DefaultBuildTargets = []BuildTarget{
	{OS: GoOSLinux, Arch: GoArchAmd64},
	{OS: GoOSLinux, Arch: GoArchArm64},
	{OS: GoOSWindows, Arch: GoArchAmd64},
	{OS: GoOSDarwin, Arch: GoArchAmd64},
	{OS: GoOSDarwin, Arch: GoArchArm64},
}
