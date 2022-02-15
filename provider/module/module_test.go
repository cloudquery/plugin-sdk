package module

import (
	"embed"
	"sort"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/*
var testdata embed.FS

func TestEmbeddedReader(t *testing.T) {
	for _, tc := range []struct {
		PreferredVersions []uint32
		ExpectedVersions  []uint32
		ExpectedFilenames map[uint32][]string
	}{
		{
			PreferredVersions: []uint32{2, 1},
			ExpectedVersions:  []uint32{2, 1},
			ExpectedFilenames: map[uint32][]string{
				2: {"file.hcl"},
				1: {"file1.hcl", "file2.hcl", "testdir/file3.hcl"},
			},
		},
		{
			PreferredVersions: []uint32{1},
			ExpectedVersions:  []uint32{1},
			ExpectedFilenames: map[uint32][]string{
				1: {"file1.hcl", "file2.hcl", "testdir/file3.hcl"},
			},
		},
		{
			PreferredVersions: []uint32{3},
		},
	} {
		info, err := EmbeddedReader(testdata, "testdata")(hclog.NewNullLogger(), "testmod", tc.PreferredVersions)
		assert.NoError(t, err)
		assert.NotNil(t, info)

		assert.EqualValues(t, []uint32{1, 2}, info.AvailableVersions)

		if len(tc.ExpectedVersions) == 0 {
			continue
		}

		assert.Equal(t, len(info.Data), len(tc.ExpectedVersions))

		for v, expected := range tc.ExpectedFilenames {
			var fnlist []string
			for _, f := range info.Data[v].Files {
				fnlist = append(fnlist, f.Name)
			}
			sort.Strings(fnlist)
			assert.Equal(t, expected, fnlist)
		}
	}
}
