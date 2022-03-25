package migration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFileStructure(t *testing.T) {
	cases := []struct {
		filenames    []string
		expecterrors []string
	}{
		{
			filenames: []string{"1_v1.0.0.up.sql", "1_v1.0.0.down.sql", "2_v1.2.0.up.sql", "2_v1.2.0.down.sql"},
		},
		{
			filenames:    []string{"1.up.sql", "1.down.sql"},
			expecterrors: []string{"less than 2 underscores"},
		},
		{
			filenames:    []string{"1_123.up.sql", "1_123.down.sql"},
			expecterrors: []string{"version should start with v"},
		},
		{
			filenames:    []string{"1_v0.0.1.up.sql", "1_v0.0.1.txt"},
			expecterrors: []string{"neither up or down migration"},
		},
		{
			filenames:    []string{"1_v1.0.0.up.sql", "1_v1.0.0.down.sql", "2_v1.2.0.up.sql", "3_v1.2.1.up.sql", "3_v1.2.1.down.sql"},
			expecterrors: []string{"missing down migration"},
		},
		{
			filenames:    []string{"1_v1.0.0.up.sql", "1_v1.0.0.down.sql", "2_v1.2.0.down.sql"},
			expecterrors: []string{"missing up migration"},
		},
		{
			filenames:    []string{"1_v0.0.1.up.sql", "2_v0.0.1.down.sql"},
			expecterrors: []string{"missing down migration", "missing up migration"},
		},
		{
			filenames:    []string{"1_v1.0.0.up.sql", "1_v1.0.0.down.sql", "2_v1.2.0.up.sql", "2_v1.2.0.down.sql", "3_v1.2.0.up.sql", "3_v1.2.0.down.sql"},
			expecterrors: []string{"mentioned in multiple versions", "mentioned in multiple versions"},
		},
		{
			filenames:    []string{"1_v1.0.0.up.sql", "1_v1.0.0.down.sql", "2_v1.2.0.up.sql", "3_v1.2.0.up.sql", "3_v1.2.0.down.sql"},
			expecterrors: []string{"missing down migration", "mentioned in multiple versions"},
		},
		{
			filenames:    []string{"21_v0.10.12.down.sql", "21_v0.10.12.up.sql", "22_v0.10.14.down.sql", "22_v0.10.14.up.sql", "23_v0.10.14.down.sql", "23_v0.10.15.up.sql"},
			expecterrors: []string{"is mentioned in multiple versions"},
		},
	}

	for _, tc := range cases {
		mf := make(map[string][]byte, len(tc.filenames))
		for _, fn := range tc.filenames {
			mf[fn] = nil
		}

		errs := checkFileStructureForDialect(mf)
		assert.Equal(t, len(tc.expecterrors), len(errs))
		if t.Failed() {
			t.Log(errs)
			t.FailNow()
		}

		matches := 0
		for _, err := range errs {
			for _, expected := range tc.expecterrors {
				if strings.Contains(err.Error(), expected) {
					matches++
					break
				}
			}
		}
		assert.Equal(t, len(tc.expecterrors), matches)
		if t.Failed() {
			t.Log(errs)
			t.FailNow()
		}
	}
}
