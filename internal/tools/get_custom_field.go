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
	return idOnlySchema("Custom field ID")
}

// Execute runs the tool and returns a formatted custom field summary.
func (t *GetCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	field, _, err := fetchByID[models.CustomField](
		ctx,
		t.client,
		args,
		"/api/custom_fields/%d/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get custom field: %w",
			err,
		)
	}

	return formatCustomField(field), nil
}
