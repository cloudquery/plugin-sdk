package plugins

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/specs"
)



func TestDestinationPlugin(t *testing.T) {
	p := NewDestinationPlugin("test", "development", NewTestDestinationMemDBClient)
	DestinationPluginTestSuiteRunner(t, p, specs.Destination{
		WriteMode: specs.WriteModeAppend,
	})
}
