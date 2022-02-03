package diag

import (
	"fmt"
	"strings"
)

type SquashedDiag struct {
	Diagnostic
	count uint64
}

func (s SquashedDiag) Description() Description {
	description := s.Diagnostic.Description()

	if _, ok := s.Diagnostic.(Countable); ok { // already squashed, don't add repeat count
		return description
	}

	switch {
	case s.count == 1:
		// no-op
	case description.Detail == "":
		description.Detail = fmt.Sprintf("[Repeated:%d]", s.count)
	case strings.HasSuffix(description.Detail, "."):
		description.Detail = fmt.Sprintf("%s [Repeated:%d]", description.Detail, s.count)
	default:
		description.Detail = fmt.Sprintf("%s. [Repeated:%d]", description.Detail, s.count)
	}

	return description
}

func (s SquashedDiag) Count() uint64 {
	return s.count
}

func (s SquashedDiag) Redacted() Diagnostic {
	rd, ok := s.Diagnostic.(Redactable)
	if !ok {
		return nil
	}

	r := rd.Redacted()
	if r == nil {
		return nil
	}

	return SquashedDiag{
		Diagnostic: r,
		count:      s.count,
	}
}

type Countable interface {
	Count() uint64
}

func CountDiag(d Diagnostic) uint64 {
	if c, ok := d.(Countable); ok {
		return c.Count()
	}

	return 1
}

var _ Countable = (*SquashedDiag)(nil)
