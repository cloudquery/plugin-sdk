package queue

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const (
	resourceCount = 10
)

type Data struct {
	Name string `json:"name"`
}
type testClient struct {
	id string
}

func (tc *testClient) ID() string {
	return tc.id
}

func testResolver(_ context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	resources := make([]*Data, 0, resourceCount)
	for i := 0; i < resourceCount; i++ {
		if parent == nil {
			resources = append(resources, &Data{Name: fmt.Sprintf("test-%d", i)})
		} else {
			item := parent.Item.(*Data)
			resources = append(resources, &Data{Name: fmt.Sprintf("%s-test-%d", item.Name, i)})
		}
	}
	res <- resources
	return nil
}

func TestScheduler(t *testing.T) {
	nopLogger := zerolog.Nop()
	m := metrics.NewMetrics()
	scheduler := NewShuffleQueueScheduler(nopLogger, m, int64(0), WithWorkerCount(1000), WithInvocationID(uuid.New().String()))
	tableClients := []WorkUnit{
		{
			Table: &schema.Table{
				Name:      "table-1",
				Resolver:  testResolver,
				Transform: transformers.TransformWithStruct(&Data{}),
				Relations: schema.Tables{
					{
						Name:      "table-1-relation-1",
						Resolver:  testResolver,
						Transform: transformers.TransformWithStruct(&Data{}),
					},
				},
			},
			Client: &testClient{id: "client-1"},
		},
		{
			Table: &schema.Table{
				Name:      "table-2",
				Resolver:  testResolver,
				Transform: transformers.TransformWithStruct(&Data{}),
				Relations: schema.Tables{
					{
						Name:      "table-2-relation-1",
						Resolver:  testResolver,
						Transform: transformers.TransformWithStruct(&Data{}),
					},
				},
			},
			Client: &testClient{id: "client-1"},
		},
	}

	for _, tc := range tableClients {
		m.InitWithClients(tc.Table, []schema.ClientMeta{tc.Client})
	}

	resolvedResources := make(chan *schema.Resource)
	msgs := make(chan message.SyncMessage, 10)
	go func() {
		defer close(resolvedResources)
		scheduler.Sync(context.Background(), tableClients, resolvedResources, msgs)
	}()

	gotResources := make([]*schema.Resource, 0)
	for r := range resolvedResources {
		gotResources = append(gotResources, r)
	}

	// 2 top level tables, each with 10 resources, and each top level has a relation with 10 resources
	require.Len(t, gotResources, 2*10+2*10*10)
}

type pushOnFirstEmptyPopQueue struct {
	*ConcurrentRandomQueue[WorkUnit]
	onFirstEmptyPop func()
}

func (q *pushOnFirstEmptyPopQueue) Pop() *WorkUnit {
	item := q.ConcurrentRandomQueue.Pop()
	if item == nil && q.onFirstEmptyPop != nil {
		onEmpty := q.onFirstEmptyPop
		q.onFirstEmptyPop = nil
		onEmpty()
	}
	return item
}

func TestDispatchHandsOffWorkQueuedByLastActiveUnit(t *testing.T) {
	activeWorkSignal := newActiveWorkSignal()
	queue := &pushOnFirstEmptyPopQueue{ConcurrentRandomQueue: NewConcurrentRandomQueue[WorkUnit](0, 1)}
	queue.Push(WorkUnit{Table: &schema.Table{Name: "parent"}})
	// The parent's worker queues a relation and finishes after the dispatcher saw an
	// empty queue but before it checked for active work.
	queue.onFirstEmptyPop = func() {
		queue.Push(WorkUnit{Table: &schema.Table{Name: "relation"}})
		activeWorkSignal.Done()
	}

	jobs := make(chan *WorkUnit)
	go dispatch(context.Background(), queue, jobs, activeWorkSignal)

	require.Equal(t, "parent", (<-jobs).Table.Name)
	relation, ok := <-jobs
	require.True(t, ok, "dispatcher stopped without handing off the queued relation")
	require.Equal(t, "relation", relation.Table.Name)
	activeWorkSignal.Done()

	_, ok = <-jobs
	require.False(t, ok)
}
