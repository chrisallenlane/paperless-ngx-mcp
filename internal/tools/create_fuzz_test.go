package tools

import (
	"encoding/json"
	"testing"
)

func FuzzValidateMatchableCreate(f *testing.F) {
	f.Add([]byte(`{"name":"test"}`))
	f.Add([]byte(`{}`))
	f.Add(
		[]byte(
			`{"name":"test","match":"foo","matching_algorithm":1,"is_insensitive":true}`,
		),
	)
	f.Add([]byte(`{"name":""}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = validateMatchableCreate(json.RawMessage(data))
	})
}

func FuzzValidateCreateTag(f *testing.F) {
	f.Add([]byte(`{"name":"inbox"}`))
	f.Add([]byte(`{}`))
	f.Add(
		[]byte(
			`{"name":"tag","color":"#ff0000","is_inbox_tag":true,"parent":5}`,
		),
	)
	f.Add([]byte(`{"name":""}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = validateCreateTag(json.RawMessage(data))
	})
}

func FuzzValidateCreateStoragePath(f *testing.F) {
	f.Add([]byte(`{"name":"archive","path":"/docs/{correspondent}"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"name":"test","path":""}`))
	f.Add([]byte(`{"name":"","path":"test"}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = validateCreateStoragePath(json.RawMessage(data))
	})
}

func FuzzValidateCreateCustomField(f *testing.F) {
	f.Add([]byte(`{"name":"field","data_type":"string"}`))
	f.Add([]byte(`{}`))
	f.Add(
		[]byte(
			`{"name":"field","data_type":"select","extra_data":{"options":["a","b"]}}`,
		),
	)
	f.Add([]byte(`{"name":"","data_type":"string"}`))
	f.Add([]byte(`{"name":"field","data_type":""}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = validateCreateCustomField(json.RawMessage(data))
	})
}

func FuzzValidateCreateSavedView(f *testing.F) {
	f.Add(
		[]byte(
			`{"name":"view","show_on_dashboard":true,"show_in_sidebar":false,"filter_rules":[{"rule_type":0,"value":"test"}]}`,
		),
	)
	f.Add([]byte(`{}`))
	f.Add(
		[]byte(
			`{"name":"view","show_on_dashboard":true,"show_in_sidebar":true,"filter_rules":[]}`,
		),
	)
	f.Add(
		[]byte(
			`{"name":"","show_on_dashboard":true,"show_in_sidebar":true,"filter_rules":[]}`,
		),
	)
	f.Add([]byte(`{"name":"view"}`))
	f.Add([]byte(`"not json"`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = validateCreateSavedView(json.RawMessage(data))
	})
}
