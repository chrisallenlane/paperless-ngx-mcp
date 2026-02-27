package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// GetNextASN retrieves the next available archive serial number.
type GetNextASN struct {
	client *client.Client
}

// NewGetNextASN creates a new GetNextASN tool instance.
func NewGetNextASN(c *client.Client) *GetNextASN {
	return &GetNextASN{client: c}
}

// Description returns a description of what this tool does.
func (t *GetNextASN) Description() string {
	return "Get the next available archive serial number (ASN)"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetNextASN) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// Execute runs the tool and returns the next ASN.
func (t *GetNextASN) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(
		ctx,
		t.client,
		"/api/documents/next_asn/",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get next ASN: %w",
			err,
		)
	}

	var asn int
	if err := json.Unmarshal(body, &asn); err != nil {
		return "", fmt.Errorf(
			"failed to parse ASN response: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"Next available ASN: %d",
		asn,
	), nil
}
