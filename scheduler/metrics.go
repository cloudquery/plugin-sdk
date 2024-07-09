package scheduler

import (
	"context"
	"sync"
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
	resources           metric.Int64Counter
	errors              metric.Int64Counter
	panics              metric.Int64Counter
	startTime           metric.Int64Counter
	started             bool
	startedLock         sync.Mutex
	endTime             metric.Int64Counter
	previousEndTime     int64
	previousEndTimeLock sync.Mutex
	attributes          []attribute.KeyValue
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

func (m *TableClientMetrics) Equal(other *TableClientMetrics) bool {
	return m.Resources == other.Resources && m.Errors == other.Errors && m.Panics == other.Panics && durationPointerEqual(m.Duration.Load(), other.Duration.Load())
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

func getOtelMeters(tableName string, clientID string) *OtelMeters {
	resources, err := otel.Meter(otelName).Int64Counter("sync.table.resources",
		metric.WithDescription("Number of resources synced for a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	errors, err := otel.Meter(otelName).Int64Counter("sync.table.errors",
		metric.WithDescription("Number of errors encountered while syncing a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	panics, err := otel.Meter(otelName).Int64Counter("sync.table.panics",
		metric.WithDescription("Number of panics encountered while syncing a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	startTime, err := otel.Meter(otelName).Int64Counter("sync.table.start_time",
		metric.WithDescription("Start time of syncing a table"),
		metric.WithUnit("ns"),
	)
	if err != nil {
		return nil
	}

	endTime, err := otel.Meter(otelName).Int64Counter("sync.table.end_time",
		metric.WithDescription("End time of syncing a table"),
		metric.WithUnit("ns"),
	)

	if err != nil {
		return nil
	}

	return &OtelMeters{
		resources: resources,
		errors:    errors,
		panics:    panics,
		startTime: startTime,
		endTime:   endTime,
		attributes: []attribute.KeyValue{
			attribute.Key("sync.client.id").String(clientID),
			attribute.Key("sync.table.name").String(tableName),
		},
	}
}

func (s *Metrics) initWithClients(table *schema.Table, clients []schema.ClientMeta) {
	s.TableClient[table.Name] = make(map[string]*TableClientMetrics, len(clients))
	for _, client := range clients {
		tableName := table.Name
		clientID := client.ID()
		s.TableClient[tableName][clientID] = &TableClientMetrics{
			otelMeters: getOtelMeters(tableName, clientID),
		}
	}
	for _, relation := range table.Relations {
		s.initWithClients(relation, clients)
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

func (m *TableClientMetrics) OtelStartTime(ctx context.Context, start time.Time) {
	if m.otelMeters == nil {
		return
	}

	// If we have already started, don't start again. This can happen for relational tables that are resolved multiple times (per parent resource)
	m.otelMeters.startedLock.Lock()
	defer m.otelMeters.startedLock.Unlock()
	if m.otelMeters.started {
		return
	}
	m.otelMeters.started = true
	m.otelMeters.startTime.Add(ctx, start.UnixNano(), metric.WithAttributes(m.otelMeters.attributes...))
}

func (m *TableClientMetrics) OtelEndTime(ctx context.Context, end time.Time) {
	if m.otelMeters == nil {
		return
	}

	m.otelMeters.previousEndTimeLock.Lock()
	defer m.otelMeters.previousEndTimeLock.Unlock()
	val := end.UnixNano()
	// If we got another end time to report, use the latest value. This can happen for relational tables that are resolved multiple times (per parent resource)
	if m.otelMeters.previousEndTime != 0 {
		m.otelMeters.endTime.Add(ctx, val-m.otelMeters.previousEndTime, metric.WithAttributes(m.otelMeters.attributes...))
	} else {
		m.otelMeters.endTime.Add(ctx, val, metric.WithAttributes(m.otelMeters.attributes...))
	}
	m.otelMeters.previousEndTime = val
}
