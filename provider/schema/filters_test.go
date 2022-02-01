package schema

import (
	"testing"
	"time"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testTableStruct struct {
	Name  string `default:"test"`
	Inner struct {
		NameNoPrefix string `default:"name_no_prefix"`
	}
	Prefix struct {
		Name string `default:"prefix_name"`
	}
}

var testTable = &Table{
	Name: "test_table",
	Columns: []Column{
		{
			Name: "name",
			Type: TypeString,
		},
		{
			Name:     "name_no_prefix",
			Type:     TypeString,
			Resolver: PathResolver("Inner.NameNoPrefix"),
		},
		{
			Name:     "prefix_name",
			Type:     TypeString,
			Resolver: PathResolver("Prefix.Name"),
		},
	},
}

func TestDeleteParentId(t *testing.T) {
	f := DeleteParentIdFilter("name")
	mockedClient := new(MockedClientMeta)
	logger := logging.New(&hclog.LoggerOptions{
		Name:   "test_log",
		Level:  hclog.Error,
		Output: nil,
	})
	mockedClient.On("Logger", mock.Anything).Return(logger)

	object := testTableStruct{}
	r := NewResourceData(PostgresDialect{}, testTable, nil, object, nil, time.Now())
	_ = r.Set("name", "test")
	assert.Equal(t, []interface{}{"name", r.Id()}, f(mockedClient, r))

	assert.Nil(t, f(mockedClient, nil))
}
