package diag

// FlatDiag is a structured diagnostic, usually can be used to create a json of diagnostics or testing.
type FlatDiag struct {
	Err         string
	Resource    string
	Type        DiagnosticType
	Severity    Severity
	Summary     string
	Description Description
}

// FlattenDiags converts Diagnostics to an array of FlatDiag
func FlattenDiags(dd Diagnostics, skipDescription bool) []FlatDiag {
	df := make([]FlatDiag, len(dd))
	for i, d := range dd {
		description := d.Description()
		df[i] = FlatDiag{
			Err:      d.Error(),
			Resource: description.Resource,
			Type:     d.Type(),
			Severity: d.Severity(),
			Summary:  description.Summary,
		}
		if !skipDescription {
			df[i].Description = description
		}
	}
	return df
}
