package diag

type RedactedDiagnostic struct {
	Diagnostic
	redacted Diagnostic
}

func NewRedactedDiagnostic(vanilla, redacted Diagnostic) RedactedDiagnostic {
	return RedactedDiagnostic{
		Diagnostic: vanilla,
		redacted:   redacted,
	}
}

func (p RedactedDiagnostic) Redacted() Diagnostic {
	return p.redacted
}

type Redactable interface {
	Redacted() Diagnostic
}

var (
	_ Redactable = (*RedactedDiagnostic)(nil)
)
