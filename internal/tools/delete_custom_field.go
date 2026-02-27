package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// DeleteCustomField deletes a custom field from Paperless-NGX.
type DeleteCustomField struct {
	client *client.Client
}

// NewDeleteCustomField creates a new DeleteCustomField tool instance.
func NewDeleteCustomField(c *client.Client) *DeleteCustomField {
	return &DeleteCustomField{client: c}
}

// Description returns a description of what this tool does.
func (t *DeleteCustomField) Description() string {
	return "Delete a custom field from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *DeleteCustomField) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Custom field ID to delete",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.ID <= 0 {
		return "", fmt.Errorf("id must be a positive integer")
	}

	path := fmt.Sprintf("/api/custom_fields/%d/", params.ID)

	if err := doDeleteRequest(ctx, t.client, path); err != nil {
		return "", fmt.Errorf(
			"failed to delete custom field: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"Custom field %d deleted successfully.",
		params.ID,
	), nil
}
