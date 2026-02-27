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

func TestCreateTag_MissingName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateTag(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
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

func TestCreateTag_EmptyName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateTag(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": ""}`),
	)
	if err == nil {
		t.Fatal("Expected error for empty name")
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

func TestCreateTag_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/tags/" {
					t.Errorf(
						"Expected /api/tags/, got %s",
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

				if req["name"] != "Important" {
					t.Errorf(
						"name = %v, want Important",
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
					"slug": "important",
					"name": "Important",
					"color": "#a6cee3",
					"text_color": "#000000",
					"match": "",
					"matching_algorithm": 1,
					"is_insensitive": true,
					"is_inbox_tag": false,
					"document_count": 0,
					"parent": null,
					"children": []
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
	tool := NewCreateTag(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Important"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Tag (ID: 1)",
		"Name: Important",
		"Color: #a6cee3",
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

func TestCreateTag_WithOptionalFields(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
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

				if req["name"] != "Inbox" {
					t.Errorf(
						"name = %v, want Inbox",
						req["name"],
					)
				}
				if req["color"] != "#ff0000" {
					t.Errorf(
						"color = %v, want #ff0000",
						req["color"],
					)
				}
				if req["is_inbox_tag"] != true {
					t.Errorf(
						"is_inbox_tag = %v, want true",
						req["is_inbox_tag"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"id": 2,
					"slug": "inbox",
					"name": "Inbox",
					"color": "#ff0000",
					"text_color": "#ffffff",
					"match": "",
					"matching_algorithm": 1,
					"is_insensitive": true,
					"is_inbox_tag": true,
					"document_count": 0,
					"parent": null,
					"children": []
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
	tool := NewCreateTag(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Inbox", "color": "#ff0000", "is_inbox_tag": true}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Tag (ID: 2)",
		"Name: Inbox",
		"Color: #ff0000",
		"Is Inbox Tag: true",
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
