//go:build darwin || linux

package limit

import (
	"syscall"
)

func GetUlimit() (Rlimit, error) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return Rlimit{rLimit.Cur, rLimit.Max}, err
}
