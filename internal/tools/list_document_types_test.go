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

const documentTypeListResponse = `{
	"count": 2,
	"next": null,
	"previous": null,
	"all": [1, 2],
	"results": [
		{
			"id": 1,
			"slug": "invoice",
			"name": "Invoice",
			"match": "invoice",
			"matching_algorithm": 1,
			"is_insensitive": true,
			"document_count": 10
		},
		{
			"id": 2,
			"slug": "receipt",
			"name": "Receipt",
			"match": "",
			"matching_algorithm": 6,
			"is_insensitive": true,
			"document_count": 3
		}
	]
}`

func TestListDocumentTypes_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/document_types/" {
					t.Errorf(
						"Expected /api/document_types/, got %s",
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
				w.Write([]byte(documentTypeListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListDocumentTypes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document Types: 2 total",
		"Invoice (ID: 1)",
		"10 documents",
		"Receipt (ID: 2)",
		"3 documents",
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

func TestListDocumentTypes_Execute_WithNameFilter(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("name__icontains") != "invoice" {
					t.Errorf(
						"Expected name__icontains=invoice, got %s",
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
						"slug": "invoice",
						"name": "Invoice",
						"match": "invoice",
						"matching_algorithm": 1,
						"is_insensitive": true,
						"document_count": 10
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
	tool := NewListDocumentTypes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "invoice"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "Invoice") {
		t.Errorf("Output missing Invoice.\nGot:\n%s", result)
	}
}

func TestListDocumentTypes_Execute_Empty(t *testing.T) {
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
	tool := NewListDocumentTypes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No document types found." {
		t.Errorf(
			"Expected empty message, got: %s",
			result,
		)
	}
}

func TestListDocumentTypes_Execute_WithPagination(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"count": 30,
					"next": "http://example.com/api/document_types/?page=2",
					"previous": null,
					"all": [1],
					"results": [{
						"id": 1,
						"slug": "invoice",
						"name": "Invoice",
						"match": "",
						"matching_algorithm": 1,
						"is_insensitive": true,
						"document_count": 10
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
	tool := NewListDocumentTypes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "more results available") {
		t.Errorf(
			"Output should show pagination hint.\nGot:\n%s",
			result,
		)
	}
}
