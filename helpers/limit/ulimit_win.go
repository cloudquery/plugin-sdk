//go:build windows

package limit

func getUlimit() (uint64, error) {
	return 0, nil
}
