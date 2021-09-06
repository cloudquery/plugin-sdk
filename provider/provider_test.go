package provider

import (
	"testing"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"

	"github.com/stretchr/testify/assert"
)

var (
	provider = Provider{
		ResourceMap: map[string]*schema.Table{
			"test": {
				Name: "sdk_test",
				Relations: []*schema.Table{
					{
						Name: "sdk_test_test1",
						Relations: []*schema.Table{
							{Name: "sdk_test_test1_test"},
						},
					},
				},
			},
			"test1": {
				Name:      "sdk_test1",
				Relations: []*schema.Table{},
			},
		},
	}

	failProvider = Provider{
		ResourceMap: map[string]*schema.Table{
			"test": {
				Name: "sdk_test",
				Relations: []*schema.Table{
					{
						Name: "sdk_test1",
					},
				},
			},
			"test1": {
				Name: "sdk_test1",
			},
		},
	}
)

func TestProviderInterpolate(t *testing.T) {
	r, err := provider.interpolateAllResources([]string{"test"})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test"}, r)

	r, err = provider.interpolateAllResources([]string{"test", "test1"})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test", "test1"}, r)

	r, err = provider.interpolateAllResources([]string{"test", "test1", "*"})
	assert.Error(t, err)
	assert.Nil(t, r)
	r, err = provider.interpolateAllResources([]string{"*"})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test", "test1"}, r)

}

func TestTableDuplicates(t *testing.T) {
	tables := make(map[string]string)
	var err error
	for r, t := range provider.ResourceMap {
		err = getTableDuplicates(r, t, tables)
	}
	assert.Nil(t, err)

	for r, t := range failProvider.ResourceMap {
		err = getTableDuplicates(r, t, tables)
	}
	assert.Error(t, err)
}
