package testdata

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenTestData(t *testing.T) {
	sch := TestSourceSchema("test", TestSourceOptions{})
	records := GenTestData(sch, GenTestDataOptions{
		SourceName: "test",
		SyncTime:   time.Now().UTC().Round(1 * time.Second),
		MaxRows:    2,
		StableUUID: uuid.New(),
		StableTime: time.Now().UTC().Round(1 * time.Second),
	})
	if len(records) == 0 {
		t.Fatal("no records")
	}
}
