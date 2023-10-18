package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

// These variables are set as part of the `package` command
var (
	Name    = ""
	Kind    = "source"
	Team    = ""
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
