package tools

import (
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
		format: func(_ int, v *models.Correspondent) string {
			return formatCorrespondent(v)
		},
	}
}

// NewGetCustomField creates a tool to get a custom field by ID.
func NewGetCustomField(c *client.Client) Tool {
	return &getTool[models.CustomField]{
		client:  c,
		desc:    "Get a custom field by ID from Paperless-NGX",
		schema:  idOnlySchema("Custom field ID"),
		pathFmt: "/api/custom_fields/%d/",
		format: func(_ int, v *models.CustomField) string {
			return formatCustomField(v)
		},
	}
}

// NewGetDocumentType creates a tool to get a document type by ID.
func NewGetDocumentType(c *client.Client) Tool {
	return &getTool[models.DocumentType]{
		client:  c,
		desc:    "Get a document type by ID from Paperless-NGX",
		schema:  idOnlySchema("Document type ID"),
		pathFmt: "/api/document_types/%d/",
		format: func(_ int, v *models.DocumentType) string {
			return formatDocumentType(v)
		},
	}
}

// NewGetDocument creates a tool to get a document by ID.
func NewGetDocument(c *client.Client) Tool {
	return &getTool[models.Document]{
		client:  c,
		desc:    "Get a document by ID from Paperless-NGX",
		schema:  idOnlySchema("Document ID"),
		pathFmt: "/api/documents/%d/",
		format: func(_ int, v *models.Document) string {
			return formatDocument(v)
		},
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
		format: func(_ int, v *models.Task) string {
			return formatTask(v)
		},
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
