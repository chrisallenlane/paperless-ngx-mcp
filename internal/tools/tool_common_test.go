package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// toolTestEntry defines a tool for shared cross-cutting tests.
type toolTestEntry struct {
	name       string
	newTool    func(*client.Client) Tool
	serverArgs string   // JSON args for ServerError test; "" to skip
	idArgsFmt  string   // fmt template for InvalidID/NegativeID; "" to skip
	required   []string // expected required schema fields; nil to skip
}

// allToolTests covers every tool except deletes (tested in
// delete_test.go).
var allToolTests = []toolTestEntry{
	// No-arg tools
	{
		name: "GetStatus",
		newTool: func(c *client.Client) Tool {
			return NewGetStatus(c)
		},
		serverArgs: `{}`,
	},
	{
		name: "GetConfig",
		newTool: func(c *client.Client) Tool {
			return NewGetConfig(c)
		},
		serverArgs: `{}`,
	},
	{
		name: "GetNextASN",
		newTool: func(c *client.Client) Tool {
			return NewGetNextASN(c)
		},
		serverArgs: `{}`,
	},

	// Get (ID-based) tools
	{
		name: "GetCorrespondent",
		newTool: func(c *client.Client) Tool {
			return NewGetCorrespondent(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "GetCustomField",
		newTool: func(c *client.Client) Tool {
			return NewGetCustomField(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "GetDocumentType",
		newTool: func(c *client.Client) Tool {
			return NewGetDocumentType(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "GetDocument",
		newTool: func(c *client.Client) Tool {
			return NewGetDocument(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "GetDocumentMetadata",
		newTool: func(c *client.Client) Tool {
			return NewGetDocumentMetadata(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "GetDocumentSuggestions",
		newTool: func(c *client.Client) Tool {
			return NewGetDocumentSuggestions(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},

	// List tools
	{
		name: "ListCorrespondents",
		newTool: func(c *client.Client) Tool {
			return NewListCorrespondents(c)
		},
		serverArgs: `{}`,
	},
	{
		name: "ListCustomFields",
		newTool: func(c *client.Client) Tool {
			return NewListCustomFields(c)
		},
		serverArgs: `{}`,
	},
	{
		name: "ListDocumentTypes",
		newTool: func(c *client.Client) Tool {
			return NewListDocumentTypes(c)
		},
		serverArgs: `{}`,
	},
	{
		name: "ListDocuments",
		newTool: func(c *client.Client) Tool {
			return NewListDocuments(c)
		},
		serverArgs: `{}`,
	},

	// Create tools
	{
		name: "CreateCorrespondent",
		newTool: func(c *client.Client) Tool {
			return NewCreateCorrespondent(c)
		},
		serverArgs: `{"name": "Test"}`,
		required:   []string{"name"},
	},
	{
		name: "CreateCustomField",
		newTool: func(c *client.Client) Tool {
			return NewCreateCustomField(c)
		},
		serverArgs: `{"name": "Test", "data_type": "string"}`,
		required:   []string{"name", "data_type"},
	},
	{
		name: "CreateDocumentType",
		newTool: func(c *client.Client) Tool {
			return NewCreateDocumentType(c)
		},
		serverArgs: `{"name": "Test"}`,
		required:   []string{"name"},
	},

	// Update tools
	{
		name: "UpdateCorrespondent",
		newTool: func(c *client.Client) Tool {
			return NewUpdateCorrespondent(c)
		},
		serverArgs: `{"id": 1, "name": "Test"}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "UpdateCustomField",
		newTool: func(c *client.Client) Tool {
			return NewUpdateCustomField(c)
		},
		serverArgs: `{"id": 1, "name": "Test"}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "UpdateDocumentType",
		newTool: func(c *client.Client) Tool {
			return NewUpdateDocumentType(c)
		},
		serverArgs: `{"id": 1, "name": "Test"}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "UpdateDocument",
		newTool: func(c *client.Client) Tool {
			return NewUpdateDocument(c)
		},
		serverArgs: `{"id": 1, "title": "Test"}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "UpdateConfig",
		newTool: func(c *client.Client) Tool {
			return NewUpdateConfig(c)
		},
		serverArgs: `{"id": 1, "deskew": true}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},

	// Delete tools
	{
		name: "DeleteCorrespondent",
		newTool: func(c *client.Client) Tool {
			return NewDeleteCorrespondent(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "DeleteCustomField",
		newTool: func(c *client.Client) Tool {
			return NewDeleteCustomField(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "DeleteDocumentType",
		newTool: func(c *client.Client) Tool {
			return NewDeleteDocumentType(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},
	{
		name: "DeleteDocument",
		newTool: func(c *client.Client) Tool {
			return NewDeleteDocument(c)
		},
		serverArgs: `{"id": 1}`,
		idArgsFmt:  `{"id": %d}`,
		required:   []string{"id"},
	},

	// Special tools (Description + InputSchema only)
	{
		name: "UploadDocument",
		newTool: func(c *client.Client) Tool {
			return NewUploadDocument(c)
		},
		required: []string{"file_path"},
	},
	{
		name: "DownloadDocument",
		newTool: func(c *client.Client) Tool {
			return NewDownloadDocument(c)
		},
		required: []string{"id", "save_path"},
	},
}

func TestAllTools_Description(t *testing.T) {
	for _, tt := range allToolTests {
		t.Run(tt.name, func(t *testing.T) {
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
	}
}

func TestAllTools_InputSchema(t *testing.T) {
	for _, tt := range allToolTests {
		t.Run(tt.name, func(t *testing.T) {
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

			if tt.required == nil {
				return
			}

			required, ok := schema["required"].([]string)
			if !ok {
				t.Fatal(
					"Schema should have required field",
				)
			}

			for _, want := range tt.required {
				found := false
				for _, r := range required {
					if r == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf(
						"%s should be in required fields",
						want,
					)
				}
			}
		})
	}
}

func TestAllTools_ServerError(t *testing.T) {
	for _, tt := range allToolTests {
		if tt.serverArgs == "" {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
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
				json.RawMessage(tt.serverArgs),
			)
			if err == nil {
				t.Fatal(
					"Expected error for server error",
				)
			}

			if !strings.Contains(
				err.Error(),
				"500",
			) {
				t.Errorf(
					"Error should mention 500, got: %s",
					err.Error(),
				)
			}
		})
	}
}

func TestAllTools_InvalidAndNegativeID(t *testing.T) {
	for _, tt := range allToolTests {
		if tt.idArgsFmt == "" {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Run("InvalidID", func(t *testing.T) {
				c := client.New(
					"http://localhost",
					"test-token",
				)
				tool := tt.newTool(c)

				_, err := tool.Execute(
					context.Background(),
					json.RawMessage(
						fmt.Sprintf(
							tt.idArgsFmt,
							0,
						),
					),
				)
				if err == nil {
					t.Fatal(
						"Expected error for invalid id",
					)
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
					json.RawMessage(
						fmt.Sprintf(
							tt.idArgsFmt,
							-1,
						),
					),
				)
				if err == nil {
					t.Fatal(
						"Expected error for negative id",
					)
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
		})
	}
}
