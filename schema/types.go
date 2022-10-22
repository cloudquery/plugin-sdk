package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/cqtypes"
)

type Status byte

const (
	Undefined Status = iota
	Null
	Present
)

type CQType interface {
	// Type() ValueType
	Set(v interface{}) error
	Get() interface{}
	// Equal(CQType) bool
	// IsValid() bool
}

type CQTypes []CQType

func (c CQTypes) MarshalJSON() ([]byte, error) {
	var res []map[string]interface{}
	for _, v := range c {
		if v == nil {
			res = append(res, nil)
			continue
		}
		var typ string
		switch v.(type) {
		case *cqtypes.Bool:
			typ = "Bool"
		case *cqtypes.Int8:
			typ = "Int8"
		case *cqtypes.Float8:
			typ = "Float8"
		case *cqtypes.UUID:
			typ = "UUID"
		case *cqtypes.Text:
			typ = "Text"
		case *cqtypes.Bytea:
			typ = "Bytea"
		case *cqtypes.TextArray:
			typ = "TextArray"
		case *cqtypes.Int8Array:
			typ = "Int8Array"
		case *cqtypes.Timestamptz:
			typ = "Timestamptz"
		case *cqtypes.JSON:
			typ = "JSON"
		case *cqtypes.UUIDArray:
			typ = "UUIDArray"
		case *cqtypes.Inet:
			typ = "Inet"
		case *cqtypes.InetArray:
			typ = "InetArray"
		case *cqtypes.CIDR:
			typ = "CIDR"
		case *cqtypes.CIDRArray:
			typ = "CIDRArray"
		case *cqtypes.Macaddr:
			typ = "Macaddr"
		case *cqtypes.MacaddrArray:
			typ = "MacaddrArray"
		default:
			return nil, fmt.Errorf("unknown type %T", v)
		}

		res = append(res, map[string]interface{}{
			"type":  typ,
			"value": v,
		})
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
			(*c)[i] = &cqtypes.Bool{}
		case "Int8":
			(*c)[i] = &cqtypes.Int8{}
		case "Float8":
			(*c)[i] = &cqtypes.Float8{}
		case "UUID":
			(*c)[i] = &cqtypes.UUID{}
		case "Text":
			(*c)[i] = &cqtypes.Text{}
		case "Bytea":
			(*c)[i] = &cqtypes.Bytea{}
		case "TextArray":
			(*c)[i] = &cqtypes.TextArray{}
		case "Int8Array":
			(*c)[i] = &cqtypes.Int8Array{}
		case "Timestamptz":
			(*c)[i] = &cqtypes.Timestamptz{}
		case "JSON":
			(*c)[i] = &cqtypes.JSON{}
		case "UUIDArray":
			(*c)[i] = &cqtypes.UUIDArray{}
		case "Inet":
			(*c)[i] = &cqtypes.Inet{}
		case "InetArray":
			(*c)[i] = &cqtypes.InetArray{}
		case "CIDR":
			(*c)[i] = &cqtypes.CIDR{}
		case "CIDRArray":
			(*c)[i] = &cqtypes.CIDRArray{}
		case "Macaddr":
			(*c)[i] = &cqtypes.Macaddr{}
		case "MacaddrArray":
			(*c)[i] = &cqtypes.MacaddrArray{}
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

// func (c CQTypes) Equal(other CQTypes) bool {
// 	if other == nil {
// 		return false
// 	}
// 	if len(c) != len(other) {
// 		return false
// 	}
// 	for i := range c {
// 		if c[i] == nil {
// 			if other[i] != nil {
// 				return false
// 			}
// 		} else {
// 			if !c[i].Equal(other[i]) {
// 				return false
// 			}
// 		}

// 	}
// 	return true
// }

func (c CQTypes) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range c {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString("]")
	return sb.String()
}
