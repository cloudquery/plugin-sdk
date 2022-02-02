package diag

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiagnostics_Squash(t *testing.T) {
	testCases := []struct {
		Name  string
		Value Diagnostics
		Want  []FlatDiag
	}{
		{
			Name: "simple squash no details",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", ""),
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", ""),
			},
			Want: []FlatDiag{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary",
					Description: Description{
						Resource: "a",
						Summary:  "some summary",
						Detail:   "Repeated[2]",
					},
				},
			},
		},
		{
			Name: "simple squash w/details",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", "some details"),
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", "some details"),
			},
			Want: []FlatDiag{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary",
					Description: Description{
						Resource: "a",
						Summary:  "some summary",
						Detail:   "Repeated[2]: some details",
					},
				},
			},
		},
		{
			Name: "different resource no squash",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", ""),
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "b", "some summary", ""),
			},
			Want: []FlatDiag{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary",
					Description: Description{
						Resource: "a",
						Summary:  "some summary",
						Detail:   "",
					},
				},
				{
					Err:      "error test",
					Resource: "b",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary",
					Description: Description{
						Resource: "b",
						Summary:  "some summary",
						Detail:   "",
					},
				},
			},
		},
		{
			Name: "different severity no squash",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), WARNING, RESOLVING, "a", "some summary", ""),
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", ""),
			},
			Want: []FlatDiag{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: WARNING,
					Summary:  "some summary",
					Description: Description{
						Resource: "a",
						Summary:  "some summary",
						Detail:   "",
					},
				},
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary",
					Description: Description{
						Resource: "a",
						Summary:  "some summary",
						Detail:   "",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Want, FlattenDiags(tc.Value.Squash(), false))
		})
	}
}
