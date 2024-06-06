package destination

import (
	"bytes"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/ipc"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// Legacy conversion functions to and from Arrow bytes. From plugin v3 onwards
// this responsibility is handled by plugin-pb-go.

func NewFromBytes(b []byte) (*arrow.Schema, error) {
	rdr, err := ipc.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return rdr.Schema(), nil
}

func NewSchemasFromBytes(b [][]byte) (schema.Schemas, error) {
	var err error
	ret := make([]*arrow.Schema, len(b))
	for i, buf := range b {
		ret[i], err = NewFromBytes(buf)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
