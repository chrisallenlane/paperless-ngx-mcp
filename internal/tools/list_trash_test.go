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

const trashListResponse = `{
	"count": 1,
	"next": null,
	"previous": null,
	"all": [1],
	"results": [
		{
			"id": 1,
			"title": "Deleted Invoice",
			"content": "",
			"correspondent": null,
			"document_type": null,
			"storage_path": null,
			"tags": [],
			"created": "2024-01-15T10:30:00Z",
			"created_date": "2024-01-15",
			"added": "2024-01-16T08:00:00Z",
			"modified": "2024-01-16T08:00:00Z",
			"archive_serial_number": null,
			"original_file_name": null,
			"archived_file_name": null,
			"mime_type": "application/pdf",
			"page_count": null,
			"custom_fields": []
		}
	]
}`

func TestListTrash_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/trash/" {
					t.Errorf(
						"Expected /api/trash/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "GET" {
					t.Errorf(
						"Expected GET, got %s",
						r.Method,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(trashListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTrash(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Documents: 1 total",
		"Deleted Invoice (ID: 1)",
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

func TestListTrash_Empty(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
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
	tool := NewListTrash(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No documents found." {
		t.Errorf(
			"Expected 'No documents found.', got: %s",
			result,
		)
	}
}
