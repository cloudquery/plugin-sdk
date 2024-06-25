package scheduler

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metrics is deprecated as we move toward open telemetry for tracing and metrics
type Metrics struct {
	TableClient map[string]map[string]*TableClientMetrics
}

type OtelMeters struct {
	resources  metric.Int64Counter
	errors     metric.Int64Counter
	panics     metric.Int64Counter
	duration   metric.Int64Counter
	attributes []attribute.KeyValue
}

type TableClientMetrics struct {
	Resources uint64
	Errors    uint64
	Panics    uint64
	Duration  atomic.Pointer[time.Duration]

	otelMeters *OtelMeters
}

func durationPointerEqual(a, b *time.Duration) bool {
	if a == nil {
		return b == nil
	}
	return b != nil && *a == *b
}

func (s *TableClientMetrics) Equal(other *TableClientMetrics) bool {
	return s.Resources == other.Resources && s.Errors == other.Errors && s.Panics == other.Panics && durationPointerEqual(s.Duration.Load(), other.Duration.Load())
}

// Equal compares to stats. Mostly useful in testing
func (s *Metrics) Equal(other *Metrics) bool {
	for table, clientStats := range s.TableClient {
		for client, stats := range clientStats {
			if _, ok := other.TableClient[table]; !ok {
				return false
			}
			if _, ok := other.TableClient[table][client]; !ok {
				return false
			}
			if !stats.Equal(other.TableClient[table][client]) {
				return false
			}
		}
	}
	for table, clientStats := range other.TableClient {
		for client, stats := range clientStats {
			if _, ok := s.TableClient[table]; !ok {
				return false
			}
			if _, ok := s.TableClient[table][client]; !ok {
				return false
			}
			if !stats.Equal(s.TableClient[table][client]) {
				return false
			}
		}
	}
	return true
}

func getOtelMeters(tableName string, clientID string, invocationID string) *OtelMeters {
	resources, err := otel.Meter(otelName).Int64Counter("sync.table.resources."+tableName,
		metric.WithDescription("Number of resources synced for the "+tableName+" table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	errors, err := otel.Meter(otelName).Int64Counter("sync.table.errors."+tableName,
		metric.WithDescription("Number of errors encountered while syncing the "+tableName+" table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	panics, err := otel.Meter(otelName).Int64Counter("sync.table.panics."+tableName,
		metric.WithDescription("Number of panics encountered while syncing the "+tableName+" table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	duration, err := otel.Meter(otelName).Int64Counter("sync.table.duration."+tableName,
		metric.WithDescription("Duration of syncing the "+tableName+" table"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil
	}

	return &OtelMeters{
		resources: resources,
		errors:    errors,
		panics:    panics,
		duration:  duration,
		attributes: []attribute.KeyValue{
			attribute.Key("sync.client.id").String(clientID),
			attribute.Key("sync.invocation.id").String(invocationID),
		},
	}
}

func (s *Metrics) initWithClients(table *schema.Table, clients []schema.ClientMeta, invocationID string) {
	s.TableClient[table.Name] = make(map[string]*TableClientMetrics, len(clients))
	for _, client := range clients {
		tableName := table.Name
		clientID := client.ID()
		s.TableClient[tableName][clientID] = &TableClientMetrics{
			otelMeters: getOtelMeters(tableName, clientID, invocationID),
		}
	}
	for _, relation := range table.Relations {
		s.initWithClients(relation, clients, invocationID)
	}
}

func (s *Metrics) TotalErrors() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Errors
		}
	}
	return total
}

func (s *Metrics) TotalErrorsAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.Errors)
		}
	}
	return total
}

func (s *Metrics) TotalPanics() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Panics
		}
	}
	return total
}

func (s *Metrics) TotalPanicsAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.Panics)
		}
	}
	return total
}

func (s *Metrics) TotalResources() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Resources
		}
	}
	return total
}

func (s *Metrics) TotalResourcesAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.Resources)
		}
	}
	return total
}

func (m *TableClientMetrics) OtelResourcesAdd(ctx context.Context, count int64) {
	if m.otelMeters == nil {
		return
	}

	m.otelMeters.resources.Add(ctx, count, metric.WithAttributes(m.otelMeters.attributes...))
}

func (m *TableClientMetrics) OtelErrorsAdd(ctx context.Context, count int64) {
	if m.otelMeters == nil {
		return
	}

	m.otelMeters.errors.Add(ctx, count, metric.WithAttributes(m.otelMeters.attributes...))
}

func (m *TableClientMetrics) OtelPanicsAdd(ctx context.Context, count int64) {
	if m.otelMeters == nil {
		return
	}

	m.otelMeters.panics.Add(ctx, count, metric.WithAttributes(m.otelMeters.attributes...))
}

func (m *TableClientMetrics) OtelDurationRecord(ctx context.Context, duration time.Duration) {
	if m.otelMeters == nil {
		return
	}

	m.otelMeters.duration.Add(ctx, duration.Milliseconds(), metric.WithAttributes(m.otelMeters.attributes...))
}
