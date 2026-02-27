package tools

import (
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

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

// NewUpdateSavedView creates a tool to update a saved view.
func NewUpdateSavedView(c *client.Client) Tool {
	return &patchTool[models.SavedView]{
		client:  c,
		desc:    "Update a saved view in Paperless-NGX",
		schema:  savedViewSchema(true),
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
