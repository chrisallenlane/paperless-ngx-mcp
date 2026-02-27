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

func TestDeleteDocument_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/1/" {
					t.Errorf(
						"Expected /api/documents/1/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "DELETE" {
					t.Errorf(
						"Expected DELETE, got %s",
						r.Method,
					)
				}

				w.WriteHeader(http.StatusNoContent)
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewDeleteDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "Document 1 deleted successfully."
	if result != expected {
		t.Errorf("Result = %q, want %q", result, expected)
	}
}

func TestDeleteDocument_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDeleteDocument(c)

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

func TestDeleteDocument_Execute_InvalidID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDeleteDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": -1}`),
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

func TestDeleteDocument_Execute_ServerError(t *testing.T) {
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
	tool := NewDeleteDocument(c)

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

func TestDeleteDocument_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDeleteDocument(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestDeleteDocument_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDeleteDocument(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	schemaType, ok := schema["type"].(string)
	if !ok || schemaType != "object" {
		t.Errorf("Schema type = %v, want object", schema["type"])
	}
}
