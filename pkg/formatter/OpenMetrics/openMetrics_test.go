package OpenMetrics

import (
	"testing"
)

func TestMetricTypes_String(t *testing.T) {
	tests := []struct {
		name string
		m    MetricTypes
		want string
	}{
		{name: "Gauge", m: Gauge, want: "gauge"},
		{name: "Counter", m: Counter, want: "counter"},
		{name: "Unknown", m: -2, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_String(t *testing.T) {
	tests := []struct {
		name   string
		metric Metric
		want   string
	}{
		{name: "Test 1", metric: Metric{
			Name: "test_gauge",
			Type: Gauge,
			Labels: map[string]string{
				"path": "/lololo",
			},
			Value: "5.3",
			Help:  "This is a test",
		}, want: `# TYPE test_gauge gauge
# HELP test_gauge This is a test
test_gauge{path="/lololo"} 5.3
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{
				Name:   tt.metric.Name,
				Type:   tt.metric.Type,
				Labels: tt.metric.Labels,
				Value:  tt.metric.Value,
				Help:   tt.metric.Help,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
