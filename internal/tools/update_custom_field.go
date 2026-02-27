package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// UpdateCustomField updates an existing custom field in Paperless-NGX.
type UpdateCustomField struct {
	client *client.Client
}

// NewUpdateCustomField creates a new UpdateCustomField tool instance.
func NewUpdateCustomField(c *client.Client) *UpdateCustomField {
	return &UpdateCustomField{client: c}
}

// Description returns a description of what this tool does.
func (t *UpdateCustomField) Description() string {
	return "Update a custom field in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UpdateCustomField) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Custom field ID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Custom field name",
			},
			"data_type": map[string]interface{}{
				"type": "string",
				"description": "Data type: string, url, " +
					"date, boolean, integer, float, " +
					"monetary, documentlink, " +
					"select, longtext",
			},
			"extra_data": map[string]interface{}{
				"type": "object",
				"description": "Additional field " +
					"configuration (JSON object)",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted custom field summary.
func (t *UpdateCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, patchBody, err := parsePatchArgs(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/custom_fields/%d/", id)

	body, err := doPatchRequest(ctx, t.client, path, patchBody)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update custom field: %w",
			err,
		)
	}

	var field models.CustomField
	if err := json.Unmarshal(body, &field); err != nil {
		return "", fmt.Errorf(
			"failed to parse custom field response: %w",
			err,
		)
	}

	return formatCustomField(&field), nil
}
