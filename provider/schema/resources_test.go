package schema

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testPrimaryKeyTable = &Table{
	Name:    "test_pk_table",
	Options: TableCreationOptions{PrimaryKeys: []string{"primary_key_str"}},
	Columns: []Column{
		{
			Name: "primary_key_str",
			Type: TypeString,
		},
	},
	Relations: []*Table{
		{
			Name:    "test_pk_rel_table",
			Options: TableCreationOptions{PrimaryKeys: []string{"primary_rel_key_str"}},
			Columns: []Column{
				{
					Name: "rel_key_str",
					Type: TypeString,
				},
			},
		},
	},
}

// TestResourcePrimaryKey checks resource id generation when primary key is set on table
func TestResourcePrimaryKey(t *testing.T) {
	r := NewResourceData(PostgresDialect{}, testPrimaryKeyTable, nil, nil, nil, time.Now())
	// save random id
	randomId := r.cqId
	// test primary table no pk
	assert.Error(t, r.GenerateCQId(), "Error expected, primary key value not set")
	// Id shouldn't change
	assert.Equal(t, randomId, r.cqId)
	err := r.Set("primary_key_str", "test")
	assert.Nil(t, err)
	assert.Nil(t, r.GenerateCQId())
	assert.NotEqual(t, randomId, r.cqId)
	randomId = r.cqId
	// validate consistency
	assert.Nil(t, r.GenerateCQId())
	assert.Equal(t, randomId, r.cqId)
	// check key length of array is as expected
	assert.Len(t, r.PrimaryKeyValues(), 1)
	var strPtr = "primary_key_str"
	assert.Nil(t, r.Set("primary_key_str", &strPtr))
	assert.Equal(t, r.PrimaryKeyValues(), []string{"primary_key_str"})
	// check stringer interface
	uuidPK := uuid.New()
	assert.Nil(t, r.Set("primary_key_str", uuidPK))
	assert.Equal(t, r.PrimaryKeyValues(), []string{uuidPK.String()})
}

// TestResourcePrimaryKey checks resource id generation when primary key is set on table
func TestResourceAddColumns(t *testing.T) {
	r := NewResourceData(PostgresDialect{}, testPrimaryKeyTable, nil, nil, nil, time.Now())
	assert.Equal(t, []string{"cq_id", "cq_meta", "primary_key_str"}, r.columns)
}

func TestResourceColumns(t *testing.T) {
	r := NewResourceData(PostgresDialect{}, testTable, nil, nil, nil, time.Now())
	errf := r.Set("name", "test")
	assert.Nil(t, errf)
	assert.Equal(t, r.Get("name"), "test")
	v, err := r.Values()
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{nil, nil, "test", nil, nil}, v)
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
	assert.Equal(t, []interface{}{nil, nil, "test", "name_no_prefix", "prefix_name"}, v)

	// check non existing col
	err = r.Set("non_exist_col", "test")
	assert.Error(t, err)
}

func TestResources(t *testing.T) {
	r1 := NewResourceData(PostgresDialect{}, testPrimaryKeyTable, nil, nil, nil, time.Now())
	r2 := NewResourceData(PostgresDialect{}, testPrimaryKeyTable, nil, nil, nil, time.Now())
	assert.Equal(t, []string{"cq_id", "cq_meta", "primary_key_str"}, r1.columns)
	assert.Equal(t, []string{"cq_id", "cq_meta", "primary_key_str"}, r2.columns)

	rr := Resources{r1, r2}
	assert.Equal(t, []string{"cq_id", "cq_meta", "primary_key_str"}, rr.ColumnNames())
	assert.Equal(t, testPrimaryKeyTable.Name, rr.TableName())
	_ = r1.Set("primary_key_str", "test")
	_ = r2.Set("primary_key_str", "test2")
	_ = r1.GenerateCQId()
	_ = r2.GenerateCQId()
	assert.Equal(t, []uuid.UUID{r1.Id(), r2.Id()}, rr.GetIds())
}
