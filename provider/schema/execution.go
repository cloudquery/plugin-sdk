package schema

import (
	"context"

	"github.com/cloudquery/go-funk"
	"github.com/hashicorp/go-hclog"
	"github.com/iancoleman/strcase"
	"golang.org/x/sync/errgroup"
)

type ClientMeta interface {
	Logger() hclog.Logger
}

// ExecutionData marks all the related execution info passed to TableResolver and ColumnResolver giving access to the Runner's meta
type ExecutionData struct {
	// The table this execution is associated with
	Table *Table
	// Database connection to insert data into
	Db Database
	// Logger associated with this execution
	Logger hclog.Logger
	// Column is set if execution is passed to ColumnResolver
	Column *Column
}

// Create a new execution data
func NewExecutionData(db Database, logger hclog.Logger, table *Table) ExecutionData {
	return ExecutionData{
		Table:  table,
		Db:     db,
		Logger: logger,
		Column: nil,
	}
}

func (e ExecutionData) ResolveTable(ctx context.Context, meta ClientMeta, parent *Resource) error {
	var clients []ClientMeta
	clients = append(clients, meta)
	if e.Table.Multiplex != nil {
		clients = e.Table.Multiplex(meta)
	}
	g, ctx := errgroup.WithContext(ctx)
	for _, client := range clients {
		client := client
		g.Go(func() error {
			err := e.callTableResolve(ctx, client, parent)
			if err != nil && !(e.Table.IgnoreError != nil && e.Table.IgnoreError(err)) {
				return err
			}
			return nil
		})
	}
	return g.Wait()
}

func (e ExecutionData) WithTable(t *Table) ExecutionData {
	return ExecutionData{
		Table:  t,
		Db:     e.Db,
		Logger: e.Logger,
	}
}

func (e ExecutionData) callTableResolve(ctx context.Context, client ClientMeta, parent *Resource) error {

	if parent == nil && e.Table.DeleteFilter != nil {
		// Delete previous fetch
		if err := e.Db.Delete(ctx, e.Table, e.Table.DeleteFilter(client)); err != nil {
			return err
		}
	}

	res := make(chan interface{})
	var resolverErr error
	go func() {
		defer close(res)
		resolverErr = e.Table.Resolver(ctx, client, parent, res)
	}()

	nc := 0
	for elem := range res {
		objects := interfaceSlice(elem)
		if len(objects) == 0 {
			continue
		}
		if err := e.resolveResources(ctx, client, parent, objects); err != nil {
			return err
		}
		nc += len(objects)
	}
	// check if channel iteration stopped because of resolver failure
	if resolverErr != nil {
		return resolverErr
	}
	// Print only parent resources
	if parent == nil {
		client.Logger().Info("fetched successfully", "table", e.Table.Name, "count", nc)
	}
	return nil
}

func (e ExecutionData) resolveResources(ctx context.Context, meta ClientMeta, parent *Resource, objects []interface{}) error {
	var resources = make([]*Resource, len(objects))
	for i, o := range objects {
		resources[i] = NewResourceData(e.Table, parent, o)
		if err := e.resolveResourceValues(ctx, meta, resources[i]); err != nil {
			return err
		}
	}

	// Before inserting resolve all table column resolvers
	if err := e.Db.Insert(ctx, e.Table, resources); err != nil {
		e.Logger.Error("failed to insert to db", "error", err)
		return err
	}

	// Finally resolve relations of each resource
	for _, rel := range e.Table.Relations {
		for _, r := range resources {
			err := e.WithTable(rel).ResolveTable(ctx, meta, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e ExecutionData) resolveResourceValues(ctx context.Context, meta ClientMeta, resource *Resource) error {
	if err := e.resolveColumns(ctx, meta, resource, resource.table.Columns); err != nil {
		return err
	}
	// call PostRowResolver if defined after columns have been resolved
	if resource.table.PostResourceResolver == nil {
		return nil
	}
	if err := resource.table.PostResourceResolver(ctx, meta, resource); err != nil {
		return err
	}
	return nil
}

func (e ExecutionData) resolveColumns(ctx context.Context, meta ClientMeta, resource *Resource, cols []Column) error {
	for _, c := range cols {
		if c.Resolver != nil {
			if err := c.Resolver(ctx, meta, resource, c); err != nil {
				return err
			}
			continue
		}
		// base use case: try to get column with CamelCase name
		v := funk.GetAllowZero(resource.Item, strcase.ToCamel(c.Name))
		if v == nil {
			v = c.Default
		}
		resource.Set(c.Name, v)
	}
	return nil
}
