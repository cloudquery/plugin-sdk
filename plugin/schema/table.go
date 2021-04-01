package schema

import (
	"context"
)

// TableResolver is the main entry point when a table fetch is called.
//
// Table resolver has 3 main arguments:
// - meta(ClientMeta): is the client returned by the plugin.Provider Configure call
// - parent(Resource): resource is the parent resource in case this table is called via parent table (i.e relation)
// - res(chan interface{}): is a channel to pass results fetched by the TableResolver
//
type TableResolver func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error

// IgnoreErrorFunc checks if returned error from table resolver should be ignored.
type IgnoreErrorFunc func(err error) bool

type RowResolver func(ctx context.Context, meta ClientMeta, resource *Resource) error

type Table struct {
	// Name of table
	Name string
	// table description
	Description string
	// Columns are the set of fields that are part of this table
	Columns []Column
	// Relations are a set of related tables defines
	Relations []*Table
	// Resolver is the main entry point to fetching table data and
	Resolver TableResolver
	// Ignore errors checks if returned error from table resolver should be ignored.
	IgnoreError IgnoreErrorFunc
	// Multiplex returns re-purposed meta clients. The sdk will execute the table with each of them
	Multiplex func(meta ClientMeta) []ClientMeta
	// DeleteFilter returns a list of key/value pairs to add when truncating this table's data from the database.
	DeleteFilter func(meta ClientMeta) []interface{}
	// Post resource resolver is called after all columns have been resolved, and before resource is inserted to database.
	PostResourceResolver RowResolver
}

// ColumnNames returns all collected columns name of table (including all inner embedded columns)
func (t Table) ColumnNames() []string {
	var cn = make([]string, len(t.Columns))
	for i, c := range t.Columns {
		cn[i] = c.Name
	}
	return append([]string{"id"}, cn...)
}

func (t Table) Column(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
