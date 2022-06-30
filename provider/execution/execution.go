package execution

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/cq-provider-sdk/helpers"
	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/cq-provider-sdk/stats"
	"github.com/hashicorp/go-hclog"
	"github.com/iancoleman/strcase"
	segmentStats "github.com/segmentio/stats/v4"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/semaphore"
)

// executionJitter adds a -1 minute to execution of fetch, so if a user fetches only 1 resources and it finishes
// faster than the <1s it won't be deleted by remove stale.
const executionJitter = -1 * time.Minute

// TableExecutor marks all the related execution info passed to TableResolver and ColumnResolver giving access to the Runner's meta
type TableExecutor struct {
	// ResourceName name of top-level resource associated with table
	ResourceName string
	// ParentExecutor is the parent executor, useful for nested tables to propagate up and use IgnoreError and so forth.
	ParentExecutor *TableExecutor
	// Table this execution is associated with
	Table *schema.Table
	// Database connection to insert data into
	Db Storage
	// Logger associated with this execution
	Logger hclog.Logger
	// classifiers
	classifiers []ErrorClassifier
	// metadata to be passed to each created resource in the execution, used by cq* resolvers.
	metadata map[string]interface{}
	// When the execution started
	executionStart time.Time
	// columns of table, this is to reduce calls to sift each time
	columns [2]schema.ColumnList
	// goroutinesSem to limit number of goroutines (clients fetched) concurrently
	goroutinesSem *semaphore.Weighted
	// timeout for each parent resource resolve call
	timeout time.Duration
}

// NewTableExecutor creates a new TableExecutor for given schema.Table
func NewTableExecutor(resourceName string, db Storage, logger hclog.Logger, table *schema.Table, metadata map[string]interface{}, classifier ErrorClassifier, goroutinesSem *semaphore.Weighted, timeout time.Duration) TableExecutor {
	var classifiers = []ErrorClassifier{defaultErrorClassifier}
	if classifier != nil {
		classifiers = append([]ErrorClassifier{classifier}, classifiers...)
	}
	var c [2]schema.ColumnList
	c[0], c[1] = db.Dialect().Columns(table).Sift()

	return TableExecutor{
		ResourceName:   resourceName,
		Table:          table,
		Db:             db,
		Logger:         logger,
		metadata:       metadata,
		classifiers:    classifiers,
		executionStart: time.Now().Add(executionJitter),
		columns:        c,
		goroutinesSem:  goroutinesSem,
		timeout:        timeout,
	}
}

// Resolve is the root function of table executor which starts an execution of a Table resolving it, and it's relations.
func (e TableExecutor) Resolve(ctx context.Context, meta schema.ClientMeta) (uint64, diag.Diagnostics) {
	var clients []schema.ClientMeta

	clients = append(clients, meta)

	if e.Table.Multiplex != nil {
		clients = e.Table.Multiplex(meta)
	}

	return e.doMultiplexResolve(ctx, clients)
}

// withTable allows to create a new TableExecutor for received *schema.Table
func (e TableExecutor) withTable(t *schema.Table, kv ...interface{}) *TableExecutor {
	var c [2]schema.ColumnList
	c[0], c[1] = e.Db.Dialect().Columns(t).Sift()
	cpy := e
	cpy.ParentExecutor = &e
	cpy.Table = t
	cpy.Logger = cpy.Logger.With(kv...)
	cpy.columns = c

	return &cpy
}

func (e TableExecutor) withLogger(kv ...interface{}) *TableExecutor {
	cpy := e
	cpy.Logger = cpy.Logger.With(kv...)
	return &cpy
}

