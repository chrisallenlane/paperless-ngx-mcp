package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetConfig retrieves the application configuration from Paperless-NGX.
type GetConfig struct {
	client *client.Client
}

// NewGetConfig creates a new GetConfig tool instance.
func NewGetConfig(c *client.Client) *GetConfig {
	return &GetConfig{client: c}
}

// Description returns a description of what this tool does.
func (t *GetConfig) Description() string {
	return "Get the current application configuration of the Paperless-NGX server"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetConfig) InputSchema() map[string]interface{} {
	return emptySchema()
}

// Execute runs the tool and returns a formatted configuration summary.
func (t *GetConfig) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(ctx, t.client, "/api/config/")
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	var configs []models.ApplicationConfiguration
	if err := json.Unmarshal(body, &configs); err != nil {
		return "", fmt.Errorf(
			"failed to parse config response: %w",
			err,
		)
	}

	if len(configs) == 0 {
		return "No configuration found.", nil
	}

	return formatConfig(&configs[0]), nil
}
