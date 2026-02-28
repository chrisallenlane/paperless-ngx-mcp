package tools

import (
	"encoding/json"
	"math"
	"testing"
)

func FuzzFormatDate(f *testing.F) {
	f.Add("2024-01-15T10:30:00Z")
	f.Add("")
	f.Add("short")
	f.Add("2024-01-15")
	f.Add("x")
	f.Add("0123456789abcdef")

	f.Fuzz(func(_ *testing.T, s string) {
		_ = formatDate(s)
	})
}

func FuzzFormatFileSize(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(1024))
	f.Add(int64(1048576))
	f.Add(int64(1073741824))
	f.Add(int64(-1))
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))

	f.Fuzz(func(_ *testing.T, bytes int64) {
		_ = formatFileSize(bytes)
	})
}

func FuzzFormatStatistics(f *testing.F) {
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"documents_total":100,"inbox_count":5}`))
	f.Add([]byte(`{"items":[1,2,3]}`))
	f.Add([]byte(`{"nested":{"a":1}}`))
	f.Add([]byte(`{"float_val":3.14,"neg":-1}`))
	f.Add([]byte(`{"arr":[{"k":"v"},{"k2":"v2"}]}`))
	f.Add([]byte(`{"empty_arr":[]}`))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var stats map[string]interface{}
		if err := json.Unmarshal(data, &stats); err != nil {
			return
		}
		_ = formatStatistics(stats)
	})
}
