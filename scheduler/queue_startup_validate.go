package scheduler

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// ValidateTablesForQueue verifies every relation-having table has an
// itemSample registered (via TransformWithStruct or SetItemSample). Only
// checked when cfg selects a non-in-memory backend.
func ValidateTablesForQueue(tables schema.Tables, cfg *QueueConfig) error {
	if cfg == nil || cfg.Type == "" || cfg.Type == QueueTypeInMemory {
		return nil
	}
	var walk func(ts []*schema.Table) error
	walk = func(ts []*schema.Table) error {
		for _, t := range ts {
			if len(t.Relations) == 0 {
				continue
			}
			if t.ItemSampleType() == nil {
				return fmt.Errorf("queue: table %q has relations but no itemSample; ensure it calls transformers.TransformWithStruct or Table.SetItemSample", t.Name)
			}
			if err := walk(t.Relations); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(tables)
}
