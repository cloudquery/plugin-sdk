package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
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

type testPluginClient struct {
	logger zerolog.Logger
}

func (c *testPluginClient) Logger() *zerolog.Logger {
	return &c.logger
}

func testPluginConfigure(ctx context.Context, p *SourcePlugin, spec specs.SourceSpec) (schema.ClientMeta, error) {
	return &testPluginClient{
		logger: p.Logger,
	}, nil
}

func testTable() *schema.Table {
	return &schema.Table{
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
	}
}

const testSourceCfg = `
kind: source
spec:
  name: testSourcePlugin
  version: 1.0.0
  spec:
    accounts:
    - name: testAccount
`

func TestSync(t *testing.T) {
	// ctx := context.Background()
	// testSourcePlugin := NewSourcePlugin(
	// 	"test",
	// 	"v1.0.0",
	// 	[]*schema.Table{
	// 		testTable(),
	// 	},
	// 	testPluginConfigure,
	// 	WithSourceLogger(zerolog.New(zerolog.NewTestWriter(t))),
	// )
	// yaml.Unmarshal()
	// testSourcePlugin.Configure(ctx)
	// 	cfg := `
	// tables:
	//   - "*"
	// configuration:
	//   regions:
	//   - "us-east-1"
	//   accounts:
	//   - name: "testAccount"
	//     regions:
	//     - "us-east-2"
	// `
	// 	resources := make(chan *schema.Resource)
	// 	var fetchErr error
	// 	var result *gojsonschema.Result
	// 	go func() {
	// 		defer close(resources)
	// 		result, fetchErr = testSourcePlugin.Fetch(context.Background(), []byte(cfg), resources)
	// 	}()
	// 	for resource := range resources {
	// 		t.Logf("%+v", resource)
	// 	}
	// 	if fetchErr != nil {
	// 		t.Errorf("fetch error: %v", fetchErr)
	// 	}
}
