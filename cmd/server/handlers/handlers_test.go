package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/lexizz/metcoll/internal/repository/metricmemoryrepository"
	"github.com/lexizz/metcoll/internal/server"
)

func TestUpdateMetric(t *testing.T) {
	metricRepository := metricmemoryrepository.New()

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "positive test #1 with type gauge",
			method: http.MethodPost,
			url:    "/update/gauge/Alloc/5",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "positive test #2 with type counter",
			method: http.MethodPost,
			url:    "/update/counter/PollCounter/5",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "positive test #3 - wrong name of metric and right value",
			method: http.MethodPost,
			url:    "/update/counter/none/5",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #1 - wrong type name",
			method: http.MethodPost,
			url:    "/update/unknown/Alloc/5",
			want: want{
				code:        http.StatusNotImplemented,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2 - wrong value of metric",
			method: http.MethodPost,
			url:    "/update/counter/Alloc/sdf",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #3 - wrong name of metric and wrong value",
			method: http.MethodPost,
			url:    "/update/counter/none/none",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #4 - without value of metric",
			method: http.MethodPost,
			url:    "/update/counter/Alloc/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #5 - without value of metric and name",
			method: http.MethodPost,
			url:    "/update/counter/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #6 - wrong url like general page",
			method: http.MethodPost,
			url:    "/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #7 - wrong http method with right url",
			method: http.MethodGet,
			url:    "/update/gauge/Alloc/5",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			writer := httptest.NewRecorder()

			routes := server.GetRoutes(metricRepository)
			routes.ServeHTTP(writer, request)

			res := writer.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("content-type"))
		})
	}
}

func TestShowValueMetric(t *testing.T) {
	metricRepository := metricmemoryrepository.New()

	firstTestMetricName := "Alloc"
	secondTestMetricName := "PollCount"

	metricRepository.InsertValue(firstTestMetricName, metrics.Gauge(5))
	metricRepository.InsertValue(secondTestMetricName, metrics.Counter(5))

	_, errValFirst := metricRepository.GetValue(firstTestMetricName)
	require.NoError(t, errValFirst)

	_, errValSecond := metricRepository.GetValue(secondTestMetricName)
	require.NoError(t, errValSecond)

	type want struct {
		code        int
		contentType string
		response    string
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodGet,
			url:    "/value/gauge/Alloc",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				response:    "5",
			},
		},
		{
			name:   "positive test #2",
			method: http.MethodGet,
			url:    "/value/counter/PollCount",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				response:    "5",
			},
		},
		{
			name:   "negative test #1 - wrong type of metric #1",
			method: http.MethodGet,
			url:    "/value/counter/Alloc",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				response:    "Wrong metric type\n",
			},
		},
		{
			name:   "negative test #2 - wrong type of metric #2",
			method: http.MethodGet,
			url:    "/value/none/Alloc",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				response:    "Wrong metric type\n",
			},
		},
		{
			name:   "negative test #3 - wrong name of metric",
			method: http.MethodGet,
			url:    "/value/gauge/none",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    "not found metric\n",
			},
		},
		{
			name:   "negative test #4 - without name of metric",
			method: http.MethodGet,
			url:    "/value/gauge",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    "Not Found\n",
			},
		},
		{
			name:   "negative test #5 - without name and type of metric",
			method: http.MethodGet,
			url:    "/value",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    "Not Found\n",
			},
		},
		{
			name:   "negative test #6 - with right name and type of metric and some additional param",
			method: http.MethodGet,
			url:    "/value/gauge/Alloc/none",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    "404 page not found\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			writer := httptest.NewRecorder()

			routes := server.GetRoutes(metricRepository)
			routes.ServeHTTP(writer, request)

			res := writer.Result()
			defer res.Body.Close()

			resultBody, errBody := io.ReadAll(res.Body)
			require.NoError(t, errBody)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("content-type"))
			assert.Equal(t, tt.want.response, string(resultBody))
		})
	}
}

