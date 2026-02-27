package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateSavedView creates a new saved view in Paperless-NGX.
type CreateSavedView struct {
	client *client.Client
}

// NewCreateSavedView creates a new CreateSavedView tool
// instance.
func NewCreateSavedView(
	c *client.Client,
) *CreateSavedView {
	return &CreateSavedView{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateSavedView) Description() string {
	return "Create a new saved view in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *CreateSavedView) InputSchema() map[string]interface{} {
	return savedViewCreateSchema()
}

// Execute runs the tool and returns a formatted saved view.
func (t *CreateSavedView) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Name            string `json:"name"`
		ShowOnDashboard *bool  `json:"show_on_dashboard"`
		ShowInSidebar   *bool  `json:"show_in_sidebar"`
		SortField       string `json:"sort_field,omitempty"`
		SortReverse     *bool  `json:"sort_reverse,omitempty"`
		FilterRules     []struct {
			RuleType int     `json:"rule_type"`
			Value    *string `json:"value"`
		} `json:"filter_rules"`
		PageSize    *int    `json:"page_size,omitempty"`
		DisplayMode *string `json:"display_mode,omitempty"`
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

	if params.ShowOnDashboard == nil {
		return "", fmt.Errorf(
			"show_on_dashboard is required",
		)
	}

	if params.ShowInSidebar == nil {
		return "", fmt.Errorf(
			"show_in_sidebar is required",
		)
	}

	if params.FilterRules == nil {
		return "", fmt.Errorf(
			"filter_rules is required",
		)
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		"/api/saved_views/",
		params,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create saved view: %w",
			err,
		)
	}

	var view models.SavedView
	if err := json.Unmarshal(body, &view); err != nil {
		return "", fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return formatSavedView(&view), nil
}
