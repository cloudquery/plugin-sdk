package scalar

import (
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/stretchr/testify/assert"
)

func TestDurationMultiplierUp(t *testing.T) {
	s := NewScalar(arrow.FixedWidthTypes.Duration_us)
	dur := time.Duration(1234) * time.Millisecond
	assert.NoError(t, s.Set(dur))
	assert.Equal(t, "1234000us", s.String())
}

func TestDurationMultiplierDown(t *testing.T) {
	s := NewScalar(arrow.FixedWidthTypes.Duration_us)
	dur := time.Duration(1234) * time.Nanosecond
	assert.NoError(t, s.Set(dur))
	assert.Equal(t, "1us", s.String())
}
