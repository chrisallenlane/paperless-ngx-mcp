package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatDocument(d *models.Document) string {
	out := fmt.Sprintf("Document (ID: %d)\n", d.ID)
	out += fmt.Sprintf("  Title: %s\n", d.Title)
	out += fmt.Sprintf(
		"  Correspondent: %s\n",
		formatOptInt(d.Correspondent),
	)
	out += fmt.Sprintf(
		"  Document Type: %s\n",
		formatOptInt(d.DocumentType),
	)
	out += fmt.Sprintf(
		"  Storage Path: %s\n",
		formatOptInt(d.StoragePath),
	)

	out += fmt.Sprintf("  Tags: %s\n", formatIntSlice(d.Tags))

	out += fmt.Sprintf("  Created: %s\n", formatDate(d.Created))
	out += fmt.Sprintf("  Added: %s\n", formatDate(d.Added))
	out += fmt.Sprintf("  Modified: %s\n", formatDate(d.Modified))
	out += fmt.Sprintf(
		"  ASN: %s\n",
		formatOptInt(d.ArchiveSerialNumber),
	)
	out += fmt.Sprintf(
		"  Original File: %s\n",
		formatOptStr(d.OriginalFileName),
	)
	out += fmt.Sprintf(
		"  Archived File: %s\n",
		formatOptStr(d.ArchivedFileName),
	)
	out += fmt.Sprintf("  MIME Type: %s\n", d.MimeType)
	out += fmt.Sprintf(
		"  Page Count: %s\n",
		formatOptInt(d.PageCount),
	)

	if len(d.CustomFields) > 0 {
		out += "  Custom Fields:\n"
		for _, cf := range d.CustomFields {
			out += fmt.Sprintf(
				"    Field %d: %s\n",
				cf.Field,
				string(cf.Value),
			)
		}
	}

	contentPreview := d.Content
	if len(contentPreview) > 500 {
		contentPreview = contentPreview[:500] + "..."
	}
	if contentPreview != "" {
		out += fmt.Sprintf("  Content: %s\n", contentPreview)
	}

	return out
}

func formatDocumentList(
	list *models.PaginatedList[models.Document],
) string {
	return formatPaginatedList(
		list,
		"No documents found.",
		"Documents",
		func(d models.Document) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d)\n"+
					"   Correspondent: %s | Type: %s"+
					" | ASN: %s | Created: %s\n",
				d.ID,
				d.Title,
				d.ID,
				formatOptInt(d.Correspondent),
				formatOptInt(d.DocumentType),
				formatOptInt(d.ArchiveSerialNumber),
				formatDate(d.Created),
			)
		},
	)
}

func formatDocumentMetadata(
	m *models.DocumentMetadata,
) string {
	out := "Document Metadata\n"

	out += "\nOriginal File:\n"
	out += fmt.Sprintf("  Filename: %s\n", m.OriginalFilename)
	out += fmt.Sprintf("  MIME Type: %s\n", m.OriginalMimeType)
	out += fmt.Sprintf(
		"  Size: %s\n",
		formatFileSize(int64(m.OriginalSize)),
	)
	out += fmt.Sprintf("  Checksum: %s\n", m.OriginalChecksum)

	out += "\nArchive File:\n"
	out += fmt.Sprintf(
		"  Has Archive Version: %v\n",
		m.HasArchiveVersion,
	)
	if m.HasArchiveVersion {
		out += fmt.Sprintf(
			"  Filename: %s\n",
			m.ArchiveMediaFilename,
		)
		out += fmt.Sprintf(
			"  Size: %s\n",
			formatFileSize(int64(m.ArchiveSize)),
		)
		out += fmt.Sprintf(
			"  Checksum: %s\n",
			m.ArchiveChecksum,
		)
	}

	out += fmt.Sprintf(
		"\nMedia Filename: %s\n",
		m.MediaFilename,
	)
	out += fmt.Sprintf("OCR Language: %s\n", m.Lang)

	return out
}

func formatDocumentSuggestions(
	s *models.DocumentSuggestions,
) string {
	out := "Document Suggestions\n"

	out += fmt.Sprintf(
		"  Correspondents: %s\n",
		formatIntSlice(s.Correspondents),
	)
	out += fmt.Sprintf(
		"  Document Types: %s\n",
		formatIntSlice(s.DocumentTypes),
	)
	out += fmt.Sprintf(
		"  Storage Paths: %s\n",
		formatIntSlice(s.StoragePaths),
	)
	out += fmt.Sprintf(
		"  Tags: %s\n",
		formatIntSlice(s.Tags),
	)
	out += fmt.Sprintf(
		"  Dates: %s\n",
		formatStringSlice(s.Dates),
	)

	return out
}
