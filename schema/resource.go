package schema

import (
	"crypto"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/thoas/go-funk"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item interface{}
	// Set if this is an embedded table
	Parent *Resource `msgpack:"parent"`
	// internal fields
	Table    *Table                 `msgpack:"table"`
	Data     map[string]interface{} `msgpack:"data"`
	cqId     uuid.UUID              `msgpack:"cq_id"`
	metadata map[string]interface{}
	CColumns []string `msgpack:"columns"`
}

func NewResourceData(t *Table, parent *Resource, item interface{}) *Resource {
	return &Resource{
		Item:     item,
		Parent:   parent,
		Table:    t,
		Data:     make(map[string]interface{}),
		cqId:     uuid.New(),
		CColumns: t.Columns.Names(),
		// metadata: metadata,
	}
}

// func (r *Resource) PrimaryKeyValues() []string {
// 	tablePrimKeys := r.dialect.PrimaryKeys(r.table)
// 	if len(tablePrimKeys) == 0 {
// 		return []string{}
// 	}
// 	results := make([]string, 0)
// 	for _, primKey := range tablePrimKeys {
// 		data := r.Get(primKey)
// 		if data == nil {
// 			continue
// 		}
// 		// we can have more types, but PKs are usually either ints, strings or a structure
// 		// hopefully supporting Stringer interface, otherwise we fallback
// 		switch v := data.(type) {
// 		case fmt.Stringer:
// 			results = append(results, v.String())
// 		case *string:
// 			results = append(results, *v)
// 		case *int:
// 			results = append(results, fmt.Sprintf("%d", *v))
// 		default:
// 			results = append(results, fmt.Sprintf("%v", v))
// 		}
// 	}
// 	return results
// }

func (r *Resource) Get(key string) interface{} {
	return r.Data[key]
}

func (r *Resource) Set(key string, value interface{}) error {
	columnExists := funk.ContainsString(r.CColumns, key)
	if !columnExists {
		return fmt.Errorf("column %s does not exist", key)
	}
	r.Data[key] = value
	return nil
}

func (r *Resource) Id() uuid.UUID {
	return r.cqId
}

func (r *Resource) Columns() []string {
	return r.CColumns
}

// func (r *Resource) Values() ([]interface{}, error) {
// 	values := make([]interface{}, 0)
// 	for _, c := range r.dialect.Columns(r.table) {
// 		v := r.Get(c.Name)
// 		if err := c.ValidateType(v); err != nil {
// 			return nil, err
// 		}
// 		values = append(values, v)
// 	}
// 	return values, nil
// }

// func (r *Resource) GenerateCQId() error {
// 	if len(r.table.Options.PrimaryKeys) == 0 {
// 		return nil
// 	}
// 	pks := r.dialect.PrimaryKeys(r.table)
// 	objs := make([]interface{}, 0, len(pks))
// 	for _, pk := range pks {
// 		if col := r.getColumnByName(pk); col == nil {
// 			return fmt.Errorf("failed to generate cq_id for %s, pk column missing %s", r.table.Name, pk)
// 		} else if col.internal {
// 			continue
// 		}

// 		value := r.Get(pk)
// 		if value == nil {
// 			return fmt.Errorf("failed to generate cq_id for %s, pk field missing %s", r.table.Name, pk)
// 		}
// 		objs = append(objs, value)
// 	}
// 	id, err := hashUUID(objs)
// 	if err != nil {
// 		return err
// 	}
// 	r.cqId = id
// 	return nil
// }

func (r *Resource) TableName() string {
	if r.Table == nil {
		return ""
	}
	return r.Table.Name
}

func (r Resource) GetMeta(key string) (interface{}, bool) {
	if r.metadata == nil {
		return nil, false
	}
	v, ok := r.metadata[key]
	return v, ok
}

func (r Resource) getColumnByName(column string) *Column {
	for _, c := range r.Table.Columns {
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
	return rr[0].Table.Name
}

func (rr Resources) ColumnNames() []string {
	if len(rr) == 0 {
		return []string{}
	}
	return rr[0].CColumns
}

func hashUUID(objs interface{}) (uuid.UUID, error) {
	// Use SHA1 because it's fast and is reasonably enough protected against accidental collisions.
	// There is no scenario here where intentional created collisions could do harm.
	digester := crypto.SHA1.New()
	hash, err := hashstructure.Hash(objs, hashstructure.FormatV2, nil)
	if err != nil {
		return uuid.Nil, err
	}
	if _, err := fmt.Fprint(digester, hash); err != nil {
		return uuid.Nil, err
	}
	data := digester.Sum(nil)
	return uuid.NewSHA1(uuid.Nil, data), nil
}
