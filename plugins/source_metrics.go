package plugins

import (
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
)

type SourceMetrics struct {
	TableClient map[string]map[string]*TableClientMetrics
}

type TableClientMetrics struct {
	Resources uint64
	Errors    uint64
	Panics    uint64
	StartTime time.Time
	EndTime   time.Time
}

func (s *TableClientMetrics) Equal(other *TableClientMetrics) bool {
	return s.Resources == other.Resources && s.Errors == other.Errors && s.Panics == other.Panics
}

// Equal compares to stats. Mostly useful in testing
func (s *SourceMetrics) Equal(other *SourceMetrics) bool {
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

func (s *SourceMetrics) initWithTables(tables schema.Tables) {
	s.TableClient = make(map[string]map[string]*TableClientMetrics)
	for _, table := range tables {
		s.TableClient[table.Name] = make(map[string]*TableClientMetrics)
		s.initWithTables(table.Relations)
	}
}

func (s *SourceMetrics) initWithClients(table *schema.Table, clients []schema.ClientMeta) {
	s.TableClient[table.Name] = make(map[string]*TableClientMetrics)
	for _, client := range clients {
		s.TableClient[table.Name][client.ID()] = &TableClientMetrics{}
		for _, relation := range table.Relations {
			s.initWithClients(relation, clients)
		}
	}
}

func (s *SourceMetrics) TotalErrors() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Errors
		}
	}
	return total
}

func (s *SourceMetrics) TotalPanics() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Panics
		}
	}
	return total
}

func (s *SourceMetrics) TotalResources() uint64 {
	var total uint64
	for _, clientMetrics := range s.TableClient {
		for _, metrics := range clientMetrics {
			total += metrics.Resources
		}
	}
	return total
}
