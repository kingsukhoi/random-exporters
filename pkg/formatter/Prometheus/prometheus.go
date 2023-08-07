package Prometheus

import (
	"errors"
	"fmt"
)

type Gauge struct {
	Name   string
	Labels map[string]string
	Value  string
	Help   string
}

/*
sample

# HELP http_request_duration_seconds A histogram of the request duration.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.05"} 24054

*/

func (g Gauge) PrintPrometheusGauge() (rtnMe string, err error) {
	if g.Name == "" || g.Value == "" {
		return "", errors.New("name and value are required")
	}
	if g.Help != "" {
		rtnMe += fmt.Sprintf("# HELP %s %s\n", g.Name, g.Help)
	}
	rtnMe += fmt.Sprintf("# TYPE %s GAUGE\n", g.Name)

	labels := ""
	for k, v := range g.Labels {
		labels += fmt.Sprintf(`%s="%s"`, k, v)
	}

	if labels != "" {
		rtnMe += fmt.Sprintf("%s{%s} %s", g.Name, labels, g.Value)
	} else {
		rtnMe += fmt.Sprintf("%s %s", g.Name, g.Value)
	}

	return
}
