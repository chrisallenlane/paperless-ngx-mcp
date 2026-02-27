package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

const bytesPerTB = 1024 * 1024 * 1024 * 1024

var matchingAlgorithmNames = map[int]string{
	0: "None",
	1: "Any word",
	2: "All words",
	3: "Exact match",
	4: "Regex",
	5: "Fuzzy word",
	6: "Automatic",
}

func formatStatus(s *models.SystemStatus) string {
	totalTB := float64(s.Storage.Total) / bytesPerTB
	availTB := float64(s.Storage.Available) / bytesPerTB

	out := fmt.Sprintf(
		"Paperless-NGX Status\nVersion: %s\nOS: %s\nInstall: %s\n\n"+
			"Storage: %.2f TB available of %.2f TB\n\n"+
			"Database: %s - %s\n"+
			"Redis: %s\n"+
			"Celery: %s\n",
		s.PNGXVersion,
		s.ServerOS,
		s.InstallType,
		availTB,
		totalTB,
		s.Database.Type,
		s.Database.Status,
		s.Tasks.RedisStatus,
		s.Tasks.CeleryStatus,
	)

	out += formatTaskLine(
		"Index",
		s.Tasks.IndexStatus,
		"last modified",
		s.Tasks.IndexLastModified,
	)
	out += formatTaskLine(
		"Classifier",
		s.Tasks.ClassifierStatus,
		"last trained",
		s.Tasks.ClassifierLastTrained,
	)
	out += formatTaskLine(
		"Sanity Check",
		s.Tasks.SanityCheckStatus,
		"last run",
		s.Tasks.SanityCheckLastRun,
	)

	return out
}

func formatTaskLine(
	name, status, dateLabel string,
	dateValue *string,
) string {
	if dateValue != nil && *dateValue != "" {
		date := formatDate(*dateValue)
		return fmt.Sprintf(
			"%s: %s (%s: %s)\n",
			name,
			status,
			dateLabel,
			date,
		)
	}
	return fmt.Sprintf("%s: %s\n", name, status)
}

func formatDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func formatConfig(c *models.ApplicationConfiguration) string {
	out := fmt.Sprintf(
		"Paperless-NGX Configuration (ID: %d)\n",
		c.ID,
	)

	out += "\nOCR Settings:\n"
	out += formatOpt("  Output Type", c.OutputType)
	out += formatOpt("  Pages", c.Pages)
	out += formatOpt("  Language", c.Language)
	out += formatOpt("  Mode", c.Mode)
	out += formatOpt("  Skip Archive File", c.SkipArchiveFile)
	out += formatOpt("  Image DPI", c.ImageDPI)
	out += formatOpt("  Unpaper Clean", c.UnpaperClean)
	out += formatOpt("  Deskew", c.Deskew)
	out += formatOpt("  Rotate Pages", c.RotatePages)
	out += formatOpt("  Rotate Pages Threshold", c.RotatePagesThreshold)
	out += formatOpt("  Max Image Pixels", c.MaxImagePixels)
	out += formatOpt(
		"  Color Conversion Strategy",
		c.ColorConversionStrategy,
	)
	out += formatOptJSON("  User Args", c.UserArgs)

	out += "\nApp Settings:\n"
	out += formatOpt("  Title", c.AppTitle)
	out += formatOpt("  Logo", c.AppLogo)

	out += "\nBarcode Settings:\n"
	out += formatOpt("  Enabled", c.BarcodesEnabled)
	out += formatOpt(
		"  TIFF Support",
		c.BarcodeEnableTiffSupport,
	)
	out += formatOpt("  String", c.BarcodeString)
	out += formatOpt(
		"  Retain Split Pages",
		c.BarcodeRetainSplitPages,
	)
	out += formatOpt("  Enable ASN", c.BarcodeEnableASN)
	out += formatOpt("  ASN Prefix", c.BarcodeASNPrefix)
	out += formatOpt("  Upscale", c.BarcodeUpscale)
	out += formatOpt("  DPI", c.BarcodeDPI)
	out += formatOpt("  Max Pages", c.BarcodeMaxPages)
	out += formatOpt("  Enable Tag", c.BarcodeEnableTag)
	out += formatOptJSON("  Tag Mapping", c.BarcodeTagMapping)

	return out
}

