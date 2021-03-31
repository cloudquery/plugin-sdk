package schema

import (
	"context"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestResourceColumns(t *testing.T) {

	r := NewResourceData(testTable, nil, nil)
	r.Set("name", "test")
	assert.Equal(t, r.Get("name"), "test")
	v, err := r.Values()
	assert.Nil(t, err)
	assert.Equal(t, v, []interface{}{r.id, "test", nil, nil})
	// Set invalid type to resource
	r.Set("name", 5)
	v, err = r.Values()
	assert.Error(t, err)
	assert.Nil(t, v)

	// Set resource fully
	r.Set("name", "test")
	r.Set("name_no_prefix", "name_no_prefix")
	r.Set("prefix_name", "prefix_name")
	v, err = r.Values()
	assert.Nil(t, err)
	assert.Equal(t, v, []interface{}{r.id, "test", "name_no_prefix", "prefix_name"})
}

func TestResourceResolveColumns(t *testing.T) {
	object := testTableStruct{}
	_ = defaults.Set(&object)
	r := NewResourceData(testTable, nil, object)
	assert.Equal(t, r.id, r.Id())
	// columns should be resolved from ColumnResolver functions or default functions
	logger := logging.New(&hclog.LoggerOptions{
		Name:   "test_log",
		Level:  hclog.Error,
		Output: nil,
	})
	exec := NewExecutionData(nil, logger, testTable)
	err := exec.resolveColumns(context.TODO(), nil, r, testTable.Columns)
	assert.Nil(t, err)
	v, err := r.Values()
	assert.Nil(t, err)
	assert.Equal(t, v, []interface{}{r.id, "test", "name_no_prefix", "prefix_name"})
}
