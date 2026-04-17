package scheduler

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/stretchr/testify/require"
)

type goodItem struct{ ID string }
type childOfGood struct{ ID string }

func TestValidateTablesForQueue_InMemorySkipsCheck(t *testing.T) {
	tables := schema.Tables{{
		Name:      "no-sample",
		Relations: []*schema.Table{{Name: "child"}},
	}}
	require.NoError(t, ValidateTablesForQueue(tables, nil))
}

func TestValidateTablesForQueue_BadgerRequiresItemSample(t *testing.T) {
	tbl := &schema.Table{
		Name:      "no-sample",
		Relations: []*schema.Table{{Name: "child"}},
	}
	tables := schema.Tables{tbl}
	err := ValidateTablesForQueue(tables, &QueueConfig{Type: QueueTypeBadger, Path: "/tmp"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no-sample")
	require.Contains(t, err.Error(), "itemSample")
}

func TestValidateTablesForQueue_BadgerWithItemSamplePasses(t *testing.T) {
	child := &schema.Table{Name: "child", Transform: transformers.TransformWithStruct(&childOfGood{})}
	root := &schema.Table{
		Name:      "root",
		Transform: transformers.TransformWithStruct(&goodItem{}),
		Relations: []*schema.Table{child},
	}
	require.NoError(t, root.Transform(root))
	require.NoError(t, child.Transform(child))
	require.NoError(t, ValidateTablesForQueue(schema.Tables{root}, &QueueConfig{Type: QueueTypeBadger, Path: "/tmp"}))
}

func TestValidateTablesForQueue_LeafWithoutItemSampleOK(t *testing.T) {
	// A leaf table (no Relations) doesn't need itemSample.
	tables := schema.Tables{{Name: "leaf"}}
	require.NoError(t, ValidateTablesForQueue(tables, &QueueConfig{Type: QueueTypeBadger, Path: "/tmp"}))
}
