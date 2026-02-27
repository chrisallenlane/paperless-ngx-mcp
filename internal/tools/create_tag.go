package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateTag creates a new tag in Paperless-NGX.
type CreateTag struct {
	client *client.Client
}

// NewCreateTag creates a new CreateTag tool instance.
func NewCreateTag(c *client.Client) *CreateTag {
	return &CreateTag{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateTag) Description() string {
	return "Create a new tag in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *CreateTag) InputSchema() map[string]interface{} {
	return tagSchema(false)
}

// Execute runs the tool and returns a formatted tag summary.
func (t *CreateTag) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Name              string `json:"name"`
		Color             string `json:"color,omitempty"`
		Match             string `json:"match,omitempty"`
		MatchingAlgorithm *int   `json:"matching_algorithm,omitempty"`
		IsInsensitive     *bool  `json:"is_insensitive,omitempty"`
		IsInboxTag        *bool  `json:"is_inbox_tag,omitempty"`
		Parent            *int   `json:"parent,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		"/api/tags/",
		params,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create tag: %w",
			err,
		)
	}

	var tag models.Tag
	if err := json.Unmarshal(body, &tag); err != nil {
		return "", fmt.Errorf(
			"failed to parse tag response: %w",
			err,
		)
	}

	return formatTag(&tag), nil
}
