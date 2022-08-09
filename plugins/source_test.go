package plugins

import (
	"context"
	"os"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/schema"
	"github.com/rs/zerolog"
	"github.com/xeipuuv/gojsonschema"
)

type Account struct {
	Name    string   `yaml:"name"`
	Regions []string `yaml:"regions"`
}

type TestConfig struct {
	Accounts []Account `yaml:"accounts"`
	Regions  []string  `yaml:"regions"`
}

func (TestConfig) Example() string {
	return ""
}

type testSourcePluginClient struct {
	logger zerolog.Logger
}

func (t testSourcePluginClient) Logger() *zerolog.Logger {
	return &t.logger
}

var testSourcePlugin = SourcePlugin{
	Name:    "testSourcePlugin",
	Version: "1.0.0",
	Configure: func(l zerolog.Logger, i interface{}) (schema.ClientMeta, error) {
		return testSourcePluginClient{logger: l}, nil
	},
	Tables: []*schema.Table{
		{
			Name: "testTable",
			Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
				res <- map[string]interface{}{
					"testColumn": 3,
				}
				return nil
			},
			Columns: []schema.Column{
				{
					Name: "testColumn",
					Type: schema.TypeInt,
				},
			},
		},
	},
	Config: func() interface{} {
		return &TestConfig{}
	},
	Logger: zerolog.New(os.Stderr),
}

func TestFetch(t *testing.T) {
	cfg := `
tables:
  - "*"
configuration:
  regions:
  - "us-east-1"
  accounts:
  - name: "testAccount"
    regions:
    - "us-east-2"
`
	resources := make(chan *schema.Resource)
	var fetchErr error
	var result *gojsonschema.Result
	go func() {
		defer close(resources)
		result, fetchErr = testSourcePlugin.Fetch(context.Background(), []byte(cfg), resources)
	}()
	for resource := range resources {
		t.Logf("%+v", resource)
	}
	if fetchErr != nil {
		t.Errorf("fetch error: %v", fetchErr)
	}
}
