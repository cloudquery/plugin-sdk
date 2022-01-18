package schema

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cloudquery/cq-provider-sdk/provider/schema/diag"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/modern-go/reflect2"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"github.com/hashicorp/go-hclog"
	"github.com/iancoleman/strcase"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/errgroup"
)

// executionJitter adds a -1 minute to execution of fetch, so if a user fetches only 1 resources and it finishes
// faster than the <1s it won't be deleted by remove stale.
const executionJitter = -1 * time.Minute

//go:generate mockgen -package=mock -destination=./mock/mock_storage.go . Storage
type Storage interface {
	QueryExecer

	Insert(ctx context.Context, t *Table, instance Resources) error
	Delete(ctx context.Context, t *Table, kvFilters []interface{}) error
	RemoveStaleData(ctx context.Context, t *Table, executionStart time.Time, kvFilters []interface{}) error
	CopyFrom(ctx context.Context, resources Resources, shouldCascade bool, CascadeDeleteFilters map[string]interface{}) error
	Close()
	Dialect() Dialect
}

type QueryExecer interface {
	pgxscan.Querier

	Exec(ctx context.Context, query string, args ...interface{}) error
}

type ClientMeta interface {
	Logger() hclog.Logger
}

// ExecutionData marks all the related execution info passed to TableResolver and ColumnResolver giving access to the Runner's meta
type ExecutionData struct {
	// ResourceName name of top-level resource associated with table
	ResourceName string
	// Table this execution is associated with
	Table *Table
	// Database connection to insert data into
	Db Storage
	// Logger associated with this execution
	Logger hclog.Logger
	// disableDelete allows disabling deletion of table data for this execution
	disableDelete bool
	// extraFields to be passed to each created resource in the execution
	extraFields map[string]interface{}
	// partialFetch if true allows partial fetching of resources
	partialFetch bool
	// PartialFetchFailureResult is a an array of resources where the fetch process failed
	PartialFetchFailureResult []ResourceFetchError
	// partialFetchChan is the channel that is used to send failed resource fetches
	partialFetchChan chan ResourceFetchError
	// When the execution started
	executionStart time.Time
}

// ResourceFetchError represents a single partial fetch failed resource
type ResourceFetchError struct {
	// table name of the failed resource fetch
	TableName string
	// root/parent table name
	RootTableName string
	// root/parent primary key values
	RootPrimaryKeyValues []string
	// error message for this resource fetch failure
	Err error
}

func (p ResourceFetchError) Error() string {
	return p.Err.Error()
}

const (
	// partialFetchFailureBufferLength defines the buffer length for the partialFetchChan.
	partialFetchFailureBufferLength = 10

	// fdLimitMessage defines the message for when a client isn't able to fetch because the open fd limit is hit
	fdLimitMessage = "try increasing number of available file descriptors via `ulimit -n 10240` or by increasing timeout via provider specific parameters"
)

