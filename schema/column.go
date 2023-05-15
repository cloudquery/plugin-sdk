package schema

import (
	"context"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
)

type ColumnList []Column

// ColumnResolver is called for each row received in TableResolver's data fetch.
// execution holds all relevant information regarding execution as well as the Column called.
// resource holds the current row we are resolving the column for.
type ColumnResolver func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error

// CreationOptions allow modification of how column is defined when table is created
type CreationOptions struct {
	PrimaryKey bool
	NotNull    bool
	// IncrementalKey is a flag that indicates if the column is used as part of an incremental key.
	// It is mainly used for documentation purposes, but may also be used as part of ensuring that
	// migrations are done correctly.
	IncrementalKey bool
	Unique         bool
}

// Column definition for Table
type Column struct {
	// Name of column
	Name string
	// Value Type of column i.e String, UUID etc'
	Type arrow.DataType
	// Description about column, this description is added as a comment in the database
	Description string
	// Column Resolver allows to set your own data for a column; this can be an API call, setting multiple embedded values, etc
	Resolver ColumnResolver
	// CreationOptions allow modifying how column is defined when table is created
	CreationOptions
	// IgnoreInTests is used to skip verifying the column is non-nil in integration tests.
	// By default, integration tests perform a fetch for all resources in cloudquery's test account, and
	// verify all columns are non-nil.
	// If IgnoreInTests is true, verification is skipped for this column.
	// Used when it is hard to create a reproducible environment with this column being non-nil (e.g. various error columns).
	IgnoreInTests bool
}

// NewColumnFromArrowField creates a new Column from an arrow.Field
// arrow.Field is a low-level representation of a CloudQuery column
// that can be sent over the wire in a cross-language way.
func NewColumnFromArrowField(f arrow.Field) Column {
	creationOptions := CreationOptions{NotNull: !f.Nullable}

	v, ok := f.Metadata.GetValue(MetadataPrimaryKey)
	creationOptions.PrimaryKey = ok && v == MetadataTrue

	v, ok = f.Metadata.GetValue(MetadataUnique)
	creationOptions.Unique = ok && v == MetadataTrue

	v, ok = f.Metadata.GetValue(MetadataIncremental)
	creationOptions.IncrementalKey = ok && v == MetadataTrue

	return Column{
		Name:            f.Name,
		Type:            f.Type,
		CreationOptions: creationOptions,
	}
}

func (c Column) ToArrowField() arrow.Field {
	mdKV := map[string]string{
		MetadataPrimaryKey:  MetadataFalse,
		MetadataUnique:      MetadataFalse,
		MetadataIncremental: MetadataFalse,
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

	return arrow.Field{
		Name:     c.Name,
		Type:     c.Type,
		Nullable: !c.CreationOptions.NotNull,
		Metadata: arrow.MetadataFrom(mdKV),
	}
}

func (c Column) String() string {
	var sb strings.Builder
	sb.WriteString(c.Name)
	sb.WriteString(":")
	sb.WriteString(c.Type.String())
	if c.CreationOptions.PrimaryKey {
		sb.WriteString(":PK")
	}
	if c.CreationOptions.NotNull {
		sb.WriteString(":NotNull")
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
