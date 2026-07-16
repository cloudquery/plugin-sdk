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

type testClient struct{}

func (testClient) ID() string { return "test" }

var _ schema.ClientMeta = testClient{}

// TestResolveResourcesChunk_PreResourceResolverPartialFailure verifies that a
// PreResourceResolver error on a single resource only drops that resource from
// the batch, while the remaining resources are still resolved and returned.
func TestResolveResourcesChunk_PreResourceResolverPartialFailure(t *testing.T) {
	for _, tc := range []struct {
		name          string
		failItems     map[int]bool
		expectedItems []int
	}{
		{
			name:          "no failures keeps all resources",
			failItems:     nil,
			expectedItems: []int{0, 1, 2, 3, 4},
		},
		{
			name:          "single failure drops only that resource",
			failItems:     map[int]bool{2: true},
			expectedItems: []int{0, 1, 3, 4},
		},
		{
			name:          "multiple failures drop only the failing resources",
			failItems:     map[int]bool{0: true, 3: true},
			expectedItems: []int{1, 2, 4},
		},
		{
			name:          "all failures drop the whole batch but do not panic",
			failItems:     map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true},
			expectedItems: []int{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			table := &schema.Table{
				Name: "test_table",
				PreResourceResolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource) error {
					if tc.failItems[resource.Item.(int)] {
						return errors.New("pre resource resolver boom")
					}
					return nil
				},
				Columns: []schema.Column{
					{
						Name: "test_column",
						Type: arrow.PrimitiveTypes.Int64,
						Resolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
							return resource.Set(c.Name, int64(resource.Item.(int)))
						},
					},
				},
			}

			client := testClient{}
			m := metrics.NewMetrics()
			m.InitWithClients(table, []schema.ClientMeta{client})

			chunk := []any{0, 1, 2, 3, 4}
			logger := zerolog.New(zerolog.NewTestWriter(t))

			resources := ResolveResourcesChunk(context.Background(), logger, m, table, client, nil, chunk, caser.New())

			gotItems := make([]int, len(resources))
			for i, r := range resources {
				gotItems[i] = r.Item.(int)
				// surviving resources should have been fully resolved through the column resolvers
				col := r.Get("test_column")
				require.True(t, col.IsValid(), "surviving resource should have its column resolved")
				require.Equal(t, int64(r.Item.(int)), col.Get(), "resolved column value should match the item")
			}
			require.ElementsMatch(t, tc.expectedItems, gotItems)

			selector := m.NewSelector(client.ID(), table.Name)
			require.Equal(t, uint64(len(tc.failItems)), m.GetErrors(selector), "expected one error per failing resource")
			require.Equal(t, uint64(len(tc.expectedItems)), m.GetResources(selector), "only surviving resources should be counted")
		})
	}
}

// TestResolveResourcesChunk_PreResourceResolverContextCancelled verifies that
// once the context is cancelled, the chunk is dropped immediately with a single
// error instead of emitting one error per remaining resource.
func TestResolveResourcesChunk_PreResourceResolverContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := 0
	table := &schema.Table{
		Name: "test_table",
		PreResourceResolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource) error {
			calls++
			if calls == 2 {
				cancel()
				return errors.New("pre resource resolver boom")
			}
			return nil
		},
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: arrow.PrimitiveTypes.Int64,
				Resolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
					return resource.Set(c.Name, int64(resource.Item.(int)))
				},
			},
		},
	}

	client := testClient{}
	m := metrics.NewMetrics()
	m.InitWithClients(table, []schema.ClientMeta{client})

	chunk := []any{0, 1, 2, 3, 4}
	logger := zerolog.New(zerolog.NewTestWriter(t))

	resources := ResolveResourcesChunk(ctx, logger, m, table, client, nil, chunk, caser.New())

	require.Empty(t, resources, "cancelled chunk should not return resources")
	require.Equal(t, 2, calls, "resolver should not be called for resources after cancellation")

	selector := m.NewSelector(client.ID(), table.Name)
	require.Equal(t, uint64(1), m.GetErrors(selector), "cancellation should be counted as a single error, not one per remaining resource")
	require.Equal(t, uint64(0), m.GetResources(selector), "no resources should be counted for a cancelled chunk")
}
