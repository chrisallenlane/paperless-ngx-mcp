package tools

import (
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

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
