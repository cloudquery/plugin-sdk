package schema

import (
	"context"
	"fmt"
	"regexp"
	"sync/atomic"
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
}

var reValidTableName = regexp.MustCompile(`^[a-z_][a-z\d_]*$`)
var reValidColumnName = regexp.MustCompile(`^[a-z_][a-z\d_]*$`)

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

func (tt Tables) ValidateTableNames() error {
	for _, t := range tt {
		if err := t.ValidateName(); err != nil {
			return err
		}
	}
	return nil
}

func (tt Tables) ValidateColumnNames() error {
	for _, t := range tt {
		if err := t.ValidateColumnNames(); err != nil {
			return err
		}
	}
	return nil
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

func (t *Table) ValidateName() error {
	ok := reValidTableName.MatchString(t.Name)
	if !ok {
		return fmt.Errorf("table name %q is not valid: table names must contain only lower-case letters, numbers and underscores, and must start with a lower-case letter or underscore", t.Name)
	}
	return nil
}

func (t *Table) ValidateDuplicateColumns() error {
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

func (t *Table) ValidateColumnNames() error {
	for _, c := range t.Columns {
		ok := reValidColumnName.MatchString(c.Name)
		if !ok {
			return fmt.Errorf("column name %q on table %q is not valid: column names must contain only lower-case letters, numbers and underscores, and must start with a lower-case letter or underscore", c.Name, t.Name)
		}
	}
	return nil
}

func (t *Table) Column(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// If the column with the same name exists, overwrites it.
// Otherwise, adds the column to the beginning of the table.
func (t *Table) OverwriteOrAddColumn(column *Column) {
	for i, c := range t.Columns {
		if c.Name == column.Name {
			t.Columns[i] = *column
			return
		}
	}
	t.Columns = append([]Column{*column}, t.Columns...)
}

func (t *Table) PrimaryKeys() []string {
	var primaryKeys []string
	for _, c := range t.Columns {
		if c.CreationOptions.PrimaryKey {
			primaryKeys = append(primaryKeys, c.Name)
		}
	}

	return primaryKeys
}

func (t *Table) TableNames() []string {
	ret := []string{t.Name}
	for _, rel := range t.Relations {
		ret = append(ret, rel.TableNames()...)
	}
	return ret
}

// Get return table by name
func (tt Tables) Get(name string) *Table {
	for _, t := range tt {
		if t.Name == name {
			return t
		}
		table := t.Relations.Get(name)
		if table != nil {
			return table
		}
	}
	return nil
}
