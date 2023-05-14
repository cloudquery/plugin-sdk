package schema

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
)

func TestSchemaEncode(t *testing.T) {
	md := arrow.NewMetadata([]string{"true"}, []string{"false"})
	md1 := arrow.NewMetadata([]string{"false"}, []string{"true"})
	schemas := Schemas{
		arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "name", Type: arrow.BinaryTypes.String},
			},
			&md,
		),
		arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "name", Type: arrow.BinaryTypes.String},
			},
			&md1,
		),
	}
	b, err := schemas.Encode()
	if err != nil {
		t.Fatal(err)
	}
	decodedSchemas, err := NewSchemasFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(decodedSchemas) != len(schemas) {
		t.Fatalf("expected %d schemas, got %d", len(schemas), len(decodedSchemas))
	}
	for i := range schemas {
		if !schemas[i].Equal(decodedSchemas[i]) {
			t.Fatalf("expected schema %d to be %v, got %v", i, schemas[i], decodedSchemas[i])
		}
	}
}