package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListCustomFields lists custom fields from Paperless-NGX.
type ListCustomFields struct {
	client *client.Client
}

// NewListCustomFields creates a new ListCustomFields tool instance.
func NewListCustomFields(c *client.Client) *ListCustomFields {
	return &ListCustomFields{client: c}
}

// Description returns a description of what this tool does.
func (t *ListCustomFields) Description() string {
	return "List custom fields in Paperless-NGX " +
		"with optional filtering by name"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *ListCustomFields) InputSchema() map[string]interface{} {
	return paginatedListSchema()
}

// Execute runs the tool and returns a formatted custom field list.
func (t *ListCustomFields) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	list, err := listResources[models.CustomField](
		ctx,
		t.client,
		"/api/custom_fields/",
		args,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list custom fields: %w",
			err,
		)
	}

	return formatCustomFieldList(list), nil
}
