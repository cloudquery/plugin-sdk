package tablefinishlogger

import (
	"os"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type testLogMessage struct {
	table       string
	client      string
	parentTable string
}

func TestTableFinishLogger(t *testing.T) {
	logger := zerolog.New(os.Stdout)
	logHook := &logHook{}
	logger.Hook(logHook)

	finishLogger := New(logger)

	// Create test tables
	parentTable := &schema.Table{Name: "parent"}
	child1Table := &schema.Table{Name: "child1"}
	child2Table := &schema.Table{Name: "child2"}

	// Create test client
	testClient := &testClient{id: "test_client"}

	// Create parent resource
	parentResource := &schema.Resource{Table: parentTable}

	tests := []struct {
		name           string
		steps          func()
		wantLogCount   int
		wantLastTable  string
		wantLastParent string
	}{
		{
			name: "parent table finishes",
			steps: func() {
				finishLogger.TableFinished(parentTable, testClient, nil)
			},
			wantLogCount:  1,
			wantLastTable: "parent",
		},
		{
			name: "child table with unfinished parent",
			steps: func() {
				finishLogger.TableStarted(child1Table, parentResource)
				finishLogger.TableFinished(child1Table, testClient, parentResource)
			},
			wantLogCount: 1, // still just the parent from before
		},
		{
			name: "multiple child instances finish after parent",
			steps: func() {
				// Parent finishes
				finishLogger.TableFinished(parentTable, testClient, nil)

				// First child instance
				finishLogger.TableStarted(child1Table, parentResource)
				finishLogger.TableFinished(child1Table, testClient, parentResource)

				// Second child instance
				finishLogger.TableStarted(child1Table, parentResource)
				finishLogger.TableFinished(child1Table, testClient, parentResource)
			},
			wantLogCount:   3,
			wantLastTable:  "child1",
			wantLastParent: "parent",
		},
		{
			name: "multiple child tables",
			steps: func() {
				// Start and finish instances of child2
				finishLogger.TableStarted(child2Table, parentResource)
				finishLogger.TableStarted(child2Table, parentResource)
				finishLogger.TableFinished(child2Table, testClient, parentResource)
				finishLogger.TableFinished(child2Table, testClient, parentResource)
			},
			wantLogCount:   4,
			wantLastTable:  "child2",
			wantLastParent: "parent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logHook.logEvents = nil // Reset messages
			tt.steps()

			if got := len(logHook.logEvents); got != tt.wantLogCount {
				t.Errorf("got %d log messages, want %d", got, tt.wantLogCount)
			}

			// if tt.wantLastTable != "" && len(logHook.logEvents) > 0 {
			// 	lastMsg := logHook.logEvents[len(logHook.logEvents)-1]
			// 	if lastMsg.table != tt.wantLastTable {
			// 		t.Errorf("got last table %q, want %q", lastMsg.table, tt.wantLastTable)
			// 	}
			// 	if tt.wantLastParent != "" && lastMsg.parentTable != tt.wantLastParent {
			// 		t.Errorf("got last parent table %q, want %q", lastMsg.parentTable, tt.wantLastParent)
			// 	}
			// }
		})
	}
}

// Mock client for testing
type testClient struct {
	id string
}

func (c *testClient) ID() string {
	return c.id
}

type logHook struct {
	logEvents []zerolog.Event
}

func (logHook *logHook) Run(logEvent *zerolog.Event, level zerolog.Level, message string) {
	// logEvent
	logHook.logEvents = append(logHook.logEvents, *logEvent)
}
