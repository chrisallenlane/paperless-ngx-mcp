package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListDocumentTypes lists document types from Paperless-NGX.
type ListDocumentTypes struct {
	client *client.Client
}

// NewListDocumentTypes creates a new ListDocumentTypes tool instance.
func NewListDocumentTypes(c *client.Client) *ListDocumentTypes {
	return &ListDocumentTypes{client: c}
}

// Description returns a description of what this tool does.
func (t *ListDocumentTypes) Description() string {
	return "List document types in Paperless-NGX " +
		"with optional filtering by name"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *ListDocumentTypes) InputSchema() map[string]interface{} {
	return paginatedListSchema()
}

// Execute runs the tool and returns a formatted document type list.
func (t *ListDocumentTypes) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	path, err := buildListPath("/api/document_types/", args)
	if err != nil {
		return "", err
	}

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list document types: %w",
			err,
		)
	}

	var list models.PaginatedList[models.DocumentType]
	if err := json.Unmarshal(body, &list); err != nil {
		return "", fmt.Errorf(
			"failed to parse document types response: %w",
			err,
		)
	}

	return formatDocumentTypeList(&list), nil
}
