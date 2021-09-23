package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/cloudquery/faker/v3"
	"github.com/hashicorp/go-hclog"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"

	"github.com/stretchr/testify/assert"
)

type (
	testStruct struct {
		Id   int
		Name string
	}
	testConfig struct{}
)

func (t testConfig) Example() string {
	return ""
}

var (
	provider = Provider{
		ResourceMap: map[string]*schema.Table{
			"test": {
				Name: "sdk_test",
				Relations: []*schema.Table{
					{
						Name: "sdk_test_test1",
						Relations: []*schema.Table{
							{Name: "sdk_test_test1_test"},
						},
					},
				},
			},
			"test1": {
				Name:      "sdk_test1",
				Relations: []*schema.Table{},
			},
		},
	}
	testResolverFunc = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
		for i := 0; i < 10; i++ {
			t := testStruct{}
			res <- faker.FakeData(&t)
		}
		return nil
	}

	testProviderCreatorFunc = func() Provider {
		return Provider{
			Name: "unitest",
			Config: func() Config {
				return &testConfig{}
			},
			ResourceMap: map[string]*schema.Table{
				"test": {
					Name:     "test_resource",
					Resolver: testResolverFunc,
					Columns: []schema.Column{
						{
							Name: "id",
							Type: schema.TypeBigInt,
						},
						{
							Name: "name",
							Type: schema.TypeString,
						},
					},
					Relations: []*schema.Table{
						{
							Name:     "test_resource_relation",
							Resolver: testResolverFunc,
							Columns: []schema.Column{
								{
									Name: "id",
									Type: schema.TypeInt,
								},
								{
									Name: "name",
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			},
		}
	}

	failProvider = Provider{
		ResourceMap: map[string]*schema.Table{
			"test": {
				Name: "sdk_test",
				Relations: []*schema.Table{
					{
						Name: "sdk_test1",
					},
				},
			},
			"test1": {
				Name: "sdk_test1",
			},
		},
	}
)

func TestProviderInterpolate(t *testing.T) {
	r, err := provider.interpolateAllResources([]string{"test"})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test"}, r)

	r, err = provider.interpolateAllResources([]string{"test", "test1"})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test", "test1"}, r)

	r, err = provider.interpolateAllResources([]string{"test", "test1", "*"})
	assert.Error(t, err)
	assert.Nil(t, r)
	r, err = provider.interpolateAllResources([]string{"*"})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test", "test1"}, r)

}

func TestTableDuplicates(t *testing.T) {
	tables := make(map[string]string)
	var err error
	for r, t := range provider.ResourceMap {
		err = getTableDuplicates(r, t, tables)
	}
	assert.Nil(t, err)

	for r, t := range failProvider.ResourceMap {
		err = getTableDuplicates(r, t, tables)
	}
	assert.Error(t, err)
}

func TestProvider_ConfigureProvider(t *testing.T) {
	tp := testProviderCreatorFunc()
	tp.Configure = func(logger hclog.Logger, i interface{}) (schema.ClientMeta, error) {
		return nil, errors.New("test error")
	}
	resp, err := tp.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "dev",
		Connection: cqproto.ConnectionDetails{
			DSN: "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable",
		},
		Config:        nil,
		DisableDelete: true,
		ExtraFields:   nil,
	})
	assert.Equal(t, "provider unitest logger not defined, make sure to run it with serve", resp.Error)
	assert.Nil(t, err)
	// set logger this time
	tp.Logger = hclog.Default()
	resp, err = tp.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "dev",
		Connection: cqproto.ConnectionDetails{
			DSN: "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable",
		},
		Config:        nil,
		DisableDelete: true,
		ExtraFields:   nil,
	})
	assert.Equal(t, "", resp.Error)
	assert.Error(t, err)
}
