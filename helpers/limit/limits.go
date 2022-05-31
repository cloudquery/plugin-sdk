package limit

import (
	"github.com/pbnjay/memory"
)

const (
	gbInBytes       int     = 1024 * 1024 * 1024
	goroutinesPerGB float64 = 250000
)

type Rlimit struct {
	Cur uint64
	Max uint64
}

func GetMaxGoRoutines() uint64 {
	limit := calculateGoRoutines(getMemory())
	ulimit, err := GetUlimit()
	if err != nil || ulimit.Max == 0 {
		return limit
	}
	if ulimit.Max > limit {
		return limit
	}
	return ulimit.Max
}

func getMemory() uint64 {
	return memory.TotalMemory()
}

func calculateGoRoutines(totalMemory uint64) uint64 {
	if totalMemory == 0 {
		// assume we have 2 GB RAM
		return uint64(goroutinesPerGB * 2)
	}
	return uint64(goroutinesPerGB * float64(totalMemory) / float64(gbInBytes))
}
