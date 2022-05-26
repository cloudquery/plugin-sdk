package diag

import "sort"

// FlatDiag is a structured diagnostic, usually can be used to create a json of diagnostics or testing.
type FlatDiag struct {
	Err         string
	Resource    string
	ResourceID  []string
	Type        Type
	Severity    Severity
	Summary     string
	Description Description
}

type FlatDiags []FlatDiag

var _ sort.Interface = (*FlatDiags)(nil)

// FlattenDiags converts Diagnostics to an array of FlatDiag
func FlattenDiags(dd Diagnostics, skipDescription bool) FlatDiags {
	df := make(FlatDiags, len(dd))
	for i, d := range dd {
		description := d.Description()
		df[i] = FlatDiag{
			Err:      d.Error(),
			Resource: description.Resource,
			Type:     d.Type(),
			Severity: d.Severity(),
			Summary:  description.Summary,
		}
		if len(description.ResourceID) > 0 {
			df[i].ResourceID = description.ResourceID
		}
		if !skipDescription {
			df[i].Description = description
		}
	}
	return df
}

func (diags FlatDiags) Len() int      { return len(diags) }
func (diags FlatDiags) Swap(i, j int) { diags[i], diags[j] = diags[j], diags[i] }
func (diags FlatDiags) Less(i, j int) bool {
	if diags[i].Severity > diags[j].Severity {
		return true
	} else if diags[i].Severity < diags[j].Severity {
		return false
	}

	if diags[i].Type > diags[j].Type {
		return true
	} else if diags[i].Type < diags[j].Type {
		return false
	}
	return diags[i].Description.Resource < diags[j].Description.Resource
}
