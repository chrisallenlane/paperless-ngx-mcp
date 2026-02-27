package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// UpdateCorrespondent updates an existing correspondent in Paperless-NGX.
type UpdateCorrespondent struct {
	client *client.Client
}

// NewUpdateCorrespondent creates a new UpdateCorrespondent tool instance.
func NewUpdateCorrespondent(
	c *client.Client,
) *UpdateCorrespondent {
	return &UpdateCorrespondent{client: c}
}

// Description returns a description of what this tool does.
func (t *UpdateCorrespondent) Description() string {
	return "Update a correspondent in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UpdateCorrespondent) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Correspondent ID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Correspondent name",
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
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted correspondent summary.
func (t *UpdateCorrespondent) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(args, &raw); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	idRaw, ok := raw["id"]
	if !ok {
		return "", fmt.Errorf("id is required")
	}

	var id int64
	if err := json.Unmarshal(idRaw, &id); err != nil {
		return "", fmt.Errorf("failed to parse id: %w", err)
	}

	if id <= 0 {
		return "", fmt.Errorf("id must be a positive integer")
	}

	patchBody := make(map[string]json.RawMessage)
	for k, v := range raw {
		if k != "id" {
			patchBody[k] = v
		}
	}

	path := fmt.Sprintf("/api/correspondents/%d/", id)

	body, err := doPatchRequest(ctx, t.client, path, patchBody)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update correspondent: %w",
			err,
		)
	}

	var corr models.Correspondent
	if err := json.Unmarshal(body, &corr); err != nil {
		return "", fmt.Errorf(
			"failed to parse correspondent response: %w",
			err,
		)
	}

	return formatCorrespondent(&corr), nil
}
