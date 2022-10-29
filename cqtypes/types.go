package cqtypes

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Status byte

const (
	Undefined Status = iota
	Null
	Present
)

type InfinityModifier int8

const (
	Infinity         InfinityModifier = 1
	None             InfinityModifier = 0
	NegativeInfinity InfinityModifier = -Infinity
)

type CQType interface {
	// Type() ValueType
	Set(v interface{}) error
	Get() interface{}
	String() string
	Equal(CQType) bool
	// IsValid() bool
}

type CQTypes []CQType

func (c CQTypes) MarshalJSON() ([]byte, error) {
	res := make([]map[string]interface{}, len(c))
	for i, v := range c {
		if v == nil {
			res[i] = nil
			continue
		}
		var typ string
		switch v.(type) {
		case *Bool:
			typ = "Bool"
		case *Int8:
			typ = "Int8"
		case *Float8:
			typ = "Float8"
		case *UUID:
			typ = "UUID"
		case *Text:
			typ = "Text"
		case *Bytea:
			typ = "Bytea"
		case *TextArray:
			typ = "TextArray"
		case *Int8Array:
			typ = "Int8Array"
		case *Timestamptz:
			typ = "Timestamptz"
		case *JSON:
			typ = "JSON"
		case *UUIDArray:
			typ = "UUIDArray"
		case *Inet:
			typ = "Inet"
		case *InetArray:
			typ = "InetArray"
		case *CIDR:
			typ = "CIDR"
		case *CIDRArray:
			typ = "CIDRArray"
		case *Macaddr:
			typ = "Macaddr"
		case *MacaddrArray:
			typ = "MacaddrArray"
		default:
			return nil, fmt.Errorf("unknown type %T", v)
		}
		res[i] = map[string]interface{}{
			"type":  typ,
			"value": v,
		}
	}
	return json.Marshal(res)
}

func (c *CQTypes) UnmarshalJSON(b []byte) error {
	var res []map[string]interface{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	*c = make(CQTypes, len(res))
	for i := range res {
		if res[i] == nil {
			(*c)[i] = nil
			continue
		}
		b, err := json.Marshal(res[i]["value"])
		if err != nil {
			return err
		}
		typ := res[i]["type"].(string)
		switch typ {
		case "Bool":
			(*c)[i] = &Bool{}
		case "Int8":
			(*c)[i] = &Int8{}
		case "Float8":
			(*c)[i] = &Float8{}
		case "UUID":
			(*c)[i] = &UUID{}
		case "Text":
			(*c)[i] = &Text{}
		case "Bytea":
			(*c)[i] = &Bytea{}
		case "TextArray":
			(*c)[i] = &TextArray{}
		case "Int8Array":
			(*c)[i] = &Int8Array{}
		case "Timestamptz":
			(*c)[i] = &Timestamptz{}
		case "JSON":
			(*c)[i] = &JSON{}
		case "UUIDArray":
			(*c)[i] = &UUIDArray{}
		case "Inet":
			(*c)[i] = &Inet{}
		case "InetArray":
			(*c)[i] = &InetArray{}
		case "CIDR":
			(*c)[i] = &CIDR{}
		case "CIDRArray":
			(*c)[i] = &CIDRArray{}
		case "Macaddr":
			(*c)[i] = &Macaddr{}
		case "MacaddrArray":
			(*c)[i] = &MacaddrArray{}
		default:
			return fmt.Errorf("unknown type %v", typ)
		}
		if err := json.Unmarshal(b, (*c)[i]); err != nil {
			return err
		}
	}
	return nil
}

func (c CQTypes) Len() int {
	return len(c)
}

func (c CQTypes) Equal(other CQTypes) bool {
	if other == nil {
		return false
	}
	if len(c) != len(other) {
		return false
	}
	for i := range c {
		if c[i] == nil && other[i] == nil {
			continue
		}
		if c[i] == nil || other[i] == nil {
			return false
		}
		if !c[i].Equal(other[i]) {
			return false
		}
	}
	return true
}

func (c CQTypes) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range c {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}
