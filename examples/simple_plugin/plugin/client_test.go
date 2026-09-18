package plugin

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

const eventSyncSpec = `{
  "num_clients": 1,
  "changes": [
    {"table": "test_some_table", "keys": {"column1": "a", "client_id": 0}, "op": "upsert", "event_id": "e1"},
    {"table": "test_some_table", "keys": {"column1": "b", "client_id": 0}, "op": "delete", "event_id": "e2"}
  ],
  "refresh_data": {
    "test_some_table": [
      {"column1": "a", "column2": 7, "client_id": 0}
    ]
  }
}`

// syncMessages runs a sync on the client and collects everything it emits.
func syncMessages(ctx context.Context, c *Client, options plugin.SyncOptions) (message.SyncMessages, error) {
	res := make(chan message.SyncMessage)
	var syncErr error
	go func() {
		defer close(res)
		syncErr = c.Sync(ctx, options, res)
	}()
	var msgs message.SyncMessages
	for msg := range res {
		msgs = append(msgs, msg)
	}
	return msgs, syncErr
}

func newEventSyncClient(t *testing.T) *Client {
	t.Helper()
	cl, err := Configure(context.Background(), zerolog.Nop(), []byte(eventSyncSpec), plugin.NewClientOptions{})
	if err != nil {
		t.Fatalf("failed to configure the client: %v", err)
	}
	return cl.(*Client)
}

func TestWatchReplaysConfiguredChanges(t *testing.T) {
	c := newEventSyncClient(t)
	ctx := context.Background()

	// No backend options: the state client falls back to the no-op client, so
	// there is no checkpoint and every change is replayed.
	ready, err := c.WatchReady(ctx, plugin.WatchOptions{Tables: []string{"*"}, Drain: true})
	if err != nil {
		t.Fatalf("WatchReady failed: %v", err)
	}
	if ready.HasCheckpoint {
		t.Fatal("expected no checkpoint without a state backend")
	}

	changes := make(chan plugin.Change, 8)
	if err := c.Watch(ctx, plugin.Position{}, changes); err != nil {
		t.Fatalf("Watch failed: %v", err)
	}
	close(changes)

	var got []plugin.Change
	for change := range changes {
		got = append(got, change)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}

	if got[0].EventID != "e1" || got[0].Op != plugin.OpUpsert {
		t.Errorf("unexpected first change: %+v", got[0])
	}
	if keys := schema.KeysString(got[0].Keys); keys != "column1=a,client_id=0" {
		t.Errorf("unexpected keys on the first change: %s", keys)
	}
	if got[1].EventID != "e2" || got[1].Op != plugin.OpDelete {
		t.Errorf("unexpected second change: %+v", got[1])
	}

	if err := c.Ack(ctx, "e1"); err != nil {
		t.Fatalf("Ack failed: %v", err)
	}
}

func TestSyncRefreshFindsAndMissesRows(t *testing.T) {
	c := newEventSyncClient(t)
	ctx := context.Background()

	table := c.tables.Get("test_some_table")
	if table == nil {
		t.Fatal("test_some_table is missing")
	}

	found, err := schema.NewKeysRecord(table, map[string]any{"column1": "a", "client_id": 0})
	if err != nil {
		t.Fatalf("failed to build keys: %v", err)
	}
	missing, err := schema.NewKeysRecord(table, map[string]any{"column1": "zz", "client_id": 0})
	if err != nil {
		t.Fatalf("failed to build keys: %v", err)
	}

	msgs, err := syncMessages(ctx, c, plugin.SyncOptions{Refresh: []plugin.RefreshResource{
		{Table: "test_some_table", Keys: found, Op: plugin.OpUpsert},
		{Table: "test_some_table", Keys: missing, Op: plugin.OpUpsert},
	}})
	if err != nil {
		t.Fatalf("refresh sync failed: %v", err)
	}

	var inserts, deletes int
	for _, msg := range msgs {
		switch msg.(type) {
		case *message.SyncInsert:
			inserts++
		case *message.SyncDeleteRecord:
			deletes++
		case *message.SyncMigrateTable:
			t.Fatal("a refresh sync must not emit migrate messages")
		}
	}
	if inserts != 1 {
		t.Errorf("expected 1 insert, got %d", inserts)
	}
	if deletes != 1 {
		t.Errorf("expected 1 delete record, got %d", deletes)
	}
}

