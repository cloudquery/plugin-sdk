package plugin

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type MigrateMode int

const (
	MigrateModeSafe MigrateMode = iota
	MigrateModeForce
)

var (
	migrateModeStrings = []string{"safe", "force"}
)

func (m MigrateMode) String() string {
	return migrateModeStrings[m]
}

type Registry int

const (
	RegistryGithub Registry = iota
	RegistryLocal
	RegistryGrpc
)

func (r Registry) String() string {
	return [...]string{"github", "local", "grpc"}[r]
}

func RegistryFromString(s string) (Registry, error) {
	switch s {
	case "github":
		return RegistryGithub, nil
	case "local":
		return RegistryLocal, nil
	case "grpc":
		return RegistryGrpc, nil
	default:
		return RegistryGithub, fmt.Errorf("unknown registry %s", s)
	}
}

type WriteMode int

const (
	WriteModeOverwriteDeleteStale WriteMode = iota
	WriteModeOverwrite
	WriteModeAppend
)

var (
	writeModeStrings = []string{"overwrite-delete-stale", "overwrite", "append"}
)

func (m WriteMode) String() string {
	return writeModeStrings[m]
}

type Scheduler int

const (
	SchedulerDFS Scheduler = iota
	SchedulerRoundRobin
)

var AllSchedulers = Schedulers{SchedulerDFS, SchedulerRoundRobin}
var AllSchedulerNames = [...]string{
	SchedulerDFS:        "dfs",
	SchedulerRoundRobin: "round-robin",
}

type Schedulers []Scheduler

func (s Schedulers) String() string {
	var buffer bytes.Buffer
	for i, scheduler := range s {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(scheduler.String())
	}
	return buffer.String()
}

func (s Scheduler) String() string {
	return AllSchedulerNames[s]
}

type GetTables func(ctx context.Context, c Client) (schema.Tables, error)

type Option func(*Plugin)

// WithDynamicTable allows the plugin to return list of tables after call to New
func WithDynamicTable(getDynamicTables GetTables) Option {
	return func(p *Plugin) {
		p.getDynamicTables = getDynamicTables
	}
}

// WithNoInternalColumns won't add internal columns (_cq_id, _cq_parent_cq_id) to the plugin tables
func WithNoInternalColumns() Option {
	return func(p *Plugin) {
		p.internalColumns = false
	}
}

func WithUnmanagedSync() Option {
	return func(p *Plugin) {
		p.unmanagedSync = true
	}
}

// WithTitleTransformer allows the plugin to control how table names get turned into titles for the
// generated documentation.
func WithTitleTransformer(t func(*schema.Table) string) Option {
	return func(p *Plugin) {
		p.titleTransformer = t
	}
}

func WithStaticTables(tables schema.Tables) Option {
	return func(p *Plugin) {
		p.staticTables = tables
	}
}

func WithManagedWriter() Option {
	return func(p *Plugin) {
		p.managedWriter = true
	}
}

func WithBatchTimeout(seconds int) Option {
	return func(p *Plugin) {
		p.batchTimeout = time.Duration(seconds) * time.Second
	}
}

func WithDefaultBatchSize(defaultBatchSize int) Option {
	return func(p *Plugin) {
		p.defaultBatchSize = defaultBatchSize
	}
}

func WithDefaultBatchSizeBytes(defaultBatchSizeBytes int) Option {
	return func(p *Plugin) {
		p.defaultBatchSizeBytes = defaultBatchSizeBytes
	}
}
