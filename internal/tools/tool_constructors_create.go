package tools

import (
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// NewCreateCorrespondent creates a tool to create a correspondent.
func NewCreateCorrespondent(c *client.Client) Tool {
	return &createTool[models.Correspondent]{
		client:   c,
		desc:     "Create a new correspondent in Paperless-NGX",
		schema:   matchableResourceSchema("Correspondent", false),
		path:     "/api/correspondents/",
		validate: validateMatchableCreate,
		format:   formatCorrespondent,
	}
}

// NewCreateDocumentType creates a tool to create a document type.
func NewCreateDocumentType(c *client.Client) Tool {
	return &createTool[models.DocumentType]{
		client:   c,
		desc:     "Create a new document type in Paperless-NGX",
		schema:   matchableResourceSchema("Document type", false),
		path:     "/api/document_types/",
		validate: validateMatchableCreate,
		format:   formatDocumentType,
	}
}

// NewCreateTag creates a tool to create a tag.
func NewCreateTag(c *client.Client) Tool {
	return &createTool[models.Tag]{
		client:   c,
		desc:     "Create a new tag in Paperless-NGX",
		schema:   tagSchema(false),
		path:     "/api/tags/",
		validate: validateCreateTag,
		format:   formatTag,
	}
}

// NewCreateStoragePath creates a tool to create a storage path.
func NewCreateStoragePath(c *client.Client) Tool {
	return &createTool[models.StoragePath]{
		client:   c,
		desc:     "Create a new storage path in Paperless-NGX",
		schema:   storagePathSchema(false),
		path:     "/api/storage_paths/",
		validate: validateCreateStoragePath,
		format:   formatStoragePath,
	}
}

// NewCreateCustomField creates a tool to create a custom field.
func NewCreateCustomField(c *client.Client) Tool {
	return &createTool[models.CustomField]{
		client:   c,
		desc:     "Create a new custom field in Paperless-NGX",
		schema:   customFieldSchema(false),
		path:     "/api/custom_fields/",
		validate: validateCreateCustomField,
		format:   formatCustomField,
	}
}

// NewCreateSavedView creates a tool to create a saved view.
func NewCreateSavedView(c *client.Client) Tool {
	return &createTool[models.SavedView]{
		client:   c,
		desc:     "Create a new saved view in Paperless-NGX",
		schema:   savedViewSchema(false),
		path:     "/api/saved_views/",
		validate: validateCreateSavedView,
		format:   formatSavedView,
	}
}

// --- Validate functions ---

// validateMatchableCreate validates arguments for creating a matchable
// resource (correspondent, document type).
func validateMatchableCreate(
	args json.RawMessage,
) (interface{}, error) {
	var params struct {
		Name              string `json:"name"`
		Match             string `json:"match,omitempty"`
		MatchingAlgorithm *int   `json:"matching_algorithm,omitempty"`
		IsInsensitive     *bool  `json:"is_insensitive,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return params, nil
}

// validateCreateTag validates arguments for creating a tag.
func validateCreateTag(
	args json.RawMessage,
) (interface{}, error) {
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
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return params, nil
}

// validateCreateStoragePath validates arguments for creating a storage
// path.
func validateCreateStoragePath(
	args json.RawMessage,
) (interface{}, error) {
	var params struct {
		Name              string `json:"name"`
		Path              string `json:"path"`
		Match             string `json:"match,omitempty"`
		MatchingAlgorithm *int   `json:"matching_algorithm,omitempty"`
		IsInsensitive     *bool  `json:"is_insensitive,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if params.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	return params, nil
}

// validateCreateCustomField validates arguments for creating a custom
// field.
func validateCreateCustomField(
	args json.RawMessage,
) (interface{}, error) {
	var params struct {
		Name     string          `json:"name"`
		DataType string          `json:"data_type"`
		Extra    json.RawMessage `json:"extra_data,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if params.DataType == "" {
		return nil, fmt.Errorf("data_type is required")
	}
	reqBody := map[string]interface{}{
		"name":      params.Name,
		"data_type": params.DataType,
	}
	if params.Extra != nil && string(params.Extra) != "null" {
		reqBody["extra_data"] = params.Extra
	}
	return reqBody, nil
}

// validateCreateSavedView validates arguments for creating a saved view.
func validateCreateSavedView(
	args json.RawMessage,
) (interface{}, error) {
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
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if params.ShowOnDashboard == nil {
		return nil, fmt.Errorf(
			"show_on_dashboard is required",
		)
	}
	if params.ShowInSidebar == nil {
		return nil, fmt.Errorf(
			"show_in_sidebar is required",
		)
	}
	if params.FilterRules == nil {
		return nil, fmt.Errorf(
			"filter_rules is required",
		)
	}
	return params, nil
}
