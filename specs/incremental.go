package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Incremental int

const (
	IncrementalNone Incremental = iota
	IncrementalTablesOnly
	IncrementalBoth
)

func (r Incremental) String() string {
	return [...]string{"none", "only", "both"}[r]
}
func (r Incremental) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (r *Incremental) UnmarshalJSON(data []byte) (err error) {
	var registry string
	if err := json.Unmarshal(data, &registry); err != nil {
		return err
	}
	if *r, err = IncrementalFromString(registry); err != nil {
		return err
	}
	return nil
}

func IncrementalFromString(s string) (Incremental, error) {
	switch s {
	case "none":
		return IncrementalNone, nil
	case "only":
		return IncrementalTablesOnly, nil
	case "both":
		return IncrementalBoth, nil
	default:
		return IncrementalNone, fmt.Errorf("unknown registry %s", s)
	}
}
