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

func TestCreateCustomField_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/custom_fields/" {
					t.Errorf(
						"Expected /api/custom_fields/, got %s",
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

				if req["name"] != "Invoice Number" {
					t.Errorf(
						"name = %v, want Invoice Number",
						req["name"],
					)
				}

				if req["data_type"] != "string" {
					t.Errorf(
						"data_type = %v, want string",
						req["data_type"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"id": 1,
					"name": "Invoice Number",
					"data_type": "string",
					"extra_data": null,
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
	tool := NewCreateCustomField(c)

	args := `{
		"name": "Invoice Number",
		"data_type": "string"
	}`

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Custom Field (ID: 1)",
		"Name: Invoice Number",
		"Data Type: string",
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

func TestCreateCustomField_Execute_WithExtraData(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
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

				if req["extra_data"] == nil {
					t.Error("extra_data should be present")
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"id": 1,
					"name": "Category",
					"data_type": "select",
					"extra_data": {
						"select_options": ["a", "b"]
					},
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
	tool := NewCreateCustomField(c)

	args := `{
		"name": "Category",
		"data_type": "select",
		"extra_data": {"select_options": ["a", "b"]}
	}`

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "select_options") {
		t.Errorf(
			"Output missing extra_data.\nGot:\n%s",
			result,
		)
	}
}

func TestCreateCustomField_Execute_MissingName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateCustomField(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"data_type": "string"}`),
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

func TestCreateCustomField_Execute_MissingDataType(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateCustomField(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Test"}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing data_type")
	}

	if !strings.Contains(err.Error(), "data_type is required") {
		t.Errorf(
			"Error should mention data_type is required, got: %s",
			err.Error(),
		)
	}
}

func TestCreateCustomField_Execute_ServerError(t *testing.T) {
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
	tool := NewCreateCustomField(c)

	args := `{"name": "Test", "data_type": "string"}`

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
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

func TestCreateCustomField_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateCustomField(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestCreateCustomField_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateCustomField(c)

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

	requiredFields := map[string]bool{}
	for _, r := range required {
		requiredFields[r] = true
	}

	if !requiredFields["name"] {
		t.Error("name should be in required fields")
	}

	if !requiredFields["data_type"] {
		t.Error("data_type should be in required fields")
	}
}
