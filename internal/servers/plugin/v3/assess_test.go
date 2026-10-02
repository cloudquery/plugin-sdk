package plugin

import (
	"context"
	"net"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

type assessorClient struct {
	mockSourceColumnAdderPluginClient
	gotTables  []plugin.TablePair
	gotOptions plugin.AssessOptions
}

func (c *assessorClient) AssessTables(_ context.Context, tables []plugin.TablePair, options plugin.AssessOptions) ([]plugin.TableFinding, error) {
	c.gotTables, c.gotOptions = tables, options
	return []plugin.TableFinding{{
		TableName:          "test_table",
		Category:           plugin.AssessCategoryManualMigrationRequired,
		SafeModeBehavior:   "rejects this change",
		ForcedModeBehavior: "drops and recreates the table",
		Columns: []plugin.ColumnFinding{{
			ColumnName:         "tags",
			Category:           plugin.AssessCategoryManualMigrationRequired,
			OldType:            "text[]",
			NewType:            "jsonb",
			SafeModeBehavior:   "rejects this change",
			ForcedModeBehavior: "drops and recreates the table",
			Evidence:           []plugin.Evidence{{SyntheticValue: `["env:prod"]`, Before: `{"tags":["env:prod"]}`, After: `{"tags":["env:prod"]}`}},
		}},
		Evidence:                 []plugin.Evidence{{SyntheticValue: "header", Before: "tags", After: "tags"}},
		CoverageIncomplete:       true,
		CoverageIncompleteReason: "nested values not compared",
	}}, nil
}

func newAssessClient(t *testing.T, newClient plugin.NewClientFunc) pb.PluginClient {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	pb.RegisterPluginServer(srv, &Server{Plugin: plugin.NewPlugin("test", "development", newClient), Logger: zerolog.Nop()})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := pb.NewPluginClient(conn)
	_, err = client.Init(context.Background(), &pb.Init_Request{NoConnection: true})
	require.NoError(t, err)
	return client
}

func tableBytes(t *testing.T, table *schema.Table) []byte {
	t.Helper()
	b, err := pb.SchemaToBytes(table.ToArrowSchema())
	require.NoError(t, err)
	return b
}

func TestAssessTablesRoundTrip(t *testing.T) {
	assessor := &assessorClient{}
	client := newAssessClient(t, func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return assessor, nil
	})
	oldTable := &schema.Table{Name: "test_table", Columns: []schema.Column{{Name: "tags", Type: arrow.ListOf(arrow.BinaryTypes.String)}}}
	newTable := &schema.Table{Name: "test_table", Columns: []schema.Column{{Name: "tags", Type: arrow.BinaryTypes.String}}}

	resp, err := client.AssessTables(context.Background(), &pb.AssessTables_Request{
		Tables:       []*pb.AssessTables_TablePair{{OldTable: tableBytes(t, oldTable), NewTable: tableBytes(t, newTable)}, {NewTable: tableBytes(t, newTable)}},
		MigrateForce: true,
	})
	require.NoError(t, err)

	require.True(t, assessor.gotOptions.MigrateForce)
	require.Len(t, assessor.gotTables, 2)
	require.Equal(t, "test_table", assessor.gotTables[0].Old.Name)
	require.True(t, arrow.TypeEqual(arrow.ListOf(arrow.BinaryTypes.String), assessor.gotTables[0].Old.Columns[0].Type))
	require.True(t, arrow.TypeEqual(arrow.BinaryTypes.String, assessor.gotTables[0].New.Columns[0].Type))
	require.Nil(t, assessor.gotTables[1].Old)

	want := &pb.AssessTables_Response{Tables: []*pb.AssessTables_TableFinding{{
		TableName:          "test_table",
		Category:           pb.AssessTables_CATEGORY_MANUAL_MIGRATION_REQUIRED,
		SafeModeBehavior:   "rejects this change",
		ForcedModeBehavior: "drops and recreates the table",
		Columns: []*pb.AssessTables_ColumnFinding{{
			ColumnName:         "tags",
			Category:           pb.AssessTables_CATEGORY_MANUAL_MIGRATION_REQUIRED,
			OldType:            "text[]",
			NewType:            "jsonb",
			SafeModeBehavior:   "rejects this change",
			ForcedModeBehavior: "drops and recreates the table",
			Evidence:           []*pb.AssessTables_Evidence{{SyntheticValue: `["env:prod"]`, Before: `{"tags":["env:prod"]}`, After: `{"tags":["env:prod"]}`}},
		}},
		Evidence:                 []*pb.AssessTables_Evidence{{SyntheticValue: "header", Before: "tags", After: "tags"}},
		CoverageIncomplete:       true,
		CoverageIncompleteReason: "nested values not compared",
	}}}
	require.Empty(t, cmp.Diff(want, resp, protocmp.Transform()))
}

func TestAssessTablesWithoutAssessorReturnsUnknown(t *testing.T) {
	client := newAssessClient(t, getColumnAdderPlugin())
	table := &schema.Table{Name: "test_table", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}

	resp, err := client.AssessTables(context.Background(), &pb.AssessTables_Request{
		Tables: []*pb.AssessTables_TablePair{{OldTable: tableBytes(t, table)}},
	})
	require.NoError(t, err)

	want := &pb.AssessTables_Response{Tables: []*pb.AssessTables_TableFinding{{
		TableName:                "test_table",
		Category:                 pb.AssessTables_CATEGORY_UNKNOWN,
		CoverageIncomplete:       true,
		CoverageIncompleteReason: plugin.AssessNotSupportedReason,
	}}}
	require.Empty(t, cmp.Diff(want, resp, protocmp.Transform()))
}

func TestAssessTablesRejectsEmptyPair(t *testing.T) {
	client := newAssessClient(t, getColumnAdderPlugin())
	_, err := client.AssessTables(context.Background(), &pb.AssessTables_Request{Tables: []*pb.AssessTables_TablePair{{}}})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAssessCategoryMatchesProto(t *testing.T) {
	for category, want := range map[plugin.AssessCategory]pb.AssessTables_Category{
		plugin.AssessCategoryUnknown:                 pb.AssessTables_CATEGORY_UNKNOWN,
		plugin.AssessCategoryNoChange:                pb.AssessTables_CATEGORY_NO_CHANGE,
		plugin.AssessCategoryAutomaticallyMigratable: pb.AssessTables_CATEGORY_AUTOMATICALLY_MIGRATABLE,
		plugin.AssessCategoryManualMigrationRequired: pb.AssessTables_CATEGORY_MANUAL_MIGRATION_REQUIRED,
		plugin.AssessCategoryTableRemoved:            pb.AssessTables_CATEGORY_TABLE_REMOVED,
		plugin.AssessCategoryFileSchemaChanged:       pb.AssessTables_CATEGORY_FILE_SCHEMA_CHANGED,
	} {
		require.Equal(t, want, tableFindingToPB(plugin.TableFinding{Category: category}).Category)
	}
	require.Len(t, pb.AssessTables_Category_name, 6)
}
