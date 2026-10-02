package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

var ErrUnableToCompare = errors.New("unable to compare")

type SyntheticCase string

const (
	SyntheticCasePopulated       SyntheticCase = "populated"
	SyntheticCaseNull            SyntheticCase = "null"
	SyntheticCaseEmptyCollection SyntheticCase = "empty collection"
	SyntheticCaseNested          SyntheticCase = "nested"
	SyntheticCaseEscaping        SyntheticCase = "escaping"
)

var syntheticCases = []SyntheticCase{
	SyntheticCasePopulated,
	SyntheticCaseNull,
	SyntheticCaseEmptyCollection,
	SyntheticCaseNested,
	SyntheticCaseEscaping,
}

const syntheticEscapingString = `"comma, \"quote\", back\\slash\nnew line\ttab é"`

// SyntheticPair is a value that means the same under the old and the new field.
// Value is its Arrow JSON form, which both fields accept and read back as equal values.
type SyntheticPair struct {
	Case  SyntheticCase
	Value string
}

// SyntheticPairs returns equivalent values for each case both fields support.
// It returns ErrUnableToCompare when no populated value means the same under both fields.
func SyntheticPairs(oldField, newField arrow.Field) ([]SyntheticPair, error) {
	var pairs []SyntheticPair
	for _, c := range syntheticCases {
		if c == SyntheticCaseNull {
			if oldField.Nullable && newField.Nullable {
				pairs = append(pairs, SyntheticPair{Case: c, Value: "null"})
			}
			continue
		}
		candidates := append(syntheticCandidates(oldField.Type, c), syntheticCandidates(newField.Type, c)...)
		for _, candidate := range candidates {
			if isEquivalent(oldField.Type, newField.Type, candidate) {
				pairs = append(pairs, SyntheticPair{Case: c, Value: candidate})
				break
			}
		}
	}
	if len(pairs) == 0 || pairs[0].Case != SyntheticCasePopulated {
		return nil, fmt.Errorf("%w: no equivalent value for %s and %s", ErrUnableToCompare, oldField.Type, newField.Type)
	}
	return pairs, nil
}

// SyntheticRecords builds one-column records for the old and the new field, with one row per pair.
func SyntheticRecords(oldField, newField arrow.Field, pairs []SyntheticPair) (oldRecord, newRecord arrow.RecordBatch, err error) {
	values := make([]string, len(pairs))
	for i, pair := range pairs {
		values[i] = pair.Value
	}
	rows := "[" + strings.Join(values, ",") + "]"
	oldRecord, err = syntheticRecord(oldField, rows)
	if err != nil {
		return nil, nil, err
	}
	newRecord, err = syntheticRecord(newField, rows)
	if err != nil {
		oldRecord.Release()
		return nil, nil, err
	}
	return oldRecord, newRecord, nil
}

func syntheticRecord(field arrow.Field, rows string) (arrow.RecordBatch, error) {
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, arrow.NewSchema([]arrow.Field{field}, nil))
	defer bldr.Release()
	if err := bldr.Field(0).UnmarshalJSON([]byte(rows)); err != nil {
		return nil, fmt.Errorf("failed to build synthetic values for field %s: %w", field.Name, err)
	}
	return bldr.NewRecordBatch(), nil
}

func isEquivalent(oldType, newType arrow.DataType, value string) bool {
	oldValue, ok := readBack(oldType, value)
	if !ok {
		return false
	}
	newValue, ok := readBack(newType, value)
	return ok && reflect.DeepEqual(oldValue, newValue)
}

func readBack(dataType arrow.DataType, value string) (any, bool) {
	bldr := array.NewBuilder(memory.DefaultAllocator, dataType)
	defer bldr.Release()
	if err := bldr.UnmarshalJSON([]byte("[" + value + "]")); err != nil || bldr.Len() != 1 {
		return nil, false
	}
	arr := bldr.NewArray()
	defer arr.Release()
	marshaled, err := json.Marshal(arr.GetOneForMarshal(0))
	if err != nil {
		return nil, false
	}
	var decoded any
	if err := json.NewDecoder(bytes.NewReader(marshaled)).Decode(&decoded); err != nil {
		return nil, false
	}
	return decoded, true
}

func syntheticCandidates(dataType arrow.DataType, c SyntheticCase) []string {
	switch dt := dataType.(type) {
	case *arrow.ListType:
		return syntheticListCandidates(dt.Elem(), c)
	case *arrow.LargeListType:
		return syntheticListCandidates(dt.Elem(), c)
	case *arrow.MapType:
		if c == SyntheticCaseNested && !isCollection(dt.ItemType()) {
			return nil
		}
		return syntheticListCandidates(dt.Elem(), c)
	case *arrow.StructType:
		return syntheticStructCandidates(dt, c)
	case *types.JSONType:
		return syntheticJSONCandidates(c)
	}
	if c == SyntheticCasePopulated {
		return syntheticScalarCandidates(dataType)
	}
	if c == SyntheticCaseEscaping && isStringLike(dataType) {
		return []string{syntheticEscapingString}
	}
	return nil
}

