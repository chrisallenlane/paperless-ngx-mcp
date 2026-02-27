package tools

import (
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// --- No-arg get tools ---

// NewGetStatus creates a tool to get the system status.
func NewGetStatus(c *client.Client) Tool {
	return &noArgGetTool[models.SystemStatus]{
		client: c,
		desc: "Get the current system status " +
			"of the Paperless-NGX server",
		path:   "/api/status/",
		format: formatStatus,
	}
}

// NewGetNextASN creates a tool to get the next available ASN.
func NewGetNextASN(c *client.Client) Tool {
	return &noArgGetTool[int]{
		client: c,
		desc: "Get the next available " +
			"archive serial number (ASN)",
		path: "/api/documents/next_asn/",
		format: func(v *int) string {
			return fmt.Sprintf(
				"Next available ASN: %d",
				*v,
			)
		},
	}
}

// --- Get tools ---

// NewGetCorrespondent creates a tool to get a correspondent by ID.
func NewGetCorrespondent(c *client.Client) Tool {
	return &getTool[models.Correspondent]{
		client:  c,
		desc:    "Get a correspondent by ID from Paperless-NGX",
		schema:  idOnlySchema("Correspondent ID"),
		pathFmt: "/api/correspondents/%d/",
		format:  formatCorrespondent,
	}
}

// NewGetCustomField creates a tool to get a custom field by ID.
func NewGetCustomField(c *client.Client) Tool {
	return &getTool[models.CustomField]{
		client:  c,
		desc:    "Get a custom field by ID from Paperless-NGX",
		schema:  idOnlySchema("Custom field ID"),
		pathFmt: "/api/custom_fields/%d/",
		format:  formatCustomField,
	}
}

// NewGetDocumentType creates a tool to get a document type by ID.
func NewGetDocumentType(c *client.Client) Tool {
	return &getTool[models.DocumentType]{
		client:  c,
		desc:    "Get a document type by ID from Paperless-NGX",
		schema:  idOnlySchema("Document type ID"),
		pathFmt: "/api/document_types/%d/",
		format:  formatDocumentType,
	}
}

// NewGetDocument creates a tool to get a document by ID.
func NewGetDocument(c *client.Client) Tool {
	return &getTool[models.Document]{
		client:  c,
		desc:    "Get a document by ID from Paperless-NGX",
		schema:  idOnlySchema("Document ID"),
		pathFmt: "/api/documents/%d/",
		format:  formatDocument,
	}
}

// NewGetDocumentMetadata creates a tool to get document metadata.
func NewGetDocumentMetadata(c *client.Client) Tool {
	return &getTool[models.DocumentMetadata]{
		client: c,
		desc: "Get file metadata for a document by ID, " +
			"including checksums, sizes, and OCR language",
		schema:  idOnlySchema("Document ID"),
		pathFmt: "/api/documents/%d/metadata/",
		format:  formatDocumentMetadata,
	}
}

// NewGetDocumentSuggestions creates a tool to get document suggestions.
func NewGetDocumentSuggestions(c *client.Client) Tool {
	return &getTool[models.DocumentSuggestions]{
		client: c,
		desc: "Get AI-generated suggestions for a document, " +
			"including correspondent, type, tags, and dates",
		schema:  idOnlySchema("Document ID"),
		pathFmt: "/api/documents/%d/suggestions/",
		format:  formatDocumentSuggestions,
	}
}

// --- List tools ---

// NewListCorrespondents creates a tool to list correspondents.
func NewListCorrespondents(c *client.Client) Tool {
	return &listTool[models.Correspondent]{
		client: c,
		desc: "List correspondents in Paperless-NGX " +
			"with optional filtering by name",
		schema:   paginatedListSchema(),
		basePath: "/api/correspondents/",
		format:   formatCorrespondentList,
	}
}

// NewListCustomFields creates a tool to list custom fields.
func NewListCustomFields(c *client.Client) Tool {
	return &listTool[models.CustomField]{
		client: c,
		desc: "List custom fields in Paperless-NGX " +
			"with optional filtering by name",
		schema:   paginatedListSchema(),
		basePath: "/api/custom_fields/",
		format:   formatCustomFieldList,
	}
}

// NewListDocumentTypes creates a tool to list document types.
func NewListDocumentTypes(c *client.Client) Tool {
	return &listTool[models.DocumentType]{
		client: c,
		desc: "List document types in Paperless-NGX " +
			"with optional filtering by name",
		schema:   paginatedListSchema(),
		basePath: "/api/document_types/",
		format:   formatDocumentTypeList,
	}
}

// NewGetTask creates a tool to get a background task by ID.
func NewGetTask(c *client.Client) Tool {
	return &getTool[models.Task]{
		client:  c,
		desc:    "Get a background task by ID from Paperless-NGX",
		schema:  idOnlySchema("Task ID"),
		pathFmt: "/api/tasks/%d/",
		format:  formatTask,
	}
}

// NewListTags creates a tool to list tags.
func NewListTags(c *client.Client) Tool {
	return &listTool[models.Tag]{
		client: c,
		desc: "List tags in Paperless-NGX " +
			"with optional filtering by name",
		schema:   paginatedListSchema(),
		basePath: "/api/tags/",
		format:   formatTagList,
	}
}

// NewGetTag creates a tool to get a tag by ID.
func NewGetTag(c *client.Client) Tool {
	return &getTool[models.Tag]{
		client:  c,
		desc:    "Get a tag by ID from Paperless-NGX",
		schema:  idOnlySchema("Tag ID"),
		pathFmt: "/api/tags/%d/",
		format:  formatTag,
	}
}

// NewUpdateTag creates a tool to update a tag.
func NewUpdateTag(c *client.Client) Tool {
	return &patchTool[models.Tag]{
		client:  c,
		desc:    "Update a tag in Paperless-NGX",
		schema:  tagSchema(true),
		pathFmt: "/api/tags/%d/",
		format:  formatTag,
	}
}

// NewDeleteTag creates a tool to delete a tag.
func NewDeleteTag(c *client.Client) Tool {
	return &deleteTool{
		client:       c,
		desc:         "Delete a tag from Paperless-NGX",
		schema:       idOnlySchema("Tag ID to delete"),
		pathFmt:      "/api/tags/%d/",
		resourceName: "Tag",
	}
}

// NewListStoragePaths creates a tool to list storage paths.
func NewListStoragePaths(c *client.Client) Tool {
	return &listTool[models.StoragePath]{
		client: c,
		desc: "List storage paths in Paperless-NGX " +
			"with optional filtering by name",
		schema:   paginatedListSchema(),
		basePath: "/api/storage_paths/",
		format:   formatStoragePathList,
	}
}

// NewGetStoragePath creates a tool to get a storage path by ID.
func NewGetStoragePath(c *client.Client) Tool {
	return &getTool[models.StoragePath]{
		client: c,
		desc: "Get a storage path by ID " +
			"from Paperless-NGX",
		schema:  idOnlySchema("Storage path ID"),
		pathFmt: "/api/storage_paths/%d/",
		format:  formatStoragePath,
	}
}

// NewUpdateStoragePath creates a tool to update a storage path.
func NewUpdateStoragePath(c *client.Client) Tool {
	return &patchTool[models.StoragePath]{
		client: c,
		desc: "Update a storage path " +
			"in Paperless-NGX",
		schema:  storagePathSchema(true),
		pathFmt: "/api/storage_paths/%d/",
		format:  formatStoragePath,
	}
}

// NewDeleteStoragePath creates a tool to delete a storage path.
func NewDeleteStoragePath(c *client.Client) Tool {
	return &deleteTool{
		client: c,
		desc: "Delete a storage path " +
			"from Paperless-NGX",
		schema: idOnlySchema(
			"Storage path ID to delete",
		),
		pathFmt:      "/api/storage_paths/%d/",
		resourceName: "Storage path",
	}
}

// NewListSavedViews creates a tool to list saved views.
func NewListSavedViews(c *client.Client) Tool {
	return &listTool[models.SavedView]{
		client: c,
		desc: "List saved views in Paperless-NGX " +
			"with optional pagination",
		schema:   paginationOnlySchema(),
		basePath: "/api/saved_views/",
		format:   formatSavedViewList,
	}
}

// NewGetSavedView creates a tool to get a saved view by ID.
func NewGetSavedView(c *client.Client) Tool {
	return &getTool[models.SavedView]{
		client: c,
		desc: "Get a saved view by ID " +
			"from Paperless-NGX",
		schema:  idOnlySchema("Saved view ID"),
		pathFmt: "/api/saved_views/%d/",
		format:  formatSavedView,
	}
}

// NewListTrash creates a tool to list soft-deleted documents.
func NewListTrash(c *client.Client) Tool {
	return &listTool[models.Document]{
		client: c,
		desc: "List soft-deleted documents in the " +
			"Paperless-NGX trash",
		schema:   paginationOnlySchema(),
		basePath: "/api/trash/",
		format:   formatDocumentList,
	}
}

// --- Create tools ---

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
		schema:   savedViewCreateSchema(),
		path:     "/api/saved_views/",
		validate: validateCreateSavedView,
		format:   formatSavedView,
	}
}

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

