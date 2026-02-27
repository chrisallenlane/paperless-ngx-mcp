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

const customFieldListResponse = `{
	"count": 2,
	"next": null,
	"previous": null,
	"all": [1, 2],
	"results": [
		{
			"id": 1,
			"name": "Invoice Number",
			"data_type": "string",
			"extra_data": null,
			"document_count": 10
		},
		{
			"id": 2,
			"name": "Due Date",
			"data_type": "date",
			"extra_data": null,
			"document_count": 3
		}
	]
}`

func TestListCustomFields_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/custom_fields/" {
					t.Errorf(
						"Expected /api/custom_fields/, got %s",
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
				w.Write([]byte(customFieldListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListCustomFields(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Custom Fields: 2 total",
		"Invoice Number (ID: 1)",
		"type: string",
		"10 documents",
		"Due Date (ID: 2)",
		"type: date",
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

func TestListCustomFields_Execute_WithNameFilter(t *testing.T) {
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
						"name": "Invoice Number",
						"data_type": "string",
						"extra_data": null,
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
	tool := NewListCustomFields(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "invoice"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "Invoice Number") {
		t.Errorf(
			"Output missing Invoice Number.\nGot:\n%s",
			result,
		)
	}
}

func TestListCustomFields_Execute_Empty(t *testing.T) {
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
	tool := NewListCustomFields(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No custom fields found." {
		t.Errorf(
			"Expected empty message, got: %s",
			result,
		)
	}
}
