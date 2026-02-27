package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// DeleteCorrespondent deletes a correspondent from Paperless-NGX.
type DeleteCorrespondent struct {
	client *client.Client
}

// NewDeleteCorrespondent creates a new DeleteCorrespondent tool instance.
func NewDeleteCorrespondent(
	c *client.Client,
) *DeleteCorrespondent {
	return &DeleteCorrespondent{client: c}
}

// Description returns a description of what this tool does.
func (t *DeleteCorrespondent) Description() string {
	return "Delete a correspondent from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *DeleteCorrespondent) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Correspondent ID to delete",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteCorrespondent) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.ID <= 0 {
		return "", fmt.Errorf("id must be a positive integer")
	}

	path := fmt.Sprintf("/api/correspondents/%d/", params.ID)

	if err := doDeleteRequest(ctx, t.client, path); err != nil {
		return "", fmt.Errorf(
			"failed to delete correspondent: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"Correspondent %d deleted successfully.",
		params.ID,
	), nil
}
