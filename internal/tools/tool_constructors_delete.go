package tools

import (
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

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
