package scheduler

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

type refreshCall struct {
	table string
	keys  string
}

// fakeRefresher answers Refresh from a canned map keyed by "<table>:<keys>".
type fakeRefresher struct {
	items map[string][]any
	errs  map[string]error

	mu    sync.Mutex
	calls []refreshCall
}

func (f *fakeRefresher) Refresh(_ context.Context, table *schema.Table, keys arrow.RecordBatch) (schema.ClientMeta, []any, error) {
	key := table.Name + ":" + schema.KeysString(keys)
	f.mu.Lock()
	f.calls = append(f.calls, refreshCall{table: table.Name, keys: schema.KeysString(keys)})
	f.mu.Unlock()
	if err, ok := f.errs[key]; ok {
		return nil, nil, err
	}
	return &testExecutionClient{}, f.items[key], nil
}

func (f *fakeRefresher) refreshCalls() []refreshCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.calls)
}

func refreshTestTables() schema.Tables {
	child := &schema.Table{
		Name: "test_refresh_child",
		Resolver: func(_ context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
			item := parent.Item.(map[string]any)
			res <- map[string]any{"child_id": item["id"].(string) + "-child"}
			return nil
		},
		Columns: schema.ColumnList{
			{Name: "child_id", Type: arrow.BinaryTypes.String, PrimaryKey: true, Resolver: schema.PathResolver("child_id")},
		},
	}
	parent := &schema.Table{
		Name: "test_refresh_parent",
		Resolver: func(context.Context, schema.ClientMeta, *schema.Resource, chan<- any) error {
			return errors.New("the table resolver must not run during a refresh")
		},
		Multiplex: func(schema.ClientMeta) []schema.ClientMeta {
			panic("Multiplex must not be called during a refresh")
		},
		Columns: schema.ColumnList{
			{Name: "id", Type: arrow.BinaryTypes.String, PrimaryKey: true, Resolver: schema.PathResolver("id")},
			{Name: "name", Type: arrow.BinaryTypes.String, Resolver: schema.PathResolver("name")},
		},
		Relations: schema.Tables{child},
	}
	schema.AddCqIDs(parent)
	return schema.Tables{parent.Copy(nil)}
}

func refreshKeys(t *testing.T, table *schema.Table, id string) arrow.RecordBatch {
	t.Helper()
	keys, err := schema.NewKeysRecord(table, map[string]any{"id": id})
	require.NoError(t, err)
	return keys
}

func runRefresh(t *testing.T, tables schema.Tables, r Refresher, resources []plugin.RefreshResource, opts ...SyncOption) (message.SyncMessages, error) {
	t.Helper()
	s := NewScheduler(WithConcurrency(2))
	res := make(chan message.SyncMessage, 128)
	var refreshErr error
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer close(res)
		refreshErr = s.Refresh(context.Background(), tables, r, resources, res, opts...)
	}()
	<-done
	var msgs message.SyncMessages
	for msg := range res {
		msgs = append(msgs, msg)
	}
	return msgs, refreshErr
}

func TestRefresh_InsertsItemsAndRelations(t *testing.T) {
	tables := refreshTestTables()
	parent := tables.Get("test_refresh_parent")
	keys := refreshKeys(t, parent, "abc")

	r := &fakeRefresher{items: map[string][]any{
		"test_refresh_parent:id=abc": {map[string]any{"id": "abc", "name": "hello"}},
	}}

	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_refresh_parent", Keys: keys, Op: plugin.OpUpsert},
	})
	require.NoError(t, err)

	inserts := msgs.GetInserts()
	require.True(t, inserts.Exists("test_refresh_parent"), "expected an insert for the refreshed table")
	require.True(t, inserts.Exists("test_refresh_child"), "expected an insert for the relation")
	require.EqualValues(t, 2, msgs.InsertItems())

	for _, msg := range msgs {
		_, isMigrate := msg.(*message.SyncMigrateTable)
		require.False(t, isMigrate, "Refresh must not emit migrate messages")
	}
}

