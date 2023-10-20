package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	r := require.New(t)
	tags := Tags{}
	tags.Add("test")
	r.Equal(1, tags.Len())
	r.True(tags.Contains("test"))
	tags.Add("test")
	r.Equal(1, tags.Len())
	tags.Remove("test")
	r.Equal(0, tags.Len())
	r.False(tags.Contains("test"))
	tags.Add("test")
	r.Equal(1, tags.Len())
	tags.Add("test2")
	r.Equal(2, tags.Len())
	tags.Remove("test")
	r.Equal(1, tags.Len())
	r.False(tags.Contains("test"))
	r.True(tags.Contains("test2"))

	b, err := tags.MarshalJSON()
	r.NoError(err)
	r.Equal(`["test2"]`, string(b))
}
