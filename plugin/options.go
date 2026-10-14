package plugin

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

type Option func(*Plugin)

func WithBuildTargets(targets []BuildTarget) Option {
	return func(p *Plugin) {
		p.targets = targets
	}
}

func WithJSONSchema(schema string) Option {
	return func(p *Plugin) {
		p.schema = schema
	}
}

func WithKind(kind string) Option {
	k := Kind(kind)
	err := k.Validate()
	if err != nil {
		panic(err)
	}
	return func(p *Plugin) {
		p.kind = k
	}
}

func WithTeam(team string) Option {
	return func(p *Plugin) {
		p.team = team
	}
}

// WithConnectionTester can be specified by a plugin to enable explicit connection testing, given a spec.
func WithConnectionTester(tester ConnectionTester) Option {
	return func(p *Plugin) {
		p.testConnFn = tester
	}
}

// WithFlushSupport declares that the client's Write handles message.WriteFlush: it commits every message received
// before the flush, then calls Done. The writers in the writers package do this, so a client whose Write passes the
// message channel to one of them can set this option. A Write that wraps or filters the channel must forward every
// message.WriteFlush, or call its Done itself, or the sync waits for it forever.
//
// A flush ends every open batch, like reaching the batch size does. Do not set this option if ending a batch early
// loses or corrupts data, for example a file destination that writes the whole sync to a single object. A client whose
// support depends on its spec implements FlushSupporter.
//
// A flush fails if a message since the previous flush failed to write. The CLI then stops saving state for the sync,
// so that it never passes the lost rows, and reports the failure once the sync finishes. Writers that stop on a failed
// write, such as the mixed batch writer, still end the sync. Without flush support, the batch writer only logs failed
// batches.
func WithFlushSupport() Option {
	return func(p *Plugin) {
		p.supportsFlush = true
	}
}

type TableOptions struct {
	Tables              []string
	SkipTables          []string
	SkipDependentTables bool
}
