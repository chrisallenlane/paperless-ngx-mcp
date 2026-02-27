package tools

// emptySchema returns an input schema with no parameters.
func emptySchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// idOnlySchema returns an input schema with a single required "id" field.
func idOnlySchema(desc string) map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": desc,
			},
		},
		"required": []string{"id"},
	}
}

// taskListSchema returns an input schema for the task list tool.
func taskListSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"description": "Filter by status: " +
					"FAILURE, PENDING, RECEIVED, " +
					"RETRY, REVOKED, STARTED, SUCCESS",
			},
			"task_name": map[string]interface{}{
				"type": "string",
				"description": "Filter by task name: " +
					"consume_file, train_classifier, " +
					"check_sanity, index_optimize",
			},
			"type": map[string]interface{}{
				"type": "string",
				"description": "Filter by type: " +
					"auto_task, scheduled_task, " +
					"manual_task",
			},
			"task_id": map[string]interface{}{
				"type": "string",
				"description": "Filter by " +
					"Celery task UUID",
			},
		},
	}
}

// paginationOnlySchema returns an input schema with only page and
// page_size parameters (no name filter).
func paginationOnlySchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number (default 1)",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Results per page (default 25)",
			},
		},
	}
}

// paginatedListSchema returns an input schema for paginated list endpoints.
func paginatedListSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number (default 1)",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Results per page (default 25)",
			},
			"name": map[string]interface{}{
				"type": "string",
				"description": "Filter by name " +
					"(case-insensitive contains)",
			},
		},
	}
}

// matchableResourceSchema returns an input schema for resources with matching
// fields (name, match, matching_algorithm, is_insensitive). Set includeID
// to true for update tools, false for create tools.
func matchableResourceSchema(
	resourceName string,
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": resourceName + " name",
		},
		"match": map[string]interface{}{
			"type":        "string",
			"description": "Match pattern for auto-assignment",
		},
		"matching_algorithm": map[string]interface{}{
			"type": "integer",
			"description": "Matching algorithm: " +
				"0=None, 1=Any word, 2=All words, " +
				"3=Exact match, 4=Regex, " +
				"5=Fuzzy word, 6=Automatic",
		},
		"is_insensitive": map[string]interface{}{
			"type":        "boolean",
			"description": "Case-insensitive matching",
		},
	}

	required := []string{"name"}
	if includeID {
		props["id"] = map[string]interface{}{
			"type":        "integer",
			"description": resourceName + " ID to update",
		}
		required = []string{"id"}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

// customFieldSchema returns an input schema for custom field tools.
// Set includeID to true for update tools, false for create tools.
func customFieldSchema(
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Custom field name",
		},
		"data_type": map[string]interface{}{
			"type": "string",
			"description": "Data type: string, url, " +
				"date, boolean, integer, float, " +
				"monetary, documentlink, " +
				"select, longtext",
		},
		"extra_data": map[string]interface{}{
			"type": "object",
			"description": "Additional field " +
				"configuration (JSON object)",
		},
	}

	required := []string{"name", "data_type"}
	if includeID {
		props["id"] = map[string]interface{}{
			"type":        "integer",
			"description": "Custom field ID to update",
		}
		required = []string{"id"}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

// documentUpdateSchema returns an input schema for the document update tool.
func documentUpdateSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID to update",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Document title (max 128 chars)",
				"maxLength":   128,
			},
			"correspondent": map[string]interface{}{
				"type": "integer",
				"description": "Correspondent ID " +
					"(null to clear)",
			},
			"document_type": map[string]interface{}{
				"type": "integer",
				"description": "Document type ID " +
					"(null to clear)",
			},
			"storage_path": map[string]interface{}{
				"type": "integer",
				"description": "Storage path ID " +
					"(null to clear)",
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "integer",
				},
				"description": "Replace all tags " +
					"with these tag IDs",
			},
			"archive_serial_number": map[string]interface{}{
				"type":        "integer",
				"description": "Archive serial number (null to clear)",
				"minimum":     0,
			},
			"created": map[string]interface{}{
				"type":        "string",
				"description": "Creation date (ISO format)",
			},
			"custom_fields": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"field": map[string]interface{}{
							"type":        "integer",
							"description": "Custom field ID",
						},
						"value": map[string]interface{}{
							"description": "Field value " +
								"(type varies by field)",
						},
					},
					"required": []string{
						"field",
						"value",
					},
				},
				"description": "Custom field values to set",
			},
		},
		"required": []string{"id"},
	}
}

// configUpdateSchema returns an input schema for the config update tool.
func configUpdateSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Config ID to update (typically 1)",
			},
			"output_type": map[string]interface{}{
				"type": "string",
				"description": "PDF output type: pdf, pdfa, " +
					"pdfa-1, pdfa-2, pdfa-3",
			},
			"pages": map[string]interface{}{
				"type":        "integer",
				"description": "Do OCR from page 1 to this value",
			},
			"language": map[string]interface{}{
				"type": "string",
				"description": "OCR language(s), " +
					"e.g. eng, eng+deu",
			},
			"mode": map[string]interface{}{
				"type": "string",
				"description": "OCR mode: skip, redo, " +
					"force, skip_noarchive",
			},
			"skip_archive_file": map[string]interface{}{
				"type": "string",
				"description": "Archive file generation: " +
					"never, with_text, always",
			},
			"image_dpi": map[string]interface{}{
				"type":        "integer",
				"description": "Image DPI fallback value",
			},
			"unpaper_clean": map[string]interface{}{
				"type": "string",
				"description": "Unpaper cleaning: " +
					"clean, clean-final, none",
			},
			"deskew": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable deskew",
			},
			"rotate_pages": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable page rotation",
			},
			"rotate_pages_threshold": map[string]interface{}{
				"type":        "number",
				"description": "Threshold for page rotation",
			},
			"max_image_pixels": map[string]interface{}{
				"type": "number",
				"description": "Maximum image size " +
					"for decompression",
			},
			"color_conversion_strategy": map[string]interface{}{
				"type": "string",
				"description": "Ghostscript color conversion: " +
					"LeaveColorUnchanged, RGB, " +
					"UseDeviceIndependentColor, " +
					"Gray, CMYK",
			},
			"user_args": map[string]interface{}{
				"type": "object",
				"description": "Additional OCRMyPDF " +
					"user arguments (JSON object)",
			},
			"app_title": map[string]interface{}{
				"type":        "string",
				"description": "Application title",
			},
			"barcodes_enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable barcode scanning",
			},
			"barcode_enable_tiff_support": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable barcode TIFF support",
			},
			"barcode_string": map[string]interface{}{
				"type":        "string",
				"description": "Barcode string pattern",
			},
			"barcode_retain_split_pages": map[string]interface{}{
				"type":        "boolean",
				"description": "Retain pages after barcode split",
			},
			"barcode_enable_asn": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable ASN barcode",
			},
			"barcode_asn_prefix": map[string]interface{}{
				"type":        "string",
				"description": "ASN barcode prefix",
			},
			"barcode_upscale": map[string]interface{}{
				"type":        "number",
				"description": "Barcode upscale factor",
			},
			"barcode_dpi": map[string]interface{}{
				"type":        "integer",
				"description": "Barcode DPI",
			},
			"barcode_max_pages": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum pages for barcode scan",
			},
			"barcode_enable_tag": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable tag barcode",
			},
			"barcode_tag_mapping": map[string]interface{}{
				"type": "object",
				"description": "Tag barcode mapping " +
					"(JSON object)",
			},
		},
		"required": []string{"id"},
	}
}
