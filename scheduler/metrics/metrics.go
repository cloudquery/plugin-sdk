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
	ResourceName = "io.cloudquery"

	resourcesMetricName = "sync.table.resources"
	errorsMetricName    = "sync.table.errors"
	panicsMetricName    = "sync.table.panics"
	durationMetricName  = "sync.table.duration"
)

var (
	resources metric.Int64Counter
	errors    metric.Int64Counter
	panics    metric.Int64Counter
	duration  metric.Int64Counter
	once      sync.Once
)

func NewMetrics() *Metrics {
	once.Do(func() {
		resources, _ = otel.Meter(ResourceName).Int64Counter(resourcesMetricName,
			metric.WithDescription("Number of resources synced for a table"),
			metric.WithUnit("/{tot}"),
		)

		errors, _ = otel.Meter(ResourceName).Int64Counter(errorsMetricName,
			metric.WithDescription("Number of errors encountered while syncing a table"),
			metric.WithUnit("/{tot}"),
		)

		panics, _ = otel.Meter(ResourceName).Int64Counter(panicsMetricName,
			metric.WithDescription("Number of panics encountered while syncing a table"),
			metric.WithUnit("/{tot}"),
		)

		duration, _ = otel.Meter(ResourceName).Int64Counter(durationMetricName,
			metric.WithDescription("Duration of syncing a table"),
			metric.WithUnit("ms"),
		)
	})

	return &Metrics{
		resources: resources,
		errors:    errors,
		panics:    panics,
		duration:  duration,

		measurements: make(map[string]tableMeasurements),
	}
}

type Metrics struct {
	resources metric.Int64Counter
	errors    metric.Int64Counter
	panics    metric.Int64Counter
	duration  metric.Int64Counter

	measurements map[string]tableMeasurements
}

type tableMeasurements struct {
	duration *durationMeasurement
	clients  map[string]*measurement
}

type measurement struct {
	resources uint64
	errors    uint64
	panics    uint64
	duration  *durationMeasurement
}

func (*Metrics) NewSelector(clientID, tableName string) Selector {
	return Selector{
		Set: attribute.NewSet(
			attribute.Key("sync.table.name").String(tableName),
			attribute.Key("sync.client.id").String(""),
		),
		clientID:  clientID,
		tableName: tableName,
	}
}

func (m *Metrics) InitWithClients(table *schema.Table, clients []schema.ClientMeta) {
	m.measurements[table.Name] = tableMeasurements{clients: make(map[string]*measurement), duration: &durationMeasurement{}}
	for _, client := range clients {
		m.measurements[table.Name].clients[client.ID()] = &measurement{duration: &durationMeasurement{}}
	}
	for _, relation := range table.Relations {
		m.InitWithClients(relation, clients)
	}
}

func (m *Metrics) TotalErrors() uint64 {
	var total uint64
	for _, clientMetrics := range m.measurements {
		for _, metrics := range clientMetrics.clients {
			total += atomic.LoadUint64(&metrics.errors)
		}
	}
	return total
}

// Deprecated: Use TotalErrors instead, it provides the same functionality but is more consistent with the naming of other metrics methods.
func (m *Metrics) TotalErrorsAtomic() uint64 {
	return m.TotalErrors()
}

func (m *Metrics) TotalPanics() uint64 {
	var total uint64
	for _, clientMetrics := range m.measurements {
		for _, metrics := range clientMetrics.clients {
			total += atomic.LoadUint64(&metrics.panics)
		}
	}
	return total
}

// Deprecated: Use TotalPanics instead, it provides the same functionality but is more consistent with the naming of other metrics methods.
func (m *Metrics) TotalPanicsAtomic() uint64 {
	return m.TotalPanics()
}

func (m *Metrics) TotalResources() uint64 {
	var total uint64
	for _, clientMetrics := range m.measurements {
		for _, metrics := range clientMetrics.clients {
			total += atomic.LoadUint64(&metrics.resources)
		}
	}
	return total
}

// Deprecated: Use TotalResources instead, it provides the same functionality but is more consistent with the naming of other metrics methods.
func (m *Metrics) TotalResourcesAtomic() uint64 {
	return m.TotalResources()
}

func (m *Metrics) TableDuration(tableName string) time.Duration {
	tc := m.measurements[tableName]
	return tc.duration.duration
}

func (m *Metrics) AddResources(ctx context.Context, count int64, selector Selector) {
	m.resources.Add(ctx, count, metric.WithAttributeSet(selector.Set))
	atomic.AddUint64(&m.measurements[selector.tableName].clients[selector.clientID].resources, uint64(count))
}

func (m *Metrics) GetResources(selector Selector) uint64 {
	return atomic.LoadUint64(&m.measurements[selector.tableName].clients[selector.clientID].resources)
}

func (m *Metrics) AddErrors(ctx context.Context, count int64, selector Selector) {
	m.errors.Add(ctx, count, metric.WithAttributeSet(selector.Set))
	atomic.AddUint64(&m.measurements[selector.tableName].clients[selector.clientID].errors, uint64(count))
}

func (m *Metrics) GetErrors(selector Selector) uint64 {
	return atomic.LoadUint64(&m.measurements[selector.tableName].clients[selector.clientID].errors)
}

func (m *Metrics) AddPanics(ctx context.Context, count int64, selector Selector) {
	m.panics.Add(ctx, count, metric.WithAttributeSet(selector.Set))
	atomic.AddUint64(&m.measurements[selector.tableName].clients[selector.clientID].panics, uint64(count))
}

func (m *Metrics) GetPanics(selector Selector) uint64 {
	return atomic.LoadUint64(&m.measurements[selector.tableName].clients[selector.clientID].panics)
}

func (m *Metrics) StartTime(start time.Time, selector Selector) {
	t := m.measurements[selector.tableName]
	tc := t.clients[selector.clientID]

	tc.duration.Start(start)
	t.duration.Start(start)
}

func (m *Metrics) EndTime(ctx context.Context, end time.Time, selector Selector) {
	t := m.measurements[selector.tableName]
	tc := t.clients[selector.clientID]

	_ = tc.duration.End(end)
	delta := t.duration.End(end)

	// only compute and add the total duration for per-table measurements (and not per-client)
	m.duration.Add(ctx, delta.Milliseconds(), metric.WithAttributeSet(selector.Set))
}

func (m *Metrics) GetDuration(selector Selector) time.Duration {
	tc := m.measurements[selector.tableName].clients[selector.clientID]
	return tc.duration.duration
}
