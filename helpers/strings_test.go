package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDuplicates(t *testing.T) {
	assert.False(t, HasDuplicates([]string{"A", "b", "c"}))
	assert.False(t, HasDuplicates([]string{"A", "a", "c"}))
	assert.True(t, HasDuplicates([]string{"a", "a", "c"}))
	assert.True(t, HasDuplicates([]string{"a", "a", "c", "c", "f"}))
}