// doMultiplexResolve resolves table with multiplexed clients appending all diagnostics returned from each multiplex.
func (e TableExecutor) doMultiplexResolve(ctx context.Context, clients []schema.ClientMeta) (uint64, diag.Diagnostics) {
	var (
		diagsChan       = make(chan diag.Diagnostics)
		totalResources  uint64
		allDiags        diag.Diagnostics
		doneClients     = 0
		numberOfClients = 0
	)
	// initially use client logger here
	e.Logger.Debug("multiplexing client", "count", len(clients))

	done := make(chan struct{})
	go func() {
		defer close(done)
		for dd := range diagsChan {
			allDiags = allDiags.Add(dd)
			doneClients++
		}
		e.Logger.Debug("multiplexed client finished", "done", doneClients, "total", numberOfClients)
	}()

	wg := &sync.WaitGroup{}
	for _, client := range clients {
		clientID := identifyClient(client)
		if clientID == "" {
			clientID = strconv.Itoa(numberOfClients + 1)
		}
		clientID = e.Table.Name + ":" + clientID

		// we can only limit on a granularity of a top table otherwise we can get deadlock
		e.Logger.Debug("trying acquire for new client", "next_id", clientID)
		if err := e.goroutinesSem.Acquire(ctx, 1); err != nil {
			diagsChan <- ClassifyError(err, diag.WithResourceName(e.ResourceName))
			break
		}
		numberOfClients++
		e.Logger.Debug("creating new multiplex client", "client_id", clientID)
		wg.Add(1)
		go func(c schema.ClientMeta, diags chan<- diag.Diagnostics, id string) {
			defer e.goroutinesSem.Release(1)
			defer wg.Done()
			tableCtx := ctx
			if e.timeout > 0 {
				ctx, cancel := context.WithTimeout(ctx, e.timeout)
				tableCtx = ctx
				defer cancel()
			}
			defer e.Logger.Debug("releasing multiplex client", "ctx_err", ctx.Err())
			// create client execution add all Client's implied Args to execution logger + add its unique client id, so all its execution can be
			// identified.
			count, resolveDiags := e.withLogger(append(c.Logger().ImpliedArgs(), "client_id", id)...).callTableResolve(tableCtx, c, nil)
			atomic.AddUint64(&totalResources, count)
			diags <- resolveDiags
		}(client, diagsChan, clientID)
	}
	wg.Wait()
	close(diagsChan)
	<-done

	e.Logger.Debug("table multiplex resolve completed")
	return totalResources, allDiags
}

// cleanupStaleData cleans resources in table that weren't update in the latest table resolve execution
func (e TableExecutor) cleanupStaleData(ctx context.Context, client schema.ClientMeta, parent *schema.Resource) error {
	// Only clean top level tables
	if parent != nil {
		return nil
	}
	e.Logger.Debug("cleaning table stale data", "last_update", e.executionStart)

	var filters []interface{}
	if e.Table.DeleteFilter != nil {
		filters = append(filters, e.Table.DeleteFilter(client, parent)...)
	}
	if err := e.Db.RemoveStaleData(ctx, e.Table, e.executionStart, filters); err != nil {
		e.Logger.Warn("failed to clean table stale data", "last_update", e.executionStart, "err", err)
		return err
	}
	e.Logger.Debug("cleaned table stale data successfully", "last_update", e.executionStart)
	return nil
}

