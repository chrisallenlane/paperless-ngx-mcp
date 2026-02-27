package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

func TestCreateDocumentType_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/document_types/" {
					t.Errorf(
						"Expected /api/document_types/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %v", err)
				}

				var req map[string]interface{}
				if err := json.Unmarshal(
					body,
					&req,
				); err != nil {
					t.Fatalf("Failed to parse body: %v", err)
				}

				if req["name"] != "Invoice" {
					t.Errorf(
						"name = %v, want Invoice",
						req["name"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"id": 1,
					"slug": "invoice",
					"name": "Invoice",
					"match": "",
					"matching_algorithm": 1,
					"is_insensitive": true
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
	tool := NewCreateDocumentType(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Invoice"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document Type (ID: 1)",
		"Name: Invoice",
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

func TestCreateDocumentType_Execute_MissingName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateDocumentType(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing name")
	}

	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf(
			"Error should mention name is required, got: %s",
			err.Error(),
		)
	}
}

func TestCreateDocumentType_Execute_EmptyName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateDocumentType(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": ""}`),
	)
	if err == nil {
		t.Fatal("Expected error for empty name")
	}

	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf(
			"Error should mention name is required, got: %s",
			err.Error(),
		)
	}
}
