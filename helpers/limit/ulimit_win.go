//go:build windows

package limit

import "errors"

func GetUlimit() (Rlimit, error) {
	return Rlimit{0, 0}, errors.New("ulimit not supported in windows")
}
