package schema

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Resources []*Resource

// Resource represents a row in it's associated table, it carries a reference to the original item, and automatically
// generates an Id based on Table's Columns. Resource data can be accessed by the Get and Set methods
type Resource struct {
	// Original resource item that wa from prior resolve
	Item any
	// Set if this is an embedded table
	Parent *Resource
	// internal fields
	Table *Table
	// This is sorted result data by column name
	data CQTypes
}

// This struct is what we send over the wire to destination.
// We dont want to reuse the same struct as otherwise we will have to comment on fields which don't get sent over the wire but still accessible
// code wise
type DestinationResource struct {
	TableName string  `json:"table_name"`
	Data      CQTypes `json:"data"`
}

func NewResourceData(t *Table, parent *Resource, item any) *Resource {
	r := Resource{
		Item:   item,
		Parent: parent,
		Table:  t,
		data:   make(CQTypes, len(t.Columns)),
	}
	for i := range r.data {
		r.data[i] = NewCqTypeFromValueType(t.Columns[i].Type)
	}
	return &r
}

func (r *Resource) ToDestinationResource() DestinationResource {
	dr := DestinationResource{
		TableName: r.Table.Name,
		Data:      r.data,
	}
	return dr
}

func (r *Resource) Get(columnName string) CQType {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	return r.data[index]
}

// Set sets a column with value. This does validation and conversion to
// one of concrete  it returns an error just for backward compatibility
// and panics in case it fails
func (r *Resource) Set(columnName string, value any) error {
	index := r.Table.Columns.Index(columnName)
	if index == -1 {
		// we panic because we want to distinguish between code error and api error
		// this also saves additional checks in our testing code
		panic(columnName + " column not found")
	}
	if err := r.data[index].Set(value); err != nil {
		panic(fmt.Errorf("failed to set column %s: %w", columnName, err))
	}
	return nil
}

// Override original item (this is useful for apis that follow list/details pattern)
func (r *Resource) SetItem(item any) {
	r.Item = item
}

func (r *Resource) GetItem() any {
	return r.Item
}

func (r *Resource) GetValues() CQTypes {
	return r.data
}

func (r *Resource) Columns() []string {
	return r.Table.Columns.Names()
}

//nolint:revive
func (r *Resource) CalculateCQID(deterministicCQID bool) error {
	if !deterministicCQID {
		return r.storeCQID(uuid.New())
	}
	names := r.Table.PrimaryKeys()
	if len(names) == 0 {
		return r.storeCQID(uuid.New())
	}
	slices.Sort(names)
	h := sha256.New()
	for _, name := range names {
		// We need to include the column name in the hash because the same value can be present in multiple columns and therefore lead to the same hash
		h.Write([]byte(name))
		h.Write([]byte(r.Get(name).String()))
	}
	return r.storeCQID(uuid.NewSHA1(uuid.UUID{}, h.Sum(nil)))
}

func (r *Resource) storeCQID(value uuid.UUID) error {
	b, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	return r.Set(CqIDColumn.Name, b)
}

// Validates that all primary keys have values.
func (r *Resource) Validate() error {
	var missingPks []string
	for i, c := range r.Table.Columns {
		if c.CreationOptions.PrimaryKey {
			if r.data[i].GetStatus() != Present {
				missingPks = append(missingPks, c.Name)
			}
		}
	}
	if len(missingPks) > 0 {
		return fmt.Errorf("missing primary key on columns: %v", missingPks)
	}
	return nil
}

func (r *Resource) PrimaryKeyValue() (values CQTypes) {
	for i, col := range r.Table.Columns {
		if col.CreationOptions.PrimaryKey {
			values = append(values, r.data[i])
		}
	}

	return values
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
