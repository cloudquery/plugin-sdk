package plugins

import (
	"testing"
)

func TestDestinationPlugin(t *testing.T) {
	p := NewDestinationPlugin("test", "development", NewTestDestinationMemDBClient)
	DestinationPluginTestSuiteRunner(t, p, nil,
		DestinationTestSuiteTests{
			Overwrite:   true,
			DeleteStale: true,
			Append:      true,
		})
}
