package tools

// tagSchema returns an input schema for tag tools.
// Set includeID to true for update tools, false for create.
func tagSchema(
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Tag name",
		},
		"color": map[string]interface{}{
			"type": "string",
			"description": "Hex color code " +
				"(e.g., #a6cee3)",
		},
		"is_inbox_tag": map[string]interface{}{
			"type": "boolean",
			"description": "Automatically assign to " +
				"newly consumed documents",
		},
		"parent": map[string]interface{}{
			"type": "integer",
			"description": "Parent tag ID " +
				"for hierarchical tags",
		},
	}
	addMatchableProps(props)

	required := withIDForUpdate(
		props,
		"Tag ID to update",
		includeID,
		[]string{"name"},
	)

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

// storagePathSchema returns an input schema for storage path
// tools. Set includeID to true for update tools, false for create.
func storagePathSchema(
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Storage path name",
		},
		"path": map[string]interface{}{
			"type": "string",
			"description": "Storage path template " +
				"(e.g., {correspondent}/" +
				"{document_type}/{title})",
		},
	}
	addMatchableProps(props)

	required := withIDForUpdate(
		props,
		"Storage path ID to update",
		includeID,
		[]string{"name", "path"},
	)

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
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

	required := withIDForUpdate(
		props,
		"Custom field ID to update",
		includeID,
		[]string{"name", "data_type"},
	)

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

// savedViewProps returns the shared property definitions for
// saved view schemas.
func savedViewProps() map[string]interface{} {
	return map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Saved view name",
		},
		"show_on_dashboard": map[string]interface{}{
			"type": "boolean",
			"description": "Show this view " +
				"on the dashboard",
		},
		"show_in_sidebar": map[string]interface{}{
			"type": "boolean",
			"description": "Show this view " +
				"in the sidebar",
		},
		"sort_field": map[string]interface{}{
			"type": "string",
			"description": "Field to sort by " +
				"(e.g., created, added, title)",
		},
		"sort_reverse": map[string]interface{}{
			"type": "boolean",
			"description": "Reverse sort order " +
				"(default false)",
		},
		"filter_rules": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"rule_type": map[string]interface{}{
						"type": "integer",
						"description": "Filter rule type: " +
							"0=title contains, " +
							"1=content contains, " +
							"3=correspondent is, " +
							"4=document type is, " +
							"5=is in inbox, " +
							"6=has tag, " +
							"17=does not have tag, " +
							"20=fulltext query, " +
							"25=storage path is " +
							"(additional types 0-47 exist)",
					},
					"value": map[string]interface{}{
						"type": "string",
						"description": "Filter value " +
							"(IDs as strings, " +
							"or search text)",
					},
				},
				"required": []string{
					"rule_type",
				},
			},
			"description": "Filter rules for " +
				"this saved view",
		},
		"page_size": map[string]interface{}{
			"type":        "integer",
			"description": "Results per page",
		},
		"display_mode": map[string]interface{}{
			"type": "string",
			"description": "Display mode: table, " +
				"smallCards, largeCards",
		},
	}
}

// savedViewSchema returns an input schema for saved view
// tools. Set includeID to true for update tools, false for
// create.
func savedViewSchema(
	includeID bool,
) map[string]interface{} {
	props := savedViewProps()

	required := withIDForUpdate(
		props,
		"Saved view ID to update",
		includeID,
		[]string{
			"name",
			"show_on_dashboard",
			"show_in_sidebar",
			"filter_rules",
		},
	)

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}
