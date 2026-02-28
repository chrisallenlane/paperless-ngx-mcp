package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

func TestDelete_Execute(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		newTool  func(*client.Client) Tool
		expected string
	}{
		{
			name: "Correspondent",
			path: "/api/correspondents/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteCorrespondent(c)
			},
			expected: "Correspondent 1 deleted successfully.",
		},
		{
			name: "CustomField",
			path: "/api/custom_fields/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteCustomField(c)
			},
			expected: "Custom field 1 deleted successfully.",
		},
		{
			name: "DocumentType",
			path: "/api/document_types/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteDocumentType(c)
			},
			expected: "Document type 1 deleted successfully.",
		},
		{
			name: "Document",
			path: "/api/documents/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteDocument(c)
			},
			expected: "Document 1 deleted successfully.",
		},
		{
			name: "Tag",
			path: "/api/tags/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteTag(c)
			},
			expected: "Tag 1 deleted successfully.",
		},
		{
			name: "StoragePath",
			path: "/api/storage_paths/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteStoragePath(c)
			},
			expected: "Storage path 1 deleted successfully.",
		},
		{
			name: "SavedView",
			path: "/api/saved_views/1/",
			newTool: func(c *client.Client) Tool {
				return NewDeleteSavedView(c)
			},
			expected: "Saved view 1 deleted successfully.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.URL.Path != tt.path {
							t.Errorf(
								"Expected %s, got %s",
								tt.path,
								r.URL.Path,
							)
						}
						if r.Method != "DELETE" {
							t.Errorf(
								"Expected DELETE, got %s",
								r.Method,
							)
						}
						w.WriteHeader(
							http.StatusNoContent,
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
			tool := tt.newTool(c)

			result, err := tool.Execute(
				context.Background(),
				json.RawMessage(`{"id": 1}`),
			)
			if err != nil {
				t.Fatalf(
					"Unexpected error: %v",
					err,
				)
			}

			if result != tt.expected {
				t.Errorf(
					"Result = %q, want %q",
					result,
					tt.expected,
				)
			}
		})
	}
}
