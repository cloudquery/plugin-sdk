package serve

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// newResource returns a resource describing this application.
func newResource(p *plugin.Plugin) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cloudquery-"+p.Name()),
			semconv.ServiceVersion(p.Version()),
		),
	)
	if err != nil {
		panic(err)
	}
	return r
}
