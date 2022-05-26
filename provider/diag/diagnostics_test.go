package diag

import (
	"errors"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiagnostics_Squash(t *testing.T) {
	testCases := []struct {
		Name  string
		Value Diagnostics
		Want  FlatDiags
	}{
		{
			Name: "simple squash no details",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary")),
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary")),
			},
			Want: FlatDiags{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary: error test",
					Description: Description{
						Resource: "a",
						Summary:  "some summary: error test",
						Detail:   "[Repeated:2]",
					},
				},
			},
		},
		{
			Name: "simple squash w/details",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary"), WithDetails("some details")),
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary"), WithDetails("some details")),
				NewBaseError(errors.New("error test2"), RESOLVING, WithResourceName("a"), WithSummary("some summary2"), WithDetails("some details2.")),
				NewBaseError(errors.New("error test2"), RESOLVING, WithResourceName("a"), WithSummary("some summary2"), WithDetails("some details2.")),
			},
			Want: FlatDiags{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary: error test",
					Description: Description{
						Resource: "a",
						Summary:  "some summary: error test",
						Detail:   "some details. [Repeated:2]",
					},
				},
				{
					Err:      "error test2",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary2: error test2",
					Description: Description{
						Resource: "a",
						Summary:  "some summary2: error test2",
						Detail:   "some details2. [Repeated:2]",
					},
				},
			},
		},
		{
			Name: "different resource no squash",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary")),
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("b"), WithSummary("some summary")),
			},
			Want: FlatDiags{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary: error test",
					Description: Description{
						Resource: "a",
						Summary:  "some summary: error test",
						Detail:   "",
					},
				},
				{
					Err:      "error test",
					Resource: "b",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary: error test",
					Description: Description{
						Resource: "b",
						Summary:  "some summary: error test",
						Detail:   "",
					},
				},
			},
		},
		{
			Name: "different severity no squash",
			Value: Diagnostics{
				NewBaseError(errors.New("error test"), RESOLVING, WithSeverity(WARNING), WithResourceName("a"), WithSummary("some summary")),
				NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary")),
			},
			Want: FlatDiags{
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: WARNING,
					Summary:  "some summary: error test",
					Description: Description{
						Resource: "a",
						Summary:  "some summary: error test",
						Detail:   "",
					},
				},
				{
					Err:      "error test",
					Resource: "a",
					Type:     RESOLVING,
					Severity: ERROR,
					Summary:  "some summary: error test",
					Description: Description{
						Resource: "a",
						Summary:  "some summary: error test",
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
			NewBaseError(errors.New("error test: 123"), RESOLVING, WithResourceName("a"), WithSummary("some summary with 456")),
			NewBaseError(errors.New("error test: xxx"), RESOLVING, WithResourceName("a"), WithSummary("some summary with xxx")),
		),
		NewRedactedDiagnostic(
			NewBaseError(errors.New("error test: 123"), RESOLVING, WithResourceName("a"), WithSummary("some summary with 456")),
			NewBaseError(errors.New("error test: xxx"), RESOLVING, WithResourceName("a"), WithSummary("some summary with xxx")),
		),
	}
	out := input.Squash()

	assert.Equal(t, FlatDiags{
		{
			Err:      "error test: 123",
			Resource: "a",
			Type:     RESOLVING,
			Severity: ERROR,
			Summary:  "some summary with 456: error test: 123",
			Description: Description{
				Resource: "a",
				Summary:  "some summary with 456: error test: 123",
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

	assert.Equal(t, FlatDiags{
		{
			Err:      "error test: xxx",
			Resource: "a",
			Type:     RESOLVING,
			Severity: ERROR,
			Summary:  "some summary with xxx: error test: xxx",
			Description: Description{
				Resource: "a",
				Summary:  "some summary with xxx: error test: xxx",
				Detail:   "[Repeated:2]",
			},
		},
	}, FlattenDiags(Diagnostics{r}, false))
}

func TestDiagnostics_Sort(t *testing.T) {
	resErrA := NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("a"), WithSummary("some summary"))
	resErrB := NewBaseError(errors.New("error test"), RESOLVING, WithResourceName("b"), WithSummary("some summary"))

	thrErrA := NewBaseError(errors.New("error test"), THROTTLE, WithResourceName("a"), WithSummary("some summary"))
	accErrA := NewBaseError(errors.New("error test"), ACCESS, WithResourceName("a"), WithSummary("some summary"))

	cases := []struct {
		diagsInWrongOrder Diagnostics
		expectedOrder     []int
	}{
		{
			Diagnostics{resErrB, resErrA},
			[]int{1, 0},
		},
		{
			Diagnostics{accErrA, resErrB, thrErrA, resErrA},
			[]int{2, 0, 3, 1},
		},
	}
	for caseNo, tc := range cases {
		assert.Equal(t, len(tc.diagsInWrongOrder), len(tc.expectedOrder), "bad test")
		if t.Failed() {
			t.FailNow()
		}

		sorted := make(Diagnostics, len(tc.diagsInWrongOrder))
		copy(sorted, tc.diagsInWrongOrder)
		sort.Sort(sorted)
		for i := range sorted {
			assert.Equalf(t, sorted[i], tc.diagsInWrongOrder[tc.expectedOrder[i]], "Case #%d item %d", caseNo+1, i)
		}
	}
}

func TestDiagnostics_BySeverity(t *testing.T) {
	diagSet := Diagnostics{
		NewBaseError(errors.New("warn test 1"), RESOLVING, WithSeverity(WARNING)),
		NewBaseError(errors.New("err test 1"), RESOLVING),
		NewBaseError(errors.New("ign test 1"), RESOLVING, WithSeverity(IGNORE)),
		NewBaseError(errors.New("warn test 2"), RESOLVING, WithSeverity(WARNING)),
		NewBaseError(errors.New("ign test 2"), RESOLVING, WithSeverity(IGNORE)),
		NewBaseError(errors.New("err test 2"), RESOLVING),
	}

	cases := []struct {
		sevs         []Severity
		expectedErrs []string
	}{
		{
			sevs:         []Severity{WARNING},
			expectedErrs: []string{"warn test 1", "warn test 2"},
		},
		{
			sevs:         []Severity{ERROR, WARNING},
			expectedErrs: []string{"warn test 1", "err test 1", "warn test 2", "err test 2"},
		},
		{
			sevs:         []Severity{IGNORE, PANIC},
			expectedErrs: []string{"ign test 1", "ign test 2"},
		},
		{
			sevs:         []Severity{PANIC},
			expectedErrs: []string{},
		},
	}
	for caseNo := range cases {
		t.Run("Test #"+strconv.Itoa(caseNo+1), func(t *testing.T) {
			tc := cases[caseNo]
			res := diagSet.BySeverity(tc.sevs...)
			assert.Equal(t, len(tc.expectedErrs), len(res))
			resErrs := make([]string, len(res))
			for i := range res {
				resErrs[i] = res[i].Error()
			}
			assert.Equal(t, tc.expectedErrs, resErrs)
		})
	}
}