func formatOpt[T any](label string, v *T) string {
	if v != nil {
		return fmt.Sprintf("%s: %v\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptJSON(label string, v json.RawMessage) string {
	if v != nil && string(v) != "null" {
		return fmt.Sprintf("%s: %s\n", label, string(v))
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func matchingAlgorithmName(algo int) string {
	name := matchingAlgorithmNames[algo]
	if name == "" {
		return "Unknown"
	}
	return name
}

func matchDisplayOrDefault(match string) string {
	if match == "" {
		return "(none)"
	}
	return match
}

func formatCorrespondent(c *models.Correspondent) string {
	algoName := matchingAlgorithmName(c.MatchingAlgorithm)
	matchDisplay := matchDisplayOrDefault(c.Match)

	lastCorr := "(none)"
	if c.LastCorrespondence != nil {
		lastCorr = formatDate(*c.LastCorrespondence)
	}

	out := fmt.Sprintf("Correspondent (ID: %d)\n", c.ID)
	out += fmt.Sprintf("  Name: %s\n", c.Name)
	out += fmt.Sprintf("  Slug: %s\n", c.Slug)
	out += fmt.Sprintf("  Match: %s\n", matchDisplay)
	out += fmt.Sprintf(
		"  Matching Algorithm: %d (%s)\n",
		c.MatchingAlgorithm,
		algoName,
	)
	out += fmt.Sprintf("  Case Insensitive: %v\n", c.IsInsensitive)
	out += fmt.Sprintf("  Document Count: %d\n", c.DocumentCount)
	out += fmt.Sprintf("  Last Correspondence: %s\n", lastCorr)

	return out
}

func formatCorrespondentList(
	list *models.PaginatedList[models.Correspondent],
) string {
	if list.Count == 0 {
		return "No correspondents found."
	}

	out := fmt.Sprintf("Correspondents: %d total\n\n", list.Count)
	for _, c := range list.Results {
		out += fmt.Sprintf(
			"%d. %s (ID: %d) — %d documents\n",
			c.ID,
			c.Name,
			c.ID,
			c.DocumentCount,
		)
	}

	if list.Next != nil {
		out += "\n(more results available — use page parameter)"
	}

	return out
}

func formatCustomField(f *models.CustomField) string {
	extraData := "(none)"
	if f.ExtraData != nil && string(f.ExtraData) != "null" {
		extraData = string(f.ExtraData)
	}

	out := fmt.Sprintf("Custom Field (ID: %d)\n", f.ID)
	out += fmt.Sprintf("  Name: %s\n", f.Name)
	out += fmt.Sprintf("  Data Type: %s\n", f.DataType)
	out += fmt.Sprintf("  Extra Data: %s\n", extraData)
	out += fmt.Sprintf("  Document Count: %d\n", f.DocumentCount)

	return out
}

func formatCustomFieldList(
	list *models.PaginatedList[models.CustomField],
) string {
	if list.Count == 0 {
		return "No custom fields found."
	}

	out := fmt.Sprintf("Custom Fields: %d total\n\n", list.Count)
	for _, f := range list.Results {
		out += fmt.Sprintf(
			"%d. %s (ID: %d) — type: %s, %d documents\n",
			f.ID,
			f.Name,
			f.ID,
			f.DataType,
			f.DocumentCount,
		)
	}

	if list.Next != nil {
		out += "\n(more results available — use page parameter)"
	}

	return out
}

func formatDocumentType(dt *models.DocumentType) string {
	algoName := matchingAlgorithmName(dt.MatchingAlgorithm)
	matchDisplay := matchDisplayOrDefault(dt.Match)

	out := fmt.Sprintf("Document Type (ID: %d)\n", dt.ID)
	out += fmt.Sprintf("  Name: %s\n", dt.Name)
	out += fmt.Sprintf("  Slug: %s\n", dt.Slug)
	out += fmt.Sprintf("  Match: %s\n", matchDisplay)
	out += fmt.Sprintf(
		"  Matching Algorithm: %d (%s)\n",
		dt.MatchingAlgorithm,
		algoName,
	)
	out += fmt.Sprintf("  Case Insensitive: %v\n", dt.IsInsensitive)
	out += fmt.Sprintf("  Document Count: %d\n", dt.DocumentCount)

	return out
}

func formatDocumentTypeList(
	list *models.PaginatedList[models.DocumentType],
) string {
	if list.Count == 0 {
		return "No document types found."
	}

	out := fmt.Sprintf(
		"Document Types: %d total\n\n",
		list.Count,
	)
	for _, dt := range list.Results {
		out += fmt.Sprintf(
			"%d. %s (ID: %d) — %d documents\n",
			dt.ID,
			dt.Name,
			dt.ID,
			dt.DocumentCount,
		)
	}

	if list.Next != nil {
		out += "\n(more results available — use page parameter)"
	}

	return out
}

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
	if list.Count == 0 {
		return "No documents found."
	}

	out := fmt.Sprintf("Documents: %d total\n\n", list.Count)
	for _, d := range list.Results {
		corr := formatOptInt(d.Correspondent)
		docType := formatOptInt(d.DocumentType)
		asn := formatOptInt(d.ArchiveSerialNumber)

		out += fmt.Sprintf(
			"%d. %s (ID: %d)\n",
			d.ID,
			d.Title,
			d.ID,
		)
		out += fmt.Sprintf(
			"   Correspondent: %s | Type: %s | "+
				"ASN: %s | Created: %s\n",
			corr,
			docType,
			asn,
			formatDate(d.Created),
		)
	}

	if list.Next != nil {
		out += "\n(more results available — use page parameter)"
	}

	return out
}

func formatOptInt(v *int) string {
	if v != nil {
		return fmt.Sprintf("%d", *v)
	}
	return "(none)"
}

func formatOptStr(v *string) string {
	if v != nil && *v != "" {
		return *v
	}
	return "(none)"
}

func formatDocumentMetadata(
	id int,
	m *models.DocumentMetadata,
) string {
	out := fmt.Sprintf("Document Metadata (ID: %d)\n", id)

	out += "\nOriginal File:\n"
	out += fmt.Sprintf("  Filename: %s\n", m.OriginalFilename)
	out += fmt.Sprintf("  MIME Type: %s\n", m.OriginalMimeType)
	out += fmt.Sprintf("  Size: %s\n", formatFileSize(m.OriginalSize))
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
			formatFileSize(m.ArchiveSize),
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

func formatFileSize(bytes int) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func formatDocumentSuggestions(
	id int,
	s *models.DocumentSuggestions,
) string {
	out := fmt.Sprintf("Document Suggestions (ID: %d)\n", id)

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

func formatIntSlice(ids []int) string {
	if len(ids) == 0 {
		return "(none)"
	}
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(strs, ", ")
}

func formatStringSlice(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}
	return strings.Join(items, ", ")
}