// callTableResolve does the actual resolving of the table calling the root table's resolver and for each returned resource resolves its columns and relations.
func (e TableExecutor) callTableResolve(ctx context.Context, client schema.ClientMeta, parent *schema.Resource) (uint64, diag.Diagnostics) {
	clock := stats.NewClockWithObserve("callTableResolve", segmentStats.Tag{Name: "client_id", Value: identifyClient(client)}, segmentStats.Tag{Name: "table", Value: e.Table.Name})
	defer clock.Stop()

	// set up all diagnostics to collect from resolving table
	var diags diag.Diagnostics

	if e.Table.Resolver == nil {
		return 0, diags.Add(diag.NewBaseError(nil, diag.SCHEMA, diag.WithSeverity(diag.ERROR), diag.WithResourceName(e.ResourceName), diag.WithSummary("table %q missing resolver, make sure table implements the resolver", e.Table.Name)))
	}

	res := make(chan interface{})
	var resolverErr error

	// we are not using goroutinesSem semaphore here as it's just a +1 goroutine and it might get us deadlocked
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				e.Logger.Error("table resolver recovered from panic", "stack", stack)
				resolverErr = diag.NewBaseError(fmt.Errorf("table resolver panic: %s", r), diag.RESOLVING, diag.WithResourceName(e.ResourceName), diag.WithSeverity(diag.PANIC),
					diag.WithSummary("panic on resource table %q fetch", e.Table.Name), diag.WithDetails("%s", stack))
			}
			close(res)
		}()
		if err := e.Table.Resolver(ctx, client, parent, res); err != nil {
			if e.IgnoreError(err) {
				e.Logger.Debug("ignored an error", "err", err)
				err = diag.NewBaseError(err, diag.RESOLVING, diag.WithSeverity(diag.IGNORE), diag.WithSummary("table %q resolver ignored error", e.Table.Name))
			}
			resolverErr = e.handleResolveError(client, parent, err)
		}
	}()

	nc := uint64(0)
	for elem := range res {
		objects := helpers.InterfaceSlice(elem)
		if len(objects) == 0 {
			continue
		}
		e.Logger.Debug("received resources from resolver", "count", len(objects))
		resolvedCount, dd := e.resolveResources(ctx, client, parent, objects)
		e.Logger.Debug("resolved resources", "original_count", len(objects), "resolved_count", resolvedCount)
		// append any diags from resolve resources
		diags = diags.Add(dd)
		nc += resolvedCount
	}
	// check if channel iteration stopped because of resolver failure
	if resolverErr != nil {
		diags = diags.Add(resolverErr)

		if diag.FromError(resolverErr, diag.INTERNAL).HasErrors() {
			e.Logger.Error("received resolve resources error", "error", resolverErr)
			return 0, diags
		}
	}
	// Print only parent resources
	if parent == nil {
		e.Logger.Info("fetched successfully", "count", nc)
	}

	if err := e.cleanupStaleData(ctx, client, parent); err != nil {
		return nc, diags.Add(ClassifyError(err, diag.WithType(diag.DATABASE), diag.WithSummary("failed to cleanup stale data on table %q", e.Table.Name)))
	}

	return nc, diags
}

// resolveResources resolves a list of resource objects inserting them into the database and resolving their relations based on the table.
func (e TableExecutor) resolveResources(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, objects []interface{}) (uint64, diag.Diagnostics) {
	var (
		resources = make(schema.Resources, 0, len(objects))
		diags     diag.Diagnostics
	)

	for i := range objects {
		resource := schema.NewResourceData(e.Db.Dialect(), e.Table, parent, objects[i], e.metadata, e.executionStart)
		// Before inserting resolve all table column resolvers
		resolveDiags := e.resolveResourceValues(ctx, meta, resource)
		diags = diags.Add(resolveDiags)
		if resolveDiags.HasErrors() {
			e.Logger.Warn("skipping failed resolved resource", "reason", resolveDiags.Error())
			continue
		}
		resources = append(resources, resource)
	}

	// only top level tables should cascade
	shouldCascade := parent == nil
	resources, dbDiags := e.saveToStorage(ctx, resources, shouldCascade)
	e.Logger.Debug("saved resources to storage", "resources", len(resources))
	diags = diags.Add(dbDiags)
	totalCount := uint64(len(resources))

	// Finally, resolve relations of each resource
	for _, rel := range e.Table.Relations {
		e.Logger.Debug("resolving table relation", "relation", rel.Name)
		for _, r := range resources {
			// ignore relation resource count
			if _, innerDiags := e.withTable(rel).callTableResolve(ctx, meta, r); innerDiags.HasDiags() {
				diags = diags.Add(innerDiags)
			}
		}
		e.Logger.Debug("finished resolving table relation", "relation", rel.Name)
	}
	return totalCount, diags
}

