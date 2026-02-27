package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// UpdateDocumentType updates an existing document type in Paperless-NGX.
type UpdateDocumentType struct {
	client *client.Client
}

// NewUpdateDocumentType creates a new UpdateDocumentType tool instance.
func NewUpdateDocumentType(
	c *client.Client,
) *UpdateDocumentType {
	return &UpdateDocumentType{client: c}
}

// Description returns a description of what this tool does.
func (t *UpdateDocumentType) Description() string {
	return "Update a document type in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UpdateDocumentType) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document type ID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Document type name",
			},
			"match": map[string]interface{}{
				"type":        "string",
				"description": "Match pattern for auto-assignment",
			},
			"matching_algorithm": map[string]interface{}{
				"type": "integer",
				"description": "Matching algorithm: " +
					"0=None, 1=Any word, 2=All words, " +
					"3=Exact match, 4=Regex, " +
					"5=Fuzzy word, 6=Automatic",
			},
			"is_insensitive": map[string]interface{}{
				"type":        "boolean",
				"description": "Case-insensitive matching",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted document type summary.
func (t *UpdateDocumentType) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, patchBody, err := parsePatchArgs(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/document_types/%d/", id)

	body, err := doPatchRequest(ctx, t.client, path, patchBody)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update document type: %w",
			err,
		)
	}

	var dt models.DocumentType
	if err := json.Unmarshal(body, &dt); err != nil {
		return "", fmt.Errorf(
			"failed to parse document type response: %w",
			err,
		)
	}

	return formatDocumentType(&dt), nil
}
