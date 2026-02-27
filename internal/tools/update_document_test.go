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

func TestUpdateDocument_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/1/" {
					t.Errorf(
						"Expected /api/documents/1/, got %s",
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

				if patch["title"] != "Updated Title" {
					t.Errorf(
						"title = %v, want Updated Title",
						patch["title"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": 1,
					"title": "Updated Title",
					"content": "Invoice from ACME Corp",
					"correspondent": 1,
					"document_type": 2,
					"storage_path": null,
					"tags": [1, 3],
					"created": "2024-01-15T10:30:00Z",
					"created_date": "2024-01-15",
					"added": "2024-01-16T08:00:00Z",
					"modified": "2024-01-17T10:00:00Z",
					"archive_serial_number": 42,
					"original_file_name": "invoice.pdf",
					"archived_file_name": null,
					"mime_type": "application/pdf",
					"page_count": 2,
					"custom_fields": []
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
	tool := NewUpdateDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 1, "title": "Updated Title"}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document (ID: 1)",
		"Title: Updated Title",
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

func TestUpdateDocument_Execute_WithTags(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %v", err)
				}

				var patch map[string]json.RawMessage
				if err := json.Unmarshal(
					body,
					&patch,
				); err != nil {
					t.Fatalf("Failed to parse body: %v", err)
				}

				if _, ok := patch["tags"]; !ok {
					t.Error("Body should contain tags")
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": 1,
					"title": "Tagged Doc",
					"content": "",
					"correspondent": null,
					"document_type": null,
					"storage_path": null,
					"tags": [1, 2, 3],
					"created": "2024-01-15T10:30:00Z",
					"created_date": "2024-01-15",
					"added": "2024-01-16T08:00:00Z",
					"modified": "2024-01-17T10:00:00Z",
					"archive_serial_number": null,
					"original_file_name": null,
					"archived_file_name": null,
					"mime_type": "application/pdf",
					"page_count": null,
					"custom_fields": []
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
	tool := NewUpdateDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 1, "tags": [1, 2, 3]}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "Tags: 1, 2, 3") {
		t.Errorf(
			"Output missing tags.\nGot:\n%s",
			result,
		)
	}
}

func TestUpdateDocument_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"title": "Test"}`),
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

func TestUpdateDocument_Execute_InvalidID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateDocument(c)

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

func TestUpdateDocument_Execute_ServerError(t *testing.T) {
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
	tool := NewUpdateDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1, "title": "Test"}`),
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

func TestUpdateDocument_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateDocument(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestUpdateDocument_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUpdateDocument(c)

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
