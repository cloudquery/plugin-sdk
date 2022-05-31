//go:build windows

package limit

func GetUlimit() (Rlimit, error) {
	return Rlimit{0, 0}, nil
}
