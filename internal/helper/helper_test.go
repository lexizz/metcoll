package helper

import (
	"testing"

	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetType(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "test with type `int`",
			value: 5,
			want:  "int",
		},
		{
			name:  "test with type `int32`",
			value: int32(5),
			want:  "int32",
		},
		{
			name:  "test with type `int64`",
			value: int64(5),
			want:  "int64",
		},
		{
			name:  "test with type `metrics.Gauge`",
			value: metrics.Gauge(5),
			want:  "gauge",
		},
		{
			name:  "test with type `metrics.Counter`",
			value: metrics.Counter(5),
			want:  "counter",
		},
		{
			name:  "test with type `metrics.Counter`",
			value: int64(5),
			want:  "int64",
		},
		{
			name:  "test with type `struct`",
			value: struct{}{},
			want:  "struct {}",
		},
		{
			name:  "test with type `float32`",
			value: float32(5),
			want:  "float32",
		},
		{
			name:  "test with type `float64`",
			value: float64(5),
			want:  "float64",
		},
		{
			name:  "test with type `string`",
			value: "test string",
			want:  "string",
		},
		{
			name:  "test with type `map`",
			value: map[string]int{},
			want:  "map[string]int",
		},
		{
			name:  "test with type `slice`",
			value: []string{},
			want:  "[]string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultType, err := GetType(tt.value)

			assert.Nil(t, err)
			require.NotEmpty(t, resultType)

			assert.Equal(t, tt.want, resultType)
		})
	}
}