// --- Update tools ---

// NewUpdateCorrespondent creates a tool to update a correspondent.
func NewUpdateCorrespondent(c *client.Client) Tool {
	return &patchTool[models.Correspondent]{
		client:  c,
		desc:    "Update a correspondent in Paperless-NGX",
		schema:  matchableResourceSchema("Correspondent", true),
		pathFmt: "/api/correspondents/%d/",
		format:  formatCorrespondent,
	}
}

// NewUpdateCustomField creates a tool to update a custom field.
func NewUpdateCustomField(c *client.Client) Tool {
	return &patchTool[models.CustomField]{
		client:  c,
		desc:    "Update a custom field in Paperless-NGX",
		schema:  customFieldSchema(true),
		pathFmt: "/api/custom_fields/%d/",
		format:  formatCustomField,
	}
}

// NewUpdateDocumentType creates a tool to update a document type.
func NewUpdateDocumentType(c *client.Client) Tool {
	return &patchTool[models.DocumentType]{
		client:  c,
		desc:    "Update a document type in Paperless-NGX",
		schema:  matchableResourceSchema("Document type", true),
		pathFmt: "/api/document_types/%d/",
		format:  formatDocumentType,
	}
}

// NewUpdateSavedView creates a tool to update a saved view.
func NewUpdateSavedView(c *client.Client) Tool {
	return &patchTool[models.SavedView]{
		client:  c,
		desc:    "Update a saved view in Paperless-NGX",
		schema:  savedViewUpdateSchema(),
		pathFmt: "/api/saved_views/%d/",
		format:  formatSavedView,
	}
}

