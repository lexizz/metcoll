package metricDbRepository

import (
	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/lexizz/metcoll/internal/repository/interfaces/metricRepository"
)

type storage []string

type metricMemory struct {
	data storage
}

var _ metricRepository.Interface = &metricMemory{}

func New() *metricMemory {
	return &metricMemory{}
}

func (m *metricMemory) IncreaseValue(fieldName string, value metrics.Counter) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *metricMemory) GetValue(nameMetric string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (m *metricMemory) InsertValue(nameMetric string, value interface{}) {
	//TODO implement me
	panic("implement me")
}
