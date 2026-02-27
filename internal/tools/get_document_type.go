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
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/document_types/%d/", id)

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get document type: %w",
			err,
		)
	}

	var dt models.DocumentType
	if err := json.Unmarshal(body, &dt); err != nil {
		return "", fmt.Errorf(
			"failed to parse document type response: %w",
			err,
		)
	}

	return formatDocumentType(&dt), nil
}
