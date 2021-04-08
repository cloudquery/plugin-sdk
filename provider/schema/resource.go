package schema

import (
	"github.com/google/uuid"
)

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item interface{}
	// Set if this is an embedded table
	Parent *Resource
	// internal fields
	table *Table
	data  map[string]interface{}
	id    uuid.UUID
}

func NewResourceData(t *Table, parent *Resource, item interface{}) *Resource {
	return &Resource{
		Item:   item,
		Parent: parent,
		table:  t,
		data:   make(map[string]interface{}),
		id:     uuid.New(),
	}
}

func (r Resource) Get(key string) interface{} {
	return r.data[key]
}

func (r Resource) Set(key string, value interface{}) {
	r.data[key] = value
}

func (r Resource) Id() uuid.UUID {
	return r.id
}

func (r Resource) Values() ([]interface{}, error) {
	values := make([]interface{}, 0)
	values = append(values, r.id)
	for _, c := range r.table.Columns {
		v := r.Get(c.Name)
		if err := c.ValidateType(v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}
