//go:build freebsd

package limit

import (
	"syscall"
)

func GetUlimit() (Rlimit, error) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return Rlimit{uint64(rLimit.Cur), uint64(rLimit.Max)}, err
}