func TestRefresh_DeterministicCQID(t *testing.T) {
	run := func(opts ...SyncOption) string {
		tables := refreshTestTables()
		parent := tables.Get("test_refresh_parent")
		r := &fakeRefresher{items: map[string][]any{
			"test_refresh_parent:id=abc": {map[string]any{"id": "abc", "name": "hello"}},
		}}
		msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
			{Table: "test_refresh_parent", Keys: refreshKeys(t, parent, "abc"), Op: plugin.OpUpsert},
		}, opts...)
		require.NoError(t, err)

		records := msgs.GetInserts().GetRecordsForTable(tables.Get("test_refresh_parent"))
		require.Len(t, records, 1)
		idx := records[0].Schema().FieldIndices(schema.CqIDColumn.Name)
		require.Len(t, idx, 1)
		return records[0].Column(idx[0]).ValueStr(0)
	}

	first := run(WithSyncDeterministicCQID(true))
	second := run(WithSyncDeterministicCQID(true))
	require.Equal(t, first, second, "deterministic _cq_id must be stable across runs")

	require.NotEqual(t, first, run(), "random _cq_id must not match the deterministic one")
}

func TestRefresh_NotFoundEmitsDeleteRecord(t *testing.T) {
	tables := refreshTestTables()
	parent := tables.Get("test_refresh_parent")

	r := &fakeRefresher{errs: map[string]error{
		"test_refresh_parent:id=gone": ErrNotFound,
	}}

	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_refresh_parent", Keys: refreshKeys(t, parent, "gone"), Op: plugin.OpUpsert},
	})
	require.NoError(t, err)
	require.Len(t, msgs, 1)

	del, ok := msgs[0].(*message.SyncDeleteRecord)
	require.True(t, ok, "expected a delete record, got %T", msgs[0])
	require.Equal(t, "test_refresh_parent", del.TableName)
	require.Len(t, del.WhereClause, 1)
	require.Equal(t, "AND", del.WhereClause[0].GroupingType)
	require.Len(t, del.WhereClause[0].Predicates, 1)
	require.Equal(t, "EQ", del.WhereClause[0].Predicates[0].Operator)
	require.Equal(t, "id", del.WhereClause[0].Predicates[0].Column)
	require.Equal(t, "gone", del.WhereClause[0].Predicates[0].Record.Column(0).ValueStr(0))
	require.Equal(t, message.TableRelations{
		{TableName: "test_refresh_child", ParentTable: "test_refresh_parent"},
	}, del.TableRelations)
}

func TestRefresh_OpDeleteSkipsRefresh(t *testing.T) {
	tables := refreshTestTables()
	parent := tables.Get("test_refresh_parent")

	r := &fakeRefresher{}
	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_refresh_parent", Keys: refreshKeys(t, parent, "abc"), Op: plugin.OpDelete},
	})
	require.NoError(t, err)
	require.Empty(t, r.refreshCalls(), "OpDelete must not call Refresh")
	require.Len(t, msgs, 1)
	require.IsType(t, &message.SyncDeleteRecord{}, msgs[0])
}

// componentsRefreshTable is shaped like an AWS table: _cq_id is the only primary
// key and identity lives in the primary key components.
func componentsRefreshTable() schema.Tables {
	table := &schema.Table{
		Name: "test_components",
		Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- map[string]any{"arn": "arn:aws:ec2:i-123", "name": "hello"}
			return nil
		},
		Columns: schema.ColumnList{
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true, Resolver: schema.PathResolver("arn")},
			{Name: "name", Type: arrow.BinaryTypes.String, Resolver: schema.PathResolver("name")},
		},
	}
	schema.AddCqIDs(table)
	return schema.Tables{table.Copy(nil)}
}

