package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

func TestUpdateCustomField_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/custom_fields/1/" {
					t.Errorf(
						"Expected /api/custom_fields/1/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "PATCH" {
					t.Errorf("Expected PATCH, got %s", r.Method)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %v", err)
				}

				var patch map[string]interface{}
				if err := json.Unmarshal(
					body,
					&patch,
				); err != nil {
					t.Fatalf("Failed to parse body: %v", err)
				}

				if _, ok := patch["id"]; ok {
					t.Error("Body should not contain id")
				}

				if patch["name"] != "Updated Field" {
					t.Errorf(
						"name = %v, want Updated Field",
						patch["name"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": 1,
					"name": "Updated Field",
					"data_type": "string",
					"extra_data": null,
					"document_count": 10
				}`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewUpdateCustomField(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 1, "name": "Updated Field"}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Custom Field (ID: 1)",
		"Name: Updated Field",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf(
				"Output missing %q.\nGot:\n%s",
				check,
				result,
			)
		}
	}
}

func TestUpdateCustomField_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateCustomField(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Test"}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing id")
	}

	if !strings.Contains(err.Error(), "id is required") {
		t.Errorf(
			"Error should mention id is required, got: %s",
			err.Error(),
		)
	}
}