// NewUpdateDocument creates a tool to update a document.
func NewUpdateDocument(c *client.Client) Tool {
	return &patchTool[models.Document]{
		client:  c,
		desc:    "Update a document in Paperless-NGX",
		schema:  documentUpdateSchema(),
		pathFmt: "/api/documents/%d/",
		format:  formatDocument,
	}
}

// NewUpdateConfig creates a tool to update the application
// configuration.
func NewUpdateConfig(c *client.Client) Tool {
	return &patchTool[models.ApplicationConfiguration]{
		client: c,
		desc: "Update the application configuration " +
			"of the Paperless-NGX server",
		schema:  configUpdateSchema(),
		pathFmt: "/api/config/%d/",
		format:  formatConfig,
	}
}

// --- Delete tools ---

// NewDeleteCorrespondent creates a tool to delete a correspondent.
func NewDeleteCorrespondent(c *client.Client) Tool {
	return &deleteTool{
		client:       c,
		desc:         "Delete a correspondent from Paperless-NGX",
		schema:       idOnlySchema("Correspondent ID to delete"),
		pathFmt:      "/api/correspondents/%d/",
		resourceName: "Correspondent",
	}
}

// NewDeleteCustomField creates a tool to delete a custom field.
func NewDeleteCustomField(c *client.Client) Tool {
	return &deleteTool{
		client:       c,
		desc:         "Delete a custom field from Paperless-NGX",
		schema:       idOnlySchema("Custom field ID to delete"),
		pathFmt:      "/api/custom_fields/%d/",
		resourceName: "Custom field",
	}
}

// NewDeleteDocumentType creates a tool to delete a document type.
func NewDeleteDocumentType(c *client.Client) Tool {
	return &deleteTool{
		client: c,
		desc:   "Delete a document type from Paperless-NGX",
		schema: idOnlySchema(
			"Document type ID to delete",
		),
		pathFmt:      "/api/document_types/%d/",
		resourceName: "Document type",
	}
}

// NewDeleteSavedView creates a tool to delete a saved view.
func NewDeleteSavedView(c *client.Client) Tool {
	return &deleteTool{
		client: c,
		desc: "Delete a saved view " +
			"from Paperless-NGX",
		schema: idOnlySchema(
			"Saved view ID to delete",
		),
		pathFmt:      "/api/saved_views/%d/",
		resourceName: "Saved view",
	}
}

// NewDeleteDocument creates a tool to delete a document.
func NewDeleteDocument(c *client.Client) Tool {
	return &deleteTool{
		client:       c,
		desc:         "Delete a document from Paperless-NGX",
		schema:       idOnlySchema("Document ID to delete"),
		pathFmt:      "/api/documents/%d/",
		resourceName: "Document",
	}
}
