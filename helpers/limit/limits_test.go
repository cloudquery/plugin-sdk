package limit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateGoRoutines(t *testing.T) {
	cases := []struct {
		Name       string
		Memory     uint64
		GoRoutines uint64
	}{
		{Name: "Zero", Memory: uint64(0), GoRoutines: uint64(500000)},
		{Name: "Below 1073741824", Memory: uint64(990498816), GoRoutines: uint64(230618)},
		{Name: "At 1073741824", Memory: uint64(1073741824), GoRoutines: uint64(250000)},
		{Name: "Above 1073741824", Memory: uint64(1573741824), GoRoutines: uint64(366415)},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, int(tc.GoRoutines), int(calculateGoRoutines(tc.Memory)))
		})
	}
}
