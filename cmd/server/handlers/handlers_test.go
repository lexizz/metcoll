package handlers_test

import (
	"github.com/lexizz/metcoll/cmd/server/handlers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetric(t *testing.T) {
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
				contentType: "text/plain",
			},
		},
		{
			name:   "positive test #2 with type counter",
			method: http.MethodPost,
			url:    "/update/counter/Alloc/5",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:   "wrong type name",
			method: http.MethodPost,
			url:    "/update/unknown/Alloc/5",
			want: want{
				code:        http.StatusNotImplemented,
				contentType: "text/plain",
			},
		},
		{
			name:   "wrong value of metric",
			method: http.MethodPost,
			url:    "/update/counter/Alloc/sdf",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		//{
		//	name:   "wrong name of metric and right value",
		//	method: http.MethodPost,
		//	url:    "/update/counter/none/5",
		//	want: want{
		//		code:        http.StatusBadRequest,
		//		contentType: "text/plain",
		//	},
		//},
		//{
		//	name:   "wrong name of metric and wrong value",
		//	method: http.MethodPost,
		//	url:    "/update/counter/none/none",
		//	want: want{
		//		code:        http.StatusBadRequest,
		//		contentType: "text/plain",
		//	},
		//},
		{
			name:   "without value of metric",
			method: http.MethodPost,
			url:    "/update/counter/Alloc/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "without value of metric and name",
			method: http.MethodPost,
			url:    "/update/counter/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "without value of metric and name and type",
			method: http.MethodPost,
			url:    "/update",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "wrong url like general page",
			method: http.MethodPost,
			url:    "/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "wrong http method with right url",
			method: http.MethodGet,
			url:    "/update/gauge/Alloc/5",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
		{
			name:   "wrong http method without right url",
			method: http.MethodGet,
			url:    "/",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			writer := httptest.NewRecorder()

			handler := http.HandlerFunc(handlers.UpdateMetric())
			handler.ServeHTTP(writer, request)

			res := writer.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			//assert.Equal(t, tt.want.contentType, res.Header.Get("content-type"))
		})
	}
}
