package main

import (
	"context"
	"log"

	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/serve"
)

func main() {
	p := plugin.NewPlugin("deadcode", "v1.0.0", memdb.NewMemDBClient)
	if err := serve.Plugin(p).Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
