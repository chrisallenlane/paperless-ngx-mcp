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

const documentResponse = `{
	"id": 1,
	"title": "Invoice 2024-001",
	"content": "Invoice from ACME Corp for services rendered in Q1 2024.",
	"correspondent": 1,
	"document_type": 2,
	"storage_path": 3,
	"tags": [1, 3, 5],
	"created": "2024-01-15T10:30:00Z",
	"created_date": "2024-01-15",
	"added": "2024-01-16T08:00:00Z",
	"modified": "2024-01-16T08:00:00Z",
	"archive_serial_number": 42,
	"original_file_name": "invoice-2024-001.pdf",
	"archived_file_name": "invoice-2024-001-archived.pdf",
	"mime_type": "application/pdf",
	"page_count": 2,
	"custom_fields": [
		{"field": 1, "value": "important"},
		{"field": 2, "value": 100.50}
	]
}`

func TestGetDocument_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/1/" {
					t.Errorf(
						"Expected /api/documents/1/, got %s",
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
				w.Write([]byte(documentResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document (ID: 1)",
		"Title: Invoice 2024-001",
		"Correspondent: 1",
		"Document Type: 2",
		"Storage Path: 3",
		"Tags: 1, 3, 5",
		"ASN: 42",
		"Original File: invoice-2024-001.pdf",
		"MIME Type: application/pdf",
		"Page Count: 2",
		"Custom Fields:",
		"Field 1:",
		"Content: Invoice from ACME Corp",
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

func TestGetDocument_Execute_NullableFields(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": 2,
					"title": "Bare Document",
					"content": "",
					"correspondent": null,
					"document_type": null,
					"storage_path": null,
					"tags": [],
					"created": "2024-03-10T14:00:00Z",
					"created_date": "2024-03-10",
					"added": "2024-03-11T09:00:00Z",
					"modified": "2024-03-11T09:00:00Z",
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
	tool := NewGetDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 2}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Correspondent: (none)",
		"Document Type: (none)",
		"Storage Path: (none)",
		"Tags: (none)",
		"ASN: (none)",
		"Original File: (none)",
		"Page Count: (none)",
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

func TestGetDocument_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing id")
	}

	if !strings.Contains(err.Error(), "positive integer") {
		t.Errorf(
			"Error should mention positive integer, got: %s",
			err.Error(),
		)
	}
}
