package metricMemoryRepository

import (
	"errors"
	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/lexizz/metcoll/internal/repository/interfaces/metricRepository"
)

type storage map[string]interface{}

type metricMemory struct {
	data storage
}

var _ metricRepository.Interface = &metricMemory{}

func New() *metricMemory {
	return &metricMemory{data: make(storage, 50)}
}

func (m *metricMemory) IncreaseValue(fieldName string, value metrics.Counter) (bool, error) {
	recordedValue, ok := m.data[fieldName]

	if !ok {
		m.data[fieldName] = value
		return true, nil
	}

	var oldValue metrics.Counter

	oldValue, ok = recordedValue.(metrics.Counter)
	if !ok {
		return false, errors.New("type change error")
	}

	m.data[fieldName] = oldValue + value

	return true, nil
}

func (m *metricMemory) InsertValue(nameMetric string, value interface{}) {
	m.data[nameMetric] = value
}

func (m *metricMemory) GetValue(nameMetric string) (interface{}, error) {
	value, ok := m.data[nameMetric]

	if !ok {
		return nil, errors.New("not found metric")
	}

	return value, nil
}
