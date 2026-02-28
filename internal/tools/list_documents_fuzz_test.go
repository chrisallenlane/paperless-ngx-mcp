package tools

import (
	"encoding/json"
	"testing"
)

func FuzzBuildDocumentListPath(f *testing.F) {
	f.Add([]byte(`{}`))
	f.Add(
		[]byte(
			`{"page":1,"page_size":25,"search":"invoice","correspondent":5,"document_type":3,"tags":[1,2,3],"is_in_inbox":true}`,
		),
	)
	f.Add([]byte(`{"tags":[]}`))
	f.Add([]byte(`{"search":"","correspondent":null}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = buildDocumentListPath(json.RawMessage(data))
	})
}
