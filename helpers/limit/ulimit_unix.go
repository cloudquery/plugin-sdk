//go:build darwin || linux

package limit

import (
	"syscall"
)

func getUlimit() (uint64, error) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return rLimit.Max, err
}
