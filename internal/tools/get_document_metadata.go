package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetDocumentMetadata retrieves metadata for a document from Paperless-NGX.
type GetDocumentMetadata struct {
	client *client.Client
}

// NewGetDocumentMetadata creates a new GetDocumentMetadata tool instance.
func NewGetDocumentMetadata(
	c *client.Client,
) *GetDocumentMetadata {
	return &GetDocumentMetadata{client: c}
}

// Description returns a description of what this tool does.
func (t *GetDocumentMetadata) Description() string {
	return "Get file metadata for a document by ID, " +
		"including checksums, sizes, and OCR language"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetDocumentMetadata) InputSchema() map[string]interface{} {
	return idOnlySchema("Document ID")
}

// Execute runs the tool and returns formatted document metadata.
func (t *GetDocumentMetadata) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/documents/%d/metadata/", id)

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get document metadata: %w",
			err,
		)
	}

	var meta models.DocumentMetadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return "", fmt.Errorf(
			"failed to parse metadata response: %w",
			err,
		)
	}

	return formatDocumentMetadata(id, &meta), nil
}
