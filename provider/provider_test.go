package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudquery/cq-provider-sdk/provider/execution"

	"github.com/cloudquery/cq-provider-sdk/provider/diag"

	"github.com/cloudquery/cq-provider-sdk/provider/schema/mock"
	"github.com/golang/mock/gomock"

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
	testClient struct{}
)

func (t testConfig) Example() string {
	return ""
}

func (t testClient) Logger() hclog.Logger {
	return hclog.Default()
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
	testResolverFunc = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		for i := 0; i < 10; i++ {
			t := testStruct{}
			time.Sleep(50 * time.Millisecond)
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
				"bad_resource": {
					Name: "bad_resource",
					Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
						return errors.New("bad error")
					},
				},
				"bad_resource_ignore_error": {
					Name: "bad_resource_ignore_error",
					IgnoreError: func(err error) bool {
						return true
					},
					Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
						return errors.New("bad error")
					},
				},
				"very_slow_resource": {
					Name: "very_slow_resource",
					Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
						select {
						case <-time.After(10 * time.Second):
							return nil
						case <-ctx.Done():
							return errors.New("context canceled")
						}
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

	parallelCheckProvider = Provider{
		Name: "parallel",
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
			},
			"test1": {
				Name:     "test1_resource",
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
			},
			"test2": {
				Name:     "test2_resource",
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
			},
			"test3": {
				Name:     "test3_resource",
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
			},
			"test4": {
				Name:     "test4_resource",
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
		Config:      nil,
		ExtraFields: nil,
	})
	assert.Equal(t, "provider unitest logger not defined, make sure to run it with serve", resp.Error)
	assert.NotNil(t, err)
	// set logger this time
	tp.Logger = hclog.Default()
	resp, err = tp.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "dev",
		Connection: cqproto.ConnectionDetails{
			DSN: "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable",
		},
		Config:      nil,
		ExtraFields: nil,
	})
	assert.Equal(t, "", resp.Error)
	assert.Error(t, err)
}

type FetchResourceTableTest struct {
	Name                   string
	ExpectedFetchResponses []*cqproto.FetchResourcesResponse
	ExpectedError          error
	MockStorageFunc        func(ctrl *gomock.Controller) *mock.MockStorage
	PartialFetch           bool
	ResourcesToFetch       []string
	Context                func() context.Context
}