// saveToStorage copies resource data to source, it has ways of inserting, first it tries the most performant CopyFrom if that does work it bulk inserts,
// finally it inserts each resource separately, appending errors for each failed resource, only successfully inserted resources are returned
func (e TableExecutor) saveToStorage(ctx context.Context, resources schema.Resources, shouldCascade bool) (schema.Resources, diag.Diagnostics) {
	var diags diag.Diagnostics
	if l := len(resources); l > 0 {
		e.Logger.Debug("storing resources", "count", l)
	}
	err := e.Db.CopyFrom(ctx, resources, shouldCascade)
	if err == nil {
		return resources, diags
	}
	e.Logger.Warn("failed copy-from to db", "error", err)
	diags = diags.Add(diag.TelemetryFromError(err, diag.CopyFromFailed))

	// fallback insert, copy from sometimes does problems, so we fall back with bulk insert
	err = e.Db.Insert(ctx, e.Table, resources, shouldCascade)
	if err == nil {
		return resources, diags
	}
	e.Logger.Error("failed insert to db", "error", err)
	diags = diags.Add(diag.TelemetryFromError(err, diag.BulkInsertFailed))
	// Setup diags, adding first diagnostic that bulk insert failed
	diags = diags.Add(ClassifyError(err, diag.WithType(diag.DATABASE), diag.WithSummary("failed bulk insert on table %q", e.Table.Name)))
	// Try to insert resource by resource if partial fetch is enabled and an error occurred
	partialFetchResources := make(schema.Resources, 0)
	var failed error
	failedCount := 0
	for id := range resources {
		if err := e.Db.Insert(ctx, e.Table, schema.Resources{resources[id]}, shouldCascade); err != nil {
			failed = err
			failedCount++
			e.Logger.Error("failed to insert resource into db", "error", err, "resource_keys", resources[id].PrimaryKeyValues())
			diags = diags.Add(ClassifyError(err, diag.WithType(diag.DATABASE)))
			continue
		}
		// If there is no error we add the resource to the final result
		partialFetchResources = append(partialFetchResources, resources[id])
	}
	if failed != nil {
		msg := "all resources"
		if failedCount < len(resources) {
			msg = "some resources"
		}
		diags = diags.Add(diag.TelemetryFromError(
			failed,
			diag.InsertFailed,
			diag.WithSummary("%s failed to insert into table %q", msg, e.Table.Name),
		))
	}
	return partialFetchResources, diags
}

// resolveResourceValues does the actual resolve of all the columns of table for said resource.
func (e TableExecutor) resolveResourceValues(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) (diags diag.Diagnostics) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			e.Logger.Error("resolve table recovered from panic", "panic_msg", r, "stack", stack)
			diags = fromError(fmt.Errorf("column resolve panic: %s", r), diag.WithResourceName(e.ResourceName), diag.WithSeverity(diag.PANIC),
				diag.WithSummary("resolve table %q recovered from panic", e.Table.Name), diag.WithDetails("%s", stack))
		}
	}()

	diags = diags.Add(e.resolveColumns(ctx, meta, resource, e.columns[0]))
	if diags.HasErrors() {
		return diags
	}

	// call PostRowResolver if defined after columns have been resolved
	if e.Table.PostResourceResolver != nil {
		if err := e.Table.PostResourceResolver(ctx, meta, resource); err != nil {
			diags = diags.Add(e.handleResolveError(meta, resource, err, diag.WithSummary("post resource resolver failed for %q", e.Table.Name)))

			if diags.HasErrors() {
				return diags
			}
		}
	}
	// Finally, resolve columns internal to the SDK
	for _, c := range e.columns[1] {
		if err := c.Resolver(ctx, meta, resource, c); err != nil {
			return diags.Add(fromError(err, diag.WithResourceName(e.ResourceName), WithResource(resource), diag.WithType(diag.INTERNAL), diag.WithSummary("default column %q resolver execution", c.Name)))
		}
	}
	return diags
}

