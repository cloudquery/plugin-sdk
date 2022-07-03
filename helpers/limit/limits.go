package limit

import (
	"errors"
	"math"

	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/pbnjay/memory"
)

const (
	gbInBytes         int     = 1024 * 1024 * 1024
	goroutinesPerGB   float64 = 250000
	minimalGoRoutines float64 = 100
	goroutineReducer          = 0.8
	// only use 75% of the available file descriptors, so to leave for other processes file descriptors
	mfoReducer = 0.75
)

type Rlimit struct {
	Cur uint64
	Max uint64
}

func GetMaxGoRoutines() uint64 {
	limit := calculateGoRoutines(getMemory())
	ulimit, err := GetUlimit()
	if err != nil || ulimit.Cur == 0 {
		return limit
	}
	if ulimit.Cur > limit {
		return limit
	}
	return ulimit.Cur
}

// DiagnoseLimits verifies if user should increase ulimit or max file descriptors to improve number of expected
// goroutines in CQ to improve performance
func DiagnoseLimits() (diags diag.Diagnostics) {
	// the amount of goroutines we want based on machine memory
	want := calculateGoRoutines(getMemory())
	// calculate file descriptor limit
	fds, err := calculateFileLimit()
	if err == nil && fds < want {
		diags = diags.Add(diag.NewBaseError(errors.New("available file descriptors capacity lower than expected"),
			diag.USER,
			diag.WithSeverity(diag.WARNING),
			diag.WithDetails("available descriptor capacity is %d want %d to run optimally, consider increasing max file descriptors on machine.", fds, want)))
	}
	ulimit, err := GetUlimit()
	if err == nil && ulimit.Cur < want {
		diags = diags.Add(diag.NewBaseError(errors.New("ulimit available for CloudQuery process lower than expected"),
			diag.USER,
			diag.WithSeverity(diag.WARNING),
			diag.WithDetails("set ulimit capacity is %d want %d to run optimally, consider increasing ulimit on this machine.", ulimit.Cur, want)))
	}
	return diags
}

func getMemory() uint64 {
	return memory.TotalMemory()
}

func calculateMemoryGoRoutines(totalMemory uint64) uint64 {
	if totalMemory == 0 {
		// assume we have 2 GB RAM
		return uint64(math.Max(minimalGoRoutines, goroutinesPerGB*2*goroutineReducer))
	}
	return uint64(math.Max(minimalGoRoutines, (goroutinesPerGB*float64(totalMemory)/float64(gbInBytes))*goroutineReducer))
}

func calculateGoRoutines(totalMemory uint64) uint64 {
	total := calculateMemoryGoRoutines(totalMemory)
	mfo, err := calculateFileLimit()
	if err != nil {
		return total
	}

	if mfo < total {
		return uint64(float64(mfo) * mfoReducer)
	}
	return total
}
