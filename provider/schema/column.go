package schema

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"strings"
	"time"

	gofrs "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"github.com/modern-go/reflect2"
	"github.com/thoas/go-funk"
)

type ValueType int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeSmallInt
	TypeInt
	TypeBigInt
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
)

func (v ValueType) String() string {
	switch v {
	case TypeBool:
		return "TypeBool"
	case TypeBigInt:
		return "TypeBigInt"
	case TypeSmallInt:
		return "TypeSmallInt"
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

func ValueTypeFromString(s string) ValueType {
	switch strings.ToLower(s) {
	case "bool", "TypeBool":
		return TypeBool
	case "int", "TypeInt":
		return TypeInt
	case "bigint", "TypeBigInt":
		return TypeBigInt
	case "smallint", "TypeSmallInt":
		return TypeSmallInt
	case "float", "TypeFloat":
		return TypeFloat
	case "uuid", "TypeUUID":
		return TypeUUID
	case "string", "TypeString":
		return TypeString
	case "json", "TypeJSON":
		return TypeJSON
	case "intarray", "TypeIntArray":
		return TypeIntArray
	case "stringarray", "TypeStringArray":
		return TypeStringArray
	case "bytearray":
		return TypeByteArray
	case "timestamp", "TypeTimestamp":
		return TypeTimestamp
	case "uuidarray", "TypeUUIDArray":
		return TypeUUIDArray
	case "inet", "TypeInet":
		return TypeInet
	case "inetrarray", "TypeInetArray":
		return TypeInetArray
	case "macaddr", "TypeMacAddr":
		return TypeMacAddr
	case "macaddrarray", "TypeMacAddrArray":
		return TypeMacAddrArray
	case "cidr", "TypeCIDR":
		return TypeCIDR
	case "cidrarray", "TypeCIDRArray":
		return TypeCIDRArray
	case "invalid", "TypeInvalid":
		return TypeInvalid
	default:
		return TypeInvalid
	}
}

// ColumnResolver is called for each row received in TableResolver's data fetch.
// execution holds all relevant information regarding execution as well as the Column called.
// resource holds the current row we are resolving the column for.
type ColumnResolver func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error

// ColumnCreationOptions allow modification of how column is defined when table is created
type ColumnCreationOptions struct {
	Nullable bool
	Unique   bool
}

// Column definition for Table
type Column struct {
	// Name of column
	Name string
	// Value Type of column i.e String, UUID etc'
	Type ValueType
	// Description about column, this description is added as a comment in the database
	Description string
	// Default value if the resolver/default getting gets a nil value
	Default interface{}
	// Column Resolver allows to set you own data based on resolving this can be an API call or setting multiple embedded values etc'
	Resolver ColumnResolver
	// Ignore errors checks if returned error from column resolver should be ignored.
	IgnoreError IgnoreErrorFunc
	// Creation options allow modifying how column is defined when table is created
	CreationOptions ColumnCreationOptions
	// IgnoreInTests if true this skips this column in tests as sometimes it might be hard
	// to create a reproducible test environment with this column being non nill. For example various error columns and so on
	IgnoreInTests bool

	// meta holds serializable information about the column's resolvers and functions
	meta *ColumnMeta
}

func (c Column) ValidateType(v interface{}) error {
	if !c.checkType(v) {
		return fmt.Errorf("column %s expected %s got %T", c.Name, c.Type.String(), v)
	}
	return nil
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
	case int8, *int8, uint8, *uint8, int16, *int16:
		return c.Type == TypeSmallInt
	case uint16, int32, *int32:
		return c.Type == TypeInt
	case int, *int, uint32, *uint32, int64, *int64:
		return c.Type == TypeBigInt
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
	case []int, []*int, *[]int:
		return c.Type == TypeIntArray || c.Type == TypeJSON
	case []interface{}:
		return c.Type == TypeJSON
	case time.Time, *time.Time:
		return c.Type == TypeTimestamp
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
		if kindName == reflect.Slice {
			itemKind := reflect2.TypeOf(v).Type1().Elem().Kind()
			if c.Type == TypeStringArray && reflect.String == itemKind {
				return true
			}
			if c.Type == TypeIntArray && reflect.Int == itemKind {
				return true
			}
			if c.Type == TypeJSON && reflect.Struct == itemKind {
				return true
			}
			if c.Type == TypeUUIDArray && reflect2.TypeOf(v).String() == "uuid.UUID" || reflect2.TypeOf(v).String() == "*uuid.UUID" {
				return c.Type == TypeUUIDArray
			}
		}
		if kindName == reflect.Struct {
			return c.Type == TypeJSON
		}
		if c.Type == TypeSmallInt && (kindName == reflect.Int8 || kindName == reflect.Int16 || kindName == reflect.Uint8) {
			return true
		}

		if c.Type == TypeInt && (kindName == reflect.Uint16 || kindName == reflect.Int32) {
			return true
		}
		if c.Type == TypeBigInt && (kindName == reflect.Int || kindName == reflect.Int64 || kindName == reflect.Uint || kindName == reflect.Uint32 || kindName == reflect.Uint64) {
			return true
		}
	}

	return false
}

func (c Column) Meta() *ColumnMeta {
	if c.meta != nil {
		return c.meta
	}
	if c.Resolver == nil {
		return &ColumnMeta{
			Resolver:     nil,
			IgnoreExists: c.IgnoreError != nil,
		}
	}
	fnName := runtime.FuncForPC(reflect.ValueOf(c.Resolver).Pointer()).Name()
	return &ColumnMeta{
		Resolver: &ResolverMeta{
			Name:    strings.TrimPrefix(fnName, "github.com/cloudquery/cq-provider-sdk/provider/"),
			Builtin: strings.HasPrefix(fnName, "github.com/cloudquery/cq-provider-sdk/"),
		},
		IgnoreExists: c.IgnoreError != nil,
	}
}

type ResolverMeta struct {
	Name    string
	Builtin bool
}

type ColumnMeta struct {
	Resolver     *ResolverMeta
	IgnoreExists bool
}

func SetColumnMeta(c Column, m *ColumnMeta) Column {
	c.meta = m
	return c
}
