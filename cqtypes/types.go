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
	Set(v interface{}) error
	Get() interface{}
	String() string
	Equal(CQType) bool
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
	cqTypes := make(CQTypes, 0, len(res))
	for i := range res {
		if res[i] == nil {
			cqTypes = append(cqTypes, nil)
			continue
		}
		b, err := json.Marshal(res[i]["value"])
		if err != nil {
			return err
		}
		typ := res[i]["type"].(string)
		switch typ {
		case "Bool":
			cqTypes = append(cqTypes, &Bool{})
		case "Int8":
			cqTypes = append(cqTypes, &Int8{})
		case "Float8":
			cqTypes = append(cqTypes, &Float8{})
		case "UUID":
			cqTypes = append(cqTypes, &UUID{})
		case "Text":
			cqTypes = append(cqTypes, &Text{})
		case "Bytea":
			cqTypes = append(cqTypes, &Bytea{})
		case "TextArray":
			cqTypes = append(cqTypes, &TextArray{})
		case "Int8Array":
			cqTypes = append(cqTypes, &Int8Array{})
		case "Timestamptz":
			cqTypes = append(cqTypes, &Timestamptz{})
		case "JSON":
			cqTypes = append(cqTypes, &JSON{})
		case "UUIDArray":
			cqTypes = append(cqTypes, &UUIDArray{})
		case "Inet":
			cqTypes = append(cqTypes, &Inet{})
		case "InetArray":
			cqTypes = append(cqTypes, &InetArray{})
		case "CIDR":
			cqTypes = append(cqTypes, &CIDR{})
		case "CIDRArray":
			cqTypes = append(cqTypes, &CIDRArray{})
		case "Macaddr":
			cqTypes = append(cqTypes, &Macaddr{})
		case "MacaddrArray":
			cqTypes = append(cqTypes, &MacaddrArray{})
		default:
			// This means a new type was added at the SDK but not yet available
			// add the destination
			continue
		}
		if err := json.Unmarshal(b, cqTypes[i]); err != nil {
			return err
		}
	}
	*c = cqTypes
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
		if v == nil {
			sb.WriteString("nil")
		} else {
			sb.WriteString(v.String())
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// var namedValues map[string]CQType

// func init() {
// 	namedValues = map[string]CQType{
// 		"Bool":       &Bool{},
// 		"Bytea": 			&Bytea{},
// 		"CIDRArray": &CIDRArray{},
// 		"CIDR": 			&CIDR{},
// 		"Float8": 		&Float8{},
// 		"InetArray": 	&InetArray{},
// 		"Inet": 			&Inet{},
// 		"Int8Array": 	&Int8Array{},
// 		"Int8": 			&Int8{},
// 		"JSON": 			&JSON{},
// 		"MacaddrArray": &MacaddrArray{},
// 		"Macaddr": 		&Macaddr{},
// 		"TextArray": 	&TextArray{},
// 		"Text": 			&Text{},
// 		"Timestamptz": &Timestamptz{},
// 		"UUIDArray": 	&UUIDArray{},
// 		"UUID": 			&UUID{},
// 	}
// }
