package schema

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	"github.com/iancoleman/strcase"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/errgroup"
)

type ClientMeta interface {
	Logger() hclog.Logger
}

// ExecutionData marks all the related execution info passed to TableResolver and ColumnResolver giving access to the Runner's meta
type ExecutionData struct {
	// Table this execution is associated with
	Table *Table
	// Database connection to insert data into
	Db Database
	// Logger associated with this execution
	Logger hclog.Logger
	// disableDelete allows to disable deletion of table data for this execution
	disableDelete bool
	// extraFields to be passed to each created resource in the execution
	extraFields map[string]interface{}
}

// NewExecutionData Create a new execution data
func NewExecutionData(db Database, logger hclog.Logger, table *Table, disableDelete bool, extraFields map[string]interface{}) ExecutionData {
	return ExecutionData{
		Table:         table,
		Db:            db,
		Logger:        logger,
		disableDelete: disableDelete,
		extraFields:   extraFields,
	}
}

func (e ExecutionData) ResolveTable(ctx context.Context, meta ClientMeta, parent *Resource) (uint64, error) {
	var clients []ClientMeta
	clients = append(clients, meta)
	if e.Table.Multiplex != nil {
		clients = e.Table.Multiplex(meta)
		meta.Logger().Debug("multiplexing client", "count", len(clients))
	}
	g, ctx := errgroup.WithContext(ctx)
	var totalResources uint64
	for _, client := range clients {
		client := client
		g.Go(func() error {
			count, err := e.callTableResolve(ctx, client, parent)
			if err != nil && !(e.Table.IgnoreError != nil && e.Table.IgnoreError(err)) {
				return err
			}
			atomic.AddUint64(&totalResources, count)
			return nil
		})
	}
	return totalResources, g.Wait()
}

func (e ExecutionData) WithTable(t *Table) ExecutionData {
	return ExecutionData{
		Table:         t,
		Db:            e.Db,
		Logger:        e.Logger,
		disableDelete: e.disableDelete,
		extraFields:   e.extraFields,
	}
}

func (e ExecutionData) truncateTable(ctx context.Context, client ClientMeta, parent *Resource) error {
	if e.Table.DeleteFilter == nil {
		return nil
	}
	if e.disableDelete && !e.Table.AlwaysDelete {
		client.Logger().Debug("skipping table truncate", "table", e.Table.Name)
		return nil
	}
	// Delete previous fetch
	if err := e.Db.Delete(ctx, e.Table, e.Table.DeleteFilter(client, parent)); err != nil {
		client.Logger().Debug("cleaning table previous fetch", "table", e.Table.Name, "always_delete", e.Table.AlwaysDelete)
		return err
	}
	return nil
}

func (e ExecutionData) callTableResolve(ctx context.Context, client ClientMeta, parent *Resource) (uint64, error) {

	if err := e.truncateTable(ctx, client, parent); err != nil {
		return 0, err
	}

	res := make(chan interface{})
	var resolverErr error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Fetch task exited with panic:\n%s\n", debug.Stack())
				e.Logger.Error("Fetch task exited with panic", "table", e.Table.Name, "stack", string(debug.Stack()))
			}
			close(res)
		}()

		resolverErr = e.Table.Resolver(ctx, client, parent, res)
	}()

	nc := uint64(0)
	for elem := range res {
		objects := interfaceSlice(elem)
		if len(objects) == 0 {
			continue
		}
		if err := e.resolveResources(ctx, client, parent, objects); err != nil {
			return 0, err
		}
		nc += uint64(len(objects))
	}
	// check if channel iteration stopped because of resolver failure
	if resolverErr != nil {
		client.Logger().Error("received resolve resources error", "table", e.Table.Name, "error", resolverErr)
		return 0, resolverErr
	}
	// Print only parent resources
	if parent == nil {
		client.Logger().Info("fetched successfully", "table", e.Table.Name, "count", nc)
	}
	return nc, nil
}

func (e ExecutionData) resolveResources(ctx context.Context, meta ClientMeta, parent *Resource, objects []interface{}) error {
	var resources = make([]*Resource, len(objects))
	for i, o := range objects {
		resources[i] = NewResourceData(e.Table, parent, o, e.extraFields)
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
		meta.Logger().Debug("resolving table relation", "table", e.Table.Name, "relation", rel.Name)
		for _, r := range resources {
			// ignore relation resource count
			_, err := e.WithTable(rel).ResolveTable(ctx, meta, r)
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
	if resource.table.PostResourceResolver != nil {
		if err := resource.table.PostResourceResolver(ctx, meta, resource); err != nil {
			return err
		}
	}
	// Finally generate cq_id for resource
	for _, c := range GetDefaultSDKColumns() {
		if err := c.Resolver(ctx, meta, resource, c); err != nil {
			return err
		}
	}
	return nil
}

func (e ExecutionData) resolveColumns(ctx context.Context, meta ClientMeta, resource *Resource, cols []Column) error {
	for _, c := range cols {
		if c.Resolver != nil {
			meta.Logger().Trace("using custom column resolver", "column", c.Name)
			if err := c.Resolver(ctx, meta, resource, c); err != nil {
				return err
			}
			continue
		}
		meta.Logger().Trace("resolving column value", "column", c.Name)
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.Item, strcase.ToCamel(c.Name), funk.WithAllowZero())
		if v == nil {
			meta.Logger().Trace("using column default value", "column", c.Name, "default", c.Default)
			v = c.Default
		}
		meta.Logger().Trace("setting column value", "column", c.Name, "value", v)
		if err := resource.Set(c.Name, v); err != nil {
			return err
		}
	}
	return nil
}

func interfaceSlice(slice interface{}) []interface{} {
	// if value is nil return nil
	if slice == nil {
		return nil
	}
	s := reflect.ValueOf(slice)
	// Keep the distinction between nil and empty slice input
	if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Slice && s.Elem().IsNil() {
		return nil
	}
	if s.Kind() != reflect.Slice {
		return []interface{}{slice}
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
