package schema

import (
	"testing"
	"time"
)

func TestTimestampTz(t *testing.T) {
	v := &Timestamptz{}
	if err := v.Scan(time.Now()); err != nil {
		t.Fatal(err)
	}
}
