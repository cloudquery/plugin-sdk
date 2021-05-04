package schema

import (
	"context"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

var testZeroTable = &Table{
	Name: "test_zero_table",
	Columns: []Column{
		{
			Name: "zero_bool",
			Type: TypeBool,
		},
		{
			Name: "zero_int",
			Type: TypeBigInt,
		},
		{
			Name: "not_zero_bool",
			Type: TypeBool,
		},
		{
			Name: "not_zero_int",
			Type: TypeBigInt,
		},
		{
			Name: "zero_int_ptr",
			Type: TypeBigInt,
		},
		{
			Name: "not_zero_int_ptr",
			Type: TypeBigInt,
		},
		{
			Name: "zero_string",
			Type: TypeString,
		},
	},
}

type zeroValuedStruct struct {
	ZeroBool      bool   `default:"false"`
	ZeroInt       int    `default:"0"`
	NotZeroInt    int    `default:"5"`
	NotZeroBool   bool   `default:"true"`
	ZeroIntPtr    *int   `default:"0"`
	NotZeroIntPtr *int   `default:"5"`
	ZeroString    string `default:""`
}

func TestResourceColumns(t *testing.T) {

	r := NewResourceData(testTable, nil, nil)
	errf := r.Set("name", "test")
	assert.Nil(t, errf)
	assert.Equal(t, r.Get("name"), "test")
	v, err := r.Values()
	assert.Nil(t, err)
	assert.Equal(t, v, []interface{}{r.id, "test", nil, nil})
	// Set invalid type to resource
	errf = r.Set("name", 5)
	assert.Nil(t, errf)
	v, err = r.Values()
	assert.Error(t, err)
	assert.Nil(t, v)

	// Set resource fully
	errf = r.Set("name", "test")
	assert.Nil(t, errf)
	errf = r.Set("name_no_prefix", "name_no_prefix")
	assert.Nil(t, errf)
	errf = r.Set("prefix_name", "prefix_name")
	assert.Nil(t, errf)
	v, err = r.Values()
	assert.Nil(t, err)
	assert.Equal(t, v, []interface{}{r.id, "test", "name_no_prefix", "prefix_name"})

	// check non existing col
	err = r.Set("non_exist_col", "test")
	assert.Error(t, err)
}

func TestResourceResolveColumns(t *testing.T) {
	t.Run("test resolve column normal", func(t *testing.T) {
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
	})

	t.Run("test resolve zero columns", func(t *testing.T) {
		object := zeroValuedStruct{}
		_ = defaults.Set(&object)
		r := NewResourceData(testZeroTable, nil, object)
		assert.Equal(t, r.id, r.Id())
		// columns should be resolved from ColumnResolver functions or default functions
		logger := logging.New(&hclog.LoggerOptions{
			Name:   "test_log",
			Level:  hclog.Error,
			Output: nil,
		})
		exec := NewExecutionData(nil, logger, testZeroTable)
		err := exec.resolveColumns(context.TODO(), nil, r, testZeroTable.Columns)
		assert.Nil(t, err)
		v, err := r.Values()
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{r.id, false, 0, true}, v[:4])
		assert.Equal(t, 0, *v[5].(*int))
		assert.Equal(t, 5, *v[6].(*int))
		assert.Equal(t, "", v[7].(string))

		object.ZeroIntPtr = nil
		r = NewResourceData(testZeroTable, nil, object)
		err = exec.resolveColumns(context.TODO(), nil, r, testZeroTable.Columns)
		assert.Nil(t, err)
		v, _ = r.Values()
		assert.Equal(t, nil, v[5])
	})

}
