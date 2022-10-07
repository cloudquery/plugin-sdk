package schema

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/caser"
	"github.com/cloudquery/plugin-sdk/helpers"
	"github.com/getsentry/sentry-go"
	"github.com/thoas/go-funk"
)

// TableResolver is the main entry point when a table is sync is called.
//
// Table resolver has 3 main arguments:
// - meta(ClientMeta): is the client returned by the plugin.Provider Configure call
// - parent(Resource): resource is the parent resource in case this table is called via parent table (i.e. relation)
// - res(chan interface{}): is a channel to pass results fetched by the TableResolver
type TableResolver func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error

type RowResolver func(ctx context.Context, meta ClientMeta, resource *Resource) error

type Multiplexer func(meta ClientMeta) []ClientMeta

type Tables []*Table

type SyncSummary struct {
	Resources uint64
	Errors    uint64
	Panics    uint64
}

type Table struct {
	// Name of table
	Name string `json:"name"`
	// table description
	Description string `json:"description"`
	// Columns are the set of fields that are part of this table
	Columns ColumnList `json:"columns"`
	// Relations are a set of related tables defines
	Relations Tables `json:"relations"`
	// Resolver is the main entry point to fetching table data and
	Resolver TableResolver `json:"-"`
	// Multiplex returns re-purposed meta clients. The sdk will execute the table with each of them
	Multiplex Multiplexer `json:"-"`
	// PostResourceResolver is called after all columns have been resolved, but before the Resource is sent to be inserted. The ordering of resolvers is:
	//  (Table) Resolver → PreResourceResolver → ColumnResolvers → PostResourceResolver
	PostResourceResolver RowResolver `json:"-"`
	// PreResourceResolver is called before all columns are resolved but after Resource is created. The ordering of resolvers is:
	//  (Table) Resolver → PreResourceResolver → ColumnResolvers → PostResourceResolver
	PreResourceResolver RowResolver `json:"-"`

	// IgnoreInTests is used to exclude a table from integration tests.
	// By default, integration tests fetch all resources from cloudquery's test account, and verify all tables
	// have at least one row.
	// When IgnoreInTests is true, integration tests won't fetch from this table.
	// Used when it is hard to create a reproducible environment with a row in this table.
	IgnoreInTests bool `json:"ignore_in_tests"`

	// Parent is the parent table in case this table is called via parent table (i.e. relation)
	Parent *Table `json:"-"`

	// Serial is used to force a signature change, which forces new table creation and cascading removal of old table and relations
	Serial string `json:"-"`

	columnsMap map[string]int
}

func (s *SyncSummary) Merge(other SyncSummary) {
	atomic.AddUint64(&s.Resources, other.Resources)
	atomic.AddUint64(&s.Errors, other.Errors)
	atomic.AddUint64(&s.Panics, other.Panics)
}

func (tt Tables) TableNames() []string {
	ret := []string{}
	for _, t := range tt {
		ret = append(ret, t.TableNames()...)
	}
	return ret
}

func (tt Tables) ValidateDuplicateColumns() error {
	for _, t := range tt {
		if err := t.ValidateDuplicateColumns(); err != nil {
			return err
		}
	}
	return nil
}

func (tt Tables) ValidateDuplicateTables() error {
	tables := make(map[string]bool, len(tt))
	for _, t := range tt {
		if _, ok := tables[t.Name]; ok {
			return fmt.Errorf("duplicate table %s", t.Name)
		}
		tables[t.Name] = true
	}
	return nil
}

func (t Table) ValidateDuplicateColumns() error {
	columns := make(map[string]bool, len(t.Columns))
	for _, c := range t.Columns {
		if _, ok := columns[c.Name]; ok {
			return fmt.Errorf("duplicate column %s in table %s", c.Name, t.Name)
		}
		columns[c.Name] = true
	}
	for _, rel := range t.Relations {
		if err := rel.ValidateDuplicateColumns(); err != nil {
			return err
		}
	}
	return nil
}