func TestFullSyncStillWorks(t *testing.T) {
	c := newEventSyncClient(t)

	msgs, err := syncMessages(context.Background(), c, plugin.SyncOptions{Tables: []string{"*"}})
	if err != nil {
		t.Fatalf("full sync failed: %v", err)
	}
	if items := msgs.InsertItems(); items != 1 {
		t.Errorf("expected 1 row from the full sync, got %d", items)
	}
}

func TestSpecValidation(t *testing.T) {
	_, err := Configure(context.Background(), zerolog.Nop(), []byte(`{"changes":[{"table":"test_some_table"}]}`), plugin.NewClientOptions{})
	if err == nil || !strings.Contains(err.Error(), "event_id is required") {
		t.Fatalf("expected an event_id validation error, got %v", err)
	}
}

// TestWatchAndAckRunConcurrently exercises the WatchClient contract: Watch and
// Ack are called from different goroutines for the lifetime of the stream, so
// everything they share has to be guarded. Run it under -race.
func TestWatchAndAckRunConcurrently(t *testing.T) {
	c := newEventSyncClient(t)
	ctx := context.Background()

	if _, err := c.WatchReady(ctx, plugin.WatchOptions{Tables: []string{"*"}}); err != nil {
		t.Fatalf("WatchReady failed: %v", err)
	}

	// Several episodes so the two goroutines actually overlap rather than
	// happening to run one after the other.
	for range 50 {
		changes := make(chan plugin.Change)
		var (
			wg       sync.WaitGroup
			watchErr error
			ackErr   error
		)

		wg.Add(3)
		go func() {
			defer wg.Done()
			defer close(changes)
			watchErr = c.Watch(ctx, plugin.Position{}, changes)
		}()
		go func() {
			defer wg.Done()
			for range changes { //revive:disable-line:empty-block // the acker below drives the checkpoint
			}
		}()
		go func() {
			defer wg.Done()
			for _, eventID := range []string{"e1", "e2"} {
				if err := c.Ack(ctx, eventID); err != nil {
					ackErr = err
					return
				}
			}
		}()
		wg.Wait()

		if watchErr != nil {
			t.Fatalf("Watch failed: %v", watchErr)
		}
		if ackErr != nil {
			t.Fatalf("Ack failed: %v", ackErr)
		}
	}

	if checkpoint := c.watchCheckpoint(); checkpoint != "e2" {
		t.Errorf("expected the checkpoint to be e2, got %q", checkpoint)
	}
}

// TestCloseIsSafeDuringAcks closes the client while acks are still in flight,
// which is what happens when a follow-mode sync is cancelled.
func TestCloseIsSafeDuringAcks(t *testing.T) {
	c := newEventSyncClient(t)
	ctx := context.Background()

	if _, err := c.WatchReady(ctx, plugin.WatchOptions{Tables: []string{"*"}}); err != nil {
		t.Fatalf("WatchReady failed: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// An ack that lands after Close is rejected, not a panic or a race.
		for range 100 {
			if err := c.Ack(ctx, "e1"); err != nil && !strings.Contains(err.Error(), "watch not started") {
				t.Errorf("unexpected ack error: %v", err)
				return
			}
		}
	}()

	if err := c.Close(ctx); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	wg.Wait()

	// Closing twice is a no-op.
	if err := c.Close(ctx); err != nil {
		t.Fatalf("second Close failed: %v", err)
	}
}
