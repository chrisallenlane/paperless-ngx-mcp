package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateDocumentType creates a new document type in Paperless-NGX.
type CreateDocumentType struct {
	client *client.Client
}

// NewCreateDocumentType creates a new CreateDocumentType tool instance.
func NewCreateDocumentType(
	c *client.Client,
) *CreateDocumentType {
	return &CreateDocumentType{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateDocumentType) Description() string {
	return "Create a new document type in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *CreateDocumentType) InputSchema() map[string]interface{} {
	return matchableResourceSchema("Document type", false)
}

// Execute runs the tool and returns a formatted document type summary.
func (t *CreateDocumentType) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	dt, err := createMatchable[models.DocumentType](
		ctx,
		t.client,
		args,
		"/api/document_types/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create document type: %w",
			err,
		)
	}

	return formatDocumentType(dt), nil
}
