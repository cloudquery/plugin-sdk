package schema

import (
	"github.com/google/uuid"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	item interface{}
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
		item:   item,
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

func (r *Resource) Set(columnName string, value interface{}) {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	var err error
	switch r.Table.Columns[index].Type {
	case TypeBool:
		r.data[index] = &Bool{}
		err = r.data[index].(*Bool).Scan(value)
	case TypeInt:
		r.data[index] = &Int64{}
		err = r.data[index].(*Int64).Scan(value)
	case TypeJSON:
		r.data[index] = &Json{}
		err = r.data[index].(*Json).Scan(value)
	case TypeTimestamp:
		r.data[index] = &Timestamptz{}
		err = r.data[index].(*Timestamptz).Scan(value)
	case TypeString:
		r.data[index] = &String{}
		err = r.data[index].(*String).Scan(value)
	case TypeUUID:
		r.data[index] = &UUID{}
		err = r.data[index].(*UUID).Scan(value)
	case TypeByteArray:
		r.data[index] = &ByteArray{}
		err = r.data[index].(*ByteArray).Scan(value)
	}
	if err != nil {
		panic(err)
	}
}

// Override original item (this is useful for apis that follow list/details pattern)
func (r *Resource) SetItem(item interface{}) {
	r.item = item
}

func (r *Resource) GetItem() interface{} {
	return r.item
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
