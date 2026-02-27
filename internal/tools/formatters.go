package tools

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

const bytesPerTB = 1024 * 1024 * 1024 * 1024

const paginationHint = "\n(more results available — use page parameter)"

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

func formatMatchableFields(
	label string,
	id int,
	name, slug, match string,
	algo int,
	isInsensitive bool,
	docCount int,
) string {
	algoName := matchingAlgorithmName(algo)
	matchDisplay := match
	if matchDisplay == "" {
		matchDisplay = "(none)"
	}

	out := fmt.Sprintf("%s (ID: %d)\n", label, id)
	out += fmt.Sprintf("  Name: %s\n", name)
	out += fmt.Sprintf("  Slug: %s\n", slug)
	out += fmt.Sprintf("  Match: %s\n", matchDisplay)
	out += fmt.Sprintf(
		"  Matching Algorithm: %d (%s)\n",
		algo,
		algoName,
	)
	out += fmt.Sprintf("  Case Insensitive: %v\n", isInsensitive)
	out += fmt.Sprintf("  Document Count: %d\n", docCount)

	return out
}

func formatCorrespondent(c *models.Correspondent) string {
	out := formatMatchableFields(
		"Correspondent",
		c.ID,
		c.Name,
		c.Slug,
		c.Match,
		c.MatchingAlgorithm,
		c.IsInsensitive,
		c.DocumentCount,
	)

	lastCorr := "(none)"
	if c.LastCorrespondence != nil {
		lastCorr = formatDate(*c.LastCorrespondence)
	}
	out += fmt.Sprintf("  Last Correspondence: %s\n", lastCorr)

	return out
}

func formatPaginatedList[T any](
	list *models.PaginatedList[T],
	emptyMsg, header string,
	formatItem func(T) string,
) string {
	if list.Count == 0 {
		return emptyMsg
	}

	out := fmt.Sprintf("%s: %d total\n\n", header, list.Count)
	for _, item := range list.Results {
		out += formatItem(item)
	}

	if list.Next != nil {
		out += paginationHint
	}

	return out
}

func formatCorrespondentList(
	list *models.PaginatedList[models.Correspondent],
) string {
	return formatPaginatedList(
		list,
		"No correspondents found.",
		"Correspondents",
		func(c models.Correspondent) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents\n",
				c.ID,
				c.Name,
				c.ID,
				c.DocumentCount,
			)
		},
	)
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
	return formatPaginatedList(
		list,
		"No custom fields found.",
		"Custom Fields",
		func(f models.CustomField) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — type: %s, %d documents\n",
				f.ID,
				f.Name,
				f.ID,
				f.DataType,
				f.DocumentCount,
			)
		},
	)
}

func formatDocumentType(dt *models.DocumentType) string {
	return formatMatchableFields(
		"Document Type",
		dt.ID,
		dt.Name,
		dt.Slug,
		dt.Match,
		dt.MatchingAlgorithm,
		dt.IsInsensitive,
		dt.DocumentCount,
	)
}

func formatDocumentTypeList(
	list *models.PaginatedList[models.DocumentType],
) string {
	return formatPaginatedList(
		list,
		"No document types found.",
		"Document Types",
		func(dt models.DocumentType) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents\n",
				dt.ID,
				dt.Name,
				dt.ID,
				dt.DocumentCount,
			)
		},
	)
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

func formatTag(tag *models.Tag) string {
	out := formatMatchableFields(
		"Tag",
		tag.ID,
		tag.Name,
		tag.Slug,
		tag.Match,
		tag.MatchingAlgorithm,
		tag.IsInsensitive,
		tag.DocumentCount,
	)
	out += fmt.Sprintf("  Color: %s\n", tag.Color)
	out += fmt.Sprintf(
		"  Text Color: %s\n",
		tag.TextColor,
	)
	out += fmt.Sprintf(
		"  Is Inbox Tag: %v\n",
		tag.IsInboxTag,
	)
	out += fmt.Sprintf(
		"  Parent: %s\n",
		formatOptInt(tag.Parent),
	)
	out += fmt.Sprintf(
		"  Children: %s\n",
		formatIntSlice(tag.Children),
	)

	return out
}

func formatTagList(
	list *models.PaginatedList[models.Tag],
) string {
	return formatPaginatedList(
		list,
		"No tags found.",
		"Tags",
		func(tag models.Tag) string {
			extra := ""
			if tag.IsInboxTag {
				extra = " [inbox]"
			}
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents%s\n",
				tag.ID,
				tag.Name,
				tag.ID,
				tag.DocumentCount,
				extra,
			)
		},
	)
}

