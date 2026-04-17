package queue

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// Codec serializes *schema.Resource values so the queue scheduler can spill
// them to an external Storage backend and reconstruct them on the consumer
// side. Uses JSON for Item payloads (exploiting existing API-response JSON
// tags) and reflects on Table.ItemSampleType() to round-trip to the concrete
// Go type.
type Codec struct {
	tablesByName map[string]*schema.Table
}

// NewCodec builds a codec with a table lookup. tables may be the flattened
// list of all tables a plugin handles — only table names are used as keys.
func NewCodec(tables schema.Tables) *Codec {
	m := make(map[string]*schema.Table, len(tables))
	walk(tables, func(t *schema.Table) { m[t.Name] = t })
	return &Codec{tablesByName: m}
}

func walk(tables schema.Tables, f func(*schema.Table)) {
	for _, t := range tables {
		f(t)
		walk(t.Relations, f)
	}
}

type serializedResource struct {
	TableName string          `json:"table_name"`
	Item      json.RawMessage `json:"item"`
	ParentID  string          `json:"parent_id,omitempty"`
}

// EncodeResource serializes r with an explicit parentID (caller-chosen UUID
// of the resource's parent in the Storage, or "" for root resources).
func (c *Codec) EncodeResource(r *schema.Resource, parentID string) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("codec: nil resource")
	}
	if r.Table == nil {
		return nil, fmt.Errorf("codec: resource has nil table")
	}
	itemBytes, err := json.Marshal(r.Item)
	if err != nil {
		return nil, fmt.Errorf("codec: marshal item for table %q: %w", r.Table.Name, err)
	}
	return json.Marshal(serializedResource{
		TableName: r.Table.Name,
		Item:      itemBytes,
		ParentID:  parentID,
	})
}

// DecodeResource reconstructs a *schema.Resource from bytes. Returns the
// resource, the parentID it references (for callers that want to chain-load
// ancestors), and any error. The returned Resource has Parent=nil; callers
// wanting the ancestor chain attached should use a higher-level helper.
func (c *Codec) DecodeResource(data []byte) (*schema.Resource, string, error) {
	var sr serializedResource
	if err := json.Unmarshal(data, &sr); err != nil {
		return nil, "", fmt.Errorf("codec: unmarshal envelope: %w", err)
	}
	tbl, ok := c.tablesByName[sr.TableName]
	if !ok {
		return nil, "", fmt.Errorf("codec: unknown table %q", sr.TableName)
	}
	sampleType := tbl.ItemSampleType()
	if sampleType == nil {
		return nil, "", fmt.Errorf("codec: table %q has no itemSample; configure TransformWithStruct or SetItemSample", sr.TableName)
	}
	ptr := reflect.New(sampleType).Interface()
	if err := json.Unmarshal(sr.Item, ptr); err != nil {
		return nil, "", fmt.Errorf("codec: unmarshal item for table %q: %w", sr.TableName, err)
	}
	item := reflect.ValueOf(ptr).Elem().Interface()
	return schema.NewResourceData(tbl, nil, item), sr.ParentID, nil
}
