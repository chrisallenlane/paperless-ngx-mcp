package tools

import (
	"context"
	"encoding/json"

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
	return idOnlySchema("Correspondent ID to delete")
}

// Execute runs the tool and returns a confirmation message.
func (t *DeleteCorrespondent) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	return deleteByID(
		ctx,
		t.client,
		args,
		"/api/correspondents/%d/",
		"Correspondent",
	)
}
