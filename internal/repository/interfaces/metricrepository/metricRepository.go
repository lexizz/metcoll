package metricrepository

import "github.com/lexizz/metcoll/internal/metrics"

type Interface interface {
	GetAll() (metrics.Type, error)
	GetValue(nameMetric string) (interface{}, error)
	InsertValue(nameMetric string, value interface{})
	IncreaseValue(fieldName string, value metrics.Counter) (bool, error)
}
