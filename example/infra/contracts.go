package infra

import "context"

// ApplicationMetric ...
type ApplicationMetric interface {
	Name() ApplicationMetricName
	WithLabelValues(...string) ApplicationMetricLabeled
}

// ApplicationMetricLabeled ...
type ApplicationMetricLabeled interface {
	Measure(context.Context, interface{}) *Error
}

// ApplicationMonitor ...
type ApplicationMonitor interface {
	Register(context.Context, ApplicationMetric) *Error
	GetMetric(name ApplicationMetricName) ApplicationMetric
	Expose(context.Context) *Error
}
