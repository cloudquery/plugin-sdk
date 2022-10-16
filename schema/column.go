package schema

import (
	"context"
	"net"
	"reflect"
	"time"

	gofrs "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"github.com/modern-go/reflect2"
	"github.com/thoas/go-funk"
)

type ValueType int


type ColumnList []Column

// ColumnResolver is called for each row received in TableResolver's data fetch.
// execution holds all relevant information regarding execution as well as the Column called.
// resource holds the current row we are resolving the column for.
type ColumnResolver func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error

// ColumnCreationOptions allow modification of how column is defined when table is created
type ColumnCreationOptions struct {
	PrimaryKey bool `json:"primary_key,omitempty"`
}

// Column definition for Table
type Column struct {
	// Name of column
	Name string `json:"name"`
	// Value Type of column i.e String, UUID etc'
	Type ValueType `json:"type"`
	// Description about column, this description is added as a comment in the database
	Description string `json:"-"`
	// Column Resolver allows to set you own data based on resolving this can be an API call or setting multiple embedded values etc'
	Resolver ColumnResolver `json:"-"`
	// Creation options allow modifying how column is defined when table is created
	CreationOptions ColumnCreationOptions
	// IgnoreInTests is used to skip verifying the column is non-nil in integration tests.
	// By default, integration tests perform a fetch for all resources in cloudquery's test account, and
	// verify all columns are non-nil.
	// If IgnoreInTests is true, verification is skipped for this column.
	// Used when it is hard to create a reproducible environment with this column being non-nil (e.g. various error columns).
	IgnoreInTests bool `json:"-"`
}

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeUUID
	TypeString
	TypeByteArray
	TypeStringArray
	TypeIntArray
	TypeTimestamp
	TypeJSON
	TypeUUIDArray
	TypeInet
	TypeInetArray
	TypeCIDR
	TypeCIDRArray
	TypeMacAddr
	TypeMacAddrArray
	TypeTimeInterval
)

func (v ValueType) String() string {
	switch v {
	case TypeBool:
		return "TypeBool"
	case TypeInt:
		return "TypeInt"
	case TypeFloat:
		return "TypeFloat"
	case TypeUUID:
		return "TypeUUID"
	case TypeString:
		return "TypeString"
	case TypeJSON:
		return "TypeJSON"
	case TypeIntArray:
		return "TypeIntArray"
	case TypeStringArray:
		return "TypeStringArray"
	case TypeTimestamp:
		return "TypeTimestamp"
	case TypeTimeInterval:
		return "TypeTimeInterval"
	case TypeByteArray:
		return "TypeByteArray"
	case TypeUUIDArray:
		return "TypeUUIDArray"
	case TypeInetArray:
		return "TypeInetArray"
	case TypeInet:
		return "TypeInet"
	case TypeMacAddrArray:
		return "TypeMacAddrArray"
	case TypeMacAddr:
		return "TypeMacAddr"
	case TypeCIDRArray:
		return "TypeCIDRArray"
	case TypeCIDR:
		return "TypeCIDR"
	case TypeInvalid:
		fallthrough
	default:
		return "TypeInvalid"
	}
}



func (c Column) checkType(v interface{}) bool {
	if reflect2.IsNil(v) {
		return true
	}

	if reflect2.TypeOf(v).Kind() == reflect.Ptr {
		return c.checkType(funk.GetOrElse(v, nil))
	}

	// Maps or slices are jsons
	if reflect2.TypeOf(v).Kind() == reflect.Map {
		return c.Type == TypeJSON
	}

	switch val := v.(type) {
	case int8, *int8, uint8, *uint8, int16, *int16, uint16, *uint16, int32, *int32, int, *int, uint32, *uint32, int64, *int64:
		return c.Type == TypeInt
	case []byte:
		if c.Type == TypeUUID {
			if _, err := uuid.FromBytes(val); err != nil {
				return false
			}
		}
		return c.Type == TypeByteArray || c.Type == TypeJSON
	case bool, *bool:
		return c.Type == TypeBool
	case string:
		if c.Type == TypeUUID {
			if _, err := uuid.Parse(val); err == nil {
				return true
			}
		}
		if c.Type == TypeJSON {
			return true
		}
		return c.Type == TypeString
	case *string:
		if c.Type == TypeJSON {
			return true
		}
		return c.Type == TypeString
	case *float32, float32, *float64, float64:
		return c.Type == TypeFloat
	case []string, []*string, *[]string:
		return c.Type == TypeStringArray || c.Type == TypeJSON
	case []int, []*int, *[]int, []int32, []*int32, []int64, []*int64, *[]int64:
		return c.Type == TypeIntArray || c.Type == TypeJSON
	case []interface{}:
		return c.Type == TypeJSON
	case time.Time, *time.Time:
		return c.Type == TypeTimestamp
	case time.Duration, *time.Duration:
		return c.Type == TypeTimeInterval
	case uuid.UUID, *uuid.UUID:
		return c.Type == TypeUUID
	case gofrs.UUID, *gofrs.UUID:
		return c.Type == TypeUUID
	case [16]byte:
		return c.Type == TypeUUID
	case net.HardwareAddr, *net.HardwareAddr:
		return c.Type == TypeMacAddr
	case []net.HardwareAddr, []*net.HardwareAddr:
		return c.Type == TypeMacAddrArray
	case net.IPAddr, *net.IPAddr, *net.IP, net.IP:
		return c.Type == TypeInet
	case []net.IPAddr, []*net.IPAddr, []*net.IP, []net.IP:
		return c.Type == TypeInetArray
	case net.IPNet, *net.IPNet:
		return c.Type == TypeCIDR
	case []net.IPNet, []*net.IPNet:
		return c.Type == TypeCIDRArray
	case interface{}:
		kindName := reflect2.TypeOf(v).Kind()
		if kindName == reflect.String && c.Type == TypeString {
			return true
		}
		switch kindName {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return c.Type == TypeInt
		}
		if kindName == reflect.Slice {
			itemKind := reflect2.TypeOf(v).Type1().Elem().Kind()
			if c.Type == TypeStringArray && reflect.String == itemKind {
				return true
			}
			if c.Type == TypeIntArray {
				switch itemKind {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					return true
				}
			}
			if c.Type == TypeJSON && (reflect.Struct == itemKind || reflect.Ptr == itemKind) {
				return true
			}
			if c.Type == TypeUUIDArray && reflect2.TypeOf(v).String() == "uuid.UUID" || reflect2.TypeOf(v).String() == "*uuid.UUID" {
				return c.Type == TypeUUIDArray
			}
		}
		if kindName == reflect.Struct {
			return c.Type == TypeJSON
		}
	}

	return false
}


func (c ColumnList) Index(col string) int {
	for i, c := range c {
		if c.Name == col {
			return i
		}
	}
	return -1
}


func (c ColumnList) Names() []string {
	ret := make([]string, len(c))
	for i := range c {
		ret[i] = c[i].Name
	}
	return ret
}

func (c ColumnList) Get(name string) *Column {
	for i := range c {
		if c[i].Name == name {
			return &c[i]
		}
	}
	return nil
}

