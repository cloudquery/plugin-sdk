package schema

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/cqtypes"
	"github.com/google/uuid"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item interface{}
	// Set if this is an embedded table
	Parent *Resource
	// internal fields
	Table *Table
	// This is sorted result data by column name
	data []CQType
}

// This struct is what we send over the wire to destination.
// We dont want to reuse the same struct as otherwise we will have to comment on fields which don't get sent over the wire but still accessible
// code wise
type DestinationResource struct {
	TableName string  `json:"table_name"`
	Data      CQTypes `json:"data"`
}

func NewResourceData(t *Table, parent *Resource, item interface{}) *Resource {
	r := Resource{
		Item:   item,
		Parent: parent,
		Table:  t,
		data:   make(CQTypes, len(t.Columns)),
	}
	return &r
}

func (r *Resource) ToDestinationResource() DestinationResource {
	dr := DestinationResource{
		TableName: r.Table.Name,
		Data:      r.data,
	}
	return dr
}

func (r *Resource) Get(columnName string) CQType {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	return r.data[index]
}

// Set sets a column with value. This does validation and conversion to
// one of concrete CQTypes. it returns an error just for backward compatability
// and panics in case it fails
func (r *Resource) Set(columnName string, value interface{}) error {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	var err error
	switch r.Table.Columns[index].Type {
	case TypeBool:
		r.data[index] = &cqtypes.Bool{}
	case TypeInt:
		r.data[index] = &cqtypes.Int8{}
	case TypeFloat:
		r.data[index] = &cqtypes.Float8{}
	case TypeUUID:
		r.data[index] = &cqtypes.UUID{}
	case TypeString:
		r.data[index] = &cqtypes.Text{}
	case TypeByteArray:
		r.data[index] = &cqtypes.Bytea{}
	case TypeStringArray:
		r.data[index] = &cqtypes.TextArray{}
	case TypeIntArray:
		r.data[index] = &cqtypes.Int8Array{}
	case TypeTimestamp:
		r.data[index] = &cqtypes.Timestamptz{}
	case TypeJSON:
		r.data[index] = &cqtypes.JSON{}
	case TypeUUIDArray:
		r.data[index] = &cqtypes.UUIDArray{}
	case TypeInet:
		r.data[index] = &cqtypes.Inet{}
	case TypeInetArray:
		r.data[index] = &cqtypes.InetArray{}
	case TypeCIDR:
		r.data[index] = &cqtypes.CIDR{}
	case TypeCIDRArray:
		r.data[index] = &cqtypes.CIDRArray{}
	case TypeMacAddr:
		r.data[index] = &cqtypes.Macaddr{}
	case TypeMacAddrArray:
		r.data[index] = &cqtypes.MacaddrArray{}
	default:
		panic(fmt.Errorf("unsupported type %s", r.Table.Columns[index].Type.String()))
	}
	err = r.data[index].Set(value)
	if err != nil {
		panic(err)
	}
	return nil
}

// Override original item (this is useful for apis that follow list/details pattern)
func (r *Resource) SetItem(item interface{}) {
	r.Item = item
}

func (r *Resource) GetItem() interface{} {
	return r.Item
}

func (r *Resource) GetValues() CQTypes {
	return r.data
}

func (r *Resource) ID() uuid.UUID {
	index := r.Table.Columns.Index(CqIDColumn.Name)
	if index == -1 {
		return uuid.UUID{}
	}
	return uuid.UUID{}
}

func (r *Resource) Columns() []string {
	return r.Table.Columns.Names()
}

func (rr Resources) GetIds() []uuid.UUID {
	rids := make([]uuid.UUID, len(rr))
	for i, r := range rr {
		rids[i] = r.ID()
	}
	return rids
}
func (rr Resources) TableName() string {
	if len(rr) == 0 {
		return ""
	}
	return rr[0].Table.Name
}

func (rr Resources) ColumnNames() []string {
	if len(rr) == 0 {
		return []string{}
	}
	return rr[0].Table.Columns.Names()
}