// resolveColumns resolves each column in the table and adds them to the resource.
func (e TableExecutor) resolveColumns(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, cols []schema.Column) (diags diag.Diagnostics) {
	var col string

	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			e.Logger.Error("resolve columns recovered from panic", "panic_msg", r, "stack", stack, "column_name", col)
			diags = fromError(fmt.Errorf("column resolve panic: %s", r), diag.WithResourceName(e.ResourceName), diag.WithSeverity(diag.PANIC),
				diag.WithSummary("resolve column %q in table %q recovered from panic", col, e.Table.Name), diag.WithDetails("%s", stack))
		}
	}()

	for _, c := range cols {
		col = c.Name
		if c.Resolver != nil {
			e.Logger.Trace("using custom column resolver", "column", c.Name)
			err := c.Resolver(ctx, meta, resource, c)
			if err == nil {
				continue
			}
			// Not allowed ignoring PK resolver errors
			if funk.ContainsString(e.Db.Dialect().PrimaryKeys(e.Table), c.Name) {
				return diags.Add(ClassifyError(err, diag.WithResourceName(e.ResourceName), WithResource(resource), diag.WithSummary("failed to resolve column %s@%s", e.Table.Name, c.Name)))
			}
			diags = diags.Add(e.handleResolveError(meta, resource, err, diag.WithSummary("column resolver %q failed for table %q", c.Name, e.Table.Name)))
			continue
		}
		e.Logger.Trace("resolving column value with path", "column", c.Name)
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.Item, strcase.ToCamel(c.Name), funk.WithAllowZero())
		e.Logger.Trace("setting column value", "column", c.Name, "value", v)
		if err := resource.Set(c.Name, v); err != nil {
			diags = diags.Add(fromError(err, diag.WithResourceName(e.ResourceName), diag.WithType(diag.INTERNAL),
				diag.WithSummary("failed to set resource value for column %s@%s", e.Table.Name, c.Name)))
		}
	}
	return diags
}

// handleResolveError handles errors returned by user defined functions, using the ErrorClassifiers if defined.
func (e TableExecutor) handleResolveError(meta schema.ClientMeta, r *schema.Resource, err error, opts ...diag.BaseErrorOption) diag.Diagnostics {
	errAsDiags := fromError(err, append(opts,
		diag.WithResourceName(e.ResourceName),
		WithResource(r),
		diag.WithOptionalSeverity(diag.ERROR),
		diag.WithType(diag.RESOLVING),
		diag.WithSummary("failed to resolve table %q", e.Table.Name),
	)...)

	classifiedDiags := make(diag.Diagnostics, 0, len(errAsDiags))
	for _, c := range e.classifiers {
		// fromError gives us diag.Diagnostics, but we need to make sure to pass one diag at a time to the classifiers and collect results,
		// mostly because Unwrap()/errors.As() can't work on multiple diags
		for _, d := range errAsDiags {
			if diags := c(meta, e.ResourceName, d); diags != nil {
				classifiedDiags = classifiedDiags.Add(diags)
			}
		}
	}
	if classifiedDiags.HasDiags() {
		return classifiedDiags
	}

	return errAsDiags
}

// IgnoreError returns true if the error is ignored via the current table IgnoreError function or in any other parent table (in that ordered)
// it stops checking the moment one of them exists and not until it returns true or fals
func (e TableExecutor) IgnoreError(err error) bool {
	// first priority is to check the tables IgnoreError function
	if e.Table.IgnoreError != nil {
		return e.Table.IgnoreError(err)
	}
	// secondy priority is to check the parent tables IgnoreError recursively
	if e.ParentExecutor != nil {
		return e.ParentExecutor.IgnoreError(err)
	}

	return false
}

func identifyClient(meta schema.ClientMeta) string {
	ider, ok := meta.(schema.ClientIdentifier)
	if ok {
		return ider.Identify()
	}
	return ""
}
