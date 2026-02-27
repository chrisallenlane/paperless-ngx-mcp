package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetCorrespondent retrieves a single correspondent from Paperless-NGX.
type GetCorrespondent struct {
	client *client.Client
}

// NewGetCorrespondent creates a new GetCorrespondent tool instance.
func NewGetCorrespondent(c *client.Client) *GetCorrespondent {
	return &GetCorrespondent{client: c}
}

// Description returns a description of what this tool does.
func (t *GetCorrespondent) Description() string {
	return "Get a correspondent by ID from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetCorrespondent) InputSchema() map[string]interface{} {
	return idOnlySchema("Correspondent ID")
}

// Execute runs the tool and returns a formatted correspondent summary.
func (t *GetCorrespondent) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/correspondents/%d/", id)

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get correspondent: %w",
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