// NewExecutionData Create a new execution data
func NewExecutionData(db Storage, logger hclog.Logger, table *Table, disableDelete bool, extraFields map[string]interface{}, partialFetch bool) ExecutionData {
	return ExecutionData{
		Table:                     table,
		Db:                        db,
		Logger:                    logger,
		disableDelete:             disableDelete,
		extraFields:               extraFields,
		PartialFetchFailureResult: []ResourceFetchError{},
		partialFetch:              partialFetch,
		executionStart:            time.Now().Add(executionJitter),
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
		e.partialFetchChan = make(chan ResourceFetchError, partialFetchFailureBufferLength)
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
		ResourceName:              e.ResourceName,
		Db:                        e.Db,
		Logger:                    e.Logger,
		disableDelete:             e.disableDelete,
		extraFields:               e.extraFields,
		partialFetch:              e.partialFetch,
		PartialFetchFailureResult: []ResourceFetchError{},
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

// cleanupStaleData cleans resources in table that weren't update in the latest table resolve execution
func (e ExecutionData) cleanupStaleData(ctx context.Context, client ClientMeta, parent *Resource) error {
	// Only clean top level tables
	if parent != nil {
		return nil
	}
	if !e.disableDelete {
		client.Logger().Debug("skipping stale data removal", "table", e.Table.Name)
		return nil
	}
	client.Logger().Debug("cleaning table table stale data", "table", e.Table.Name, "last_update", e.executionStart)

	filters := make([]interface{}, 0)
	for k, v := range e.extraFields {
		filters = append(filters, k, v)
	}
	if e.Table.DeleteFilter != nil {
		filters = append(filters, e.Table.DeleteFilter(client, parent)...)
	}
	return e.Db.RemoveStaleData(ctx, e.Table, e.executionStart, filters)
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
			// add partial fetch error, this will be passed in diagnostics, although it was ignored
			_ = e.checkPartialFetchError(err, parent, "table resolver ignored error")
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
	if err := e.cleanupStaleData(ctx, client, parent); err != nil {
		return nc, e.checkPartialFetchError(err, nil, "failed to clean stale data")
	}
	return nc, nil
}

func (e *ExecutionData) resolveResources(ctx context.Context, meta ClientMeta, parent *Resource, objects []interface{}) error {
	var resources = make(Resources, 0, len(objects))
	for _, o := range objects {
		resource := NewResourceData(e.Db.Dialect(), e.Table, parent, o, e.extraFields, e.executionStart)
		// Before inserting resolve all table column resolvers
		if err := e.resolveResourceValues(ctx, meta, resource); err != nil {
			if partialFetchErr := e.checkPartialFetchError(err, resource, "failed to resolve resource"); partialFetchErr != nil {
				return partialFetchErr
			}
			e.Logger.Warn("skipping failed resolved resource", "reason", err.Error())
			continue
		}
		resources = append(resources, resource)
	}

	// only top level tables should cascade, disable delete is turned on.
	// if we didn't disable delete all data should be wiped before resolve
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
	e.Logger.Warn("failed copy-from to db", "error", err, "table", e.Table.Name)

	// fallback insert, copy from sometimes does problems, so we fall back with bulk insert
	err = e.Db.Insert(ctx, e.Table, resources)
	if err == nil {
		return resources, nil
	}
	e.Logger.Error("failed insert to db", "error", err, "table", e.Table.Name)

	// Partial fetch check
	if partialFetchErr := e.checkPartialFetchError(err, nil, "failed to copy resources into the db"); partialFetchErr != nil {
		return nil, partialFetchErr
	}

	// Try to insert resource by resource if partial fetch is enabled and an error occurred
	partialFetchResources := make(Resources, 0)
	for id := range resources {
		if err := e.Db.Insert(ctx, e.Table, Resources{resources[id]}); err != nil {
			e.Logger.Error("failed to insert resource into db", "error", err, "resource_keys", resources[id].Keys(), "table", e.Table.Name)
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
			err = fmt.Errorf("recovered from panic: %s", r)
		}
	}()

	providerCols, internalCols := siftColumns(e.Db.Dialect().Columns(resource.table))

	if err = e.resolveColumns(ctx, meta, resource, providerCols); err != nil {
		return fmt.Errorf("resolve columns error: %w", err)
	}
	// call PostRowResolver if defined after columns have been resolved
	if resource.table.PostResourceResolver != nil {
		if err = resource.table.PostResourceResolver(ctx, meta, resource); err != nil {
			return fmt.Errorf("post resource resolver failed: %w", err)
		}
	}
	// Finally, resolve columns internal to the SDK
	for _, c := range internalCols {
		if err = c.Resolver(ctx, meta, resource, c); err != nil {
			return fmt.Errorf("default column %s resolver execution failed: %w", c.Name, err)
		}
	}
	return err
}

func (e *ExecutionData) resolveColumns(ctx context.Context, meta ClientMeta, resource *Resource, cols []Column) error {
	for _, c := range cols {
		if c.Resolver != nil {
			meta.Logger().Trace("using custom column resolver", "column", c.Name, "table", e.Table.Name)
			err := c.Resolver(ctx, meta, resource, c)
			if err == nil {
				continue
			}
			// Not allowed ignoring PK resolver errors
			if funk.ContainsString(resource.Keys(), c.Name) {
				return err
			}
			// check if column resolver defined an IgnoreError function, if it does check if ignore should be ignored.
			if c.IgnoreError == nil || !c.IgnoreError(err) {
				return fmt.Errorf("column %s: %w", c.Name, err)
			}

			if reflect2.IsNil(c.Default) {
				return err
			}
			// Set default value if defined, otherwise it will be nil
			if err := resource.Set(c.Name, c.Default); err != nil {
				return err
			}
			continue
		}
		meta.Logger().Trace("resolving column value with path", "column", c.Name, "table", e.Table.Name)
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

	partialFetchFailure := ResourceFetchError{
		TableName: e.Table.Name,
		Err:       fmt.Errorf("%s: %w", customMsg, err),
	}
	e.Logger.Debug("fetch error occurred and partial fetch is enabled", "msg", partialFetchFailure.Error(), "table", e.Table.Name)

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
			partialFetchFailure.RootPrimaryKeyValues = root.Keys()
		}
	}
	if strings.Contains(err.Error(), ": socket: too many open files") {
		// Return a Diagnostic error so that it can be properly propegated back to the user via the CLI
		partialFetchFailure.Err = diag.FromError(err, diag.WARNING, diag.THROTTLE, e.ResourceName, err.Error(), fdLimitMessage)
	}
	// Send information via our channel
	e.partialFetchChan <- partialFetchFailure
	return nil
}

// siftColumns gets a column list and returns a list of provider columns, and another list of internal columns, cqId column being the very last one
func siftColumns(cols []Column) ([]Column, []Column) {
	providerCols, internalCols := make([]Column, 0, len(cols)), make([]Column, 0, len(cols))

	cqIdColIndex := -1
	for i := range cols {
		if cols[i].internal {
			if cols[i].Name == cqIdColumn.Name {
				cqIdColIndex = len(internalCols)
			}

			internalCols = append(internalCols, cols[i])
		} else {
			providerCols = append(providerCols, cols[i])
		}
	}

	// resolve cqId last, as it would need other PKs to be resolved, some might be internal (cq_fetch_date)
	if lastIndex := len(internalCols) - 1; cqIdColIndex > -1 && cqIdColIndex != lastIndex {
		internalCols[cqIdColIndex], internalCols[lastIndex] = internalCols[lastIndex], internalCols[cqIdColIndex]
	}

	return providerCols, internalCols
}
