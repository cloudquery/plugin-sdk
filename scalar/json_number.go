package scalar

import (
	"bytes"
	"encoding/json"
)

// unmarshalJSONWithNumbers decodes JSON with UseNumber so int64/uint64 values that
// exceed float64 precision survive as json.Number instead of being rounded.
func unmarshalJSONWithNumbers(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(v)
}
