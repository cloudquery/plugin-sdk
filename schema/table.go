package schema

import (
	"context"
	"runtime/debug"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/helpers"
	"github.com/iancoleman/strcase"
	"github.com/thoas/go-funk"
)

// TableResolver is the main entry point when a table fetch is called.
//
// Table resolver has 3 main arguments:
// - meta(ClientMeta): is the client returned by the plugin.Provider Configure call
// - parent(Resource): resource is the parent resource in case this table is called via parent table (i.e. relation)
// - res(chan interface{}): is a channel to pass results fetched by the TableResolver
//
type TableResolver func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error

// IgnoreErrorFunc checks if returned error from table resolver should be ignored.
type IgnoreErrorFunc func(err error) bool

type RowResolver func(ctx context.Context, meta ClientMeta, resource *Resource) error

type Table struct {
	// Name of table
	Name string
	// table description
	Description string
	// Columns are the set of fields that are part of this table
	Columns ColumnList
	// Relations are a set of related tables defines
	Relations []*Table
	// Resolver is the main entry point to fetching table data and
	Resolver TableResolver `msgpack:"-"`
	// Ignore errors checks if returned error from table resolver should be ignored.
	IgnoreError IgnoreErrorFunc `msgpack:"-"`
	// Multiplex returns re-purposed meta clients. The sdk will execute the table with each of them
	Multiplex func(meta ClientMeta) []ClientMeta `msgpack:"-"`
	// Post resource resolver is called after all columns have been resolved, and before resource is inserted to database.
	PostResourceResolver RowResolver `msgpack:"-"`
	// Options allow modification of how the table is defined when created
	Options TableCreationOptions

	// IgnoreInTests is used to exclude a table from integration tests.
	// By default, integration tests fetch all resources from cloudquery's test account, and verify all tables
	// have at least one row.
	// When IgnoreInTests is true, integration tests won't fetch from this table.
	// Used when it is hard to create a reproducible environment with a row in this table.
	IgnoreInTests bool

	// Parent is the parent table in case this table is called via parent table (i.e. relation)
	Parent *Table

	// Serial is used to force a signature change, which forces new table creation and cascading removal of old table and relations
	Serial string
}

// TableCreationOptions allow modifying how table is created such as defining primary keys, indices, foreign keys and constraints.
type TableCreationOptions struct {
	// List of columns to set as primary keys. If this is empty, a random unique ID is generated.
	PrimaryKeys []string
}

func (t Table) Column(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func (tco TableCreationOptions) signature() string {
	return strings.Join(tco.PrimaryKeys, ";")
}

func (t Table) TableNames() []string {
	ret := []string{t.Name}
	for _, rel := range t.Relations {
		ret = append(ret, rel.TableNames()...)
	}
	return ret
}

// Call the table resolver with with all of it's relation for every reolved resource
func (t Table) Resolve(ctx context.Context, meta ClientMeta, parent *Resource, resolvedResources chan<- *Resource) {
	res := make(chan interface{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				meta.Logger().Error().Str("table_name", t.Name).Str("stack", stack).Msg("table resolver finished with panic")
			}
			close(res)
		}()
		meta.Logger().Debug().Str("table_name", t.Name).Msg("table resolver started")
		if err := t.Resolver(ctx, meta, parent, res); err != nil {
			meta.Logger().Error().Str("table_name", t.Name).Err(err).Msg("table resolver finished with error")
		}
		meta.Logger().Debug().Str("table_name", t.Name).Msg("table resolver finished successfully")
	}()

	// each result is an array of interface{}
	for elem := range res {
		objects := helpers.InterfaceSlice(elem)
		if len(objects) == 0 {
			continue
		}
		for i := range objects {
			resource := NewResourceData(&t, parent, objects[i])
			t.resolveColumns(ctx, meta, resource)
			if t.PostResourceResolver != nil {
				meta.Logger().Trace().Str("table_name", t.Name).Msg("post resource resolver started")
				if err := t.PostResourceResolver(ctx, meta, resource); err != nil {
					meta.Logger().Error().Str("table_name", t.Name).Err(err).Msg("post resource resolver finished with error")
				}
				meta.Logger().Trace().Str("table_name", t.Name).Msg("post resource resolver finished successfully")
			}
			resolvedResources <- resource
			for _, rel := range t.Relations {
				rel.Resolve(ctx, meta, resource, resolvedResources)
			}
		}
	}
}

func (t Table) resolveColumns(ctx context.Context, meta ClientMeta, resource *Resource) {
	for _, c := range t.Columns {
		if c.Resolver != nil {
			meta.Logger().Trace().Str("colum_name", c.Name).Str("table_name", t.Name).Msg("column resolver custom started")
			if err := c.Resolver(ctx, meta, resource, c); err != nil {
				meta.Logger().Error().Str("colum_name", c.Name).Str("table_name", t.Name).Err(err).Msg("column resolver finished with error")
			}
			meta.Logger().Trace().Str("colum_name", c.Name).Str("table_name", t.Name).Msg("column resolver finished successfully")
		} else {
			meta.Logger().Trace().Str("colum_name", c.Name).Str("table_name", t.Name).Msg("column resolver default started")
			// base use case: try to get column with CamelCase name
			v := funk.Get(resource.Item, strcase.ToCamel(c.Name), funk.WithAllowZero())
			if err := resource.Set(c.Name, v); err != nil {
				meta.Logger().Error().Str("colum_name", c.Name).Str("table_name", t.Name).Err(err).Msg("column resolver default finished with error")
			}
			meta.Logger().Trace().Str("colum_name", c.Name).Str("table_name", t.Name).Msg("column resolver default finished successfully")
		}
	}
}
