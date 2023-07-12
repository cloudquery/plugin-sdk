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

type TableOptions struct {
	Tables              []string
	SkipTables          []string
	SkipDependentTables bool
}
