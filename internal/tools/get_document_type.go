package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetDocumentType retrieves a single document type from Paperless-NGX.
type GetDocumentType struct {
	client *client.Client
}

// NewGetDocumentType creates a new GetDocumentType tool instance.
func NewGetDocumentType(c *client.Client) *GetDocumentType {
	return &GetDocumentType{client: c}
}

// Description returns a description of what this tool does.
func (t *GetDocumentType) Description() string {
	return "Get a document type by ID from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetDocumentType) InputSchema() map[string]interface{} {
	return idOnlySchema("Document type ID")
}

// Execute runs the tool and returns a formatted document type summary.
func (t *GetDocumentType) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	dt, _, err := fetchByID[models.DocumentType](
		ctx,
		t.client,
		args,
		"/api/document_types/%d/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get document type: %w",
			err,
		)
	}

	return formatDocumentType(dt), nil
}
