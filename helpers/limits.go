package helpers

import (
	"github.com/pbnjay/memory"
)

const GB_IN_BYTES uint64 = 1024 * 1024 * 1024
const GO_ROUTINES_PER_GB uint64 = 250000

func GetMaxGoRoutines() uint64 {
	totalMemory := memory.TotalMemory()
	if totalMemory == 0 {
		// assume we have 2 GB RAM
		return GO_ROUTINES_PER_GB * 2
	}
	gb := float64(totalMemory) / float64(GB_IN_BYTES)
	return uint64(gb) * GO_ROUTINES_PER_GB
}
