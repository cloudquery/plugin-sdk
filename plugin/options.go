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

type TableOptions struct {
	Tables              []string
	SkipTables          []string
	SkipDependentTables bool
}
