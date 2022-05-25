package diag

import (
	"fmt"
	"strings"
)

type SquashedDiag struct {
	Diagnostic
	count uint64
}

type Countable interface {
	Count() uint64
}

type Unsquashable interface {
	Unsquash() Diagnostic
}

var (
	_ Countable    = (*SquashedDiag)(nil)
	_ Unsquashable = (*SquashedDiag)(nil)
)

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

// Count returns the number of diagnostics inside the squashed diagnostic
func (s SquashedDiag) Count() uint64 {
	return s.count
}

// Redacted returns the redacted version of the first diagnostic, if there is any
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

// Unsquash returns the first diagnostic of the squashed set
func (s SquashedDiag) Unsquash() Diagnostic {
	return s.Diagnostic
}

func CountDiag(d Diagnostic) uint64 {
	if c, ok := d.(Countable); ok {
		return c.Count()
	}

	return 1
}

func UnsquashDiag(d Diagnostic) Diagnostic {
	if c, ok := d.(Unsquashable); ok {
		return c.Unsquash()
	}

	return d
}
