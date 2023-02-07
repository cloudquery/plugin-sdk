package postgresql

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

func getTestConnection() string {
	testConn := os.Getenv("CQ_PG_TEST_CONN")
	if testConn == "" {
		return "postgresql://postgres:pass@localhost:5432/postgres?sslmode=disable"
	}
	return testConn
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func TestBackend(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	spec := specs.Source{
		Name:    "test_" + t.Name() + "_" + randSeq(10),
		Version: "v1",
		Path:    "/tmp/test",
		Backend: specs.BackendPostgres,
		BackendSpec: Spec{
			ConnectionString: getTestConnection(),
			TableName:        "cloudquery_state_test",
		},
	}
	b, err := New(ctx, zerolog.Logger{}, spec)
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}
	if b == nil {
		t.Fatalf("expected backend to be not nil")
	}

	tableName := "test_table"
	clientID := "test_client"
	got, err := b.Get(ctx, tableName, clientID)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value, but got %s", got)
	}

	err = b.Set(ctx, tableName, clientID, "test_value_to_overwrite")
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	// Set again with different value, should not error
	err = b.Set(ctx, tableName, clientID, "test_value")
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	got, err = b.Get(ctx, tableName, clientID)
	if err != nil {
		t.Fatalf("failed to get value after setting it: %v", err)
	}
	if got != "test_value" {
		t.Fatalf("expected value to be test_value, but got %s", got)
	}

	err = b.Close(ctx)
	if err != nil {
		t.Fatalf("failed to close local backend: %v", err)
	}

	b, err = New(ctx, zerolog.Logger{}, spec)
	if err != nil {
		t.Fatalf("failed to open local backend after closing it: %v", err)
	}

	got, err = b.Get(ctx, tableName, clientID)
	if err != nil {
		t.Fatalf("failed to get value after closing and reopening local backend: %v", err)
	}
	if got != "test_value" {
		t.Fatalf("expected value to be test_value, but got %s", got)
	}

	got, err = b.Get(ctx, "some_other_table", clientID)
	if err != nil {
		t.Fatalf("failed to get value after closing and reopening local backend: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value for some_other_table -> test_key, but got %s", got)
	}
	err = b.Close(ctx)
	if err != nil {
		t.Fatalf("failed to close local backend the second time: %v", err)
	}

	// check that state is namespaced by source name
	spec.Name = "test2"
	local2, err := New(ctx, zerolog.Logger{}, spec)
	if err != nil {
		t.Fatalf("failed to create local backend for test2: %v", err)
	}

	got, err = local2.Get(ctx, "test_table", clientID)
	if err != nil {
		t.Fatalf("failed to get value for local backend test2: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value for test2 -> test_table -> test_key, but got %s", got)
	}
	err = local2.Close(ctx)
	if err != nil {
		t.Fatalf("failed to close second local backend: %v", err)
	}
}
