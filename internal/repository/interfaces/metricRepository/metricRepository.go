package metricRepository

import "github.com/lexizz/metcoll/internal/metrics"

type Interface interface {
	GetValue(nameMetric string) (interface{}, error)
	InsertValue(nameMetric string, value interface{})
	IncreaseValue(fieldName string, value metrics.Counter) (bool, error)
}
