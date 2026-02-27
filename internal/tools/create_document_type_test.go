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

func TestCreateDocumentType_Execute_ServerError(t *testing.T) {
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
	tool := NewCreateDocumentType(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Test"}`),
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

func TestCreateDocumentType_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateDocumentType(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestCreateDocumentType_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateDocumentType(c)

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

	foundName := false
	for _, r := range required {
		if r == "name" {
			foundName = true
		}
	}
	if !foundName {
		t.Error("name should be in required fields")
	}
}
