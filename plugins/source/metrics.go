package source

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v2/schema"
	"golang.org/x/exp/slices"
)

type Metrics struct {
	TableClient map[string]map[string]*TableClientMetrics
}

type TableClientMetrics struct {
	// These should only be accessed with 'Atomic*' methods.
	Resources uint64
	Errors    uint64
	Panics    uint64

	// These accesses must be protected by the mutex.
	startTime time.Time
	endTime   time.Time
	mutex     sync.Mutex
}

func (s *TableClientMetrics) Equal(other *TableClientMetrics) bool {
	return s.Resources == other.Resources && s.Errors == other.Errors && s.Panics == other.Panics
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

func (s *Metrics) MarkStart(table *schema.Table, clientID string) {
	now := time.Now()

	s.TableClient[table.Name][clientID].mutex.Lock()
	defer s.TableClient[table.Name][clientID].mutex.Unlock()
	s.TableClient[table.Name][clientID].startTime = now
}

// if the table is a top-level table, we need to mark all of its descendents as 'done' as well.
// This is because, when a top-level table is empty (no resources), its descendants are never actually
// synced.
func (s *Metrics) MarkEnd(table *schema.Table, clientID string) {
	now := time.Now()

	if table.Parent == nil {
		s.markEndRecursive(table, clientID, now)
		return
	}

	s.TableClient[table.Name][clientID].mutex.Lock()
	defer s.TableClient[table.Name][clientID].mutex.Unlock()
	s.TableClient[table.Name][clientID].endTime = now
}

func (s *Metrics) markEndRecursive(table *schema.Table, clientID string, now time.Time) {
	// We don't use defer with Unlock(), because we want to unlock the mutex as soon as possible.
	s.TableClient[table.Name][clientID].mutex.Lock()
	s.TableClient[table.Name][clientID].endTime = now
	s.TableClient[table.Name][clientID].mutex.Unlock()

	for _, relation := range table.Relations {
		s.markEndRecursive(relation, clientID, now)
	}
}

func (s *Metrics) InProgressTables() []string {
	var inProgressTables []string

	for table, tableMetrics := range s.TableClient {
		for _, clientMetrics := range tableMetrics {
			clientMetrics.mutex.Lock()
			endTime := clientMetrics.endTime
			clientMetrics.mutex.Unlock()
			if endTime.IsZero() {
				inProgressTables = append(inProgressTables, table)
				break
			}
		}
	}

	slices.Sort(inProgressTables)

	return inProgressTables
}
