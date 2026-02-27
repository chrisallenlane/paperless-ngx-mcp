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

const documentSuggestionsResponse = `{
	"correspondents": [1, 3],
	"document_types": [2],
	"storage_paths": [1],
	"tags": [1, 5, 7],
	"dates": ["2024-01-15", "2024-02-01"]
}`

func TestGetDocumentSuggestions_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/1/suggestions/" {
					t.Errorf(
						"Expected /api/documents/1/suggestions/, got %s",
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
				w.Write(
					[]byte(documentSuggestionsResponse),
				)
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetDocumentSuggestions(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document Suggestions (ID: 1)",
		"Correspondents: 1, 3",
		"Document Types: 2",
		"Storage Paths: 1",
		"Tags: 1, 5, 7",
		"Dates: 2024-01-15, 2024-02-01",
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

func TestGetDocumentSuggestions_Execute_Empty(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"correspondents": [],
					"document_types": [],
					"storage_paths": [],
					"tags": [],
					"dates": []
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
	tool := NewGetDocumentSuggestions(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Correspondents: (none)",
		"Document Types: (none)",
		"Storage Paths: (none)",
		"Tags: (none)",
		"Dates: (none)",
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

func TestGetDocumentSuggestions_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentSuggestions(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing id")
	}
}

func TestGetDocumentSuggestions_Execute_ServerError(t *testing.T) {
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
	tool := NewGetDocumentSuggestions(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
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

func TestGetDocumentSuggestions_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentSuggestions(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestGetDocumentSuggestions_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentSuggestions(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	schemaType, ok := schema["type"].(string)
	if !ok || schemaType != "object" {
		t.Errorf(
			"Schema type = %v, want object",
			schema["type"],
		)
	}
}
