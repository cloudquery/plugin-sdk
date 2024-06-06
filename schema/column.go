package schema

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/apache/arrow/go/v16/arrow"
)

type ColumnList []Column

// ColumnResolver is called for each row received in TableResolver's data fetch.
// execution holds all relevant information regarding execution as well as the Column called.
// resource holds the current row we are resolving the column for.
type ColumnResolver func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error

// Column definition for Table
type Column struct {
	// Name of column
	Name string `json:"name"`
	// Value Type of column i.e String, UUID etc'
	Type arrow.DataType `json:"type"`
	// Description about column, this description is added as a comment in the database
	Description string `json:"description"`
	// Column Resolver allows to set your own data for a column; this can be an API call, setting multiple embedded values, etc
	Resolver ColumnResolver `json:"-"`

	// IgnoreInTests is used to skip verifying the column is non-nil in integration tests.
	// By default, integration tests perform a fetch for all resources in cloudquery's test account, and
	// verify all columns are non-nil.
	// If IgnoreInTests is true, verification is skipped for this column.
	// Used when it is hard to create a reproducible environment with this column being non-nil (e.g. various error columns).
	IgnoreInTests bool `json:"-"`

	// PrimaryKey requires the destinations supporting this to include this column into the primary key
	PrimaryKey bool `json:"primary_key"`
	// NotNull requires the destinations supporting this to mark this column as non-nullable
	NotNull bool `json:"not_null"`
	// IncrementalKey is a flag that indicates if the column is used as part of an incremental key.
	// It is mainly used for documentation purposes, but may also be used as part of ensuring that
	// migrations are done correctly.
	IncrementalKey bool `json:"incremental_key"`
	// Unique requires the destinations supporting this to mark this column as unique
	Unique bool `json:"unique"`

	// PrimaryKeyComponent is a flag that indicates if the column is used as part of the input to calculate the value of `_cq_id`.
	PrimaryKeyComponent bool `json:"primary_key_component"`
}

// NewColumnFromArrowField creates a new Column from an arrow.Field
// arrow.Field is a low-level representation of a CloudQuery column
// that can be sent over the wire in a cross-language way.
func NewColumnFromArrowField(f arrow.Field) Column {
	column := Column{
		Name:    f.Name,
		Type:    f.Type,
		NotNull: !f.Nullable,
	}

	v, ok := f.Metadata.GetValue(MetadataPrimaryKey)
	column.PrimaryKey = ok && v == MetadataTrue

	v, ok = f.Metadata.GetValue(MetadataUnique)
	column.Unique = ok && v == MetadataTrue

	v, ok = f.Metadata.GetValue(MetadataIncremental)
	column.IncrementalKey = ok && v == MetadataTrue

	v, ok = f.Metadata.GetValue(MetadataPrimaryKeyComponent)
	column.PrimaryKeyComponent = ok && v == MetadataTrue

	return column
}

func (c Column) ToArrowField() arrow.Field {
	mdKV := map[string]string{
		MetadataPrimaryKey:          MetadataFalse,
		MetadataUnique:              MetadataFalse,
		MetadataIncremental:         MetadataFalse,
		MetadataPrimaryKeyComponent: MetadataFalse,
	}
	if c.PrimaryKey {
		mdKV[MetadataPrimaryKey] = MetadataTrue
	}
	if c.Unique {
		mdKV[MetadataUnique] = MetadataTrue
	}
	if c.IncrementalKey {
		mdKV[MetadataIncremental] = MetadataTrue
	}
	if c.PrimaryKeyComponent {
		mdKV[MetadataPrimaryKeyComponent] = MetadataTrue
	}

	return arrow.Field{
		Name:     c.Name,
		Type:     c.Type,
		Nullable: !c.NotNull,
		Metadata: arrow.MetadataFrom(mdKV),
	}
}

func (c Column) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name                string `json:"name"`
		Type                string `json:"type"`
		Description         string `json:"description"`
		PrimaryKey          bool   `json:"primary_key"`
		NotNull             bool   `json:"not_null"`
		Unique              bool   `json:"unique"`
		IncrementalKey      bool   `json:"incremental_key"`
		PrimaryKeyComponent bool   `json:"primary_key_component"`
	}
	var alias Alias
	alias.Name = c.Name
	alias.Type = c.Type.String()
	alias.Description = c.Description
	alias.PrimaryKey = c.PrimaryKey
	alias.NotNull = c.NotNull
	alias.Unique = c.Unique
	alias.IncrementalKey = c.IncrementalKey
	alias.PrimaryKeyComponent = c.PrimaryKeyComponent

	return json.Marshal(alias)
}

func (c Column) String() string {
	var sb strings.Builder
	sb.WriteString(c.Name)
	sb.WriteString(":")
	sb.WriteString(c.Type.String())
	if c.PrimaryKey {
		sb.WriteString(":PK")
	}
	if c.NotNull {
		sb.WriteString(":NotNull")
	}
	if c.Unique {
		sb.WriteString(":Unique")
	}
	if c.IncrementalKey {
		sb.WriteString(":IncrementalKey")
	}

	if c.PrimaryKeyComponent {
		sb.WriteString(":PrimaryKeyComponent")
	}
	return sb.String()
}

func (c ColumnList) Index(col string) int {
	for i, c := range c {
		if c.Name == col {
			return i
		}
	}
	return -1
}

func (c ColumnList) Names() []string {
	ret := make([]string, len(c))
	for i := range c {
		ret[i] = c[i].Name
	}
	return ret
}

func (c ColumnList) Get(name string) *Column {
	for i := range c {
		if c[i].Name == name {
			return &c[i]
		}
	}
	return nil
}

func (c ColumnList) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, col := range c {
		sb.WriteString(col.String())
		if i != len(c)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}
