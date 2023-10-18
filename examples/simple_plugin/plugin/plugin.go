package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

// These variables are set as part of the `package` command
var (
	Team    = "example-team"
	Kind    = "source"
	Name    = "example"
	Version = "development"
)

func Plugin() *plugin.Plugin {
	return plugin.NewPlugin(
		Name,
		Version,
		Configure,
		plugin.WithKind(Kind),
		plugin.WithTeam(Team),
	)
}
