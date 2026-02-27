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

const documentTypeResponse = `{
	"id": 1,
	"slug": "invoice",
	"name": "Invoice",
	"match": "invoice",
	"matching_algorithm": 1,
	"is_insensitive": true,
	"document_count": 10
}`

func TestGetDocumentType_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/document_types/1/" {
					t.Errorf(
						"Expected /api/document_types/1/, got %s",
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
				w.Write([]byte(documentTypeResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetDocumentType(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document Type (ID: 1)",
		"Name: Invoice",
		"Slug: invoice",
		"Match: invoice",
		"Matching Algorithm: 1 (Any word)",
		"Case Insensitive: true",
		"Document Count: 10",
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

func TestGetDocumentType_Execute_EmptyMatch(t *testing.T) {
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
					"slug": "receipt",
					"name": "Receipt",
					"match": "",
					"matching_algorithm": 6,
					"is_insensitive": true,
					"document_count": 0
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
	tool := NewGetDocumentType(c)

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

func TestGetDocumentType_Execute_InvalidID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentType(c)

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

func TestGetDocumentType_Execute_ServerError(t *testing.T) {
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
	tool := NewGetDocumentType(c)

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

func TestGetDocumentType_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentType(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestGetDocumentType_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentType(c)

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
