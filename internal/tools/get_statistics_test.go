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

const statisticsResponse = `{
	"documents_total": 150,
	"documents_inbox": 5,
	"inbox_tags": [1, 2],
	"document_file_type_counts": [
		{"mime_type": "application/pdf", "count": 100},
		{"mime_type": "text/plain", "count": 50}
	],
	"character_count": 500000,
	"tag_count": 20,
	"correspondent_count": 10,
	"document_type_count": 8,
	"storage_path_count": 3
}`

func TestGetStatistics_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/statistics/" {
					t.Errorf(
						"Expected /api/statistics/, got %s",
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
				w.Write([]byte(statisticsResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetStatistics(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Paperless-NGX Statistics",
		"Documents Total: 150",
		"Documents Inbox: 5",
		"Character Count: 500000",
		"Tag Count: 20",
		"Correspondent Count: 10",
		"Document Type Count: 8",
		"Storage Path Count: 3",
		"Inbox Tags:",
		"Document File Type Counts:",
		"application/pdf",
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

func TestGetStatistics_Empty(t *testing.T) {
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
				w.Write([]byte(`{}`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetStatistics(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No statistics available." {
		t.Errorf(
			"Expected 'No statistics available.', got: %s",
			result,
		)
	}
}
