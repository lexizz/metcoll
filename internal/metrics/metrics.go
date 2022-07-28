package metrics

import (
	"math/rand"
	"runtime"
	"sync"
)

type (
	Gauge   float64
	Counter int64
)

type Metrics struct {
	ID      string   `json:"id"`              // имя метрики
	MType   string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta   *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value   *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	counter Counter
	mutex   *sync.RWMutex
}

type Collection map[string]interface{}

func New() *Metrics {
	var m sync.RWMutex

	met := Metrics{}
	met.mutex = &m

	return &met
}

func (met *Metrics) GetListAvailable() Collection {
	return Collection{
		"Alloc":         nil,
		"BuckHashSys":   nil,
		"Frees":         nil,
		"GCCPUFraction": nil,
		"GCSys":         nil,
		"HeapAlloc":     nil,
		"HeapIdle":      nil,
		"HeapInuse":     nil,
		"HeapObjects":   nil,
		"HeapReleased":  nil,
		"HeapSys":       nil,
		"LastGC":        nil,
		"Lookups":       nil,
		"MCacheInuse":   nil,
		"MCacheSys":     nil,
		"MSpanInuse":    nil,
		"MSpanSys":      nil,
		"Mallocs":       nil,
		"NextGC":        nil,
		"NumForcedGC":   nil,
		"NumGC":         nil,
		"OtherSys":      nil,
		"PauseTotalNs":  nil,
		"StackInuse":    nil,
		"StackSys":      nil,
		"Sys":           nil,
		"TotalAlloc":    nil,
		"PollCount":     nil,
		"RandomValue":   nil,
	}
}

func (met *Metrics) CollectData() Collection {
	var memoryStat runtime.MemStats
	met.setCounter()

	runtime.ReadMemStats(&memoryStat)

	met.mutex.Lock()
	data := Collection{
		"Alloc":         Gauge(memoryStat.Alloc),
		"BuckHashSys":   Gauge(memoryStat.BuckHashSys),
		"Frees":         Gauge(memoryStat.Frees),
		"GCCPUFraction": Gauge(memoryStat.GCCPUFraction),
		"GCSys":         Gauge(memoryStat.GCSys),
		"HeapAlloc":     Gauge(memoryStat.HeapAlloc),
		"HeapIdle":      Gauge(memoryStat.HeapIdle),
		"HeapInuse":     Gauge(memoryStat.HeapInuse),
		"HeapObjects":   Gauge(memoryStat.HeapObjects),
		"HeapReleased":  Gauge(memoryStat.HeapReleased),
		"HeapSys":       Gauge(memoryStat.HeapSys),
		"LastGC":        Gauge(memoryStat.LastGC),
		"Lookups":       Gauge(memoryStat.Lookups),
		"MCacheInuse":   Gauge(memoryStat.MCacheInuse),
		"MCacheSys":     Gauge(memoryStat.MCacheSys),
		"MSpanInuse":    Gauge(memoryStat.MSpanInuse),
		"MSpanSys":      Gauge(memoryStat.MSpanSys),
		"Mallocs":       Gauge(memoryStat.Mallocs),
		"NextGC":        Gauge(memoryStat.NextGC),
		"NumForcedGC":   Gauge(memoryStat.NumForcedGC),
		"NumGC":         Gauge(memoryStat.NumGC),
		"OtherSys":      Gauge(memoryStat.OtherSys),
		"PauseTotalNs":  Gauge(memoryStat.PauseTotalNs),
		"StackInuse":    Gauge(memoryStat.StackInuse),
		"StackSys":      Gauge(memoryStat.StackSys),
		"Sys":           Gauge(memoryStat.Sys),
		"TotalAlloc":    Gauge(memoryStat.TotalAlloc),
		"PollCount":     met.getCallCounter(),
		"RandomValue":   Gauge(rand.Float64()), //nolint:gosec
	}
	met.mutex.Unlock()

	return data
}

func (met *Metrics) getCallCounter() Counter {
	return met.counter
}

func (met *Metrics) setCounter() {
	met.mutex.Lock()
	met.counter++
	met.mutex.Unlock()
}
