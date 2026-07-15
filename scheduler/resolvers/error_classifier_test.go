package resolvers

import (
	"context"
	"errors"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var errResolverBoom = errors.New("resolver boom")

// TestResolveResourcesChunk_ErrorClassifier verifies that the error classifier controls
// whether resolver errors at each phase are counted as errors, and that it receives an
// ErrorEvent describing the phase that failed.
func TestResolveResourcesChunk_ErrorClassifier(t *testing.T) {
	for _, tc := range []struct {
		name       string
		table      func() *schema.Table
		wantPhase  schema.ErrorPhase
		wantColumn string
	}{
		{
			name: "column resolver",
			table: func() *schema.Table {
				return &schema.Table{
					Name: "test_table",
					Columns: []schema.Column{
						{
							Name: "test_column",
							Type: arrow.PrimitiveTypes.Int64,
							Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, _ schema.Column) error {
								return errResolverBoom
							},
						},
					},
				}
			},
			wantPhase:  schema.ErrorPhaseColumnResolver,
			wantColumn: "test_column",
		},
		{
			name: "pre resource chunk resolver",
			table: func() *schema.Table {
				return &schema.Table{
					Name: "test_table",
					PreResourceChunkResolver: &schema.RowsChunkResolver{
						ChunkSize: 10,
						RowsResolver: func(_ context.Context, _ schema.ClientMeta, _ []*schema.Resource) error {
							return errResolverBoom
						},
					},
					Columns: []schema.Column{{Name: "test_column", Type: arrow.PrimitiveTypes.Int64}},
				}
			},
			wantPhase: schema.ErrorPhasePreResourceChunkResolver,
		},
		{
			name: "pre resource resolver",
			table: func() *schema.Table {
				return &schema.Table{
					Name: "test_table",
					PreResourceResolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource) error {
						return errResolverBoom
					},
					Columns: []schema.Column{{Name: "test_column", Type: arrow.PrimitiveTypes.Int64}},
				}
			},
			wantPhase: schema.ErrorPhasePreResourceResolver,
		},
		{
			name: "post resource resolver",
			table: func() *schema.Table {
				return &schema.Table{
					Name: "test_table",
					PostResourceResolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource) error {
						return errResolverBoom
					},
					Columns: []schema.Column{{Name: "test_column", Type: arrow.PrimitiveTypes.Int64}},
				}
			},
			wantPhase: schema.ErrorPhasePostResourceResolver,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := testClient{}
			chunk := []any{0}
			logger := zerolog.New(zerolog.NewTestWriter(t))

			t.Run("raised without classifier", func(t *testing.T) {
				table := tc.table()
				m := metrics.NewMetrics()
				m.InitWithClients(table, []schema.ClientMeta{client})
				ResolveResourcesChunk(context.Background(), logger, m, table, client, nil, chunk, caser.New(), nil)
				require.Equal(t, uint64(1), m.GetErrors(m.NewSelector(client.ID(), table.Name)))
			})

			t.Run("suppressed by classifier", func(t *testing.T) {
				table := tc.table()
				m := metrics.NewMetrics()
				m.InitWithClients(table, []schema.ClientMeta{client})
				var gotEvent schema.ErrorEvent
				classifier := func(_ context.Context, _ error, event schema.ErrorEvent) bool {
					gotEvent = event
					return true
				}
				ResolveResourcesChunk(context.Background(), logger, m, table, client, nil, chunk, caser.New(), classifier)
				require.Equal(t, uint64(0), m.GetErrors(m.NewSelector(client.ID(), table.Name)), "suppressed error should not be counted")
				require.Equal(t, tc.wantPhase, gotEvent.Phase)
				require.Equal(t, table, gotEvent.Table)
				if tc.wantColumn != "" {
					require.NotNil(t, gotEvent.Column)
					require.Equal(t, tc.wantColumn, gotEvent.Column.Name)
				} else {
					require.Nil(t, gotEvent.Column)
				}
			})
		})
	}
}
