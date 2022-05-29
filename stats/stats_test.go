package stats

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/segmentio/stats/v4"
	"github.com/stretchr/testify/assert"
)

func Test_getMeasurementDetails(t *testing.T) {
	type args struct {
		name string
		tags []stats.Tag
	}
	type want struct {
		id    string
		stamp bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Non stamped clock",
			args: args{name: "measurement", tags: []stats.Tag{{Name: "table", Value: "name_of_table"}}},
			want: want{id: "measurement:table:name_of_table", stamp: false},
		},
		{
			name: "Stamped clock",
			args: args{name: "measurement", tags: []stats.Tag{{Name: "table", Value: "name_of_table"}, {Name: "stamp", Value: "total"}}},
			want: want{id: "measurement:table:name_of_table", stamp: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, stamp := getMeasurementDetails(tt.args.name, tt.args.tags)
			assert.EqualValues(t, tt.want.id, id)
			assert.EqualValues(t, tt.want.stamp, stamp)
		})
	}
}

func TestLogHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	logger := hclog.NewNullLogger()
	handler := newHandler(logger)
	Start(ctx, logger, WithTick(time.Hour), WithHandler(handler))

	clock1 := NewClockWithObserve("withStop", stats.Tag{Name: "table", Value: "table1"})
	NewClockWithObserve("withoutStop", stats.Tag{Name: "table", Value: "table2"})
	Flush()

	logHandler := handler.(*durationLogger)
	assert.Len(t, logHandler.trackedOperations.Keys(), 2)
	assert.EqualValues(t, "withStop:table:table1", logHandler.trackedOperations.Keys()[0])
	assert.EqualValues(t, "withoutStop:table:table2", logHandler.trackedOperations.Keys()[1])
	clock1.Stop()

	Flush()

	assert.Len(t, logHandler.trackedOperations.Keys(), 1)
	assert.EqualValues(t, "withoutStop:table:table2", logHandler.trackedOperations.Keys()[0])

	cancel()
}
