package tools

import (
	"context"
	"encoding/json"

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
	return idOnlySchema("Custom field ID to delete")
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	return deleteByID(
		ctx,
		t.client,
		args,
		"/api/custom_fields/%d/",
		"Custom field",
	)
}
