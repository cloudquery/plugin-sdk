package diag

import "fmt"

type SquashedDiag struct {
	Diagnostic
	Count int
}

func (s SquashedDiag) Description() Description {
	description := s.Diagnostic.Description()
	if s.Count == 1 {
		return description
	}
	if description.Detail == "" {
		description.Detail = fmt.Sprintf("Repeated[%d]", s.Count)
	} else {
		description.Detail = fmt.Sprintf("Repeated[%d]: %s", s.Count, description.Detail)
	}

	return description
}
