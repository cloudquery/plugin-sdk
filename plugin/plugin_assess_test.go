package plugin

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func TestAssessTablesWithoutAssessorReturnsUnknown(t *testing.T) {
	ctx := context.Background()
	p := NewPlugin("test", "v1.0.0", newTestPluginClient)
	require.NoError(t, p.Init(ctx, nil, NewClientOptions{NoConnection: true}))

	findings, err := p.AssessTables(ctx, []TablePair{
		{Old: &schema.Table{Name: "changed"}, New: &schema.Table{Name: "changed"}},
		{New: &schema.Table{Name: "added"}},
		{Old: &schema.Table{Name: "removed"}},
	}, AssessOptions{})
	require.NoError(t, err)

	unknown := func(name string) TableFinding {
		return TableFinding{TableName: name, Category: AssessCategoryUnknown, CoverageIncomplete: true, CoverageIncompleteReason: AssessNotSupportedReason}
	}
	require.Equal(t, []TableFinding{unknown("changed"), unknown("added"), unknown("removed")}, findings)
}

func TestAssessTablesBeforeInit(t *testing.T) {
	p := NewPlugin("test", "v1.0.0", newTestPluginClient)
	_, err := p.AssessTables(context.Background(), nil, AssessOptions{})
	require.Error(t, err)
}
