package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListDocumentNotes lists notes for a document.
type ListDocumentNotes struct {
	client *client.Client
}

// NewListDocumentNotes creates a new ListDocumentNotes tool
// instance.
func NewListDocumentNotes(
	c *client.Client,
) *ListDocumentNotes {
	return &ListDocumentNotes{client: c}
}

// Description returns a description of what this tool does.
func (t *ListDocumentNotes) Description() string {
	return "List notes attached to a document " +
		"in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *ListDocumentNotes) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID",
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number (default 1)",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Results per page (default 25)",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted notes list.
func (t *ListDocumentNotes) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		ID       int  `json:"id"`
		Page     *int `json:"page"`
		PageSize *int `json:"page_size"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.ID <= 0 {
		return "", fmt.Errorf(
			"id must be a positive integer",
		)
	}

	path := fmt.Sprintf(
		"/api/documents/%d/notes/",
		params.ID,
	)

	q := url.Values{}
	addPaginationQuery(q, params.Page, params.PageSize)
	path = appendQuery(path, q)

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list document notes: %w",
			err,
		)
	}

	var list models.PaginatedList[models.Note]
	if err := json.Unmarshal(body, &list); err != nil {
		return "", fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return formatNoteList(params.ID, &list), nil
}
