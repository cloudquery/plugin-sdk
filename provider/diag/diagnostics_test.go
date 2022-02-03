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
						Detail:   "[Repeated:2]",
					},
				},
			},
		},
		{
			Name: "simple squash w/details",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", "some details"),
				NewBaseError(errors.New("error test"), ERROR, RESOLVING, "a", "some summary", "some details"),
				NewBaseError(errors.New("error test2"), ERROR, RESOLVING, "a", "some summary2", "some details2."),
				NewBaseError(errors.New("error test2"), ERROR, RESOLVING, "a", "some summary2", "some details2."),
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
						Detail:   "some details. [Repeated:2]",
					},
				},
				{
					Err:      "error test2",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary2",
					Description: Description{
						Resource: "a",
						Summary:  "some summary2",
						Detail:   "some details2. [Repeated:2]",
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
			sq := tc.Value.Squash()
			assert.Equal(t, tc.Want, FlattenDiags(sq, false))
			assert.Equal(t, tc.Want, FlattenDiags(sq.Squash(), false)) // double squash, should still work
		})
	}
}

func TestDiagnostics_SquashRedactable(t *testing.T) {
	input := Diagnostics{
		NewRedactedDiagnostic(
			NewBaseError(errors.New("error test: 123"), ERROR, RESOLVING, "a", "some summary: 123", ""),
			NewBaseError(errors.New("error test: xxx"), ERROR, RESOLVING, "a", "some summary: xxx", ""),
		),
		NewRedactedDiagnostic(
			NewBaseError(errors.New("error test: 123"), ERROR, RESOLVING, "a", "some summary: 123", ""),
			NewBaseError(errors.New("error test: xxx"), ERROR, RESOLVING, "a", "some summary: xxx", ""),
		),
	}
	out := input.Squash()

	assert.Equal(t, []FlatDiag{
		{
			Err:      "error test: 123",
			Resource: "a",
			Type:     RESOLVING,
			Severity: ERROR,
			Summary:  "some summary: 123",
			Description: Description{
				Resource: "a",
				Summary:  "some summary: 123",
				Detail:   "[Repeated:2]",
			},
		},
	}, FlattenDiags(out, false))

	assert.Len(t, out, 1)

	rd, ok := out[0].(Redactable)
	assert.True(t, ok)
	assert.NotNil(t, rd)

	if t.Failed() {
		t.FailNow()
	}

	r := rd.Redacted()
	assert.NotNil(t, r)

	assert.Equal(t, []FlatDiag{
		{
			Err:      "error test: xxx",
			Resource: "a",
			Type:     RESOLVING,
			Severity: ERROR,
			Summary:  "some summary: xxx",
			Description: Description{
				Resource: "a",
				Summary:  "some summary: xxx",
				Detail:   "[Repeated:2]",
			},
		},
	}, FlattenDiags(Diagnostics{r}, false))

}
