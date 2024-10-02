package scheduler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
)

type Strategy int

func (s *Strategy) String() string {
	if s == nil {
		return ""
	}
	return AllStrategyNames[*s]
}

// MarshalJSON implements json.Marshaler.
func (s *Strategy) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	if s == nil {
		b.Write([]byte("null"))
		return b.Bytes(), nil
	}
	b.Write([]byte{'"'})
	b.Write([]byte(s.String()))
	b.Write([]byte{'"'})
	return b.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Strategy) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}
	strategy, err := StrategyForName(name)
	if err != nil {
		return err
	}
	*s = strategy
	return nil
}

func (s *Strategy) Validate() error {
	if s == nil {
		return errors.New("scheduler strategy is nil")
	}
	for _, strategy := range AllStrategies {
		if strategy == *s {
			return nil
		}
	}
	return fmt.Errorf("unknown scheduler strategy: %d", s)
}

func (Strategy) JSONSchema() *jsonschema.Schema {
	enum := make([]any, len(AllStrategyNames))
	for i, s := range AllStrategyNames {
		enum[i] = s
	}
	return &jsonschema.Schema{
		Type:    "string",
		Enum:    enum,
		Default: AllStrategyNames[StrategyDFS],
		Title:   "CloudQuery scheduling strategy",
	}
}

var AllStrategies = Strategies{StrategyDFS, StrategyRoundRobin, StrategyShuffle, StrategyRandomQueue}
var AllStrategyNames = [...]string{
	StrategyDFS:         "dfs",
	StrategyRoundRobin:  "round-robin",
	StrategyShuffle:     "shuffle",
	StrategyRandomQueue: "random-queue",
}

func StrategyForName(s string) (Strategy, error) {
	for i, name := range AllStrategyNames {
		if name == s {
			return AllStrategies[i], nil
		}
	}
	return StrategyDFS, fmt.Errorf("unknown scheduler strategy: %s", s)
}

type Strategies []Strategy

func (s Strategies) String() string {
	var buffer bytes.Buffer
	for i, strategy := range s {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(strategy.String())
	}
	return buffer.String()
}
