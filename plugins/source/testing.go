package source

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

func TestPluginSync(t *testing.T, plugin *Plugin, spec specs.Source, opts ...TestPluginOption) {
	t.Helper()

	o := &testPluginOptions{
		parallel:   true,
		validators: []Validator{validateColumnsHaveData},
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.parallel {
		t.Parallel()
	}

	resourcesChannel := make(chan *schema.Resource)
	var syncErr error

	if err := plugin.Init(context.Background(), spec); err != nil {
		t.Fatal(err)
	}

	go func() {
		defer close(resourcesChannel)
		syncErr = plugin.Sync(context.Background(), resourcesChannel)
	}()

	syncedResources := make([]*schema.Resource, 0)
	for resource := range resourcesChannel {
		syncedResources = append(syncedResources, resource)
	}
	if syncErr != nil {
		t.Fatal(syncErr)
	}
	for _, validator := range o.validators {
		err := validator(plugin, syncedResources)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type TestPluginOption func(*testPluginOptions)

func WithTestPluginNoParallel() TestPluginOption {
	return func(f *testPluginOptions) {
		f.parallel = false
	}
}

func WithTestPluginAdditionalValidators(v Validator) TestPluginOption {
	return func(f *testPluginOptions) {
		f.validators = append(f.validators, v)
	}
}

type testPluginOptions struct {
	parallel   bool
	validators []Validator
}
