package schema

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type innerStruct struct {
	Value string
}

type testStruct struct {
	Inner      innerStruct
	Value      int
	unexported bool
}

func TestPathResolver(t *testing.T) {
	r1 := PathResolver("Inner.Value")
	r2 := PathResolver("Value")
	r3 := PathResolver("unexported")
	resource := &Resource{
		Item: testStruct{Inner: innerStruct{Value: "bla"}, Value: 5, unexported: false},
		data: map[string]interface{}{},
	}
	err := r1(context.TODO(), nil, resource, Column{Name: "test"})

	assert.Nil(t, err)
	assert.Equal(t, resource.Get("test"), "bla")

	err = r2(context.TODO(), nil, resource, Column{Name: "int_value"})

	assert.Nil(t, err)
	assert.Equal(t, resource.Get("int_value"), 5)

	err = r3(context.TODO(), nil, resource, Column{Name: "unexported"})
	assert.Nil(t, err)
	assert.Nil(t, resource.Get("unexported"))
}
