package memdb

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination/batchingwriter"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/rs/zerolog"
)

type opencloseclient struct {
	*client

	errOnOpen    bool
	errOnClose   bool
	openTables   []string
	closedTables []string
	tblLock      *sync.Mutex
}

func getNewOCWClient(base opencloseclient, options ...Option) destination.NewClientFunc {
	return func(ctx context.Context, lgr zerolog.Logger, spec specs.Destination) (destination.Client, error) {
		c, err := GetNewClient(options...)(ctx, lgr, spec)
		if err != nil {
			return nil, err
		}
		base.client = c.(*client)
		base.tblLock = &sync.Mutex{}

		return &base, nil
	}
}

func (c *opencloseclient) WriteTableBatch(_ context.Context, _ specs.Source, table *schema.Table, _ time.Time, resources []arrow.Record) error {
	tableName := table.Name
	for _, resource := range resources {
		c.memoryDBLock.Lock()
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[tableName] = append(c.memoryDB[tableName], resource)
		} else {
			c.overwrite(table, resource)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *opencloseclient) OpenTable(_ context.Context, _ specs.Source, table *schema.Table) error {
	if c.errOnOpen {
		return fmt.Errorf("errOnOpen")
	}
	c.tblLock.Lock()
	defer c.tblLock.Unlock()
	c.openTables = append(c.openTables, table.Name)
	return nil
}
func (c *opencloseclient) CloseTable(_ context.Context, _ specs.Source, table *schema.Table) error {
	if c.errOnClose {
		return fmt.Errorf("errOnClose")
	}
	c.tblLock.Lock()
	defer c.tblLock.Unlock()
	c.closedTables = append(c.closedTables, table.Name)
	return nil
}

func (c *opencloseclient) Close(_ context.Context) error {
	// run the OpenClose-specific test checks here
	c.tblLock.Lock()
	defer c.tblLock.Unlock()
	if l := len(c.openTables); l == 0 && !c.errOnOpen {
		return fmt.Errorf("no tables were opened")
	} else if l > 0 && c.errOnOpen {
		return fmt.Errorf("%d tables were opened, expected 0", l)
	}
	if l := len(c.closedTables); l == 0 && !c.errOnClose {
		return fmt.Errorf("no tables were closed")
	} else if l > 0 && c.errOnClose {
		return fmt.Errorf("%d tables were closed, expected 0", l)
	}
	return nil
}

//nolint:revive
func validateOCWClient(expectErrors bool) func(t *testing.T, p *destination.Plugin, destSpec specs.Destination) {
	return func(t *testing.T, p *destination.Plugin, destSpec specs.Destination) {
		t.Helper()
		if expectErrors && p.Metrics().Errors == 0 {
			t.Fatal("expected errors, got none")
		}
		if !expectErrors && p.Metrics().Errors > 0 {
			t.Fatalf("not expected errors, got %d", p.Metrics().Errors)
		}
	}
}

func TestPluginManagedClientWithOCW(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", getNewOCWClient(opencloseclient{}), destination.WithManagedWriter(batchingwriter.New()))
		},
		specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
			Validate:                 validateOCWClient(false),
		})
}

func TestPluginManagedClientWithOCWCloseError(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", getNewOCWClient(opencloseclient{errOnOpen: false, errOnClose: true}), destination.WithManagedWriter(batchingwriter.New()))
		},
		specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
			Validate:                 validateOCWClient(true),
		})
}
