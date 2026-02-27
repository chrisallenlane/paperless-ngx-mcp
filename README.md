# paperless-ngx-mcp

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for
[Paperless-NGX](https://docs.paperless-ngx.com/), written in Go.

This MCP server integrates Paperless-NGX with Claude and other AI assistants,
providing tools to query and manage your document management system.

## Features

- **Complete MCP Implementation**: Full JSON-RPC 2.0 over stdio
- **Type-Safe Models**: Structured Go types with JSON marshaling
- **Token Authentication**: Paperless-NGX API token-based auth
- **Tool System**: Clean interface for adding new capabilities
- **Testing Infrastructure**: Comprehensive test suite
- **Code Quality Tools**: Formatting, linting, and vetting built-in

## Project Structure

```
paperless-ngx-mcp/
├── cmd/
│   └── paperless-ngx-mcp/   # Main application entry point
├── internal/
│   ├── client/               # HTTP client for Paperless-NGX API
│   ├── models/               # Paperless-NGX data structures
│   ├── server/               # MCP JSON-RPC server
│   └── tools/                # Tool implementations
├── Makefile                  # Build automation
├── CLAUDE.md                 # Development guidance
├── README.md                 # This file
└── SETUP.md                  # Setup instructions
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)
- A running Paperless-NGX instance with an API token

### Installation

```bash
make install
```

### Configuration

Set environment variables for your Paperless-NGX instance:

```bash
export PAPERLESS_URL="https://paperless.example.com"
export PAPERLESS_TOKEN="your-api-token"
```

### Running

The MCP server communicates via stdin/stdout:

```bash
./dist/paperless-ngx-mcp
```

See `SETUP.md` for integration with Claude Code and Claude Desktop.

## Available Tools

### System

#### `get_status`
Returns the current system status of the Paperless-NGX server, including
version, OS, storage usage, database/Redis/Celery status, and task timestamps.

#### `get_config`
Returns the current application configuration, grouped by category (OCR, App,
Barcode). Null values display as "(default)".

#### `update_config`
Updates application configuration via PATCH. Accepts `id` (required) plus any
config field. Only included fields are modified.

### Correspondents

#### `list_correspondents`
Lists correspondents with optional filtering by name and pagination (page,
page_size parameters).

#### `get_correspondent`
Gets a single correspondent by ID, including name, match pattern, matching
algorithm, document count, and last correspondence date.

#### `create_correspondent`
Creates a new correspondent. Requires `name`; optionally accepts `match`,
`matching_algorithm`, and `is_insensitive`.

#### `update_correspondent`
Updates an existing correspondent via PATCH. Requires `id`; any other
correspondent field is optional.

#### `delete_correspondent`
Deletes a correspondent by ID.

### Custom Fields

#### `list_custom_fields`
Lists custom fields with optional filtering by name and pagination (page,
page_size parameters).

#### `get_custom_field`
Gets a single custom field by ID, including name, data type, extra data, and
document count.

#### `create_custom_field`
Creates a new custom field. Requires `name` and `data_type`; optionally accepts
`extra_data`.

#### `update_custom_field`
Updates an existing custom field via PATCH. Requires `id`; any other field is
optional.

#### `delete_custom_field`
Deletes a custom field by ID.

### Document Types

#### `list_document_types`
Lists document types with optional filtering by name and pagination (page,
page_size parameters).

#### `get_document_type`
Gets a single document type by ID, including name, match pattern, matching
algorithm, and document count.

#### `create_document_type`
Creates a new document type. Requires `name`; optionally accepts `match`,
`matching_algorithm`, and `is_insensitive`.

#### `update_document_type`
Updates an existing document type via PATCH. Requires `id`; any other
document type field is optional.

#### `delete_document_type`
Deletes a document type by ID.

### Documents

#### `list_documents`
Lists documents with optional filtering and full-text search. Supports
`page`, `page_size`, `search` (full-text), `correspondent` (ID),
`document_type` (ID), `tags` (array of IDs, matches all), and
`is_in_inbox` (boolean).

#### `get_document`
Gets a single document by ID, including title, correspondent, document
type, tags, dates, ASN, file names, MIME type, page count, custom
fields, and a content preview.

#### `update_document`
Updates an existing document via PATCH. Requires `id`; optionally
accepts `title`, `correspondent`, `document_type`, `storage_path`,
`tags`, `archive_serial_number`, `created`, and `custom_fields`.

#### `delete_document`
Deletes a document by ID.

#### `get_document_metadata`
Gets file-level metadata for a document, including checksums, sizes,
MIME type, archive version details, and OCR language.

#### `get_document_suggestions`
Gets AI-generated suggestions for a document, including correspondent,
document type, storage path, tags, and dates.

#### `get_next_asn`
Gets the next available archive serial number (ASN).

#### `upload_document`
Uploads a local file to Paperless-NGX for processing. Requires
`file_path` (absolute); optionally accepts `title`, `correspondent`,
`document_type`, `storage_path`, `tags`, `archive_serial_number`, and
`created`. Returns a task ID for async processing.

#### `download_document`
Downloads a document file to the local filesystem. Requires `id` and
`save_path`; optionally accepts `original` (boolean, default false) to
download the original file instead of the archive version.

### Tags

#### `list_tags`
Lists tags with optional filtering by name and pagination (page,
page_size parameters).

#### `get_tag`
Gets a single tag by ID, including name, color, text color, match
pattern, matching algorithm, inbox flag, parent, and children.

#### `create_tag`
Creates a new tag. Requires `name`; optionally accepts `color`,
`match`, `matching_algorithm`, `is_insensitive`, `is_inbox_tag`,
and `parent`.

#### `update_tag`
Updates an existing tag via PATCH. Requires `id`; any other tag
field is optional.

#### `delete_tag`
Deletes a tag by ID.

### Storage Paths

#### `list_storage_paths`
Lists storage paths with optional filtering by name and pagination
(page, page_size parameters).

#### `get_storage_path`
Gets a single storage path by ID, including name, path template,
match pattern, matching algorithm, and document count.

#### `create_storage_path`
Creates a new storage path. Requires `name` and `path` (template
string, e.g., `{correspondent}/{document_type}/{title}`);
optionally accepts `match`, `matching_algorithm`, and
`is_insensitive`.

#### `update_storage_path`
Updates an existing storage path via PATCH. Requires `id`; any
other storage path field is optional.

#### `delete_storage_path`
Deletes a storage path by ID.

### Saved Views

#### `list_saved_views`
Lists saved views with optional pagination (page, page_size
parameters).

#### `get_saved_view`
Gets a single saved view by ID, including name, dashboard/sidebar
flags, sort configuration, display mode, and all filter rules
with human-readable rule type names.

#### `create_saved_view`
Creates a new saved view. Requires `name`, `show_on_dashboard`,
`show_in_sidebar`, and `filter_rules`; optionally accepts
`sort_field`, `sort_reverse`, `page_size`, and `display_mode`.

#### `update_saved_view`
Updates an existing saved view via PATCH. Requires `id`; any
other saved view field is optional.

#### `delete_saved_view`
Deletes a saved view by ID.

### Document Notes

#### `list_document_notes`
Lists notes for a document. Requires `id`; optionally accepts
`page` and `page_size`.

#### `create_document_note`
Adds a note to a document. Requires `id` and `note` (text
content). Returns the full updated notes list for the document.

#### `delete_document_note`
Deletes a note from a document. Requires `document_id` and
`note_id`. Returns the updated notes list after deletion.

### Tasks

#### `list_tasks`
Lists background processing tasks. Optionally accepts `status`
(FAILURE, PENDING, RECEIVED, RETRY, REVOKED, STARTED, SUCCESS),
`task_name` (consume_file, train_classifier, check_sanity,
index_optimize), `type` (auto_task, scheduled_task,
manual_task), and `task_id` (Celery task UUID).

#### `get_task`
Gets a single background task by ID, including task UUID, status,
type, file name, creation and completion dates, result, and
related document.

### Statistics and Trash

#### `get_statistics`
Gets document and resource count statistics from Paperless-NGX,
including document counts, inbox counts, and per-type breakdowns.

#### `list_trash`
Lists soft-deleted documents in the Paperless-NGX trash, with
optional pagination (page, page_size parameters).

## Development

### Build

```bash
make build
```

### Test

```bash
make test
make coverage
```

### Code Quality

```bash
make check   # format, lint, vet, test
```

## Project Conventions

- **Formatting**: 80-column line wrapping with golines + gofumpt
- **Testing**: Standard library `testing` package
- **Dependencies**: Minimal - only Go stdlib for production code
- **Error Handling**: Always wrap errors with context
- **Type Safety**: Use structs, not `map[string]interface{}`

## Architecture

### MCP Protocol Flow

```
Claude → stdin → JSON-RPC Request → Tool Execution → JSON-RPC Response → stdout → Claude
```

### Key Components

- **Server** (`internal/server`): Handles JSON-RPC protocol
- **Client** (`internal/client`): Makes HTTP requests to Paperless-NGX API
- **Tools** (`internal/tools`): Implements MCP tools
- **Models** (`internal/models`): Type-safe Paperless-NGX data structures

## Resources

- [MCP Specification](https://modelcontextprotocol.io/)
- [Paperless-NGX Documentation](https://docs.paperless-ngx.com/)
- [Paperless-NGX API](https://docs.paperless-ngx.com/api/)