func TestShowPossibleValue(t *testing.T) {
	metricRepository := metricmemoryrepository.New()

	metricsObject := metrics.New()

	for metricName, metricValue := range metricsObject.CollectData() {
		metricRepository.InsertValue(metricName, metricValue)
	}

	listMetrics, err := metricRepository.GetAll()
	require.NoError(t, err)
	require.NotEmpty(t, listMetrics)

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodGet,
			url:    "/",
			want: want{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name:   "negative test #1 - wrong http method",
			method: http.MethodPost,
			url:    "/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2 - wrong general page",
			method: http.MethodPost,
			url:    "/none",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			writer := httptest.NewRecorder()

			routes := server.GetRoutes(metricRepository)
			routes.ServeHTTP(writer, request)

			res := writer.Result()
			defer res.Body.Close()

			t.Log(res.Header.Get("content-type"))

			resultBody, errBody := io.ReadAll(res.Body)
			require.NoError(t, errBody)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("content-type"))

			assert.NotEmpty(t, resultBody)
		})
	}
}

func TestUpdateMetricJSON(t *testing.T) {
	metricRepository := metricmemoryrepository.New()

	type want struct {
		code        int
		contentType string
	}

	testValue := float64(32.33)

	tests := []struct {
		name   string
		method string
		url    string
		metric metrics.Metrics
		want   want
	}{
		{
			name:   "positive test #1 add new metric",
			method: http.MethodPost,
			url:    "/update",
			metric: metrics.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: &testValue,
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json; charset=UTF-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, errJSON := json.Marshal(tt.metric)
			assert.NoError(t, errJSON)

			bodyBuff := bytes.NewBuffer(body)
			request := httptest.NewRequest(tt.method, tt.url, bodyBuff)
			writer := httptest.NewRecorder()

			routes := server.GetRoutes(metricRepository)
			routes.ServeHTTP(writer, request)

			result := writer.Result()
			defer result.Body.Close()

			expectedValue, errGetValue := metricRepository.GetValue(tt.metric.ID)
			assert.NoError(t, errGetValue)

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("content-type"))
			assert.Equal(t, metrics.Gauge(testValue), expectedValue)
		})
	}
}

func TestGetValueJSON(t *testing.T) {
	metricRepository := metricmemoryrepository.New()

	type want struct {
		code        int
		contentType string
		resultValue interface{}
	}

	tests := []struct {
		name   string
		method string
		url    string
		value  interface{}
		metric metrics.Metrics
		want   want
	}{
		{
			name:   "positive test #1 get metric with type gauge",
			method: http.MethodPost,
			url:    "/value",
			value:  metrics.Gauge(32.33),
			metric: metrics.Metrics{
				ID:    "Alloc",
				MType: "gauge",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json; charset=UTF-8",
				resultValue: float64(32.33),
			},
		},
		{
			name:   "positive test #2 get metric with type counter",
			method: http.MethodPost,
			url:    "/value",
			value:  metrics.Counter(32),
			metric: metrics.Metrics{
				ID:    "PollCount",
				MType: "counter",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json; charset=UTF-8",
				resultValue: int64(32),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricRepository.InsertValue(tt.metric.ID, tt.value)

			requestBody, errJSON := json.Marshal(tt.metric)
			assert.NoError(t, errJSON)

			bodyBuff := bytes.NewBuffer(requestBody)
			request := httptest.NewRequest(tt.method, tt.url, bodyBuff)
			writer := httptest.NewRecorder()

			routes := server.GetRoutes(metricRepository)
			routes.ServeHTTP(writer, request)

			result := writer.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("content-type"))

			resultBody, errReadALl := io.ReadAll(result.Body)
			assert.NoError(t, errReadALl)

			var resultMetric metrics.Metrics
			errDecode := json.Unmarshal(resultBody, &resultMetric)
			assert.NoError(t, errDecode)

			fmt.Printf("== resultMetric: %v\n", resultMetric)

			if tt.metric.MType == "counter" {
				assert.Equal(t, tt.want.resultValue, *resultMetric.Delta)
				return
			}

			assert.Equal(t, tt.want.resultValue, *resultMetric.Value)
		})
	}
}
