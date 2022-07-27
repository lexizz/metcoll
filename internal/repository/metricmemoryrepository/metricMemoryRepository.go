package metricmemoryrepository

import (
	"errors"
	"strings"

	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/lexizz/metcoll/internal/repository/interfaces/metricrepository"
)

type storage metrics.Collection

type metricMemory struct {
	data storage
}

var _ metricrepository.Interface = &metricMemory{}

func (m *metricMemory) GetAll() (metrics.Collection, error) {
	if len(m.data) < 1 {
		return nil, errors.New("data not found")
	}

	return metrics.Collection(m.data), nil
}

func New() *metricMemory {
	return &metricMemory{data: make(storage, 50)}
}

func (m *metricMemory) IncreaseValue(fieldName string, value metrics.Counter) (bool, error) {
	fieldName = strings.ToLower(fieldName)
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
	nameMetric = strings.ToLower(nameMetric)
	m.data[nameMetric] = value
}

func (m *metricMemory) GetValue(nameMetric string) (interface{}, error) {
	nameMetric = strings.ToLower(nameMetric)
	value, ok := m.data[nameMetric]

	if !ok {
		return nil, errors.New("not found metric")
	}

	return value, nil
}
