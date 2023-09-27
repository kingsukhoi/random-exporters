package OpenMetrics

import (
	"fmt"
	"strings"
)

type MetricTypes int

const (
	Gauge MetricTypes = iota
	Counter
)

type Metric struct {
	Name   string
	Type   MetricTypes
	Labels map[string]string
	Value  string
	Help   string
}

/**
# TYPE acme_http_router_request_seconds summary
# UNIT acme_http_router_request_seconds seconds
# HELP acme_http_router_request_seconds Latency though all of ACME's HTTP request router.
acme_http_router_request_seconds_sum{path="/api/v1",method="GET"} 9036.32
*/

func (m MetricTypes) String() string {
	switch m {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return "unknown"
}

func (m Metric) String() string {
	labels := ""
	for k, v := range m.Labels {
		labels += fmt.Sprintf(`%s="%s",`, k, v)
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`# TYPE %s %s`+"\n", m.Name, m.Type.String()))
	sb.WriteString(fmt.Sprintf(`# HELP %s %s`+"\n", m.Name, m.Help))

	if labels != "" {
		sb.WriteString(fmt.Sprintf("%s{%s} %s", m.Name, labels, m.Value))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s", m.Name, m.Value))
	}

	return sb.String() + "\n"
}
