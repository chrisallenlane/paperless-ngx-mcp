package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// DeleteDocument deletes a document from Paperless-NGX.
type DeleteDocument struct {
	client *client.Client
}

// NewDeleteDocument creates a new DeleteDocument tool instance.
func NewDeleteDocument(c *client.Client) *DeleteDocument {
	return &DeleteDocument{client: c}
}

// Description returns a description of what this tool does.
func (t *DeleteDocument) Description() string {
	return "Delete a document from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *DeleteDocument) InputSchema() map[string]interface{} {
	return idOnlySchema("Document ID to delete")
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteDocument) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/documents/%d/", id)

	if err := doDeleteRequest(ctx, t.client, path); err != nil {
		return "", fmt.Errorf(
			"failed to delete document: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"Document %d deleted successfully.",
		id,
	), nil
}
