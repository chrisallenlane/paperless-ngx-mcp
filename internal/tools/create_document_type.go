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
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Document type name",
			},
			"match": map[string]interface{}{
				"type":        "string",
				"description": "Match pattern for auto-assignment",
			},
			"matching_algorithm": map[string]interface{}{
				"type": "integer",
				"description": "Matching algorithm: " +
					"0=None, 1=Any word, 2=All words, " +
					"3=Exact match, 4=Regex, " +
					"5=Fuzzy word, 6=Automatic",
			},
			"is_insensitive": map[string]interface{}{
				"type":        "boolean",
				"description": "Case-insensitive matching",
			},
		},
		"required": []string{"name"},
	}
}

// Execute runs the tool and returns a formatted document type summary.
func (t *CreateDocumentType) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Name              string `json:"name"`
		Match             string `json:"match,omitempty"`
		MatchingAlgorithm *int   `json:"matching_algorithm,omitempty"`
		IsInsensitive     *bool  `json:"is_insensitive,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		"/api/document_types/",
		params,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create document type: %w",
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
