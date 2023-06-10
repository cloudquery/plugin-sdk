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

type Option func(*Plugin)