func TestRefresh_OpDeleteOnPrimaryKeyComponentTable(t *testing.T) {
	tables := componentsRefreshTable()
	table := tables.Get("test_components")
	require.Equal(t, []string{schema.CqIDColumn.Name}, table.PrimaryKeys())

	keys, err := schema.NewKeysRecord(table, map[string]any{"arn": "arn:aws:ec2:i-123"})
	require.NoError(t, err)

	r := &fakeRefresher{}
	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_components", Keys: keys, Op: plugin.OpDelete},
	})
	require.NoError(t, err)
	require.Empty(t, r.refreshCalls(), "OpDelete must not call Refresh")
	require.Len(t, msgs, 1)

	del, ok := msgs[0].(*message.SyncDeleteRecord)
	require.True(t, ok, "expected a delete record, got %T", msgs[0])
	require.Equal(t, "test_components", del.TableName)
	require.Len(t, del.WhereClause[0].Predicates, 1)
	require.Equal(t, "arn", del.WhereClause[0].Predicates[0].Column)
	require.Equal(t, "arn:aws:ec2:i-123", del.WhereClause[0].Predicates[0].Record.Column(0).ValueStr(0))
}

// iamPolicyLikeRefreshTable mirrors aws_iam_policies: four identity columns, of
// which an event knows only account_id and arn.
func iamPolicyLikeRefreshTable() schema.Tables {
	table := &schema.Table{
		Name: "test_iam_policies",
		Resolver: func(context.Context, schema.ClientMeta, *schema.Resource, chan<- any) error {
			return errors.New("the table resolver must not run during a refresh")
		},
		Columns: schema.ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true, Resolver: schema.PathResolver("account_id")},
			{Name: "id", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true, Resolver: schema.PathResolver("id")},
			{Name: "input_hash", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true, Resolver: schema.PathResolver("input_hash")},
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true, Resolver: schema.PathResolver("arn")},
		},
	}
	schema.AddCqIDs(table)
	return schema.Tables{table.Copy(nil)}
}

func TestRefresh_OpDeleteWithIdentityColumnSubset(t *testing.T) {
	tables := iamPolicyLikeRefreshTable()
	table := tables.Get("test_iam_policies")
	require.Equal(t, []string{"account_id", "id", "input_hash", "arn"}, table.IdentityColumns())

	// The event knows 2 of the 4 identity columns.
	keys, err := schema.NewKeysRecord(table, map[string]any{
		"account_id": "123456789012",
		"arn":        "arn:aws:iam::123456789012:policy/demo",
	})
	require.NoError(t, err)

	r := &fakeRefresher{}
	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_iam_policies", Keys: keys, Op: plugin.OpDelete},
	})
	require.NoError(t, err)
	require.Empty(t, r.refreshCalls(), "OpDelete must not call Refresh")
	require.Len(t, msgs, 1)

	del, ok := msgs[0].(*message.SyncDeleteRecord)
	require.True(t, ok, "expected a delete record, got %T", msgs[0])
	require.Len(t, del.WhereClause[0].Predicates, 2, "predicates must cover exactly the columns the event supplied")
	require.Equal(t, "account_id", del.WhereClause[0].Predicates[0].Column)
	require.Equal(t, "arn", del.WhereClause[0].Predicates[1].Column)
	require.Equal(t, "arn:aws:iam::123456789012:policy/demo", del.WhereClause[0].Predicates[1].Record.Column(0).ValueStr(0))
}

func TestRefresh_NotFoundWithIdentityColumnSubset(t *testing.T) {
	tables := iamPolicyLikeRefreshTable()
	table := tables.Get("test_iam_policies")

	keys, err := schema.NewKeysRecord(table, map[string]any{"arn": "arn:aws:iam::123456789012:policy/gone"})
	require.NoError(t, err)

	r := &fakeRefresher{errs: map[string]error{
		"test_iam_policies:arn=arn:aws:iam::123456789012:policy/gone": ErrNotFound,
	}}
	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_iam_policies", Keys: keys, Op: plugin.OpUpsert},
	})
	require.NoError(t, err)
	require.Len(t, msgs, 1)

	del, ok := msgs[0].(*message.SyncDeleteRecord)
	require.True(t, ok, "expected a delete record, got %T", msgs[0])
	require.Len(t, del.WhereClause[0].Predicates, 1)
	require.Equal(t, "arn", del.WhereClause[0].Predicates[0].Column)
}

