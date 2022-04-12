package diag

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/hcl/v2"
)

type Diagnostics []Diagnostic

func (diags Diagnostics) Error() string {
	switch {
	case len(diags) == 0:
		// should never happen, since we don't create this wrapper if
		// there are no diagnostics in the list.
		return "no errors"
	case len(diags) == 1:
		desc := diags[0].Description()
		var ret bytes.Buffer
		if len(desc.ResourceID) > 0 {
			fmt.Fprintf(&ret, "[%s] ", strings.Join(desc.ResourceID, ","))
		}
		if desc.Detail == "" {
			fmt.Fprintf(&ret, "%s", desc.Summary)
		} else {
			fmt.Fprintf(&ret, "%s: %s", desc.Summary, desc.Detail)
		}
		return ret.String()
	default:
		var ret bytes.Buffer
		fmt.Fprintf(&ret, "%d problems:\n", len(diags))
		for _, diag := range diags {
			desc := diag.Description()
			fmt.Fprintf(&ret, "\n- ")
			if len(desc.ResourceID) > 0 {
				fmt.Fprintf(&ret, "[%s] ", strings.Join(desc.ResourceID, ","))
			}
			if desc.Detail == "" {
				fmt.Fprintf(&ret, "%s", desc.Summary)
			} else {
				fmt.Fprintf(&ret, "%s: %s", desc.Summary, desc.Detail)
			}
		}
		return ret.String()
	}
}

func (diags Diagnostics) HasErrors() bool {
	for _, d := range diags {
		if d.Severity() >= ERROR {
			return true
		}
	}
	return false
}

func (diags Diagnostics) HasDiags() bool {
	return len(diags) > 0
}

func (diags Diagnostics) Add(new ...interface{}) Diagnostics {
	for _, item := range new {
		if item == nil {
			continue
		}
		switch ti := item.(type) {
		case Diagnostic:
			diags = append(diags, ti)
		case Diagnostics:
			diags = append(diags, ti...) // flatten
		case error:
			switch {
			case errwrap.ContainsType(ti, Diagnostics(nil)):
				// If we have an errwrap wrapper with a Diagnostics hiding
				// inside then we'll unpick it here to get access to the
				// individual diagnostics.
				diags = diags.Add(errwrap.GetType(ti, Diagnostics(nil)))
			case errwrap.ContainsType(ti, hcl.Diagnostics(nil)):
				// Likewise, if we have HCL diagnostics we'll unpick that too.
				diags = diags.Add(errwrap.GetType(ti, hcl.Diagnostics(nil)))
			default:
				diags = append(diags, nativeError{ti})
			}
		default:
			panic(fmt.Errorf("can't construct diagnostic(s) from %T", item))
		}
	}

	// Given the above, we should never end up with a non-nil empty slice
	// here, but we'll make sure of that so callers can rely on empty == nil
	if len(diags) == 0 {
		return nil
	}

	return diags
}

// Squash attempts to squash diagnostics
func (diags Diagnostics) Squash() Diagnostics {
	dd := make(map[string]*SquashedDiag, len(diags))
	sdd := make(Diagnostics, 0)
	for i, d := range diags {
		keygen := d
		if rd, ok := d.(Redactable); ok {
			if r := rd.Redacted(); r != nil {
				keygen = r
			}
		}

		key := fmt.Sprintf("%s_%s_%d_%d", keygen.Error(), keygen.Description().Resource, keygen.Severity(), keygen.Type())
		if sd, ok := dd[key]; ok {
			sd.count += CountDiag(d)
			continue
		}
		nsd := &SquashedDiag{
			Diagnostic: diags[i],
			count:      CountDiag(d),
		}
		dd[key] = nsd
		sdd = append(sdd, nsd)
	}
	return sdd
}

func (diags Diagnostics) Warnings() uint64 {
	return diags.CountBySeverity(WARNING, false)
}

func (diags Diagnostics) Errors() uint64 {
	return diags.CountBySeverity(ERROR, false)
}

// CountBySeverity returns number of diagnostics of the given severity. If includeSquashed is false, squashed diags are counted as a single diag.
func (diags Diagnostics) CountBySeverity(sev Severity, includeSquashed bool) uint64 {
	var count uint64 = 0

	for _, d := range diags {
		if d.Severity() != sev {
			continue
		}
		if includeSquashed {
			count += CountDiag(d)
		} else {
			count++
		}
	}
	return count
}

func (diags Diagnostics) Redacted() Diagnostics {
	res := make(Diagnostics, len(diags))
	for i := range diags {
		if rd, ok := diags[i].(Redactable); ok {
			if r := rd.Redacted(); r != nil {
				res[i] = r
				continue
			}
		}

		res[i] = diags[i]
	}
	return res
}

func (diags Diagnostics) Len() int      { return len(diags) }
func (diags Diagnostics) Swap(i, j int) { diags[i], diags[j] = diags[j], diags[i] }
func (diags Diagnostics) Less(i, j int) bool {
	if diags[i].Severity() > diags[j].Severity() {
		return true
	} else if diags[i].Severity() < diags[j].Severity() {
		return false
	}

	if diags[i].Type() > diags[j].Type() {
		return true
	} else if diags[i].Type() < diags[j].Type() {
		return false
	}
	return diags[i].Description().Resource < diags[j].Description().Resource
}