func (t Table) Column(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func (t Table) PrimaryKeys() []string {
	var primaryKeys []string
	for _, c := range t.Columns {
		if c.CreationOptions.PrimaryKey {
			primaryKeys = append(primaryKeys, c.Name)
		}
	}

	return primaryKeys
}

func (t Table) ColumnIndex(name string) int {
	var once sync.Once
	once.Do(func() {
		if t.columnsMap == nil {
			t.columnsMap = make(map[string]int)
			for i, c := range t.Columns {
				t.columnsMap[c.Name] = i
			}
		}
	})
	if index, ok := t.columnsMap[name]; ok {
		return index
	}
	return -1
}

func (t Table) TableNames() []string {
	ret := []string{t.Name}
	for _, rel := range t.Relations {
		ret = append(ret, rel.TableNames()...)
	}
	return ret
}

// Call the table resolver with with all of it's relation for every reolved resource
func (t Table) Resolve(ctx context.Context, meta ClientMeta, parent *Resource, resolvedResources chan<- *Resource) (summary SyncSummary) {
	tableStartTime := time.Now()
	meta.Logger().Info().Str("table", t.Name).Msg("table resolver started")

	res := make(chan interface{})
	startTime := time.Now()
	go func(sum *SyncSummary) {
		defer func() {
			if err := recover(); err != nil {
				stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", t.Name)
					sentry.CurrentHub().CaptureMessage(stack)
				})
				meta.Logger().Error().Interface("error", err).Str("table", t.Name).TimeDiff("duration", time.Now(), startTime).Str("stack", stack).Msg("table resolver finished with panic")
				atomic.AddUint64(&sum.Panics, 1)
			}
			close(res)
		}()
		meta.Logger().Debug().Str("table", t.Name).Msg("table resolver started")
		if err := t.Resolver(ctx, meta, parent, res); err != nil {
			meta.Logger().Error().Str("table", t.Name).TimeDiff("duration", time.Now(), startTime).Err(err).Msg("table resolver finished with error")
			atomic.AddUint64(&sum.Errors, 1)
			return
		}
		meta.Logger().Debug().Str("table", t.Name).TimeDiff("duration", time.Now(), startTime).Msg("table resolver finished successfully")
	}(&summary)
	tableResources := 0
	for elem := range res {
		objects := helpers.InterfaceSlice(elem)
		if len(objects) == 0 {
			continue
		}
		for i := range objects {
			summary.Merge(t.resolveObject(ctx, meta, parent, objects[i], resolvedResources))
		}
	}
	meta.Logger().Info().Str("table", t.Name).Int("total_resources", tableResources).TimeDiff("duration", time.Now(), tableStartTime).Msg("fetch table finished")

	return summary
}

func (t Table) resolveObject(ctx context.Context, meta ClientMeta, parent *Resource, item interface{}, resolvedResources chan<- *Resource) (summary SyncSummary) {
	resource := NewResourceData(&t, parent, item)
	objectStartTime := time.Now()
	csr := caser.New()
	meta.Logger().Info().Str("table", t.Name).Msg("object resolver started")
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("table", t.Name)
				sentry.CurrentHub().CaptureMessage(stack)
			})
			meta.Logger().Error().Interface("error", err).Str("table", t.Name).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("object resolver finished with panic")
			summary.Panics++
		}
	}()
	if t.PreResourceResolver != nil {
		meta.Logger().Trace().Str("table", t.Name).Msg("pre resource resolver started")
		if err := t.PreResourceResolver(ctx, meta, resource); err != nil {
			meta.Logger().Error().Str("table", t.Name).Err(err).Msg("pre resource resolver failed")
			summary.Errors++
			return summary
		}
		meta.Logger().Trace().Str("table", t.Name).Msg("pre resource resolver finished successfully")
	}

	for _, c := range t.Columns {
		if c.Resolver != nil {
			meta.Logger().Trace().Str("column_name", c.Name).Str("table", t.Name).Msg("column resolver custom started")
			if err := c.Resolver(ctx, meta, resource, c); err != nil {
				meta.Logger().Error().Str("column_name", c.Name).Str("table", t.Name).Err(err).Msg("column resolver finished with error")
				summary.Errors++
			}
			meta.Logger().Trace().Str("column_name", c.Name).Str("table", t.Name).Msg("column resolver finished successfully")
		} else {
			meta.Logger().Trace().Str("column_name", c.Name).Str("table", t.Name).Msg("column resolver default started")
			// base use case: try to get column with CamelCase name
			v := funk.Get(resource.Item, csr.ToPascal(c.Name), funk.WithAllowZero())
			if v != nil {
				if err := resource.Set(c.Name, v); err != nil {
					meta.Logger().Error().Str("column_name", c.Name).Str("table", t.Name).Err(err).Msg("column resolver default finished with error")
					summary.Errors++
				}
				meta.Logger().Trace().Str("column_name", c.Name).Str("table", t.Name).Msg("column resolver default finished successfully")
			} else {
				meta.Logger().Trace().Str("column_name", c.Name).Str("table", t.Name).Msg("column resolver default finished successfully with nil")
			}
		}
	}

	if t.PostResourceResolver != nil {
		meta.Logger().Trace().Str("table", t.Name).Msg("post resource resolver started")
		if err := t.PostResourceResolver(ctx, meta, resource); err != nil {
			meta.Logger().Error().Str("table", t.Name).Stack().Err(err).Msg("post resource resolver finished with error")
			summary.Errors++
		} else {
			meta.Logger().Trace().Str("table", t.Name).Msg("post resource resolver finished successfully")
		}
	}
	summary.Resources = 1
	resolvedResources <- resource

	for _, rel := range t.Relations {
		summary.Merge(rel.Resolve(ctx, meta, resource, resolvedResources))
	}

	return summary
}
