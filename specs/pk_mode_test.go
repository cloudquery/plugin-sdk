package specs

import (
	"testing"
)

func TestPKModeFromString(t *testing.T) {
	var pkMode PKMode
	if err := pkMode.UnmarshalJSON([]byte(`"cq-id"`)); err != nil {
		t.Fatal(err)
	}
	if pkMode != PKModeCQID {
		t.Fatalf("expected PKModeCQID, got %v", pkMode)
	}
	if err := pkMode.UnmarshalJSON([]byte(`"composite-keys"`)); err != nil {
		t.Fatal(err)
	}
	if pkMode != PKModeCompositeKeys {
		t.Fatalf("expected PKModeCompositeKeys, got %v", pkMode)
	}
}

func TestTestPKMode(t *testing.T) {
	for _, pkModeStr := range pkModeStrings {
		pkMode, err := PKModeFromString(pkModeStr)
		if err != nil {
			t.Fatal(err)
		}
		if pkModeStr != pkMode.String() {
			t.Fatalf("expected:%s got:%s", pkModeStr, pkMode.String())
		}
	}
}
