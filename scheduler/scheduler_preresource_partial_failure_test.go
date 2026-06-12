package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestSchedulerPreResourceResolverPartialFailureWithRelations verifies that a
// PreResourceResolver failure only drops the failing resource and its own
// subtree, at every level of a parent -> child1 -> child2 hierarchy:
//   - parent "p2" fails: p2 and all its descendants are dropped, the other
//     parents and their descendants survive
//   - child1 "p1/c0" fails: only that row and its child2 descendants are
//     dropped, sibling "p1/c1" and its descendants survive
func TestSchedulerPreResourceResolverPartialFailureWithRelations(t *testing.T) {
	for _, strategy := range AllStrategies {
		t.Run(strategy.String(), func(t *testing.T) {
			nameColumn := schema.Column{
				Name: "name",
				Type: arrow.BinaryTypes.String,
				Resolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
					return resource.Set(c.Name, resource.Item.(string))
				},
			}

			child2 := &schema.Table{
				Name: "test_child2",
				Resolver: func(_ context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
					res <- []string{parent.Item.(string) + "/g0"}
					return nil
				},
				Columns: []schema.Column{nameColumn},
			}

			child1 := &schema.Table{
				Name: "test_child1",
				Resolver: func(_ context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
					p := parent.Item.(string)
					res <- []string{p + "/c0", p + "/c1"}
					return nil
				},
				PreResourceResolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource) error {
					if resource.Item.(string) == "p1/c0" {
						return errors.New("child1 pre resource resolver boom")
					}
					return nil
				},
				Columns:   []schema.Column{nameColumn},
				Relations: schema.Tables{child2},
			}

			parentTable := &schema.Table{
				Name: "test_parent",
				Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
					// a single slice so all parents land in one chunk, like one API page
					res <- []string{"p0", "p1", "p2", "p3", "p4"}
					return nil
				},
				PreResourceResolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource) error {
					if resource.Item.(string) == "p2" {
						return errors.New("parent pre resource resolver boom")
					}
					return nil
				},
				Columns:   []schema.Column{nameColumn},
				Relations: schema.Tables{child1},
			}

			tables := schema.Tables{parentTable}
			c := testExecutionClient{}
			sc := NewScheduler(
				WithLogger(zerolog.New(zerolog.NewTestWriter(t)).Level(zerolog.DebugLevel)),
				WithStrategy(strategy),
			)
			msgs := make(chan message.SyncMessage, 500)
			require.NoError(t, sc.Sync(context.Background(), &c, tables, msgs))
			close(msgs)

			var messages message.SyncMessages
			for msg := range msgs {
				messages = append(messages, msg)
			}

			collect := func(tb *schema.Table) []string {
				values := []string{}
				for _, rec := range messages.GetInserts().GetRecordsForTable(tb) {
					idx := rec.Schema().FieldIndices("name")[0]
					col := rec.Column(idx).(*array.String)
					for i := 0; i < col.Len(); i++ {
						values = append(values, col.Value(i))
					}
				}
				return values
			}

			require.ElementsMatch(t,
				[]string{"p0", "p1", "p3", "p4"},
				collect(parentTable),
				"only the failing parent should be dropped")

			require.ElementsMatch(t,
				[]string{"p0/c0", "p0/c1", "p1/c1", "p3/c0", "p3/c1", "p4/c0", "p4/c1"},
				collect(child1),
				"children of surviving parents should sync, except the failing child; no children of the dropped parent")

			require.ElementsMatch(t,
				[]string{"p0/c0/g0", "p0/c1/g0", "p1/c1/g0", "p3/c0/g0", "p3/c1/g0", "p4/c0/g0", "p4/c1/g0"},
				collect(child2),
				"grandchildren should only sync under surviving child1 rows")
		})
	}
}
