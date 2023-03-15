package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type PKMode int

const (
	PKModeCompositeKeys PKMode = iota
	PKModeCQID
)

var (
	pkModeStrings = []string{"default", "cq-ids-only"}
)

func (m PKMode) String() string {
	return pkModeStrings[m]
}

func (m PKMode) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(m.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (m *PKMode) UnmarshalJSON(data []byte) (err error) {
	var pkMode string
	if err := json.Unmarshal(data, &pkMode); err != nil {
		return err
	}
	if *m, err = PKModeFromString(pkMode); err != nil {
		return err
	}
	return nil
}

func PKModeFromString(s string) (PKMode, error) {
	switch s {
	case "default":
		return PKModeCompositeKeys, nil
	case "cq-ids-only":
		return PKModeCQID, nil
	}
	return 0, fmt.Errorf("invalid pk mode: %s", s)
}
