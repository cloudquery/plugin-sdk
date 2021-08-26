package schema

import (
	"context"
	"fmt"
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
	// partialFetch if true allows partial fetching of resources
	partialFetch bool
	// PartialFetchFailureResult is a map of resources where the fetch process failed
	PartialFetchFailureResult []PartialFetchFailedResource
	// partialFetchChan is the channel that is used to send failed resource fetches
	partialFetchChan chan PartialFetchFailedResource
}

// PartialFetchFailedResource represents a single partial fetch failed resource
type PartialFetchFailedResource struct {
	// table name of the failed resource fetch
	TableName string
	// root/parent table name
	RootTableName string
	// root/parent primary key values
	RootPrimaryKeyValues []string
	// error message for this resource fetch failure
	Error string
}

// partialFetchFailureBufferLength defines the buffer length for the partialFetchChan.
const partialFetchFailureBufferLength = 10

// NewExecutionData Create a new execution data
func NewExecutionData(db Database, logger hclog.Logger, table *Table, disableDelete bool, extraFields map[string]interface{}, partialFetch bool) ExecutionData {
	return ExecutionData{
		Table:                     table,
		Db:                        db,
		Logger:                    logger,
		disableDelete:             disableDelete,
		extraFields:               extraFields,
		PartialFetchFailureResult: []PartialFetchFailedResource{},
		partialFetch:              partialFetch,
	}
}

func (e *ExecutionData) ResolveTable(ctx context.Context, meta ClientMeta, parent *Resource) (uint64, error) {
	var clients []ClientMeta
	clients = append(clients, meta)
	if e.Table.Multiplex != nil {
		if parent != nil {
			meta.Logger().Warn("relation client multiplexing is not allowed, skipping multiplex", "table", e.Table.Name)
		} else {
			clients = e.Table.Multiplex(meta)
			meta.Logger().Debug("multiplexing client", "count", len(clients), "table", e.Table.Name)
		}
	}
	g, ctx := errgroup.WithContext(ctx)
	// Start the partial fetch failure result channel routine
	finishedPartialFetchChan := make(chan bool)
	if e.partialFetch {
		e.partialFetchChan = make(chan PartialFetchFailedResource, partialFetchFailureBufferLength)
		go func() {
			for fetchResourceFailure := range e.partialFetchChan {
				meta.Logger().Debug("received failed partial fetch resource", "resource", fetchResourceFailure, "table", e.Table.Name)
				e.PartialFetchFailureResult = append(e.PartialFetchFailureResult, fetchResourceFailure)
			}
			finishedPartialFetchChan <- true
		}()
	}
	var totalResources uint64
	for _, client := range clients {
		client := client
		g.Go(func() error {
			count, err := e.callTableResolve(ctx, client, parent)
			atomic.AddUint64(&totalResources, count)
			return err
		})
	}
	err := g.Wait()
	if e.partialFetch {
		close(e.partialFetchChan)
		<-finishedPartialFetchChan
	}
	return totalResources, err
}

func (e *ExecutionData) WithTable(t *Table) *ExecutionData {
	return &ExecutionData{
		Table:                     t,
		Db:                        e.Db,
		Logger:                    e.Logger,
		disableDelete:             e.disableDelete,
		extraFields:               e.extraFields,
		partialFetch:              e.partialFetch,
		PartialFetchFailureResult: []PartialFetchFailedResource{},
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
	client.Logger().Debug("cleaning table previous fetch", "table", e.Table.Name, "always_delete", e.Table.AlwaysDelete)
	if err := e.Db.Delete(ctx, e.Table, e.Table.DeleteFilter(client, parent)); err != nil {
		return err
	}
	return nil
}

func (e ExecutionData) callTableResolve(ctx context.Context, client ClientMeta, parent *Resource) (uint64, error) {

	if e.Table.Resolver == nil {
		return 0, fmt.Errorf("table %s missing resolver, make sure table implements the resolver", e.Table.Name)
	}
	if err := e.truncateTable(ctx, client, parent); err != nil {
		return 0, err
	}

	res := make(chan interface{})
	var resolverErr error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				client.Logger().Error("table resolver recovered from panic", "table", e.Table.Name, "stack", string(debug.Stack()))
				resolverErr = fmt.Errorf("failed table %s fetch. Error: %s", e.Table.Name, r)
			}
			close(res)
		}()
		err := e.Table.Resolver(ctx, client, parent, res)
		if err != nil && e.Table.IgnoreError != nil && e.Table.IgnoreError(err) {
			client.Logger().Warn("ignored an error", "err", err, "table", e.Table.Name)
			return
		}
		resolverErr = err
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
		return 0, e.checkPartialFetchError(resolverErr, nil, "table resolve error")
	}
	// Print only parent resources
	if parent == nil {
		client.Logger().Info("fetched successfully", "table", e.Table.Name, "count", nc)
	}
	return nc, nil
}

func (e *ExecutionData) resolveResources(ctx context.Context, meta ClientMeta, parent *Resource, objects []interface{}) error {
	var resources = make(Resources, len(objects))
	for i, o := range objects {
		resources[i] = NewResourceData(e.Table, parent, o, e.extraFields)
		// Before inserting resolve all table column resolvers
		if err := e.resolveResourceValues(ctx, meta, resources[i]); err != nil {
			e.Logger.Error("failed to resolve resource values", "error", err)
			return err
		}
	}

	// only top level tables should cascade, disable delete is turned on.
	// if we didn't disable delete all data should be wiped before resolve)
	shouldCascade := parent == nil && e.disableDelete
	var err error
	resources, err = e.copyDataIntoDB(ctx, resources, shouldCascade)
	if err != nil {
		return err
	}

	// Finally, resolve relations of each resource
	for _, rel := range e.Table.Relations {
		meta.Logger().Debug("resolving table relation", "table", e.Table.Name, "relation", rel.Name)
		for _, r := range resources {
			// ignore relation resource count
			_, err := e.WithTable(rel).ResolveTable(ctx, meta, r)
			if err != nil {
				if partialFetchErr := e.checkPartialFetchError(err, r, "resolve relation error"); partialFetchErr != nil {
					return partialFetchErr
				}
			}
		}
	}
	return nil
}

