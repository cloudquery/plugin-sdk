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
		table: &Table{
			Columns: []Column{
				{
					Name: "test",
					Type: TypeString,
				},
				{
					Name: "int_value",
					Type: TypeInt,
				},
				{
					Name: "unexported",
					Type: TypeBool,
				},
			},
		},
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

func TestInterfaceSlice(t *testing.T) {
	var sType []interface{}
	var names []string
	names = append(names, "first", "second")
	assert.IsTypef(t, sType, interfaceSlice(names), "")
	assert.IsTypef(t, sType, interfaceSlice(&names), "")
	assert.IsTypef(t, sType, interfaceSlice(1), "")
	assert.IsTypef(t, sType, interfaceSlice(innerStruct{"asdsa"}), "")
	assert.IsTypef(t, sType, interfaceSlice(&innerStruct{"asdsa"}), "")
	pSlice := []*innerStruct{{"asdsa"}, {"asdsa"}, {"asdsa"}}
	assert.IsTypef(t, sType, interfaceSlice(pSlice), "")
	assert.IsTypef(t, sType, interfaceSlice(&pSlice), "")

}
