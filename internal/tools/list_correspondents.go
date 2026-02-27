package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListCorrespondents lists correspondents from Paperless-NGX.
type ListCorrespondents struct {
	client *client.Client
}

// NewListCorrespondents creates a new ListCorrespondents tool instance.
func NewListCorrespondents(c *client.Client) *ListCorrespondents {
	return &ListCorrespondents{client: c}
}

// Description returns a description of what this tool does.
func (t *ListCorrespondents) Description() string {
	return "List correspondents in Paperless-NGX " +
		"with optional filtering by name"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *ListCorrespondents) InputSchema() map[string]interface{} {
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

// Execute runs the tool and returns a formatted correspondent list.
func (t *ListCorrespondents) Execute(
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

	path := "/api/correspondents/"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list correspondents: %w",
			err,
		)
	}

	var list models.PaginatedCorrespondentList
	if err := json.Unmarshal(body, &list); err != nil {
		return "", fmt.Errorf(
			"failed to parse correspondents response: %w",
			err,
		)
	}

	return formatCorrespondentList(&list), nil
}

func formatCorrespondentList(
	list *models.PaginatedCorrespondentList,
) string {
	if list.Count == 0 {
		return "No correspondents found."
	}

	out := fmt.Sprintf("Correspondents: %d total\n\n", list.Count)
	for _, c := range list.Results {
		out += fmt.Sprintf(
			"%d. %s (ID: %d) — %d documents\n",
			c.ID,
			c.Name,
			c.ID,
			c.DocumentCount,
		)
	}

	if list.Next != nil {
		out += "\n(more results available — use page parameter)"
	}

	return out
}
