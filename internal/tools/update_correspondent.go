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
	return matchableResourceSchema("Correspondent", true)
}

// Execute runs the tool and returns a formatted correspondent summary.
func (t *UpdateCorrespondent) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, patchBody, err := parsePatchArgs(args)
	if err != nil {
		return "", err
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
