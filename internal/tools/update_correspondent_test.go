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

func TestUpdateCorrespondent_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/correspondents/1/" {
					t.Errorf(
						"Expected /api/correspondents/1/, got %s",
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

				if patch["name"] != "Updated Corp" {
					t.Errorf(
						"name = %v, want Updated Corp",
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
					"slug": "updated-corp",
					"name": "Updated Corp",
					"match": "acme",
					"matching_algorithm": 1,
					"is_insensitive": true,
					"document_count": 5,
					"last_correspondence": null
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
	tool := NewUpdateCorrespondent(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 1, "name": "Updated Corp"}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Correspondent (ID: 1)",
		"Name: Updated Corp",
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

func TestUpdateCorrespondent_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateCorrespondent(c)

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

func TestUpdateCorrespondent_Execute_InvalidID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateCorrespondent(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 0}`),
	)
	if err == nil {
		t.Fatal("Expected error for invalid id")
	}

	if !strings.Contains(err.Error(), "positive integer") {
		t.Errorf(
			"Error should mention positive integer, got: %s",
			err.Error(),
		)
	}
}

func TestUpdateCorrespondent_Execute_ServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewUpdateCorrespondent(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1, "name": "Test"}`),
	)
	if err == nil {
		t.Fatal("Expected error for server error response")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"Error should mention status code, got: %s",
			err.Error(),
		)
	}
}

func TestUpdateCorrespondent_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateCorrespondent(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestUpdateCorrespondent_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateCorrespondent(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	schemaType, ok := schema["type"].(string)
	if !ok || schemaType != "object" {
		t.Errorf("Schema type = %v, want object", schema["type"])
	}

	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Schema should have required field")
	}

	foundID := false
	for _, r := range required {
		if r == "id" {
			foundID = true
		}
	}
	if !foundID {
		t.Error("id should be in required fields")
	}
}
