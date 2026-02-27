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

func TestCreateStoragePath_MissingName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateStoragePath(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path": "/test/"}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing name")
	}

	if !strings.Contains(
		err.Error(),
		"name is required",
	) {
		t.Errorf(
			"Error should mention name is required, got: %s",
			err.Error(),
		)
	}
}

func TestCreateStoragePath_MissingPath(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateStoragePath(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Test"}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing path")
	}

	if !strings.Contains(
		err.Error(),
		"path is required",
	) {
		t.Errorf(
			"Error should mention path is required, got: %s",
			err.Error(),
		)
	}
}

func TestCreateStoragePath_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/storage_paths/" {
					t.Errorf(
						"Expected /api/storage_paths/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "POST" {
					t.Errorf(
						"Expected POST, got %s",
						r.Method,
					)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf(
						"Failed to read body: %v",
						err,
					)
				}

				var req map[string]interface{}
				if err := json.Unmarshal(
					body,
					&req,
				); err != nil {
					t.Fatalf(
						"Failed to parse body: %v",
						err,
					)
				}

				if req["name"] != "Invoices" {
					t.Errorf(
						"name = %v, want Invoices",
						req["name"],
					)
				}

				if req["path"] != "{correspondent}/{title}" {
					t.Errorf(
						"path = %v, want template",
						req["path"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"id": 1,
					"slug": "invoices",
					"name": "Invoices",
					"path": "{correspondent}/{title}",
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
	tool := NewCreateStoragePath(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Invoices", "path": "{correspondent}/{title}"}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Storage Path (ID: 1)",
		"Name: Invoices",
		"Path: {correspondent}/{title}",
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
