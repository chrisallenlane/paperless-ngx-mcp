package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetDocument retrieves a single document from Paperless-NGX.
type GetDocument struct {
	client *client.Client
}

// NewGetDocument creates a new GetDocument tool instance.
func NewGetDocument(c *client.Client) *GetDocument {
	return &GetDocument{client: c}
}

// Description returns a description of what this tool does.
func (t *GetDocument) Description() string {
	return "Get a document by ID from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetDocument) InputSchema() map[string]interface{} {
	return idOnlySchema("Document ID")
}

// Execute runs the tool and returns a formatted document summary.
func (t *GetDocument) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	doc, _, err := fetchByID[models.Document](
		ctx,
		t.client,
		args,
		"/api/documents/%d/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get document: %w",
			err,
		)
	}

	return formatDocument(doc), nil
}
