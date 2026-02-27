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

const correspondentResponse = `{
	"id": 1,
	"slug": "acme-corp",
	"name": "ACME Corp",
	"match": "acme",
	"matching_algorithm": 1,
	"is_insensitive": true,
	"document_count": 5,
	"last_correspondence": "2026-02-15T10:00:00Z"
}`

func TestGetCorrespondent_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/correspondents/1/" {
					t.Errorf(
						"Expected /api/correspondents/1/, got %s",
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
				w.Write([]byte(correspondentResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetCorrespondent(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Correspondent (ID: 1)",
		"Name: ACME Corp",
		"Slug: acme-corp",
		"Match: acme",
		"Matching Algorithm: 1 (Any word)",
		"Case Insensitive: true",
		"Document Count: 5",
		"Last Correspondence: 2026-02-15",
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

func TestGetCorrespondent_Execute_NullLastCorrespondence(
	t *testing.T,
) {
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
					"slug": "john-doe",
					"name": "John Doe",
					"match": "",
					"matching_algorithm": 6,
					"is_insensitive": true,
					"document_count": 0,
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
	tool := NewGetCorrespondent(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 2}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Match: (none)",
		"Matching Algorithm: 6 (Automatic)",
		"Last Correspondence: (none)",
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
