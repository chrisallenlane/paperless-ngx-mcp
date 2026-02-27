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
	return idOnlySchema("Custom field ID to delete")
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/custom_fields/%d/", id)

	if err := doDeleteRequest(ctx, t.client, path); err != nil {
		return "", fmt.Errorf(
			"failed to delete custom field: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"Custom field %d deleted successfully.",
		id,
	), nil
}
