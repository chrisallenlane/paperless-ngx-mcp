package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetCustomField retrieves a single custom field from Paperless-NGX.
type GetCustomField struct {
	client *client.Client
}

// NewGetCustomField creates a new GetCustomField tool instance.
func NewGetCustomField(c *client.Client) *GetCustomField {
	return &GetCustomField{client: c}
}

// Description returns a description of what this tool does.
func (t *GetCustomField) Description() string {
	return "Get a custom field by ID from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetCustomField) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Custom field ID",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted custom field summary.
func (t *GetCustomField) Execute(
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

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get custom field: %w",
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

func formatCustomField(f *models.CustomField) string {
	extraData := "(none)"
	if f.ExtraData != nil && string(f.ExtraData) != "null" {
		extraData = string(f.ExtraData)
	}

	out := fmt.Sprintf("Custom Field (ID: %d)\n", f.ID)
	out += fmt.Sprintf("  Name: %s\n", f.Name)
	out += fmt.Sprintf("  Data Type: %s\n", f.DataType)
	out += fmt.Sprintf("  Extra Data: %s\n", extraData)
	out += fmt.Sprintf("  Document Count: %d\n", f.DocumentCount)

	return out
}
