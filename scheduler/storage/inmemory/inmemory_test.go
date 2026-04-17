package inmemory_test

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage/inmemory"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage/storagetest"
)

func TestInMemory_Contract(t *testing.T) {
	storagetest.TestContract(t, func(t *testing.T) storage.Storage {
		return inmemory.New(1)
	})
}
