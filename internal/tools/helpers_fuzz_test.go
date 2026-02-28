package tools

import (
	"encoding/json"
	"testing"
)

func FuzzValidateFilePath(f *testing.F) {
	f.Add("/absolute/path")
	f.Add("relative/path")
	f.Add("../traversal")
	f.Add("")
	f.Add("..")
	f.Add("/foo/../bar")
	f.Add("/foo/\x00bar")
	f.Add("//double/slash")
	f.Add("/a/b/c/d/e/f/g/h")

	f.Fuzz(func(_ *testing.T, path string) {
		// must not panic; errors are expected
		_ = validateFilePath(path)
	})
}

func FuzzParseIDArg(f *testing.F) {
	f.Add([]byte(`{"id":1}`))
	f.Add([]byte(`{"id":-1}`))
	f.Add([]byte(`{"id":0}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(`{"id":"string"}`))
	f.Add([]byte(`{"id":999999999}`))
	f.Add([]byte(`{"id":1.5}`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = parseIDArg(json.RawMessage(data))
	})
}

func FuzzParsePatchArgs(f *testing.F) {
	f.Add([]byte(`{"id":1,"name":"foo"}`))
	f.Add([]byte(`{"id":1}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(`{"id":"string","name":"foo"}`))
	f.Add([]byte(`{"id":1,"a":null,"b":true,"c":[1,2]}`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _, _ = parsePatchArgs(json.RawMessage(data))
	})
}

func FuzzBuildListPath(f *testing.F) {
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"page":1}`))
	f.Add([]byte(`{"page_size":25}`))
	f.Add([]byte(`{"name":"test"}`))
	f.Add([]byte(`{"page":1,"page_size":10,"name":"foo"}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = buildListPath("/api/test/", json.RawMessage(data))
	})
}
