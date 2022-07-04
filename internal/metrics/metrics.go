package metrics

import (
	"math/rand"
	"runtime"
)

type (
	gauge   float64
	counter int64
)

type Metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

var (
	memoryStat  runtime.MemStats
	callCounter counter = 0
)

func Collect() Metrics {
	runtime.ReadMemStats(&memoryStat)

	return Metrics{
		Alloc:         gauge(memoryStat.Alloc),
		BuckHashSys:   gauge(memoryStat.BuckHashSys),
		Frees:         gauge(memoryStat.Frees),
		GCCPUFraction: gauge(memoryStat.GCCPUFraction),
		GCSys:         gauge(memoryStat.GCSys),
		HeapAlloc:     gauge(memoryStat.HeapAlloc),
		HeapIdle:      gauge(memoryStat.HeapIdle),
		HeapInuse:     gauge(memoryStat.HeapInuse),
		HeapObjects:   gauge(memoryStat.HeapObjects),
		HeapReleased:  gauge(memoryStat.HeapReleased),
		HeapSys:       gauge(memoryStat.HeapSys),
		LastGC:        gauge(memoryStat.LastGC),
		Lookups:       gauge(memoryStat.Lookups),
		MCacheInuse:   gauge(memoryStat.MCacheInuse),
		MCacheSys:     gauge(memoryStat.MCacheSys),
		MSpanInuse:    gauge(memoryStat.MSpanInuse),
		MSpanSys:      gauge(memoryStat.MSpanSys),
		Mallocs:       gauge(memoryStat.Mallocs),
		NextGC:        gauge(memoryStat.NextGC),
		NumForcedGC:   gauge(memoryStat.NumForcedGC),
		NumGC:         gauge(memoryStat.NumGC),
		OtherSys:      gauge(memoryStat.OtherSys),
		PauseTotalNs:  gauge(memoryStat.PauseTotalNs),
		StackInuse:    gauge(memoryStat.StackInuse),
		StackSys:      gauge(memoryStat.StackSys),
		Sys:           gauge(memoryStat.Sys),
		TotalAlloc:    gauge(memoryStat.TotalAlloc),
		PollCount:     getCallCounter(),
		RandomValue:   gauge(rand.Float64()),
	}
}

func getCallCounter() counter {
	return func() counter {
		callCounter++
		return callCounter
	}()
}
