package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

const documentListResponse = `{
	"count": 2,
	"next": null,
	"previous": null,
	"all": [1, 2],
	"results": [
		{
			"id": 1,
			"title": "Invoice 2024-001",
			"content": "Invoice from ACME Corp",
			"correspondent": 1,
			"document_type": 2,
			"storage_path": null,
			"tags": [1, 3],
			"created": "2024-01-15T10:30:00Z",
			"created_date": "2024-01-15",
			"added": "2024-01-16T08:00:00Z",
			"modified": "2024-01-16T08:00:00Z",
			"archive_serial_number": 42,
			"original_file_name": "invoice.pdf",
			"archived_file_name": "invoice-archived.pdf",
			"mime_type": "application/pdf",
			"page_count": 2,
			"custom_fields": []
		},
		{
			"id": 2,
			"title": "Tax Return 2023",
			"content": "Federal tax return",
			"correspondent": null,
			"document_type": null,
			"storage_path": null,
			"tags": [],
			"created": "2024-03-10T14:00:00Z",
			"created_date": "2024-03-10",
			"added": "2024-03-11T09:00:00Z",
			"modified": "2024-03-11T09:00:00Z",
			"archive_serial_number": null,
			"original_file_name": "tax-2023.pdf",
			"archived_file_name": null,
			"mime_type": "application/pdf",
			"page_count": 10,
			"custom_fields": []
		}
	]
}`

func TestListDocuments_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/" {
					t.Errorf(
						"Expected /api/documents/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(documentListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListDocuments(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Documents: 2 total",
		"Invoice 2024-001 (ID: 1)",
		"Tax Return 2023 (ID: 2)",
		"Correspondent: 1",
		"Correspondent: (none)",
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

func TestListDocuments_Execute_WithSearch(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("search") != "invoice" {
					t.Errorf(
						"Expected search=invoice, got %s",
						r.URL.RawQuery,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"count": 1,
					"next": null,
					"previous": null,
					"all": [1],
					"results": [{
						"id": 1,
						"title": "Invoice 2024-001",
						"content": "Invoice from ACME",
						"correspondent": 1,
						"document_type": 2,
						"storage_path": null,
						"tags": [1],
						"created": "2024-01-15T10:30:00Z",
						"created_date": "2024-01-15",
						"added": "2024-01-16T08:00:00Z",
						"modified": "2024-01-16T08:00:00Z",
						"archive_serial_number": null,
						"original_file_name": "invoice.pdf",
						"archived_file_name": null,
						"mime_type": "application/pdf",
						"page_count": 2,
						"custom_fields": []
					}]
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
	tool := NewListDocuments(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"search": "invoice"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "Invoice 2024-001") {
		t.Errorf(
			"Output missing Invoice 2024-001.\nGot:\n%s",
			result,
		)
	}
}

func TestListDocuments_Execute_WithFilters(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				if q.Get("correspondent__id") != "1" {
					t.Errorf(
						"Expected correspondent__id=1, got %s",
						q.Get("correspondent__id"),
					)
				}
				if q.Get("document_type__id") != "2" {
					t.Errorf(
						"Expected document_type__id=2, got %s",
						q.Get("document_type__id"),
					)
				}
				tagValues := q["tags__id__all"]
				if len(tagValues) != 2 {
					t.Errorf(
						"Expected 2 tag filters, got %d",
						len(tagValues),
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"count": 0,
					"next": null,
					"previous": null,
					"all": [],
					"results": []
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
	tool := NewListDocuments(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"correspondent": 1, "document_type": 2, "tags": [1, 3]}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestListDocuments_Execute_Empty(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"count": 0,
					"next": null,
					"previous": null,
					"all": [],
					"results": []
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
	tool := NewListDocuments(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No documents found." {
		t.Errorf(
			"Expected empty message, got: %s",
			result,
		)
	}
}

func TestListDocuments_Execute_ServerError(t *testing.T) {
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
	tool := NewListDocuments(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
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

func TestListDocuments_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewListDocuments(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestListDocuments_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewListDocuments(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	schemaType, ok := schema["type"].(string)
	if !ok || schemaType != "object" {
		t.Errorf("Schema type = %v, want object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	expectedProps := []string{
		"page", "page_size", "search",
		"correspondent", "document_type",
		"tags", "is_in_inbox",
	}
	for _, prop := range expectedProps {
		if _, ok := props[prop]; !ok {
			t.Errorf("Schema missing property %q", prop)
		}
	}
}
