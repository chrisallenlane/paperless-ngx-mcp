package tools

import (
	"context"
	"encoding/json"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// DeleteDocumentType deletes a document type from Paperless-NGX.
type DeleteDocumentType struct {
	client *client.Client
}

// NewDeleteDocumentType creates a new DeleteDocumentType tool instance.
func NewDeleteDocumentType(
	c *client.Client,
) *DeleteDocumentType {
	return &DeleteDocumentType{client: c}
}

// Description returns a description of what this tool does.
func (t *DeleteDocumentType) Description() string {
	return "Delete a document type from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *DeleteDocumentType) InputSchema() map[string]interface{} {
	return idOnlySchema("Document type ID to delete")
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteDocumentType) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	return deleteByID(
		ctx,
		t.client,
		args,
		"/api/document_types/%d/",
		"Document type",
	)
}
