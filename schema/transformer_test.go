package schema

import "testing"

func TestTransformWithTransformer(t *testing.T) {
	cqTypes := make(CQTypes, 0, TypeEnd-TypeInvalid)
	for i := TypeInvalid + 1; i < TypeEnd; i++ {
		if deprecatedTypesValues.isDeprecated(i) {
			continue
		}
		cqTypes = append(cqTypes, NewCqTypeFromValueType(i))
	}
	transformedVal, err := TransformWithTransformer(&DefaultTransformer{}, cqTypes)
	if err != nil {
		t.Fatal(err)
	}
	if transformedVal == nil {
		t.Fatal("expected not nil")
	}
}
