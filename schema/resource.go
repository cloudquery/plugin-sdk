package schema

import (
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/google/uuid"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item any
	// Set if this is an embedded table
	Parent *Resource
	// internal fields
	Table *Table
	// This is sorted result data by column name
	data map[string]any
	bldr array.RecordBuilder
}

// This struct is what we send over the wire to destination.
// We dont want to reuse the same struct as otherwise we will have to comment on fields which don't get sent over the wire but still accessible
// code wise
// type DestinationResource struct {
// 	TableName string  `json:"table_name"`
// 	Data      CQTypes `json:"data"`
// }

func NewResourceData(t *Table, parent *Resource, item any) *Resource {
	r := Resource{
		Item:   item,
		Parent: parent,
		Table:  t,
		data:   make(map[string]any, len(t.Columns)),
	}
	for _, c := range t.Columns {
		r.data[c.Name] = nil
	}
	return &r
}

func (r *Resource) Get(columnName string) any {
	return r.data[columnName]
}

// Set sets a column with value. This does validation and conversion to
// one of concrete  it returns an error just for backward compatibility
// and panics in case it fails
func (r *Resource) Set(columnName string, value any) error {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	r.data[columnName] = value
	return nil
}

// Override original item (this is useful for apis that follow list/details pattern)
func (r *Resource) SetItem(item any) {
	r.Item = item
}

func (r *Resource) GetItem() any {
	return r.Item
}

//nolint:revive
func (r *Resource) CalculateCQID(deterministicCQID bool) error {
	panic("not implemented")
}

func (r *Resource) storeCQID(value uuid.UUID) error {
	panic("not implemented")
}

// Validates that all primary keys have values.
func (r *Resource) Validate() error {
	panic("not implemented")
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
