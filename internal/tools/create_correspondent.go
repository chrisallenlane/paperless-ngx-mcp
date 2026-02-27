package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateCorrespondent creates a new correspondent in Paperless-NGX.
type CreateCorrespondent struct {
	client *client.Client
}

// NewCreateCorrespondent creates a new CreateCorrespondent tool instance.
func NewCreateCorrespondent(
	c *client.Client,
) *CreateCorrespondent {
	return &CreateCorrespondent{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateCorrespondent) Description() string {
	return "Create a new correspondent in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *CreateCorrespondent) InputSchema() map[string]interface{} {
	return matchableResourceSchema("Correspondent", false)
}

// Execute runs the tool and returns a formatted correspondent summary.
func (t *CreateCorrespondent) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params matchableCreateParams
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		"/api/correspondents/",
		params,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create correspondent: %w",
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