func formatStoragePath(sp *models.StoragePath) string {
	out := formatMatchableFields(
		"Storage Path",
		sp.ID,
		sp.Name,
		sp.Slug,
		sp.Match,
		sp.MatchingAlgorithm,
		sp.IsInsensitive,
		sp.DocumentCount,
	)
	out += fmt.Sprintf("  Path: %s\n", sp.Path)

	return out
}

func formatStoragePathList(
	list *models.PaginatedList[models.StoragePath],
) string {
	return formatPaginatedList(
		list,
		"No storage paths found.",
		"Storage Paths",
		func(sp models.StoragePath) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents\n",
				sp.ID,
				sp.Name,
				sp.ID,
				sp.DocumentCount,
			)
		},
	)
}

func formatTask(t *models.Task) string {
	out := fmt.Sprintf("Task (ID: %d)\n", t.ID)
	out += fmt.Sprintf("  Task UUID: %s\n", t.TaskID)
	out += fmt.Sprintf("  Status: %s\n", t.Status)
	out += fmt.Sprintf("  Type: %s\n", t.Type)
	out += fmt.Sprintf(
		"  Task Name: %s\n",
		formatOptStr(t.TaskName),
	)
	out += fmt.Sprintf(
		"  File Name: %s\n",
		formatOptStr(t.TaskFileName),
	)
	out += fmt.Sprintf(
		"  Created: %s\n",
		formatOptDate(t.DateCreated),
	)
	out += fmt.Sprintf(
		"  Done: %s\n",
		formatOptDate(t.DateDone),
	)
	out += fmt.Sprintf(
		"  Result: %s\n",
		formatOptStr(t.Result),
	)
	out += fmt.Sprintf(
		"  Acknowledged: %v\n",
		t.Acknowledged,
	)
	out += fmt.Sprintf(
		"  Related Document: %s\n",
		formatOptStr(t.RelatedDocument),
	)

	return out
}

func formatTaskArray(tasks []models.Task) string {
	if len(tasks) == 0 {
		return "No tasks found."
	}

	out := fmt.Sprintf("Tasks: %d total\n\n", len(tasks))
	for _, task := range tasks {
		name := formatOptStr(task.TaskName)
		out += fmt.Sprintf(
			"%d. %s — %s",
			task.ID,
			name,
			task.Status,
		)
		if task.TaskFileName != nil {
			out += fmt.Sprintf(
				" — %s",
				*task.TaskFileName,
			)
		}
		out += "\n"
	}

	return out
}