func TestRefresh_NotFoundOnPrimaryKeyComponentTable(t *testing.T) {
	tables := componentsRefreshTable()
	table := tables.Get("test_components")

	keys, err := schema.NewKeysRecord(table, map[string]any{"arn": "arn:aws:ec2:i-gone"})
	require.NoError(t, err)

	r := &fakeRefresher{errs: map[string]error{
		"test_components:arn=arn:aws:ec2:i-gone": ErrNotFound,
	}}
	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_components", Keys: keys, Op: plugin.OpUpsert},
	})
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	require.IsType(t, &message.SyncDeleteRecord{}, msgs[0])
}

func TestRefresh_UpsertOnPrimaryKeyComponentTable(t *testing.T) {
	tables := componentsRefreshTable()
	table := tables.Get("test_components")

	keys, err := schema.NewKeysRecord(table, map[string]any{"arn": "arn:aws:ec2:i-123"})
	require.NoError(t, err)

	r := &fakeRefresher{items: map[string][]any{
		"test_components:arn=arn:aws:ec2:i-123": {map[string]any{"arn": "arn:aws:ec2:i-123", "name": "hello"}},
	}}
	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_components", Keys: keys, Op: plugin.OpUpsert},
	}, WithSyncDeterministicCQID(true))
	require.NoError(t, err)
	require.EqualValues(t, 1, msgs.InsertItems())
	require.True(t, msgs.GetInserts().Exists("test_components"))
}

func TestRefresh_OneFailureDoesNotStopTheOthers(t *testing.T) {
	tables := refreshTestTables()
	parent := tables.Get("test_refresh_parent")

	boom := errors.New("boom")
	r := &fakeRefresher{
		items: map[string][]any{
			"test_refresh_parent:id=ok1": {map[string]any{"id": "ok1", "name": "one"}},
			"test_refresh_parent:id=ok2": {map[string]any{"id": "ok2", "name": "two"}},
		},
		errs: map[string]error{"test_refresh_parent:id=bad": boom},
	}

	msgs, err := runRefresh(t, tables, r, []plugin.RefreshResource{
		{Table: "test_refresh_parent", Keys: refreshKeys(t, parent, "ok1"), Op: plugin.OpUpsert},
		{Table: "test_refresh_parent", Keys: refreshKeys(t, parent, "bad"), Op: plugin.OpUpsert},
		{Table: "test_refresh_parent", Keys: refreshKeys(t, parent, "ok2"), Op: plugin.OpUpsert},
	})
	require.ErrorIs(t, err, boom)

	// Both healthy resources produced their parent row and their relation row.
	require.EqualValues(t, 4, msgs.InsertItems())

	var syncErrors []*message.SyncError
	for _, msg := range msgs {
		if e, ok := msg.(*message.SyncError); ok {
			syncErrors = append(syncErrors, e)
		}
	}
	require.Len(t, syncErrors, 1)
	require.Equal(t, "test_refresh_parent", syncErrors[0].TableName)
	require.Contains(t, syncErrors[0].Error, "boom")
}

func TestRefresh_UnknownTable(t *testing.T) {
	tables := refreshTestTables()
	parent := tables.Get("test_refresh_parent")

	_, err := runRefresh(t, tables, &fakeRefresher{}, []plugin.RefreshResource{
		{Table: "test_no_such_table", Keys: refreshKeys(t, parent, "abc"), Op: plugin.OpUpsert},
	})
	require.ErrorContains(t, err, "test_no_such_table")
}

func TestRefresh_NoResources(t *testing.T) {
	msgs, err := runRefresh(t, refreshTestTables(), &fakeRefresher{}, nil)
	require.NoError(t, err)
	require.Empty(t, msgs)
}
