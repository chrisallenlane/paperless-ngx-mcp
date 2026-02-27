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

const correspondentListResponse = `{
	"count": 2,
	"next": null,
	"previous": null,
	"all": [1, 2],
	"results": [
		{
			"id": 1,
			"slug": "acme-corp",
			"name": "ACME Corp",
			"match": "acme",
			"matching_algorithm": 1,
			"is_insensitive": true,
			"document_count": 5
		},
		{
			"id": 2,
			"slug": "john-doe",
			"name": "John Doe",
			"match": "",
			"matching_algorithm": 6,
			"is_insensitive": true,
			"document_count": 0
		}
	]
}`

func TestListCorrespondents_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/correspondents/" {
					t.Errorf(
						"Expected /api/correspondents/, got %s",
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
				w.Write([]byte(correspondentListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListCorrespondents(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Correspondents: 2 total",
		"ACME Corp (ID: 1)",
		"5 documents",
		"John Doe (ID: 2)",
		"0 documents",
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

func TestListCorrespondents_Execute_WithNameFilter(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("name__icontains") != "acme" {
					t.Errorf(
						"Expected name__icontains=acme, got %s",
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
						"slug": "acme-corp",
						"name": "ACME Corp",
						"match": "acme",
						"matching_algorithm": 1,
						"is_insensitive": true,
						"document_count": 5
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
	tool := NewListCorrespondents(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "acme"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "ACME Corp") {
		t.Errorf("Output missing ACME Corp.\nGot:\n%s", result)
	}
}

func TestListCorrespondents_Execute_Empty(t *testing.T) {
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
	tool := NewListCorrespondents(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No correspondents found." {
		t.Errorf(
			"Expected empty message, got: %s",
			result,
		)
	}
}
