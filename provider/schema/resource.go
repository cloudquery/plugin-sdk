package schema

import (
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure"
	"github.com/thoas/go-funk"
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
	table          *Table
	data           map[string]interface{}
	cqId           uuid.UUID
	metadata       map[string]interface{}
	columns        []string
	dialect        Dialect
	executionStart time.Time
}

func NewResourceData(dialect Dialect, t *Table, parent *Resource, item interface{}, metadata map[string]interface{}, startTime time.Time) *Resource {
	return &Resource{
		Item:           item,
		Parent:         parent,
		table:          t,
		data:           make(map[string]interface{}),
		cqId:           uuid.New(),
		columns:        dialect.Columns(t).Names(),
		metadata:       metadata,
		dialect:        dialect,
		executionStart: startTime,
	}
}
func (r *Resource) PrimaryKeyValues() []string {
	tablePrimKeys := r.dialect.PrimaryKeys(r.table)
	if len(tablePrimKeys) == 0 {
		return []string{}
	}
	results := make([]string, 0)
	for _, primKey := range tablePrimKeys {
		data := r.Get(primKey)
		if data == nil {
			continue
		}
		// we can have more types, but PKs are usually either ints, strings or a structure
		// hopefully supporting Stringer interface, otherwise we fallback
		switch v := data.(type) {
		case fmt.Stringer:
			results = append(results, v.String())
		case *string:
			results = append(results, *v)
		case *int:
			results = append(results, fmt.Sprintf("%d", *v))
		default:
			results = append(results, fmt.Sprintf("%v", v))
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
	for _, c := range r.dialect.Columns(r.table) {
		v := r.Get(c.Name)
		if err := c.ValidateType(v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func (r *Resource) GenerateCQId() error {
	if len(r.table.Options.PrimaryKeys) == 0 {
		return nil
	}
	pks := r.dialect.PrimaryKeys(r.table)
	objs := make([]interface{}, 0, len(pks))
	for _, pk := range pks {
		if col := r.getColumnByName(pk); col == nil {
			return fmt.Errorf("failed to generate cq_id for %s, pk column missing %s", r.table.Name, pk)
		} else if col.internal {
			continue
		}

		value := r.Get(pk)
		if value == nil {
			return fmt.Errorf("failed to generate cq_id for %s, pk field missing %s", r.table.Name, pk)
		}
		objs = append(objs, value)
	}
	id, err := hashUUID(objs)
	if err != nil {
		return err
	}
	r.cqId = id
	return nil
}

func (r *Resource) TableName() string {
	if r.table == nil {
		return ""
	}
	return r.table.Name
}

func (r Resource) getColumnByName(column string) *Column {
	for _, c := range r.dialect.Columns(r.table) {
		if strings.Compare(column, c.Name) == 0 {
			return &c
		}
	}
	return nil
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
