package metrics

import (
	"math/rand"
	"runtime"
)

type (
	Gauge   float64
	Counter int64
)

type Metrics struct{}

type Type map[string]interface{}

var callCounter Counter

func (met *Metrics) GetListAvailable() Type {
	return Type{
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

func (met *Metrics) CollectData() Type {
	var memoryStat runtime.MemStats

	runtime.ReadMemStats(&memoryStat)

	return Type{
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
		"PollCount":     getCallCounter(),
		"RandomValue":   Gauge(rand.Float64()), //nolint:gosec
	}
}

func getCallCounter() Counter {
	return func() Counter {
		callCounter++
		return callCounter
	}()
}
