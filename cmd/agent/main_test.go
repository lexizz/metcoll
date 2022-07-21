package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lexizz/metcoll/internal/metrics"
)

func TestConvertToString(t *testing.T) {
	met := metrics.New()
	exp := exporter{
		httpClient:  &http.Client{},
		metrics:     met,
		metricsData: met.CollectData(),
	}

	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "test convert from metrics.Gauge",
			value: metrics.Gauge(5.3),
			want:  "5.30",
		},
		{
			name:  "test convert from metrics.Gauge crop to 2 numbers",
			value: metrics.Gauge(5.312121212121212),
			want:  "5.31",
		},
		{
			name:  "test convert from metrics.Gauge with zero",
			value: metrics.Gauge(0),
			want:  "0.00",
		},
		{
			name:  "test convert from number with zero",
			value: 0,
			want:  "0",
		},
		{
			name:  "test convert from metrics.Counter",
			value: metrics.Counter(5),
			want:  "5",
		},
		{
			name:  "test convert from metrics.Counter with zero",
			value: metrics.Counter(0),
			want:  "0",
		},
		{
			name:  "test convert from string",
			value: "5.3",
			want:  "5.3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := exp.convertValueToString(test.value, strings.ToLower(reflect.TypeOf(test.value).Name()))

			assert.Equal(t, test.want, result)
		})
	}
}

func TestSendRequest(t *testing.T) {
	met := metrics.New()
	exp := exporter{
		httpClient:  &http.Client{},
		metrics:     met,
		metricsData: met.CollectData(),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "text/plain; charset=UTF-8", r.Header.Get("Content-Type"))

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	exp.sendRequest(ts.URL)
}

func TestGetListUrls(t *testing.T) {
	met := metrics.New()
	exp := exporter{
		httpClient:  &http.Client{},
		metrics:     met,
		metricsData: met.CollectData(),
	}

	countMetrics := len(exp.metricsData)
	countUrls := len(exp.getListUrls())

	assert.Equal(t, countMetrics, countUrls)
}