func formatOptDate(v *string) string {
	if v != nil && *v != "" {
		return formatDate(*v)
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

func formatSavedView(v *models.SavedView) string {
	out := fmt.Sprintf("Saved View (ID: %d)\n", v.ID)
	out += fmt.Sprintf("  Name: %s\n", v.Name)
	out += fmt.Sprintf(
		"  Show on Dashboard: %v\n",
		v.ShowOnDashboard,
	)
	out += fmt.Sprintf(
		"  Show in Sidebar: %v\n",
		v.ShowInSidebar,
	)
	out += fmt.Sprintf(
		"  Sort Field: %s\n",
		formatOptStr(v.SortField),
	)
	out += fmt.Sprintf(
		"  Sort Reverse: %v\n",
		v.SortReverse,
	)
	out += fmt.Sprintf(
		"  Page Size: %s\n",
		formatOptInt(v.PageSize),
	)
	out += fmt.Sprintf(
		"  Display Mode: %s\n",
		formatOptStr(v.DisplayMode),
	)

	if len(v.FilterRules) == 0 {
		out += "  Filter Rules: (none)\n"
	} else {
		out += "  Filter Rules:\n"
		for _, r := range v.FilterRules {
			val := "(null)"
			if r.Value != nil {
				val = *r.Value
			}
			out += fmt.Sprintf(
				"    - %s: %s\n",
				ruleTypeName(r.RuleType),
				val,
			)
		}
	}

	return out
}

func formatSavedViewList(
	list *models.PaginatedList[models.SavedView],
) string {
	return formatPaginatedList(
		list,
		"No saved views found.",
		"Saved Views",
		func(v models.SavedView) string {
			flags := ""
			if v.ShowOnDashboard {
				flags += " [dashboard]"
			}
			if v.ShowInSidebar {
				flags += " [sidebar]"
			}
			return fmt.Sprintf(
				"%d. %s (ID: %d) — "+
					"%d filter rules%s\n",
				v.ID,
				v.Name,
				v.ID,
				len(v.FilterRules),
				flags,
			)
		},
	)
}

// ruleTypeName returns a human-readable name for a filter
// rule type.
func ruleTypeName(ruleType int) string {
	names := map[int]string{
		0:  "Title contains",
		1:  "Content contains",
		2:  "ASN is",
		3:  "Correspondent is",
		4:  "Document type is",
		5:  "Is in inbox",
		6:  "Has tag",
		7:  "Has any tag",
		8:  "Created before",
		9:  "Created after",
		10: "Created year",
		11: "Created month",
		12: "Created day",
		13: "Added before",
		14: "Added after",
		15: "Modified before",
		16: "Modified after",
		17: "Does not have tag",
		18: "Does not have ASN",
		19: "Title or content contains",
		20: "Fulltext query",
		21: "More like document",
		22: "Has tags in",
		23: "ASN greater than",
		24: "ASN less than",
		25: "Storage path is",
		26: "Has correspondent in",
		27: "Does not have correspondent in",
		28: "Has document type in",
		29: "Does not have document type in",
		30: "Has storage path in",
		31: "Does not have storage path in",
		32: "Has tags all",
		33: "Owner is",
		34: "Owner is in",
		35: "Does not have owner in",
		36: "Correspondent starts with",
		37: "Correspondent ends with",
		38: "Title starts with",
		39: "Title ends with",
		44: "Has custom field value",
		45: "Custom field query",
		46: "Is shared by me",
		47: "Has custom fields in",
	}
	if name, ok := names[ruleType]; ok {
		return name
	}
	return fmt.Sprintf("Rule type %d", ruleType)
}

func formatNote(note *models.Note) string {
	out := fmt.Sprintf("Note (ID: %d)\n", note.ID)
	out += fmt.Sprintf(
		"  Author: %s\n",
		note.User.Username,
	)
	out += fmt.Sprintf(
		"  Created: %s\n",
		formatDate(note.Created),
	)
	out += fmt.Sprintf("  Content: %s\n", note.Note)

	return out
}

func formatNoteList(
	docID int,
	list *models.PaginatedList[models.Note],
) string {
	if list.Count == 0 {
		return fmt.Sprintf(
			"No notes found for document %d.",
			docID,
		)
	}

	out := fmt.Sprintf(
		"Notes for Document %d: %d total\n\n",
		docID,
		list.Count,
	)

	for _, note := range list.Results {
		out += fmt.Sprintf(
			"--- Note %d (by %s on %s) ---\n%s\n\n",
			note.ID,
			note.User.Username,
			formatDate(note.Created),
			note.Note,
		)
	}

	if list.Next != nil {
		out += "(More notes available — " +
			"use page parameter)\n"
	}

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

func formatStatistics(
	stats map[string]interface{},
) string {
	if len(stats) == 0 {
		return "No statistics available."
	}

	out := "Paperless-NGX Statistics\n\n"

	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		label := formatStatLabel(key)
		switch v := stats[key].(type) {
		case []interface{}:
			out += fmt.Sprintf(
				"  %s: %s\n",
				label,
				formatStatSlice(v),
			)
		default:
			out += fmt.Sprintf(
				"  %s: %s\n",
				label,
				formatStatValue(v),
			)
		}
	}

	return out
}

func formatStatLabel(key string) string {
	words := strings.Split(key, "_")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func formatStatValue(v interface{}) string {
	if f, ok := v.(float64); ok {
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		return fmt.Sprintf("%.2f", f)
	}
	return fmt.Sprintf("%v", v)
}

func formatStatSlice(items []interface{}) string {
	if len(items) == 0 {
		return "(none)"
	}

	parts := make([]string, len(items))
	for i, item := range items {
		switch v := item.(type) {
		case map[string]interface{}:
			pairs := make([]string, 0, len(v))
			for k, val := range v {
				pairs = append(
					pairs,
					fmt.Sprintf(
						"%s=%s",
						k,
						formatStatValue(val),
					),
				)
			}
			sort.Strings(pairs)
			parts[i] = strings.Join(pairs, ", ")
		default:
			parts[i] = formatStatValue(item)
		}
	}
	return strings.Join(parts, "; ")
}
