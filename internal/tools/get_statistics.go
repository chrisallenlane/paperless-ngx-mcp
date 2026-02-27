package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// GetStatistics retrieves document and resource count statistics.
type GetStatistics struct {
	client *client.Client
}

// NewGetStatistics creates a new GetStatistics tool instance.
func NewGetStatistics(c *client.Client) *GetStatistics {
	return &GetStatistics{client: c}
}

// Description returns a description of what this tool does.
func (t *GetStatistics) Description() string {
	return "Get document and resource count statistics " +
		"from Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *GetStatistics) InputSchema() map[string]interface{} {
	return emptySchema()
}

// Execute runs the tool and returns formatted statistics.
func (t *GetStatistics) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(
		ctx,
		t.client,
		"/api/statistics/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get statistics: %w",
			err,
		)
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(body, &stats); err != nil {
		return "", fmt.Errorf(
			"failed to parse statistics response: %w",
			err,
		)
	}

	return formatStatistics(stats), nil
}
