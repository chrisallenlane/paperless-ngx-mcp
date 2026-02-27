package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// UpdateDocument updates an existing document in Paperless-NGX.
type UpdateDocument struct {
	client *client.Client
}

// NewUpdateDocument creates a new UpdateDocument tool instance.
func NewUpdateDocument(c *client.Client) *UpdateDocument {
	return &UpdateDocument{client: c}
}

// Description returns a description of what this tool does.
func (t *UpdateDocument) Description() string {
	return "Update a document in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UpdateDocument) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID to update",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Document title (max 128 chars)",
				"maxLength":   128,
			},
			"correspondent": map[string]interface{}{
				"type": "integer",
				"description": "Correspondent ID " +
					"(null to clear)",
			},
			"document_type": map[string]interface{}{
				"type": "integer",
				"description": "Document type ID " +
					"(null to clear)",
			},
			"storage_path": map[string]interface{}{
				"type": "integer",
				"description": "Storage path ID " +
					"(null to clear)",
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "integer",
				},
				"description": "Replace all tags " +
					"with these tag IDs",
			},
			"archive_serial_number": map[string]interface{}{
				"type":        "integer",
				"description": "Archive serial number (null to clear)",
				"minimum":     0,
			},
			"created": map[string]interface{}{
				"type":        "string",
				"description": "Creation date (ISO format)",
			},
			"custom_fields": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"field": map[string]interface{}{
							"type":        "integer",
							"description": "Custom field ID",
						},
						"value": map[string]interface{}{
							"description": "Field value " +
								"(type varies by field)",
						},
					},
					"required": []string{
						"field",
						"value",
					},
				},
				"description": "Custom field values to set",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted document summary.
func (t *UpdateDocument) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, patchBody, err := parsePatchArgs(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/documents/%d/", id)

	body, err := doPatchRequest(ctx, t.client, path, patchBody)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update document: %w",
			err,
		)
	}

	var doc models.Document
	if err := json.Unmarshal(body, &doc); err != nil {
		return "", fmt.Errorf(
			"failed to parse document response: %w",
			err,
		)
	}

	return formatDocument(&doc), nil
}
