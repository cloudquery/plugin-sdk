package metrics

import "go.opentelemetry.io/otel/attribute"

type Selector struct {
	attribute.Set

	clientID  string
	tableName string
}
