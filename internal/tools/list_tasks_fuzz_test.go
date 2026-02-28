package tools

import (
	"encoding/json"
	"testing"
)

func FuzzBuildTaskListPath(f *testing.F) {
	f.Add([]byte(`{}`))
	f.Add(
		[]byte(
			`{"status":"SUCCESS","task_name":"consume","type":"file","task_id":"abc-123"}`,
		),
	)
	f.Add([]byte(`{"status":"FAILURE"}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = buildTaskListPath(json.RawMessage(data))
	})
}
