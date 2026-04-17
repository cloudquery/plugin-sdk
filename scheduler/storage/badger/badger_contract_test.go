package badger_test

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	bstore "github.com/cloudquery/plugin-sdk/v4/scheduler/storage/badger"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage/storagetest"
)

func TestBadger_Contract(t *testing.T) {
	storagetest.TestContract(t, func(t *testing.T) storage.Storage {
		dir := t.TempDir()
		s, err := bstore.Open(bstore.Options{Path: dir})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = s.Close(context.Background()) })
		return s
	})
}
