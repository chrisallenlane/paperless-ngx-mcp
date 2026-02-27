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
	return customFieldSchema(true)
}

// Execute runs the tool and returns a formatted custom field summary.
func (t *UpdateCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	field, err := patchByID[models.CustomField](
		ctx,
		t.client,
		args,
		"/api/custom_fields/%d/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update custom field: %w",
			err,
		)
	}

	return formatCustomField(field), nil
}
