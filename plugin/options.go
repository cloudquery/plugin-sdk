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

func WithTitle(title string) Option {
	return func(p *Plugin) {
		p.title = title
	}
}

func WithDescription(description string) Option {
	return func(p *Plugin) {
		p.description = description
	}
}

func WithShortDescription(shortDescription string) Option {
	return func(p *Plugin) {
		p.shortDescription = shortDescription
	}
}

type TableOptions struct {
	Tables              []string
	SkipTables          []string
	SkipDependentTables bool
}
