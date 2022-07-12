package metrics

import (
	"math/rand"
	"runtime"
)

type (
	Gauge   float64
	Counter int64
)

type Metrics map[string]interface{}

var callCounter Counter = 0

func CollectData() Metrics {
	var memoryStat runtime.MemStats

	runtime.ReadMemStats(&memoryStat)

	return Metrics{
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
		"RandomValue":   Gauge(rand.Float64()),
	}
}

func getCallCounter() Counter {
	return func() Counter {
		callCounter++
		return callCounter
	}()
}
