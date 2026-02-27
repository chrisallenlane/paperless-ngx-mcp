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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("Success", func(t *testing.T) {
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

			t.Run("InvalidID", func(t *testing.T) {
				c := client.New(
					"http://localhost",
					"test-token",
				)
				tool := tt.newTool(c)

				_, err := tool.Execute(
					context.Background(),
					json.RawMessage(`{"id": 0}`),
				)
				if err == nil {
					t.Fatal("Expected error for invalid id")
				}
				if !strings.Contains(
					err.Error(),
					"positive integer",
				) {
					t.Errorf(
						"Error should mention positive integer, got: %s",
						err.Error(),
					)
				}
			})

			t.Run("NegativeID", func(t *testing.T) {
				c := client.New(
					"http://localhost",
					"test-token",
				)
				tool := tt.newTool(c)

				_, err := tool.Execute(
					context.Background(),
					json.RawMessage(`{"id": -1}`),
				)
				if err == nil {
					t.Fatal("Expected error for negative id")
				}
				if !strings.Contains(
					err.Error(),
					"positive integer",
				) {
					t.Errorf(
						"Error should mention positive integer, got: %s",
						err.Error(),
					)
				}
			})

			t.Run("ServerError", func(t *testing.T) {
				server := httptest.NewServer(
					http.HandlerFunc(
						func(
							w http.ResponseWriter,
							_ *http.Request,
						) {
							w.WriteHeader(
								http.StatusInternalServerError,
							)
							w.Write(
								[]byte(
									"Internal Server Error",
								),
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

				_, err := tool.Execute(
					context.Background(),
					json.RawMessage(`{"id": 1}`),
				)
				if err == nil {
					t.Fatal(
						"Expected error for server error response",
					)
				}
				if !strings.Contains(
					err.Error(),
					"500",
				) {
					t.Errorf(
						"Error should mention status code, got: %s",
						err.Error(),
					)
				}
			})

			t.Run("Description", func(t *testing.T) {
				c := client.New(
					"http://localhost",
					"test-token",
				)
				tool := tt.newTool(c)
				if tool.Description() == "" {
					t.Error(
						"Description should not be empty",
					)
				}
			})

			t.Run("InputSchema", func(t *testing.T) {
				c := client.New(
					"http://localhost",
					"test-token",
				)
				tool := tt.newTool(c)
				schema := tool.InputSchema()
				if schema == nil {
					t.Fatal(
						"InputSchema should not be nil",
					)
				}

				schemaType, ok := schema["type"].(string)
				if !ok || schemaType != "object" {
					t.Errorf(
						"Schema type = %v, want object",
						schema["type"],
					)
				}

				required, ok := schema["required"].([]string)
				if !ok {
					t.Fatal(
						"Schema should have required field",
					)
				}

				foundID := false
				for _, r := range required {
					if r == "id" {
						foundID = true
					}
				}
				if !foundID {
					t.Error(
						"id should be in required fields",
					)
				}
			})
		})
	}
}
