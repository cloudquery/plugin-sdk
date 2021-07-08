package schema

import (
	"testing"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteParentId(t *testing.T) {
	f := DeleteParentIdFilter("name")
	mockedClient := new(mockedClientMeta)
	logger := logging.New(&hclog.LoggerOptions{
		Name:   "test_log",
		Level:  hclog.Error,
		Output: nil,
	})
	mockedClient.On("Logger", mock.Anything).Return(logger)

	object := testTableStruct{}
	r := NewResourceData(testTable, nil, object, nil)
	_ = r.Set("name", "test")
	assert.Equal(t, []interface{}{"name", r.Id()}, f(mockedClient, r))

	assert.Nil(t, f(mockedClient, nil))
}
