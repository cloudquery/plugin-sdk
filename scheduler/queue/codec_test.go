package queue

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/stretchr/testify/require"
)

type codecTestItem struct {
	ID   string
	Name string
	Age  int
}

func TestCodec_RoundTripRoot(t *testing.T) {
	tbl := &schema.Table{
		Name:      "codec_test",
		Transform: transformers.TransformWithStruct(&codecTestItem{}),
		Columns:   schema.ColumnList{{Name: "id", Type: arrow.BinaryTypes.String}},
	}
	require.NoError(t, tbl.Transform(tbl))

	item := codecTestItem{ID: "a", Name: "alice", Age: 30}
	res := schema.NewResourceData(tbl, nil, item)
	require.NoError(t, res.Set("id", "a"))

	tables := schema.Tables{tbl}
	c := NewCodec(tables)

	data, err := c.EncodeResource(res, "") // empty parentID = root
	require.NoError(t, err)
	require.NotEmpty(t, data)

	decoded, parentID, err := c.DecodeResource(data)
	require.NoError(t, err)
	require.Equal(t, "", parentID)
	require.Equal(t, "codec_test", decoded.Table.Name)
	require.Nil(t, decoded.Parent)

	typed, ok := decoded.Item.(codecTestItem)
	require.True(t, ok, "Item should round-trip to concrete type, got %T", decoded.Item)
	require.Equal(t, item, typed)
}

func TestCodec_RoundTripWithParentRef(t *testing.T) {
	parentTbl := &schema.Table{
		Name:      "parent_tbl",
		Transform: transformers.TransformWithStruct(&codecTestItem{}),
		Columns:   schema.ColumnList{{Name: "id", Type: arrow.BinaryTypes.String}},
	}
	childTbl := &schema.Table{
		Name:      "child_tbl",
		Transform: transformers.TransformWithStruct(&codecTestItem{}),
		Columns:   schema.ColumnList{{Name: "id", Type: arrow.BinaryTypes.String}},
	}
	require.NoError(t, parentTbl.Transform(parentTbl))
	require.NoError(t, childTbl.Transform(childTbl))

	child := schema.NewResourceData(childTbl, nil, codecTestItem{ID: "c"})
	c := NewCodec(schema.Tables{parentTbl, childTbl})

	data, err := c.EncodeResource(child, "parent-id-123")
	require.NoError(t, err)

	decoded, parentID, err := c.DecodeResource(data)
	require.NoError(t, err)
	require.Equal(t, "parent-id-123", parentID)
	require.Equal(t, "child_tbl", decoded.Table.Name)
}

func TestCodec_DecodeUnknownTableErrors(t *testing.T) {
	other := &schema.Table{Name: "unknown_table", Transform: transformers.TransformWithStruct(&codecTestItem{})}
	require.NoError(t, other.Transform(other))

	// Use a codec that does NOT have other/unknown_table registered.
	c := NewCodec(schema.Tables{})
	r := schema.NewResourceData(other, nil, codecTestItem{})
	blob, err := c.EncodeResource(r, "")
	require.NoError(t, err)
	_, _, err = c.DecodeResource(blob)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown_table")
}
