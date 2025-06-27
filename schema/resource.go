package schema

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"slices"

	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/google/uuid"
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
	data scalar.Vector
	// bldr array.RecordBuilder
}

func NewResourceData(t *Table, parent *Resource, item any) *Resource {
	r := Resource{
		Item:   item,
		Parent: parent,
		Table:  t,
		data:   make(scalar.Vector, len(t.Columns)),
	}
	for i := range r.data {
		r.data[i] = scalar.NewScalar(t.Columns[i].Type)
	}
	return &r
}

func (r *Resource) Get(columnName string) scalar.Scalar {
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

func (r *Resource) GetValues() scalar.Vector {
	return r.data
}

//nolint:revive
func (r *Resource) CalculateCQID(deterministicCQID bool) error {
	// if `PrimaryKeyComponent` is set, we calculate the CQID based on those components
	pkComponents := r.Table.PrimaryKeyComponents()
	if len(pkComponents) > 0 {
		return r.storeCQID(newUUID(uuid.UUID{}, calculateCqIDValue(r, pkComponents).Sum(nil)))
	}

	// If deterministicCQID is false, we generate a random CQID
	if !deterministicCQID {
		return r.storeCQID(uuid.New())
	}
	names := r.Table.PrimaryKeys()
	// If there are no primary keys or if CQID is the only PK, we generate a random CQID
	if len(names) == 0 || (len(names) == 1 && names[0] == CqIDColumn.Name) {
		return r.storeCQID(uuid.New())
	}

	return r.storeCQID(newUUID(uuid.UUID{}, calculateCqIDValue(r, names).Sum(nil)))
}

func calculateCqIDValue(r *Resource, cols []string) hash.Hash {
	h := sha256.New()
	slices.Sort(cols)
	for _, col := range cols {
		// We need to include the column name in the hash because the same value can be present in multiple columns and therefore lead to the same hash
		h.Write([]byte(col))
		h.Write([]byte(r.Get(col).String()))
	}
	return h
}

func (r *Resource) storeCQID(value uuid.UUID) error {
	// We skip if _cq_id is not present.
	// Mostly the problem here is because the transformation step is baked into the resolving step
	if r.Table.Columns.Get(CqIDColumn.Name) == nil {
		return nil
	}
	b, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	return r.Set(CqIDColumn.Name, b)
}

func (r *Resource) StoreCQClientID(clientID string) error {
	// We skip if _cq_client_id is not present.
	if r.Table.Columns.Get(CqClientIDColumn.Name) == nil {
		return nil
	}
	return r.Set(CqClientIDColumn.Name, clientID)
}

type PKError struct {
	MissingPKs []string
}

func (e *PKError) Error() string {
	return fmt.Sprintf("missing primary key on columns: %v", e.MissingPKs)
}

type PKComponentError struct {
	MissingPKComponents []string
}

func (e *PKComponentError) Error() string {
	return fmt.Sprintf("missing primary key component on columns: %v", e.MissingPKComponents)
}

// Validates that all primary keys have values.
func (r *Resource) Validate() error {
	var missingPks []string
	var missingPKComponents []string
	for i, c := range r.Table.Columns {
		switch {
		case c.PrimaryKey && !r.data[i].IsValid():
			missingPks = append(missingPks, c.Name)
		case c.PrimaryKeyComponent && !r.data[i].IsValid():
			missingPKComponents = append(missingPKComponents, c.Name)
		}
	}
	if len(missingPks) > 0 {
		return &PKError{MissingPKs: missingPks}
	}
	if len(missingPKComponents) > 0 {
		return &PKComponentError{MissingPKComponents: missingPKComponents}
	}
	return nil
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