func syntheticListCandidates(elem arrow.DataType, c SyntheticCase) []string {
	switch c {
	case SyntheticCaseEmptyCollection:
		return []string{"[]"}
	case SyntheticCaseNested:
		if !isCollection(elem) {
			return nil
		}
		return wrapInList(syntheticCandidates(elem, SyntheticCasePopulated))
	default:
		return wrapInList(syntheticCandidates(elem, c))
	}
}

func wrapInList(values []string) []string {
	wrapped := make([]string, len(values))
	for i, v := range values {
		wrapped[i] = "[" + v + "]"
	}
	return wrapped
}

func syntheticStructCandidates(dt *arrow.StructType, c SyntheticCase) []string {
	if c == SyntheticCaseEmptyCollection {
		return nil
	}
	if c == SyntheticCaseNested && !hasCollectionField(dt) {
		return nil
	}
	hasEscaping := false
	fields := make([]string, dt.NumFields())
	for i, field := range dt.Fields() {
		value := firstCandidate(field.Type, SyntheticCasePopulated)
		if c == SyntheticCaseEscaping {
			if escaping := firstCandidate(field.Type, SyntheticCaseEscaping); escaping != "" {
				value = escaping
				hasEscaping = true
			}
		}
		if value == "" {
			return nil
		}
		name, _ := json.Marshal(field.Name)
		fields[i] = string(name) + ":" + value
	}
	if c == SyntheticCaseEscaping && !hasEscaping {
		return nil
	}
	return []string{"{" + strings.Join(fields, ",") + "}"}
}

func firstCandidate(dataType arrow.DataType, c SyntheticCase) string {
	candidates := syntheticCandidates(dataType, c)
	if len(candidates) == 0 {
		return ""
	}
	return candidates[0]
}

func hasCollectionField(dt *arrow.StructType) bool {
	for _, field := range dt.Fields() {
		if isCollection(field.Type) {
			return true
		}
	}
	return false
}

func isCollection(dataType arrow.DataType) bool {
	switch dataType.(type) {
	case *arrow.ListType, *arrow.LargeListType, *arrow.MapType, *arrow.StructType, *types.JSONType:
		return true
	}
	return false
}

func syntheticJSONCandidates(c SyntheticCase) []string {
	switch c {
	case SyntheticCasePopulated:
		return []string{`{"env":"prod"}`}
	case SyntheticCaseEmptyCollection:
		return []string{`{}`, `[]`}
	case SyntheticCaseNested:
		return []string{`{"env":{"tags":["prod"]}}`}
	case SyntheticCaseEscaping:
		return []string{`{"note":` + syntheticEscapingString + `}`}
	}
	return nil
}

func isStringLike(dataType arrow.DataType) bool {
	return arrow.TypeEqual(dataType, arrow.BinaryTypes.String) || arrow.TypeEqual(dataType, arrow.BinaryTypes.LargeString)
}

func syntheticScalarCandidates(dataType arrow.DataType) []string {
	switch {
	case isStringLike(dataType):
		return []string{`"env:prod"`}
	case arrow.TypeEqual(dataType, types.ExtensionTypes.UUID):
		return []string{`"6ba7b810-9dad-11d1-80b4-00c04fd430c8"`}
	case arrow.TypeEqual(dataType, types.ExtensionTypes.Inet):
		return []string{`"192.0.2.1/24"`}
	case arrow.TypeEqual(dataType, types.ExtensionTypes.MAC):
		return []string{`"00:00:5e:00:53:01"`}
	case arrow.IsInteger(dataType.ID()), arrow.IsFloating(dataType.ID()), dataType.ID() == arrow.DURATION:
		return []string{`42`}
	case arrow.IsDecimal(dataType.ID()):
		return []string{`"12.5"`}
	}
	switch dataType.ID() {
	case arrow.BOOL:
		return []string{`true`}
	case arrow.BINARY, arrow.LARGE_BINARY:
		return []string{`"ZW52OnByb2Q="`}
	case arrow.TIMESTAMP:
		return []string{`"2024-01-02T03:04:05Z"`}
	case arrow.DATE32, arrow.DATE64:
		return []string{`"2024-01-02"`}
	case arrow.TIME32, arrow.TIME64:
		return []string{`"03:04:05"`}
	case arrow.INTERVAL_MONTHS:
		return []string{`{"months":1}`}
	case arrow.INTERVAL_DAY_TIME:
		return []string{`{"days":1,"milliseconds":1}`}
	case arrow.INTERVAL_MONTH_DAY_NANO:
		return []string{`{"months":1,"days":1,"nanoseconds":1}`}
	}
	return nil
}
