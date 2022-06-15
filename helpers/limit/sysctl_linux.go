//go:build linux

package limit

import (
	"github.com/lorenzosaino/go-sysctl"
	"github.com/spf13/cast"
)

func calculateFileLimit() (uint64, error) {
	maxFileOpen, err := sysctl.Get("fs.file-max")
	if err != nil {
		return 0, err
	}
	mfo, err := cast.ToUint64E(maxFileOpen)
	if err != nil {
		return 0, err
	}

	fileNr, err := sysctl.Get("fs.file-nr")
	if err != nil {
		return 0, err
	}
	fnr := cast.ToUint64(fileNr)

	return uint64(float64(mfo-fnr) * goroutineReducer), nil
}
