package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	m := NewMetrics()

	m.measurements["test_table_1"] = tableMeasurements{
		clients: map[string]*measurement{
			"test_client_1": {duration: &durationMeasurement{}},
			"test_client_2": {duration: &durationMeasurement{}},
		},
		duration: &durationMeasurement{},
	}
	m.measurements["test_table_2"] = tableMeasurements{
		clients: map[string]*measurement{
			"test_client_1": {duration: &durationMeasurement{}},
		},
		duration: &durationMeasurement{},
	}

	require.Equal(t, m.TotalResources(), uint64(0))
	require.Equal(t, m.TotalErrors(), uint64(0))
	require.Equal(t, m.TotalPanics(), uint64(0))

	s1 := m.NewSelector("test_client_1", "test_table_1")

	// test single table, single client
	m.StartTime(time.Now(), s1)
	m.AddResources(t.Context(), 1, s1)
	require.Equal(t, m.TotalResources(), uint64(1))
	require.Equal(t, m.GetResources(s1), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(0))
	require.Equal(t, m.TotalPanics(), uint64(0))

	m.AddErrors(t.Context(), 1, s1)
	require.Equal(t, m.TotalResources(), uint64(1))
	require.Equal(t, m.GetResources(s1), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(1))
	require.Equal(t, m.GetErrors(s1), uint64(1))
	require.Equal(t, m.TotalPanics(), uint64(0))

	m.AddPanics(t.Context(), 1, s1)
	require.Equal(t, m.TotalResources(), uint64(1))
	require.Equal(t, m.GetResources(s1), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(1))
	require.Equal(t, m.GetErrors(s1), uint64(1))
	require.Equal(t, m.TotalPanics(), uint64(1))
	require.Equal(t, m.GetPanics(s1), uint64(1))

	time.Sleep(1 * time.Millisecond)
	m.EndTime(t.Context(), time.Now(), s1)

	// test single table, multiple clients
	s2 := m.NewSelector("test_client_2", "test_table_1")

	m.StartTime(time.Now(), s2)
	m.AddResources(t.Context(), 1, s2)
	require.Equal(t, m.TotalResources(), uint64(2))
	require.Equal(t, m.GetResources(s2), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(1))
	require.Equal(t, m.GetErrors(s2), uint64(0))
	require.Equal(t, m.TotalPanics(), uint64(1))
	require.Equal(t, m.GetPanics(s2), uint64(0))

	m.AddErrors(t.Context(), 1, s2)
	require.Equal(t, m.TotalResources(), uint64(2))
	require.Equal(t, m.GetResources(s2), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(2))
	require.Equal(t, m.GetErrors(s2), uint64(1))
	require.Equal(t, m.TotalPanics(), uint64(1))
	require.Equal(t, m.GetPanics(s2), uint64(0))

	m.AddPanics(t.Context(), 1, s2)
	require.Equal(t, m.TotalResources(), uint64(2))
	require.Equal(t, m.GetResources(s2), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(2))
	require.Equal(t, m.GetErrors(s2), uint64(1))
	require.Equal(t, m.TotalPanics(), uint64(2))
	require.Equal(t, m.GetPanics(s2), uint64(1))

	time.Sleep(1 * time.Millisecond)
	m.EndTime(t.Context(), time.Now(), s2)

	// test multiple tables, multiple clients
	s3 := m.NewSelector("test_client_1", "test_table_2")

	m.StartTime(time.Now(), s3)
	m.AddResources(t.Context(), 1, s3)
	require.Equal(t, m.TotalResources(), uint64(3))
	require.Equal(t, m.GetResources(s3), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(2))
	require.Equal(t, m.GetErrors(s3), uint64(0))
	require.Equal(t, m.TotalPanics(), uint64(2))
	require.Equal(t, m.GetPanics(s3), uint64(0))

	m.AddErrors(t.Context(), 1, s3)
	require.Equal(t, m.TotalResources(), uint64(3))
	require.Equal(t, m.GetResources(s3), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(3))
	require.Equal(t, m.GetErrors(s3), uint64(1))
	require.Equal(t, m.TotalPanics(), uint64(2))
	require.Equal(t, m.GetPanics(s3), uint64(0))

	m.AddPanics(t.Context(), 1, s3)
	require.Equal(t, m.TotalResources(), uint64(3))
	require.Equal(t, m.GetResources(s3), uint64(1))
	require.Equal(t, m.TotalErrors(), uint64(3))
	require.Equal(t, m.GetErrors(s3), uint64(1))
	require.Equal(t, m.TotalPanics(), uint64(3))
	require.Equal(t, m.GetPanics(s3), uint64(1))

	time.Sleep(1 * time.Millisecond)
	m.EndTime(t.Context(), time.Now(), s3)

	require.Greater(t, m.GetDuration(s1), 0*time.Nanosecond)
	require.Greater(t, m.GetDuration(s2), 0*time.Nanosecond)

	// This should work because the 2 metrics are built sequentially; in practice though, this is probably not the case.
	require.GreaterOrEqual(t, m.TableDuration(s1.tableName), m.GetDuration(s1)+m.GetDuration(s2))
}
