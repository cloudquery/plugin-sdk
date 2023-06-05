package plugin

import (
	"fmt"

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

type Option func(*Plugin)


// WithNoInternalColumns won't add internal columns (_cq_id, _cq_parent_cq_id) to the plugin tables
func WithNoInternalColumns() Option {
	return func(p *Plugin) {
		p.internalColumns = false
	}
}

// WithTitleTransformer allows the plugin to control how table names get turned into titles for the
// generated documentation.
func WithTitleTransformer(t func(*schema.Table) string) Option {
	return func(p *Plugin) {
		p.titleTransformer = t
	}
}