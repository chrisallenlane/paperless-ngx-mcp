package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

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
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number (default 1)",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Results per page (default 25)",
			},
			"name": map[string]interface{}{
				"type": "string",
				"description": "Filter by name " +
					"(case-insensitive contains)",
			},
		},
	}
}

// Execute runs the tool and returns a formatted custom field list.
func (t *ListCustomFields) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Page     *int    `json:"page"`
		PageSize *int    `json:"page_size"`
		Name     *string `json:"name"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	q := url.Values{}
	if params.Page != nil {
		q.Set("page", fmt.Sprintf("%d", *params.Page))
	}
	if params.PageSize != nil {
		q.Set("page_size", fmt.Sprintf("%d", *params.PageSize))
	}
	if params.Name != nil {
		q.Set("name__icontains", *params.Name)
	}

	path := "/api/custom_fields/"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list custom fields: %w",
			err,
		)
	}

	var list models.PaginatedCustomFieldList
	if err := json.Unmarshal(body, &list); err != nil {
		return "", fmt.Errorf(
			"failed to parse custom fields response: %w",
			err,
		)
	}

	return formatCustomFieldList(&list), nil
}

func formatCustomFieldList(
	list *models.PaginatedCustomFieldList,
) string {
	if list.Count == 0 {
		return "No custom fields found."
	}

	out := fmt.Sprintf("Custom Fields: %d total\n\n", list.Count)
	for _, f := range list.Results {
		out += fmt.Sprintf(
			"%d. %s (ID: %d) — type: %s, %d documents\n",
			f.ID,
			f.Name,
			f.ID,
			f.DataType,
			f.DocumentCount,
		)
	}

	if list.Next != nil {
		out += "\n(more results available — use page parameter)"
	}

	return out
}