func (e *ExecutionData) copyDataIntoDB(ctx context.Context, resources Resources, shouldCascade bool) (Resources, error) {
	err := e.Db.CopyFrom(ctx, resources, shouldCascade, e.extraFields)
	if err == nil {
		return resources, nil
	}
	e.Logger.Warn("failed copy-from to db", "error", err)

	// fallback insert, copy from sometimes does problems so we fall back with insert
	err = e.Db.Insert(ctx, e.Table, resources)
	if err == nil {
		return resources, nil
	}
	e.Logger.Error("failed insert to db", "error", err)

	// Partial fetch check
	if partialFetchErr := e.checkPartialFetchError(err, nil, "failed to copy resources into the db"); partialFetchErr != nil {
		return nil, partialFetchErr
	}

	// Try to insert resource by resource if partial fetch is enabled and an error occurred
	partialFetchResources := make(Resources, 0)
	for id := range resources {
		if err := e.Db.Insert(ctx, e.Table, Resources{resources[id]}); err != nil {
			e.Logger.Error("failed to insert resource into db", "error", err, "resource", resources[id])
		} else {
			// If there is no error we add the resource to the final result
			partialFetchResources = append(partialFetchResources, resources[id])
		}
	}
	return partialFetchResources, nil
}

func (e *ExecutionData) resolveResourceValues(ctx context.Context, meta ClientMeta, resource *Resource) (err error) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger.Error("resolve resource recovered from panic", "table", e.Table.Name, "stack", string(debug.Stack()))
			if partialFetchErr := e.checkPartialFetchError(fmt.Errorf("failed resolve resource. Error: %s", r), resource, "resolve resource recovered from panic"); partialFetchErr != nil {
				err = partialFetchErr
			}
		}
	}()
	if err = e.resolveColumns(ctx, meta, resource, resource.table.Columns); err != nil {
		if partialFetchErr := e.checkPartialFetchError(err, resource, "resolve column error"); partialFetchErr != nil {
			return partialFetchErr
		}
	}
	// call PostRowResolver if defined after columns have been resolved
	if resource.table.PostResourceResolver != nil {
		if err = resource.table.PostResourceResolver(ctx, meta, resource); err != nil {
			if partialFetchErr := e.checkPartialFetchError(err, resource, "post resource resolver failed"); partialFetchErr != nil {
				return partialFetchErr
			}
		}
	}
	// Finally generate cq_id for resource
	for _, c := range GetDefaultSDKColumns() {
		if err = c.Resolver(ctx, meta, resource, c); err != nil {
			if partialFetchErr := e.checkPartialFetchError(err, resource, "column resolver execution failed"); partialFetchErr != nil {
				return partialFetchErr
			}
		}
	}
	return err
}

func (e *ExecutionData) resolveColumns(ctx context.Context, meta ClientMeta, resource *Resource, cols []Column) error {
	for _, c := range cols {
		if c.Resolver != nil {
			meta.Logger().Trace("using custom column resolver", "column", c.Name, "table", e.Table.Name)
			if err := c.Resolver(ctx, meta, resource, c); err != nil {
				return err
			}
			continue
		}
		meta.Logger().Trace("resolving column value", "column", c.Name, "table", e.Table.Name)
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.Item, strcase.ToCamel(c.Name), funk.WithAllowZero())
		if v == nil {
			meta.Logger().Trace("using column default value", "column", c.Name, "default", c.Default, "table", e.Table.Name)
			v = c.Default
		}
		meta.Logger().Trace("setting column value", "column", c.Name, "value", v, "table", e.Table.Name)
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

func (e *ExecutionData) checkPartialFetchError(err error, res *Resource, customMsg string) error {
	// Fast path if partial fetch is disabled
	if !e.partialFetch {
		return err
	}

	partialFetchFailure := PartialFetchFailedResource{
		Error: fmt.Sprintf("%s: %s", customMsg, err.Error()),
	}
	e.Logger.Debug("fetch error occurred and partial fetch is enabled", "msg", partialFetchFailure.Error, "table", e.Table.Name)

	// If resource is given
	if res != nil {
		partialFetchFailure.TableName = res.table.Name

		// Find root/parent resource if one exists
		var root *Resource
		currRes := res
		for root == nil {
			root = currRes
			if currRes != nil && currRes.Parent != nil {
				currRes = res.Parent
				root = nil
			}
		}

		if root != res {
			partialFetchFailure.RootTableName = root.table.Name
			partialFetchFailure.RootPrimaryKeyValues = getPrimaryKeyValues(root)
		}
	}

	// Send information via our channel
	e.partialFetchChan <- partialFetchFailure

	return nil
}

func getPrimaryKeyValues(res *Resource) []string {
	tablePrimKeys := res.table.Options.PrimaryKeys
	if len(tablePrimKeys) == 0 {
		return []string{}
	}
	results := make([]string, len(tablePrimKeys))
	for _, primKey := range tablePrimKeys {
		data := res.Get(primKey)
		if data != nil {
			results = append(results, fmt.Sprintf("%v", data))
		}
	}
	return results
}
