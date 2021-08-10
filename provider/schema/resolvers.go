package schema

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gofrs/uuid"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

// PathResolver resolves a field in the Resource.Item
//
// Examples:
// PathResolver("Field")
// PathResolver("InnerStruct.Field")
// PathResolver("InnerStruct.InnerInnerStruct.Field")
func PathResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, funk.Get(r.Item, path, funk.WithAllowZero()))
	}
}

// ParentIdResolver resolves the cq_id from the parent
// if you want to reference the parent's primary keys use ParentResourceFieldResolver as required.
func ParentIdResolver(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
	return r.Set(c.Name, r.Parent.Id())
}

// ParentResourceFieldResolver resolves a field from the parent's resource, the value is expected to be set
// if name isn't set the field will be set to null
func ParentResourceFieldResolver(name string) ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, r.Parent.Get(name))
	}
}

// ParentPathResolver resolves a field from the parent
func ParentPathResolver(path string) ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, funk.Get(r.Parent.Item, path, funk.WithAllowZero()))
	}
}

// DateUTCResolver resolves the different date formats (ISODate - 2011-10-05T14:48:00.000Z is default) into *time.Time and converts the date to utc timezone
//
// Examples:
// DateUTCResolver("Date") - resolves using RFC.RFC3339 as default
// DateUTCResolver("InnerStruct.Field", time.RFC822)  - resolves using time.RFC822
// DateUTCResolver("InnerStruct.Field", time.RFC822, "2011-10-05")  - resolves using a few resolvers one by one
func DateUTCResolver(path string, rfcs ...string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		data, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}
		date, err := parseDate(data, rfcs...)
		if err != nil {
			return err
		}
		return r.Set(c.Name, date.UTC())
	}
}

// DateResolver resolves the different date formats (ISODate - 2011-10-05T14:48:00.000Z is default) into *time.Time
//
// Examples:
// DateResolver("Date") - resolves using RFC.RFC3339 as default
// DateResolver("InnerStruct.Field", time.RFC822)  - resolves using time.RFC822
// DateResolver("InnerStruct.Field", time.RFC822, "2011-10-05")  - resolves using a few resolvers one by one
func DateResolver(path string, rfcs ...string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		data, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}
		date, err := parseDate(data, rfcs...)
		if err != nil {
			return err
		}
		return r.Set(c.Name, date)
	}
}

func parseDate(dateStr string, rfcs ...string) (date *time.Time, err error) {
	if dateStr == "" {
		return nil, nil
	}

	// set default rfc
	if len(rfcs) == 0 {
		rfcs = append(rfcs, time.RFC3339)
	}

	var d time.Time
	for _, rfc := range rfcs {
		d, err = time.Parse(rfc, dateStr)
		if err == nil {
			date = &d
			return
		}
	}
	return
}

// IPAddressResolver resolves the ip string value and returns net.IP
//
// Examples:
// IPAddressResolver("IP")
func IPAddressResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		ipStr, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}
		ip := net.ParseIP(ipStr)
		if ipStr != "" && ip == nil {
			return fmt.Errorf("failed to parse IP from %s", ipStr)
		}
		return r.Set(c.Name, ip)
	}
}

// MACAddressResolver resolves the mac string value and returns net.HardwareAddr
//
// Examples:
// MACAddressResolver("MAC")
func MACAddressResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		macStr, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}
		mac, err := net.ParseMAC(macStr)
		if err != nil {
			return err
		}
		return r.Set(c.Name, mac)
	}
}

// IPNetResolver resolves the network string value and returns net.IPNet
//
// Examples:
// IPNetResolver("Network")
func IPNetResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		ipStr, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}
		_, inet, err := net.ParseCIDR(ipStr)
		if err != nil {
			return err
		}
		return r.Set(c.Name, inet)
	}
}

// UUIDResolver resolves the uuid string value and returns uuid.UUID
//
// Examples:
// UUIDResolver("Resource.UUID")
func UUIDResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		uuidString, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}

		uuid, err := uuid.FromString(uuidString)
		if err != nil {
			return err
		}
		return r.Set(c.Name, uuid)
	}
}

// StringResolver tries to cast value into string
//
// Examples:
// StringResolver("Id")
func StringResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		str, err := cast.ToStringE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}

		return r.Set(c.Name, str)
	}
}

// IntResolver tries to cast value into int
//
// Examples:
// IntResolver("Id")
func IntResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		i, err := cast.ToIntE(funk.Get(r.Item, path, funk.WithAllowZero()))
		if err != nil {
			return err
		}

		return r.Set(c.Name, i)
	}
}
