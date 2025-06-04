package metrics

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

const (
	OtelName = "io.cloudquery"
)

func NewMetrics() *Metrics {
	resources, err := otel.Meter(OtelName).Int64Counter("sync.table.resources",
		metric.WithDescription("Number of resources synced for a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	errors, err := otel.Meter(OtelName).Int64Counter("sync.table.errors",
		metric.WithDescription("Number of errors encountered while syncing a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	panics, err := otel.Meter(OtelName).Int64Counter("sync.table.panics",
		metric.WithDescription("Number of panics encountered while syncing a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	startTime, err := otel.Meter(OtelName).Int64Counter("sync.table.start_time",
		metric.WithDescription("Start time of syncing a table"),
		metric.WithUnit("ns"),
	)
	if err != nil {
		return nil
	}

	endTime, err := otel.Meter(OtelName).Int64Counter("sync.table.end_time",
		metric.WithDescription("End time of syncing a table"),
		metric.WithUnit("ns"),
	)

	if err != nil {
		return nil
	}

	return &Metrics{
		TableClient: make(map[string]map[string]*TableClientMetrics),
		resources:   resources,
		errors:      errors,
		panics:      panics,
		startTime:   startTime,
		endTime:     endTime,
	}
}

// Metrics is deprecated as we move toward open telemetry for tracing and metrics
type Metrics struct {
	resources metric.Int64Counter
	errors    metric.Int64Counter
	panics    metric.Int64Counter

	startTime   metric.Int64Counter
	started     bool
	startedLock sync.Mutex

	endTime             metric.Int64Counter
	previousEndTime     int64
	previousEndTimeLock sync.Mutex

	TableClient map[string]map[string]*TableClientMetrics
}

type OtelMeters struct {
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

func (m *TableClientMetrics) Equal(other *TableClientMetrics) bool {
	return m.Resources == other.Resources && m.Errors == other.Errors && m.Panics == other.Panics && durationPointerEqual(m.Duration.Load(), other.Duration.Load())
}

// Equal compares to stats. Mostly useful in testing
func (m *Metrics) Equal(other *Metrics) bool {
	for table, clientStats := range m.TableClient {
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
			if _, ok := m.TableClient[table]; !ok {
				return false
			}
			if _, ok := m.TableClient[table][client]; !ok {
				return false
			}
			if !stats.Equal(m.TableClient[table][client]) {
				return false
			}
		}
	}
	return true
}

func GetOtelAttributeSet(tableName string, clientID string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Key("sync.client.id").String(clientID),
		attribute.Key("sync.table.name").String(tableName),
	}
}

func (m *Metrics) InitWithClients(table *schema.Table, clients []schema.ClientMeta) {
	m.TableClient[table.Name] = make(map[string]*TableClientMetrics, len(clients))
	for _, client := range clients {
		tableName := table.Name
		clientID := client.ID()
		m.TableClient[tableName][clientID] = &TableClientMetrics{
			otelMeters: &OtelMeters{attributes: GetOtelAttributeSet(tableName, clientID)},
		}
	}
	for _, relation := range table.Relations {
		m.InitWithClients(relation, clients)
	}
}

func (m *Metrics) TotalErrors() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Errors
		}
	}
	return total
}

func (m *Metrics) TotalErrorsAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.Errors)
		}
	}
	return total
}

func (m *Metrics) TotalPanics() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Panics
		}
	}
	return total
}

func (m *Metrics) TotalPanicsAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.Panics)
		}
	}
	return total
}

func (m *Metrics) TotalResources() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Resources
		}
	}
	return total
}

func (m *Metrics) TotalResourcesAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.Resources)
		}
	}
	return total
}

func (m *Metrics) OtelResourcesAdd(ctx context.Context, count int64, tc *TableClientMetrics) {
	m.resources.Add(ctx, count, metric.WithAttributes(tc.otelMeters.attributes...))
	atomic.AddUint64(&tc.Resources, uint64(count))
}

func (m *Metrics) OtelErrorsAdd(ctx context.Context, count int64, tc *TableClientMetrics) {
	m.errors.Add(ctx, count, metric.WithAttributes(tc.otelMeters.attributes...))
	atomic.AddUint64(&tc.Errors, uint64(count))
}

func (m *Metrics) OtelPanicsAdd(ctx context.Context, count int64, tc *TableClientMetrics) {
	m.panics.Add(ctx, count, metric.WithAttributes(tc.otelMeters.attributes...))
	atomic.AddUint64(&tc.Panics, uint64(count))
}

func (m *Metrics) OtelStartTime(ctx context.Context, start time.Time, tc *TableClientMetrics) {
	if m.startTime == nil {
		return
	}

	// If we have already started, don't start again. This can happen for relational tables that are resolved multiple times (per parent resource)
	m.startedLock.Lock()
	defer m.startedLock.Unlock()
	if m.started {
		return
	}

	m.started = true
	m.startTime.Add(ctx, start.UnixNano(), metric.WithAttributes(tc.otelMeters.attributes...))
}

func (m *Metrics) OtelEndTime(ctx context.Context, end time.Time, tc *TableClientMetrics) {
	if m.endTime == nil {
		return
	}

	m.previousEndTimeLock.Lock()
	defer m.previousEndTimeLock.Unlock()
	val := end.UnixNano()

	// If we got another end time to report, use the latest value. This can happen for relational tables that are resolved multiple times (per parent resource)
	if m.previousEndTime != 0 {
		m.endTime.Add(ctx, val-m.previousEndTime, metric.WithAttributes(tc.otelMeters.attributes...))
	} else {
		m.endTime.Add(ctx, val, metric.WithAttributes(tc.otelMeters.attributes...))
	}
	m.previousEndTime = val
}
