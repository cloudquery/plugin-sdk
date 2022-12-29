package schema

import "testing"

func TestTextSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result Text
	}{
		{source: "foo", result: Text{Str: "foo", Status: Present}},
		{source: _string("bar"), result: Text{Str: "bar", Status: Present}},
		{source: (*string)(nil), result: Text{Status: Null}},
	}

	for i, tt := range successfulTests {
		var d Text
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestText_LessThan(t *testing.T) {
	cases := []struct {
		a    Text
		b    Text
		want bool
	}{
		{a: Text{Str: "foo", Status: Present}, b: Text{Str: "bar", Status: Present}, want: false},
		{a: Text{Str: "bar", Status: Present}, b: Text{Str: "foo", Status: Present}, want: true},
		{a: Text{Str: "foo", Status: Undefined}, b: Text{Str: "foo", Status: Present}, want: true},
		{a: Text{Str: "foo", Status: Present}, b: Text{Str: "foo", Status: Undefined}, want: false},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
