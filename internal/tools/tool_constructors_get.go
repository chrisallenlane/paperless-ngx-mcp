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

// NewGetConfig creates a tool to get the application
// configuration.
func NewGetConfig(c *client.Client) Tool {
	return &noArgGetToolRaw{
		client: c,
		desc: "Get the current application " +
			"configuration of the " +
			"Paperless-NGX server",
		path: "/api/config/",
		process: func(body []byte) (string, error) {
			var configs []models.ApplicationConfiguration
			if err := json.Unmarshal(
				body,
				&configs,
			); err != nil {
				return "", fmt.Errorf(
					"failed to parse response: %w",
					err,
				)
			}
			if len(configs) == 0 {
				return "No configuration found.", nil
			}
			return formatConfig(&configs[0]), nil
		},
	}
}

// NewGetStatistics creates a tool to get document and resource
// count statistics.
func NewGetStatistics(c *client.Client) Tool {
	return &noArgGetToolRaw{
		client: c,
		desc: "Get document and resource count " +
			"statistics from Paperless-NGX",
		path: "/api/statistics/",
		process: func(body []byte) (string, error) {
			var stats map[string]interface{}
			if err := json.Unmarshal(
				body,
				&stats,
			); err != nil {
				return "", fmt.Errorf(
					"failed to parse response: %w",
					err,
				)
			}
			return formatStatistics(stats), nil
		},
	}
}

// --- Get-by-ID tools ---

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
