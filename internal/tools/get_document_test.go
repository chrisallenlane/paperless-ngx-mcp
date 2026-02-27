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
