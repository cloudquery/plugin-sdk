package scheduler

import (
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// Metrics is deprecated as we move toward open telemetry for tracing and metrics
type Metrics struct {
	TableClient map[string]map[string]*TableClientMetrics
}

type TableClientMetrics struct {
	Resources uint64
	Errors    uint64
	Panics    uint64
	Duration  atomic.Value
}

func (s *TableClientMetrics) Equal(other *TableClientMetrics) bool {
	return s.Resources == other.Resources &&
		s.Errors == other.Errors &&
		s.Panics == other.Panics &&
		s.Duration.Load().(time.Duration) == other.Duration.Load().(time.Duration)
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

func (s *Metrics) initWithClients(table *schema.Table, clients []schema.ClientMeta) {
	s.TableClient[table.Name] = make(map[string]*TableClientMetrics, len(clients))
	for _, client := range clients {
		s.TableClient[table.Name][client.ID()] = &TableClientMetrics{}
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
