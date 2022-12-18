package destination

import "testing"

func TestPluginUnmanagedClient(t *testing.T) {
	p := NewPlugin("test", "development", newMemDBClient)
	PluginTestSuiteRunner(t, p, nil,
		PluginTestSuiteTests{})
}

func TestPluginManagedClient(t *testing.T) {
	p := NewPlugin("test", "development", newMemDBClient, WithManagerWriter())
	PluginTestSuiteRunner(t, p, nil,
		PluginTestSuiteTests{})
}