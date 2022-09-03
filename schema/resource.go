package schema

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item interface{} `json:"-"`
	// Set if this is an embedded table
	Parent *Resource `json:"-"`
	// internal fields
	Table *Table `json:"-"`
	// This is sorted result data by column name
	Data      map[string]interface{} `json:"data"`
	TableName string                 `json:"table_name"`
}

func NewResourceData(t *Table, parent *Resource, fetchTime time.Time, item interface{}) *Resource {
	r := Resource{
		Item:      item,
		Parent:    parent,
		Table:     t,
		Data:      make(map[string]interface{}, len(t.Columns)),
		TableName: t.Name,
	}
	r.Data[CqFetchTime.Name] = fetchTime
	return &r
}

func (r *Resource) PrimaryKeyValue() string {
	pks := r.Table.PrimaryKeys()
	if len(pks) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, primKey := range pks {
		data := r.Get(primKey)
		if data == nil {
			continue
		}
		// we can have more types, but PKs are usually either ints, strings or a structure
		// hopefully supporting Stringer interface, otherwise we fallback
		switch v := data.(type) {
		case fmt.Stringer:
			sb.WriteString(v.String())
		case *string:
			sb.WriteString(*v)
		case *int:
			sb.WriteString(fmt.Sprintf("%d", *v))
		default:
			sb.WriteString(fmt.Sprintf("%d", v))
		}
	}
	return sb.String()
}

func (r *Resource) Get(key string) interface{} {
	return r.Data[key]
}

func (r *Resource) Set(key string, value interface{}) error {
	r.Data[key] = value
	return nil
}

// Override original item (this is useful for apis that follow list/details pattern)
func (r *Resource) SetItem(item interface{}) {
	r.Item = item
}

func (r *Resource) Id() uuid.UUID {
	if r.Data[CqIdColumn.Name] == nil {
		return uuid.UUID{}
	}
	return r.Data[CqIdColumn.Name].(uuid.UUID)
}

func (r *Resource) Columns() []string {
	return r.Table.Columns.Names()
}

func (rr Resources) GetIds() []uuid.UUID {
	rids := make([]uuid.UUID, len(rr))
	for i, r := range rr {
		rids[i] = r.Id()
	}
	return rids
}
func (rr Resources) TableName() string {
	if len(rr) == 0 {
		return ""
	}
	return rr[0].Table.Name
}

func (rr Resources) ColumnNames() []string {
	if len(rr) == 0 {
		return []string{}
	}
	return rr[0].Table.Columns.Names()
}
