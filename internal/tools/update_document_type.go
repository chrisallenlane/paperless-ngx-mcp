package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// UpdateDocumentType updates an existing document type in Paperless-NGX.
type UpdateDocumentType struct {
	client *client.Client
}

// NewUpdateDocumentType creates a new UpdateDocumentType tool instance.
func NewUpdateDocumentType(
	c *client.Client,
) *UpdateDocumentType {
	return &UpdateDocumentType{client: c}
}

// Description returns a description of what this tool does.
func (t *UpdateDocumentType) Description() string {
	return "Update a document type in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UpdateDocumentType) InputSchema() map[string]interface{} {
	return matchableResourceSchema("Document type", true)
}

// Execute runs the tool and returns a formatted document type summary.
func (t *UpdateDocumentType) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	dt, err := patchByID[models.DocumentType](
		ctx,
		t.client,
		args,
		"/api/document_types/%d/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update document type: %w",
			err,
		)
	}

	return formatDocumentType(dt), nil
}
