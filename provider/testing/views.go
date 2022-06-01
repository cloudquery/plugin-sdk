package testing

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/testlog"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type ViewTestCase struct {
	// Provider to configure and create tables before executing the view
	Provider *provider.Provider
	// SQLView statement that create that view
	SQLView string
}

func HelperTestView(t *testing.T, resource ViewTestCase) {
	t.Helper()

	conn, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}

	l := testlog.New(t)
	l.SetLevel(hclog.Info)
	resource.Provider.Logger = l

	for _, table := range resource.Provider.ResourceMap {
		if err := dropAndCreateTable(context.Background(), conn, table); err != nil {
			assert.FailNow(t, fmt.Sprintf("failed to create tables %s", table.Name), err)
		}
	}

	if err := conn.Exec(context.Background(), resource.SQLView); err != nil {
		t.Fatal(err)
	}
}
