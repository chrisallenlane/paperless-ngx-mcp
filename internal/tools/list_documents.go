package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListDocuments lists documents from Paperless-NGX.
type ListDocuments struct {
	client *client.Client
}

// NewListDocuments creates a new ListDocuments tool instance.
func NewListDocuments(c *client.Client) *ListDocuments {
	return &ListDocuments{client: c}
}

// Description returns a description of what this tool does.
func (t *ListDocuments) Description() string {
	return "List documents in Paperless-NGX with optional " +
		"filtering and full-text search"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *ListDocuments) InputSchema() map[string]interface{} {
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
			"search": map[string]interface{}{
				"type": "string",
				"description": "Full-text search across " +
					"title and content",
			},
			"correspondent": map[string]interface{}{
				"type": "integer",
				"description": "Filter by " +
					"correspondent ID",
			},
			"document_type": map[string]interface{}{
				"type": "integer",
				"description": "Filter by " +
					"document type ID",
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "integer",
				},
				"description": "Filter by tag IDs " +
					"(documents with ALL specified tags)",
			},
			"is_in_inbox": map[string]interface{}{
				"type":        "boolean",
				"description": "Filter inbox documents",
			},
		},
	}
}

// Execute runs the tool and returns a formatted document list.
func (t *ListDocuments) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	path, err := buildDocumentListPath(args)
	if err != nil {
		return "", err
	}

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list documents: %w",
			err,
		)
	}

	var list models.PaginatedList[models.Document]
	if err := json.Unmarshal(body, &list); err != nil {
		return "", fmt.Errorf(
			"failed to parse documents response: %w",
			err,
		)
	}

	return formatDocumentList(&list), nil
}

// buildDocumentListPath constructs the API path with query parameters
// for document listing.
func buildDocumentListPath(args json.RawMessage) (string, error) {
	var params struct {
		Page          *int   `json:"page"`
		PageSize      *int   `json:"page_size"`
		Search        string `json:"search"`
		Correspondent *int   `json:"correspondent"`
		DocumentType  *int   `json:"document_type"`
		Tags          []int  `json:"tags"`
		IsInInbox     *bool  `json:"is_in_inbox"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	q := url.Values{}
	if params.Page != nil {
		q.Set("page", fmt.Sprintf("%d", *params.Page))
	}
	if params.PageSize != nil {
		q.Set(
			"page_size",
			fmt.Sprintf("%d", *params.PageSize),
		)
	}
	if params.Search != "" {
		q.Set("search", params.Search)
	}
	if params.Correspondent != nil {
		q.Set(
			"correspondent__id",
			fmt.Sprintf("%d", *params.Correspondent),
		)
	}
	if params.DocumentType != nil {
		q.Set(
			"document_type__id",
			fmt.Sprintf("%d", *params.DocumentType),
		)
	}
	for _, tagID := range params.Tags {
		q.Add(
			"tags__id__all",
			fmt.Sprintf("%d", tagID),
		)
	}
	if params.IsInInbox != nil {
		q.Set(
			"is_in_inbox",
			fmt.Sprintf("%t", *params.IsInInbox),
		)
	}

	return appendQuery("/api/documents/", q), nil
}
