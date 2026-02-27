package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateCustomField creates a new custom field in Paperless-NGX.
type CreateCustomField struct {
	client *client.Client
}

// NewCreateCustomField creates a new CreateCustomField tool instance.
func NewCreateCustomField(c *client.Client) *CreateCustomField {
	return &CreateCustomField{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateCustomField) Description() string {
	return "Create a new custom field in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *CreateCustomField) InputSchema() map[string]interface{} {
	return customFieldSchema(false)
}

// Execute runs the tool and returns a formatted custom field summary.
func (t *CreateCustomField) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Name     string          `json:"name"`
		DataType string          `json:"data_type"`
		Extra    json.RawMessage `json:"extra_data,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if params.DataType == "" {
		return "", fmt.Errorf("data_type is required")
	}

	reqBody := map[string]interface{}{
		"name":      params.Name,
		"data_type": params.DataType,
	}
	if params.Extra != nil && string(params.Extra) != "null" {
		reqBody["extra_data"] = params.Extra
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		"/api/custom_fields/",
		reqBody,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create custom field: %w",
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
