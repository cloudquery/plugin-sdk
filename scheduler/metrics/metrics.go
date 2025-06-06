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

func NewMetrics(invocationID string) *Metrics {
	resources, err := otel.Meter(ResourceName).Int64Counter(resourcesMetricName,
		metric.WithDescription("Number of resources synced for a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	errors, err := otel.Meter(ResourceName).Int64Counter(errorsMetricName,
		metric.WithDescription("Number of errors encountered while syncing a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	panics, err := otel.Meter(ResourceName).Int64Counter(panicsMetricName,
		metric.WithDescription("Number of panics encountered while syncing a table"),
		metric.WithUnit("/{tot}"),
	)
	if err != nil {
		return nil
	}

	duration, err := otel.Meter(ResourceName).Int64Counter(durationMetricName,
		metric.WithDescription("Duration of syncing a table"),
		metric.WithUnit("ns"),
	)
	if err != nil {
		return nil
	}

	return &Metrics{
		invocationID: invocationID,

		resources: resources,
		errors:    errors,
		panics:    panics,
		duration:  duration,

		TableClient: make(map[string]map[string]*tableClientMetrics),
	}
}

type Metrics struct {
	invocationID string

	resources metric.Int64Counter
	errors    metric.Int64Counter
	panics    metric.Int64Counter
	duration  metric.Int64Counter

	TableClient map[string]map[string]*tableClientMetrics
}

type tableClientMetrics struct {
	resources uint64
	errors    uint64
	panics    uint64
	duration  atomic.Pointer[time.Duration]

	startTime   time.Time
	started     bool
	startedLock sync.Mutex
}

func durationPointerEqual(a, b *time.Duration) bool {
	if a == nil {
		return b == nil
	}
	return b != nil && *a == *b
}

func (m *tableClientMetrics) Equal(other *tableClientMetrics) bool {
	return m.resources == other.resources && m.errors == other.errors && m.panics == other.panics && durationPointerEqual(m.duration.Load(), other.duration.Load())
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

func (m *Metrics) NewSelector(clientID, tableName string) Selector {
	return Selector{
		Set: attribute.NewSet(
			attribute.Key("sync.invocation.id").String(m.invocationID),
			attribute.Key("sync.table.name").String(tableName),
		),
		clientID:  clientID,
		tableName: tableName,
	}
}

func (m *Metrics) InitWithClients(invocationID string, table *schema.Table, clients []schema.ClientMeta) {
	m.TableClient[table.Name] = make(map[string]*tableClientMetrics, len(clients))
	for _, client := range clients {
		m.TableClient[table.Name][client.ID()] = &tableClientMetrics{}
	}
	for _, relation := range table.Relations {
		m.InitWithClients(invocationID, relation, clients)
	}
}

func (m *Metrics) TotalErrors() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.errors
		}
	}
	return total
}

func (m *Metrics) TotalErrorsAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.errors)
		}
	}
	return total
}

func (m *Metrics) TotalPanics() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.panics
		}
	}
	return total
}

func (m *Metrics) TotalPanicsAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.panics)
		}
	}
	return total
}

func (m *Metrics) TotalResources() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.resources
		}
	}
	return total
}

func (m *Metrics) TotalResourcesAtomic() uint64 {
	var total uint64
	for _, clientMetrics := range m.TableClient {
		for _, metrics := range clientMetrics {
			total += atomic.LoadUint64(&metrics.resources)
		}
	}
	return total
}

func (m *Metrics) ResourcesAdd(ctx context.Context, count int64, selector Selector) {
	m.resources.Add(ctx, count, metric.WithAttributeSet(selector.Set))
	atomic.AddUint64(&m.TableClient[selector.tableName][selector.clientID].resources, uint64(count))
}

func (m *Metrics) ResourcesGet(selector Selector) uint64 {
	return atomic.LoadUint64(&m.TableClient[selector.tableName][selector.clientID].resources)
}

func (m *Metrics) ErrorsAdd(ctx context.Context, count int64, selector Selector) {
	m.errors.Add(ctx, count, metric.WithAttributeSet(selector.Set))
	atomic.AddUint64(&m.TableClient[selector.tableName][selector.clientID].errors, uint64(count))
}

func (m *Metrics) ErrorsGet(selector Selector) uint64 {
	return atomic.LoadUint64(&m.TableClient[selector.tableName][selector.clientID].errors)
}

func (m *Metrics) PanicsAdd(ctx context.Context, count int64, selector Selector) {
	m.panics.Add(ctx, count, metric.WithAttributeSet(selector.Set))
	atomic.AddUint64(&m.TableClient[selector.tableName][selector.clientID].panics, uint64(count))
}

func (m *Metrics) PanicsGet(selector Selector) uint64 {
	return atomic.LoadUint64(&m.TableClient[selector.tableName][selector.clientID].panics)
}

func (m *Metrics) StartTime(start time.Time, selector Selector) {
	tc := m.TableClient[selector.tableName][selector.clientID]

	// If we have already started, don't start again. This can happen for relational tables that are resolved multiple times (per parent resource)
	tc.startedLock.Lock()
	defer tc.startedLock.Unlock()
	if tc.started {
		return
	}

	tc.started = true
	tc.startTime = start
}

func (m *Metrics) EndTime(ctx context.Context, end time.Time, selector Selector) {
	tc := m.TableClient[selector.tableName][selector.clientID]
	duration := time.Duration(end.UnixNano() - tc.startTime.UnixNano())
	tc.duration.Store(&duration)
	m.duration.Add(ctx, duration.Nanoseconds(), metric.WithAttributeSet(selector.Set))
}

func (m *Metrics) DurationGet(selector Selector) *time.Duration {
	tc := m.TableClient[selector.tableName][selector.clientID]
	if tc == nil {
		return nil
	}
	return tc.duration.Load()
}
