package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetStatus retrieves the system status from Paperless-NGX.
type GetStatus struct {
	client *client.Client
}

// NewGetStatus creates a new GetStatus tool instance.
func NewGetStatus(c *client.Client) *GetStatus {
	return &GetStatus{client: c}
}

// Description returns a description of what this tool does.
func (t *GetStatus) Description() string {
	return "Get the current system status of the Paperless-NGX server"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetStatus) InputSchema() map[string]interface{} {
	return emptySchema()
}

// Execute runs the tool and returns a formatted status summary.
func (t *GetStatus) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(ctx, t.client, "/api/status/")
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	var status models.SystemStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return "", fmt.Errorf(
			"failed to parse status response: %w",
			err,
		)
	}

	return formatStatus(&status), nil
}
