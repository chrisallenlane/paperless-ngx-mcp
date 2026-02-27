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
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Correspondent ID",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted correspondent summary.
func (t *GetCorrespondent) Execute(
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

var matchingAlgorithmNames = map[int]string{
	0: "None",
	1: "Any word",
	2: "All words",
	3: "Exact match",
	4: "Regex",
	5: "Fuzzy word",
	6: "Automatic",
}

func formatCorrespondent(c *models.Correspondent) string {
	algoName := matchingAlgorithmNames[c.MatchingAlgorithm]
	if algoName == "" {
		algoName = "Unknown"
	}

	matchDisplay := c.Match
	if matchDisplay == "" {
		matchDisplay = "(none)"
	}

	lastCorr := "(none)"
	if c.LastCorrespondence != nil {
		lastCorr = formatDate(*c.LastCorrespondence)
	}

	out := fmt.Sprintf("Correspondent (ID: %d)\n", c.ID)
	out += fmt.Sprintf("  Name: %s\n", c.Name)
	out += fmt.Sprintf("  Slug: %s\n", c.Slug)
	out += fmt.Sprintf("  Match: %s\n", matchDisplay)
	out += fmt.Sprintf(
		"  Matching Algorithm: %d (%s)\n",
		c.MatchingAlgorithm,
		algoName,
	)
	out += fmt.Sprintf("  Case Insensitive: %v\n", c.IsInsensitive)
	out += fmt.Sprintf("  Document Count: %d\n", c.DocumentCount)
	out += fmt.Sprintf("  Last Correspondence: %s\n", lastCorr)

	return out
}
