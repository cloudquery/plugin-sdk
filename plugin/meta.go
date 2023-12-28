package plugin

import cqapi "github.com/cloudquery/cloudquery-api-go"

type Meta struct {
	Team            cqapi.PluginTeam
	Kind            cqapi.PluginKind
	Name            cqapi.PluginName
	SkipUsageClient bool
}
