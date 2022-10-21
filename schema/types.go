package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CQType interface{
	Type() ValueType
	Scan(v interface{}) error
	Equal(CQType) bool
}



type CQTypes []CQType

func (c CQTypes) MarshalJSON() ([]byte, error) {
	var res []map[string]interface{}
	for _, v := range c {
		if v == nil {
			res = append(res, nil)
			continue
		}
		res = append(res, map[string]interface{}{
			"type": v.Type(),
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
		typ := ValueType(int(res[i]["type"].(float64)))

		switch typ {
			case TypeBool:
				var r Bool
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				(*c)[i] = &r
			case TypeInt:
				var r Int64
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				(*c)[i] = &r
			case TypeUUID:
				var r UUID
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
			case TypeString:
				var r String
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown type %v", typ)
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
		if c[i] == nil {
			if other[i] != nil {
				return false
			}
		} else {
			if !c[i].Equal(other[i]) {
				return false
			}
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
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString("]")
	return sb.String()
}