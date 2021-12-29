package schema

import (
	"crypto"
	"fmt"
	"strings"

	"github.com/mitchellh/hashstructure"
	"github.com/thoas/go-funk"

	"github.com/google/uuid"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item interface{}
	// Set if this is an embedded table
	Parent *Resource
	// internal fields
	table       *Table
	data        map[string]interface{}
	cqId        uuid.UUID
	extraFields map[string]interface{}
	columns     []string
}

func NewResourceData(t *Table, parent *Resource, item interface{}, extraFields map[string]interface{}) *Resource {
	return &Resource{
		Item:        item,
		Parent:      parent,
		table:       t,
		data:        make(map[string]interface{}),
		cqId:        uuid.New(),
		columns:     getResourceColumns(t, extraFields),
		extraFields: extraFields,
	}
}
func (r *Resource) Keys() []string {
	tablePrimKeys := r.table.PrimaryKeys()
	if len(tablePrimKeys) == 0 {
		return []string{}
	}
	results := make([]string, 0)
	for _, primKey := range tablePrimKeys {
		data := r.Get(primKey)
		if data != nil {
			results = append(results, fmt.Sprintf("%v", data))
		}
	}
	return results
}

func (r *Resource) Get(key string) interface{} {
	return r.data[key]
}

func (r *Resource) Set(key string, value interface{}) error {
	columnExists := funk.ContainsString(r.columns, key)
	if !columnExists {
		return fmt.Errorf("column %s does not exist", key)
	}
	r.data[key] = value
	return nil
}

func (r *Resource) Id() uuid.UUID {
	return r.cqId
}

func (r *Resource) Values() ([]interface{}, error) {
	values := make([]interface{}, 0)
	for _, c := range append(r.table.Columns, GetDefaultSDKColumns()...) {
		v := r.Get(c.Name)
		if err := c.ValidateType(v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	for _, v := range r.extraFields {
		values = append(values, v)
	}
	return values, nil
}

func (r *Resource) GenerateCQId() error {
	if len(r.table.Options.PrimaryKeys) == 0 {
		return nil
	}
	objs := make([]interface{}, len(r.table.PrimaryKeys()))
	for i, pk := range r.table.PrimaryKeys() {
		value := r.Get(pk)
		if value == nil {
			return fmt.Errorf("failed to generate cq_id for %s, pk field missing %s", r.table.Name, pk)
		}
		objs[i] = value
	}
	id, err := hashUUID(objs)
	if err != nil {
		return err
	}
	r.cqId = id
	return nil
}

func (r Resource) getColumnByName(column string) *Column {
	for _, c := range r.table.Columns {
		if strings.Compare(column, c.Name) == 0 {
			return &c
		}
	}
	return nil
}

func hashUUID(objs interface{}) (uuid.UUID, error) {
	// Use SHA1 because it's fast and is reasonably enough protected against accidental collisions.
	// There is no scenario here where intentional created collisions could do harm.
	digester := crypto.SHA1.New()
	hash, err := hashstructure.Hash(objs, nil)
	if err != nil {
		return uuid.Nil, err
	}
	if _, err := fmt.Fprint(digester, hash); err != nil {
		return uuid.Nil, err
	}
	data := digester.Sum(nil)
	return uuid.NewSHA1(uuid.Nil, data), nil
}

func getResourceColumns(t *Table, fields map[string]interface{}) []string {
	columns := t.ColumnNames()
	for k := range fields {
		columns = append(columns, k)
	}
	return columns
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
	return rr[0].table.Name
}

func (rr Resources) ColumnNames() []string {
	if len(rr) == 0 {
		return []string{}
	}
	return rr[0].columns
}
