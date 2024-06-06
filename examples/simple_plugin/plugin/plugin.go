package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

// These variables are used by the `package` command and for checking
// usage limits.
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
		plugin.WithConnectionTester(TestConnection),
	)
}
