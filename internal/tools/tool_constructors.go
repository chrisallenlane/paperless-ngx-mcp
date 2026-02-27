package tools

import (
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

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
	return &getToolWithID[models.DocumentMetadata]{
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
	return &getToolWithID[models.DocumentSuggestions]{
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

// --- Create tools ---

// NewCreateCorrespondent creates a tool to create a correspondent.
func NewCreateCorrespondent(c *client.Client) Tool {
	return &createMatchableTool[models.Correspondent]{
		client: c,
		desc:   "Create a new correspondent in Paperless-NGX",
		schema: matchableResourceSchema("Correspondent", false),
		path:   "/api/correspondents/",
		format: formatCorrespondent,
	}
}

// NewCreateDocumentType creates a tool to create a document type.
func NewCreateDocumentType(c *client.Client) Tool {
	return &createMatchableTool[models.DocumentType]{
		client: c,
		desc:   "Create a new document type in Paperless-NGX",
		schema: matchableResourceSchema("Document type", false),
		path:   "/api/document_types/",
		format: formatDocumentType,
	}
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
