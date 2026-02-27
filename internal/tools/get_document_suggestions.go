package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetDocumentSuggestions retrieves AI suggestions for a document.
type GetDocumentSuggestions struct {
	client *client.Client
}

// NewGetDocumentSuggestions creates a new GetDocumentSuggestions tool instance.
func NewGetDocumentSuggestions(
	c *client.Client,
) *GetDocumentSuggestions {
	return &GetDocumentSuggestions{client: c}
}

// Description returns a description of what this tool does.
func (t *GetDocumentSuggestions) Description() string {
	return "Get AI-generated suggestions for a document, " +
		"including correspondent, type, tags, and dates"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetDocumentSuggestions) InputSchema() map[string]interface{} {
	return idOnlySchema("Document ID")
}

// Execute runs the tool and returns formatted suggestions.
func (t *GetDocumentSuggestions) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf(
		"/api/documents/%d/suggestions/",
		id,
	)

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get document suggestions: %w",
			err,
		)
	}

	var sugg models.DocumentSuggestions
	if err := json.Unmarshal(body, &sugg); err != nil {
		return "", fmt.Errorf(
			"failed to parse suggestions response: %w",
			err,
		)
	}

	return formatDocumentSuggestions(id, &sugg), nil
}
