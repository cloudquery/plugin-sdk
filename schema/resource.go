package schema

import (
	"fmt"

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
	data CQTypes
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
	for i := range r.data {
		r.data[i] = NewCqTypeFromValueType(t.Columns[i].Type)
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
// one of concrete  it returns an error just for backward compatibility
// and panics in case it fails
func (r *Resource) Set(columnName string, value interface{}) error {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	if err := r.data[index].Set(value); err != nil {
		panic(fmt.Errorf("failed to set column %s: %w", columnName, err))
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