func TestProvider_FetchResources(t *testing.T) {
	tp := testProviderCreatorFunc()
	tp.Logger = hclog.Default()
	tp.Configure = func(logger hclog.Logger, i interface{}) (schema.ClientMeta, error) {
		return &testClient{}, nil
	}
	_, err := tp.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "dev",
		Connection: cqproto.ConnectionDetails{
			DSN: "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable",
		},
		Config:      nil,
		ExtraFields: nil,
	})
	ctrl := gomock.NewController(t)
	var fetchCases = []FetchResourceTableTest{
		{
			Name: "ignore error resource",
			ExpectedFetchResponses: []*cqproto.FetchResourcesResponse{
				{
					ResourceName: "bad_resource_ignore_error",
					Summary: cqproto.ResourceFetchSummary{
						Status:        cqproto.ResourceFetchComplete,
						ResourceCount: 0,
						Diagnostics: diag.Diagnostics{diag.NewBaseError(errors.New("table[bad_resource_ignore_error] resolver ignored error. Error: bad error"),
							diag.IGNORE, diag.RESOLVING, "bad_resource_ignore_error", "table[bad_resource_ignore_error] resolver ignored error. Error: bad error", "")},
					},
				}},
			ExpectedError: nil,
			MockStorageFunc: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockDB := mock.NewMockStorage(ctrl)
				mockDB.EXPECT().Dialect().Return(schema.PostgresDialect{})
				mockDB.EXPECT().Close()
				return mockDB
			},
			PartialFetch:     true,
			ResourcesToFetch: []string{"bad_resource_ignore_error"},
		},
		{
			Name: "returning error resource",
			ExpectedFetchResponses: []*cqproto.FetchResourcesResponse{
				{
					ResourceName: "bad_resource",
					Summary: cqproto.ResourceFetchSummary{
						Status:        cqproto.ResourceFetchPartial,
						ResourceCount: 0,
						Diagnostics:   diag.Diagnostics{diag.NewBaseError(errors.New("bad error"), diag.ERROR, diag.RESOLVING, "bad_resource", "failed to resolve table \"bad_resource\"", "")},
					},
				},
			},
			MockStorageFunc: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockDB := mock.NewMockStorage(ctrl)
				mockDB.EXPECT().Dialect().Return(schema.PostgresDialect{})
				mockDB.EXPECT().Close()
				return mockDB
			},
			PartialFetch:     false,
			ResourcesToFetch: []string{"bad_resource"},
		},
		{
			Name: "context canceled execution",
			ExpectedFetchResponses: []*cqproto.FetchResourcesResponse{
				{
					ResourceName: "very_slow_resource",
					Summary: cqproto.ResourceFetchSummary{
						Status:        cqproto.ResourceFetchCanceled,
						ResourceCount: 0,
						Diagnostics: diag.Diagnostics{diag.NewBaseError(errors.New("context canceled"),
							diag.ERROR, diag.RESOLVING, "very_slow_resource", "failed to resolve table \"very_slow_resource\"", "")},
					},
				}},
			ExpectedError: nil,
			MockStorageFunc: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockDB := mock.NewMockStorage(ctrl)
				mockDB.EXPECT().Dialect().Return(schema.PostgresDialect{})
				mockDB.EXPECT().Close()
				return mockDB
			},
			Context: func() context.Context {
				//nolint:govet
				ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
				return ctx
			},

			ResourcesToFetch: []string{"very_slow_resource"},
		},
	}

	for _, tt := range fetchCases {
		t.Run(tt.Name, func(t *testing.T) {
			tp.storageCreator = func(ctx context.Context, logger hclog.Logger, dbURL string) (execution.Storage, error) {
				return tt.MockStorageFunc(ctrl), nil
			}
			ctx := context.Background()
			if tt.Context != nil {
				ctx = tt.Context()
			}

			err = tp.FetchResources(ctx, &cqproto.FetchResourcesRequest{
				Resources: tt.ResourcesToFetch,
			}, &testResourceSender{
				t,
				tt.ExpectedFetchResponses,
			})
			if tt.ExpectedError != nil {
				assert.Equal(t, err, tt.ExpectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

type testResourceSender struct {
	t                 *testing.T
	ExpectedResponses []*cqproto.FetchResourcesResponse
}

func (f *testResourceSender) Send(r *cqproto.FetchResourcesResponse) error {
	for _, e := range f.ExpectedResponses {
		if e.ResourceName != r.ResourceName {
			continue
		}
		assert.Equal(f.t, r.Error, e.Error)
		assert.Equal(f.t, e.Summary.Status, r.Summary.Status)
		assert.Equal(f.t, diag.FlattenDiags(e.Summary.Diagnostics, true), diag.FlattenDiags(r.Summary.Diagnostics, true))
	}
	return nil
}

func TestProvider_FetchResourcesParallelLimit(t *testing.T) {
	parallelCheckProvider.Configure = func(logger hclog.Logger, i interface{}) (schema.ClientMeta, error) {
		return testClient{}, nil
	}
	parallelCheckProvider.Logger = hclog.Default()
	resp, err := parallelCheckProvider.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "dev",
		Connection: cqproto.ConnectionDetails{
			DSN: "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable",
		},
		Config:      nil,
		ExtraFields: nil,
	})
	assert.Equal(t, "", resp.Error)
	assert.Nil(t, err)

	// it runs 5 resources at a time. each resource takes ~500ms
	start := time.Now()
	err = parallelCheckProvider.FetchResources(context.Background(), &cqproto.FetchResourcesRequest{Resources: []string{"*"}}, &testResourceSender{})
	assert.Nil(t, err)
	length := time.Since(start)
	assert.Less(t, length, 1000*time.Millisecond)

	// it runs 5 resources one by one. each resource takes ~500ms
	start = time.Now()
	err = parallelCheckProvider.FetchResources(context.Background(), &cqproto.FetchResourcesRequest{Resources: []string{"*"}, ParallelFetchingLimit: 1}, &testResourceSender{})
	assert.Nil(t, err)
	length = time.Since(start)
	assert.Greater(t, length, 2500*time.Millisecond)
}
