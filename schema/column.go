package schema

import (
	"context"
	"encoding/json"
	"fmt"
)

type ColumnList []Column

// ColumnResolver is called for each row received in TableResolver's data fetch.
// execution holds all relevant information regarding execution as well as the Column called.
// resource holds the current row we are resolving the column for.
type ColumnResolver func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error

// ColumnCreationOptions allow modification of how column is defined when table is created
type ColumnCreationOptions struct {
	PrimaryKey bool `json:"primary_key,omitempty"`
}

// Column definition for Table
type Column struct {
	// Name of column
	Name string `json:"name,omitempty"`
	// Value Type of column i.e String, UUID etc'
	Type ValueType `json:"type,omitempty"`
	// Description about column, this description is added as a comment in the database
	Description string `json:"-"`
	// Column Resolver allows to set you own data based on resolving this can be an API call or setting multiple embedded values etc'
	Resolver ColumnResolver `json:"-"`
	// Creation options allow modifying how column is defined when table is created
	CreationOptions ColumnCreationOptions `json:"creation_options,omitempty"`
	// IgnoreInTests is used to skip verifying the column is non-nil in integration tests.
	// By default, integration tests perform a fetch for all resources in cloudquery's test account, and
	// verify all columns are non-nil.
	// If IgnoreInTests is true, verification is skipped for this column.
	// Used when it is hard to create a reproducible environment with this column being non-nil (e.g. various error columns).
	IgnoreInTests bool `json:"-"`
}

func (c *ColumnList) UnmarshalJSON(data []byte) (err error) {
	var tmp []Column
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal column list: %w, %s", err, data)
	}
	res := make(ColumnList, 0, len(tmp))
	for _, column := range tmp {
		if column.Type != TypeInvalid {
			res = append(res, column)
		}
	}
	*c = res
	return nil
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
