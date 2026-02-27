# CLAUDE.md

This file provides guidance to Claude Code when working with this project.

## Project Overview

**paperless-ngx-mcp** is a Model Context Protocol (MCP) server for
[Paperless-NGX](https://docs.paperless-ngx.com/), written in Go. It provides
tools for querying and managing a Paperless-NGX document management system
through Claude and other AI assistants.

**Tech Stack:**
- **Language**: Go 1.21+
- **Protocol**: MCP (Model Context Protocol) via JSON-RPC 2.0 over stdio
- **Dependencies**: Minimal - Go stdlib only for production code

## Project Structure

```
paperless-ngx-mcp/
├── cmd/
│   └── paperless-ngx-mcp/     # Main application
│       └── main.go            # Entry point, configuration, initialization
├── internal/                  # Private application packages
│   ├── client/                # HTTP client for Paperless-NGX API
│   │   ├── client.go          # HTTP client with request helpers
│   │   └── client_test.go     # Client tests
│   ├── models/                # Paperless-NGX data structures
│   │   ├── models.go          # Domain models
│   │   └── models_test.go     # Model tests
│   ├── server/                # MCP server implementation
│   │   ├── server.go          # JSON-RPC server, request routing
│   │   ├── server_test.go     # Protocol tests
│   │   └── types.go           # JSON-RPC request/response types
│   └── tools/                           # MCP tool implementations
│       ├── tool.go                      # Tool interface definition
│       ├── tool_factory.go              # Generic data-driven tool types
│       ├── tool_constructors.go         # Constructor functions for all tools
│       ├── helpers.go                   # Shared utility/HTTP helper functions
│       ├── helpers_test.go              # Helper function tests
│       ├── schemas.go                   # JSON schema builder functions
│       ├── formatters.go                # All response formatting functions
│       ├── get_config.go                # Application config tool (array response)
│       ├── get_config_test.go           # Config tool tests
│       ├── get_statistics.go            # Statistics tool (dynamic response)
│       ├── get_statistics_test.go       # Statistics tool tests
│       ├── list_documents.go            # List documents tool (custom filters)
│       ├── list_documents_test.go       # List documents tests
│       ├── list_tasks.go                # List tasks tool (custom filters)
│       ├── list_tasks_test.go           # List tasks tests
│       ├── create_custom_field.go       # Create custom field tool
│       ├── create_custom_field_test.go  # Create custom field tests
│       ├── create_storage_path.go       # Create storage path tool
│       ├── create_storage_path_test.go  # Create storage path tests
│       ├── create_tag.go                # Create tag tool
│       ├── create_tag_test.go           # Create tag tests
│       ├── create_saved_view.go         # Create saved view tool
│       ├── create_saved_view_test.go    # Create saved view tests
│       ├── list_document_notes.go       # List document notes tool
│       ├── create_document_note.go      # Create document note tool
│       ├── delete_document_note.go      # Delete document note tool
│       ├── document_notes_test.go       # Document notes tool tests
│       ├── upload_document.go           # Upload document tool
│       ├── upload_document_test.go      # Upload document tests
│       ├── download_document.go         # Download document tool
│       ├── download_document_test.go    # Download document tests
│       ├── tool_common_test.go          # Cross-cutting tests (description, schema)
│       ├── delete_test.go               # Delete tool tests
│       ├── get_correspondent_test.go    # Get correspondent tests
│       ├── get_custom_field_test.go     # Get custom field tests
│       ├── get_document_metadata_test.go # Get document metadata tests
│       ├── get_document_suggestions_test.go # Get document suggestions tests
│       ├── get_document_test.go         # Get document tests
│       ├── get_document_type_test.go    # Get document type tests
│       ├── get_next_asn_test.go         # Get next ASN tests
│       ├── get_status_test.go           # Status tool tests
│       ├── list_document_types_test.go  # List document types tests
│       ├── list_trash_test.go           # List trash tests
│       ├── update_config_test.go        # Config update tests
│       └── update_document_test.go      # Update document tests
├── Makefile                   # Build automation
├── CLAUDE.md                  # This file
├── README.md                  # User-facing documentation
└── SETUP.md                   # Setup instructions
```

This follows the **standard Go project layout**:
- `cmd/` - Main application entry points
- `internal/` - Private packages that cannot be imported by external projects

## Architecture

### MCP Protocol Implementation

The server implements MCP via **JSON-RPC 2.0 over stdio**:

1. **Stdin** - JSON-RPC requests from Claude
2. **Process** - Route to handlers, execute tools
3. **Stdout** - JSON-RPC responses back to Claude

**Key Methods:**
- `initialize` - Handshake, declare capabilities
- `tools/list` - Return available tools and their schemas
- `tools/call` - Execute a specific tool

**Flow:**
```
Claude → stdin → Scanner → JSON unmarshal → handleRequest() → execute tool → JSON marshal → stdout → Claude
```

### HTTP Client (`internal/client/client.go`)

HTTP client for Paperless-NGX API requests:

**HTTP Methods:**
- `Get(ctx, path)` - GET requests
- `Post(ctx, path, body)` - POST requests with JSON body
- `Patch(ctx, path, body)` - PATCH requests with JSON body
- `Delete(ctx, path)` - DELETE requests
- `PostMultipart(ctx, path, body, contentType)` - POST requests with multipart/form-data body

**Testing Support:**
- `HTTPDoer` interface allows mocking HTTP requests
- `NewWithHTTPClient(baseURL, token, httpClient)` - Test constructor
- Use `httptest.Server` for testing without real API calls

### Tool Interface (`internal/tools/tool.go`)

Every tool must implement:

```go
type Tool interface {
    Execute(ctx context.Context, args json.RawMessage) (string, error)
    Description() string
    InputSchema() map[string]interface{}
}
```

**Execute** - Runs the tool with parsed arguments, returns formatted string response
**Description** - Human-readable description for Claude
**InputSchema** - JSON Schema defining required/optional parameters

### Data-Driven Tool Factory (`internal/tools/tool_factory.go`)

Generic tool types that implement the `Tool` interface. Each type is parameterized
by the response model type `T` and holds configuration (description, schema, path,
formatter) rather than behavior. This eliminates per-tool boilerplate.

**Factory types:**
- **`noArgGetTool[T]`** - GET with no input; unmarshals response into `T`, calls `format(*T)`
- **`getTool[T]`** - GET by ID; calls `fetchByID[T]`, calls `format(id, *T)`
- **`listTool[T]`** - Paginated list; calls `listResources[T]`, calls `format(*PaginatedList[T])`
- **`patchTool[T]`** - PATCH by ID; calls `patchByID[T]`, calls `format(*T)`
- **`createMatchableTool[T]`** - POST for matchable resources; calls `createMatchable[T]`, calls `format(*T)`
- **`deleteTool`** - DELETE by ID; calls `deleteByID`, returns confirmation string

**Constructor functions** live in `tool_constructors.go` and instantiate these types
with the appropriate configuration for each named tool.

### Tool Registration (`internal/server/server.go`)

Tools are registered in `registerTools()`:

```go
s.tools["tool_name"] = tools.NewToolName(s.client)
```

The server automatically discovers and exposes all registered tools via `tools/list`.

### Type-Safe Models (`internal/models/models.go`)

Domain models for Paperless-NGX API responses:

```go
type SystemStatus struct {
    PNGXVersion string `json:"pngx_version"`
    // ...
}
```

**Benefits:**
- Compile-time type checking (no `map[string]interface{}`)
- IDE autocomplete support
- Self-documenting code

### Helper Functions (`internal/tools/helpers.go`)

Shared utility functions to eliminate code duplication:

**HTTP helpers:**
- **`doAPIRequest(ctx, client, path)`** - Common GET request pattern (expects 200)
- **`doPostRequest(ctx, client, path, body)`** - Common POST request pattern (expects 201)
- **`doPatchRequest(ctx, client, path, body)`** - Common PATCH request pattern (expects 200)
- **`doDeleteRequest(ctx, client, path)`** - Common DELETE request pattern (expects 204)
- **`readResponse(resp, expectedStatus)`** - Read and validate HTTP response body

**Argument parsing helpers:**
- **`parseIDArg(args)`** - Parse and validate a required `id` field from JSON args
- **`parsePatchArgs(args)`** - Extract `id` and build patch body (all fields except `id`)

**Generic operation helpers:**
- **`fetchByID[T](ctx, client, args, pathFmt)`** - Parse ID, fetch resource, unmarshal response
- **`patchByID[T](ctx, client, args, pathFmt)`** - Parse patch args, PATCH resource, unmarshal response
- **`listResources[T](ctx, client, basePath, args)`** - Build list path, fetch, unmarshal paginated list
- **`createMatchable[T](ctx, client, args, path)`** - Parse matchable params, POST, unmarshal response
- **`deleteByID(ctx, client, args, pathFmt, resourceName)`** - Parse ID, DELETE resource, return confirmation

**Shared types:**
- **`matchableCreateParams`** - Common params struct for creating matchable resources
- **`listParams`** - Common pagination and filter parameters

**File path helpers:**
- **`validateFilePath(path)`** - Validates that a file path is absolute and contains no `..` traversal sequences

**List query helpers:**
- **`buildListPath(basePath, args)`** - Build URL path with query parameters from list args

### Schema Builders (`internal/tools/schemas.go`)

JSON schema builder functions, separated from helpers to keep concerns distinct:

- **`emptySchema()`** - Schema with no parameters (for tools like `get_status`, `get_next_asn`)
- **`idOnlySchema(desc)`** - Schema with a single required `id` integer field
- **`paginatedListSchema()`** - Schema for list tools (page, page_size, name filter)
- **`paginationOnlySchema()`** - Schema for list tools without a name filter (page, page_size only); used by `list_saved_views` and `list_trash`
- **`matchableResourceSchema(resourceName, includeID)`** - Schema for matchable resources (correspondents, document types) with name, match, matching_algorithm, is_insensitive fields; set `includeID` true for update tools
- **`tagSchema(includeID)`** - Schema for tag tools with name, color, match, matching_algorithm, is_insensitive, is_inbox_tag, parent fields; set `includeID` true for update tools
- **`storagePathSchema(includeID)`** - Schema for storage path tools with name, path, match, matching_algorithm, is_insensitive fields; set `includeID` true for update tools
- **`taskListSchema()`** - Schema for the task list tool with status, task_name, type, task_id filters
- **`savedViewCreateSchema()`** - Schema for creating a saved view with name, show_on_dashboard, show_in_sidebar, filter_rules (required) plus sort/display options
- **`savedViewUpdateSchema()`** - Schema for updating a saved view; adds required `id` to the create schema fields
- **`customFieldSchema(includeID)`** - Schema for custom field tools with name, data_type, extra_data fields
- **`documentUpdateSchema()`** - Schema for the document update tool with all document fields
- **`configUpdateSchema()`** - Schema for the config update tool with all config fields
- **`withIDForUpdate(props, idDesc, includeID, createRequired)`** - Helper that conditionally adds an `id` field and adjusts required fields; used by schema builders that serve both create and update tools

### Response Formatters (`internal/tools/formatters.go`)

All response formatting functions are centralized here:

- **`formatStatus`** - System status summary
- **`formatConfig`** - Application configuration grouped by category
- **`formatStatistics`** - Document and resource count statistics (dynamic key-value map)
- **`formatMatchableFields`** - Shared formatter for resources with matching fields (name, slug, match, algorithm, document count); used by correspondents, document types, tags, and storage paths
- **`formatPaginatedList[T]`** - Generic paginated list formatter; handles empty message, header with count, per-item formatting, and pagination hint
- **`formatCorrespondent`** / **`formatCorrespondentList`** - Correspondent details and lists
- **`formatCustomField`** / **`formatCustomFieldList`** - Custom field details and lists
- **`formatDocumentType`** / **`formatDocumentTypeList`** - Document type details and lists
- **`formatDocument`** / **`formatDocumentList`** - Document details and lists
- **`formatDocumentMetadata`** - Document file metadata (checksums, sizes, OCR language)
- **`formatDocumentSuggestions`** - AI-generated document suggestions
- **`formatTag`** / **`formatTagList`** - Tag details and lists (includes color, inbox flag, parent/children)
- **`formatStoragePath`** / **`formatStoragePathList`** - Storage path details and lists
- **`formatTask`** / **`formatTaskArray`** - Single task details and task array list (note: tasks use an array response, not paginated)
- **`formatNote`** / **`formatNoteList`** - Document note details and lists
- **`formatSavedView`** / **`formatSavedViewList`** - Saved view details and lists with filter rule display
- **`ruleTypeName`** - Filter rule type integer-to-name lookup (types 0–47)
- **`formatOpt[T]`** / **`formatOptJSON`** - Nullable field formatting helpers
- **`formatOptInt`** / **`formatOptStr`** / **`formatOptDate`** - Nullable int, string, and date formatting helpers
- **`formatFileSize`** - Human-readable byte size formatting (B/KB/MB/GB)
- **`formatIntSlice`** / **`formatStringSlice`** - Slice-to-string formatting helpers
- **`matchingAlgorithmName`** - Matching algorithm integer-to-name lookup
- **`formatDate`** / **`formatTaskLine`** - Date and task status line formatting

## Development Workflow

### Building

```bash
make build
# Output: dist/paperless-ngx-mcp
```

### Testing

```bash
make fmt       # Format code
make lint      # Lint code
make vet       # Run vet
make test      # Run tests
make coverage  # Tests with coverage report
make check     # All checks (format, lint, vet, test)
```

### Installing

```bash
make install   # Install to $GOPATH/bin
```

## Adding a New Tool

Most tools are implemented using the data-driven factory types in `tool_factory.go`.
Choose the appropriate factory type based on what the tool does, then add a constructor
in `tool_constructors.go`. Only create a dedicated file if the tool has logic that
cannot be expressed with the factory types (see `get_config.go` and `list_documents.go`
for examples).

### 1. Choose the factory type

- **`noArgGetTool[T]`** - GET with no parameters (e.g., `get_status`, `get_next_asn`)
- **`getTool[T]`** - GET by ID (e.g., `get_correspondent`, `get_document`, `get_task`)
- **`listTool[T]`** - Paginated list with optional name filter (e.g., `list_correspondents`, `list_tags`, `list_saved_views`)
- **`patchTool[T]`** - PATCH by ID (e.g., `update_correspondent`, `update_config`, `update_saved_view`)
- **`createMatchableTool[T]`** - POST for resources with matching fields (e.g., `create_correspondent`)
- **`deleteTool`** - DELETE by ID (e.g., `delete_correspondent`, `delete_tag`, `delete_storage_path`)

### 2. Add a constructor in `internal/tools/tool_constructors.go`

```go
// NewGetMyResource creates a tool to get a my-resource by ID.
func NewGetMyResource(c *client.Client) Tool {
    return &getTool[models.MyResource]{
        client:  c,
        desc:    "Get a my-resource by ID from Paperless-NGX",
        schema:  idOnlySchema("My resource ID"),
        pathFmt: "/api/my_resources/%d/",
        format: func(_ int, v *models.MyResource) string {
            return formatMyResource(v)
        },
    }
}
```

Add a formatter in `internal/tools/formatters.go` if needed.

### 3. Register in `internal/server/server.go`

Add to `registerTools()`:
```go
s.tools["get_my_resource"] = tools.NewGetMyResource(s.client)
```

### 4. Write tests

Add the tool to the cross-cutting tests in `tool_common_test.go` (description and
schema coverage). Add a tool-specific test file for input validation and response
parsing if needed.

### 5. Rebuild and test

```bash
make check
make build
```

### When to create a dedicated file

Create a dedicated tool file (e.g., `my_tool.go`) only when the tool cannot use a
factory type, such as when:
- The API response is not a standard JSON object (e.g., array response like `get_config`, dynamic map like `get_statistics`)
- The URL construction requires custom logic beyond standard pagination (e.g., `list_documents`, `list_tasks`)
- The tool has multi-step behavior or non-standard HTTP semantics (e.g., `create_document_note` returns 200 not 201; `delete_document_note` returns 200 not 204)
- The resource has required fields that don't fit the `matchableCreateParams` struct (e.g., `create_storage_path` requires both `name` and `path`; `create_saved_view` has required boolean fields)

## Code Quality Standards

### Input Validation
Always validate input before making API calls.

### Use Helper Functions
Prefer `doAPIRequest`/`doPatchRequest` over duplicating HTTP boilerplate.

### Type Safety
Use models package instead of `map[string]interface{}`.

### Error Messages
Include context in error messages:
```go
return "", fmt.Errorf("descriptive context: %w", err)
```

### Testing Requirements
Every new tool should have:
- Input validation tests
- Description and schema tests
- Tests run in `make check`

### Code Organization
- Keep it simple - prefer standard library over dependencies
- Use factory types in `tool_factory.go` + constructors in `tool_constructors.go` for standard CRUD tools
- Create dedicated files only for tools with non-standard logic (array responses, custom URL construction)
- Shared HTTP/operation logic in `helpers.go`
- JSON schema builders in `schemas.go`
- Response formatting in `formatters.go`
- Type definitions in `models.go`

## Current Tools

Most tools are constructed via factory types in `tool_constructors.go`. Exceptions
are noted below.

### `get_status` (`tool_constructors.go` — `NewGetStatus`)

- **Endpoint**: `GET /api/status/`
- **Input**: None
- **Output**: Formatted status summary
- **Model**: `models.SystemStatus`

### `get_config` (`get_config.go` — dedicated file)

- **Endpoint**: `GET /api/config/`
- **Input**: None
- **Output**: Config summary grouped by category (OCR, App, Barcode)
- **Model**: `models.ApplicationConfiguration`
- **Note**: Dedicated file because response is a JSON array; tool takes the first element

### `update_config` (`tool_constructors.go` — `NewUpdateConfig`)

- **Endpoint**: `PATCH /api/config/{id}/`
- **Input**: `id` (required) + any config fields
- **Output**: Updated config summary
- **Note**: Only included fields are modified; `app_logo` skipped (binary upload)

### `list_correspondents` (`tool_constructors.go` — `NewListCorrespondents`)

- **Endpoint**: `GET /api/correspondents/`
- **Input**: `page`, `page_size`, `name` (all optional)
- **Output**: Paginated correspondent list
- **Model**: `models.PaginatedList[models.Correspondent]`

### `get_correspondent` (`tool_constructors.go` — `NewGetCorrespondent`)

- **Endpoint**: `GET /api/correspondents/{id}/`
- **Input**: `id` (required)
- **Output**: Correspondent details with matching algorithm name
- **Model**: `models.Correspondent`

### `create_correspondent` (`tool_constructors.go` — `NewCreateCorrespondent`)

- **Endpoint**: `POST /api/correspondents/`
- **Input**: `name` (required), `match`, `matching_algorithm`, `is_insensitive` (optional)
- **Output**: Created correspondent details

### `update_correspondent` (`tool_constructors.go` — `NewUpdateCorrespondent`)

- **Endpoint**: `PATCH /api/correspondents/{id}/`
- **Input**: `id` (required) + any correspondent fields
- **Output**: Updated correspondent details

### `delete_correspondent` (`tool_constructors.go` — `NewDeleteCorrespondent`)

- **Endpoint**: `DELETE /api/correspondents/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `list_custom_fields` (`tool_constructors.go` — `NewListCustomFields`)

- **Endpoint**: `GET /api/custom_fields/`
- **Input**: `page`, `page_size`, `name` (all optional)
- **Output**: Paginated custom field list
- **Model**: `models.PaginatedList[models.CustomField]`

### `get_custom_field` (`tool_constructors.go` — `NewGetCustomField`)

- **Endpoint**: `GET /api/custom_fields/{id}/`
- **Input**: `id` (required)
- **Output**: Custom field details with extra data
- **Model**: `models.CustomField`

### `create_custom_field` (`create_custom_field.go` — dedicated file)

- **Endpoint**: `POST /api/custom_fields/`
- **Input**: `name`, `data_type` (required), `extra_data` (optional)
- **Output**: Created custom field details
- **Note**: Dedicated file because custom fields use a different schema than matchable resources

### `update_custom_field` (`tool_constructors.go` — `NewUpdateCustomField`)

- **Endpoint**: `PATCH /api/custom_fields/{id}/`
- **Input**: `id` (required) + any custom field fields
- **Output**: Updated custom field details

### `delete_custom_field` (`tool_constructors.go` — `NewDeleteCustomField`)

- **Endpoint**: `DELETE /api/custom_fields/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `list_document_types` (`tool_constructors.go` — `NewListDocumentTypes`)

- **Endpoint**: `GET /api/document_types/`
- **Input**: `page`, `page_size`, `name` (all optional)
- **Output**: Paginated document type list
- **Model**: `models.PaginatedList[models.DocumentType]`

### `get_document_type` (`tool_constructors.go` — `NewGetDocumentType`)

- **Endpoint**: `GET /api/document_types/{id}/`
- **Input**: `id` (required)
- **Output**: Document type details with matching algorithm name
- **Model**: `models.DocumentType`

### `create_document_type` (`tool_constructors.go` — `NewCreateDocumentType`)

- **Endpoint**: `POST /api/document_types/`
- **Input**: `name` (required), `match`, `matching_algorithm`, `is_insensitive` (optional)
- **Output**: Created document type details

### `update_document_type` (`tool_constructors.go` — `NewUpdateDocumentType`)

- **Endpoint**: `PATCH /api/document_types/{id}/`
- **Input**: `id` (required) + any document type fields
- **Output**: Updated document type details

### `delete_document_type` (`tool_constructors.go` — `NewDeleteDocumentType`)

- **Endpoint**: `DELETE /api/document_types/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `list_documents` (`list_documents.go` — dedicated file)

- **Endpoint**: `GET /api/documents/`
- **Input**: `page`, `page_size`, `search`, `correspondent` (ID), `document_type` (ID), `tags` (array of IDs), `is_in_inbox` (all optional)
- **Output**: Paginated document list with concise summaries
- **Model**: `models.PaginatedList[models.Document]`
- **Note**: Dedicated file because document filtering uses custom URL parameters beyond standard pagination

### `get_document` (`tool_constructors.go` — `NewGetDocument`)

- **Endpoint**: `GET /api/documents/{id}/`
- **Input**: `id` (required)
- **Output**: Full document details including custom fields and content preview
- **Model**: `models.Document`

### `update_document` (`tool_constructors.go` — `NewUpdateDocument`)

- **Endpoint**: `PATCH /api/documents/{id}/`
- **Input**: `id` (required) + `title`, `correspondent`, `document_type`, `storage_path`, `tags`, `archive_serial_number`, `created`, `custom_fields` (all optional)
- **Output**: Updated document details

### `delete_document` (`tool_constructors.go` — `NewDeleteDocument`)

- **Endpoint**: `DELETE /api/documents/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `get_document_metadata` (`tool_constructors.go` — `NewGetDocumentMetadata`)

- **Endpoint**: `GET /api/documents/{id}/metadata/`
- **Input**: `id` (required)
- **Output**: File metadata (checksums, sizes, MIME type, archive version, OCR language)
- **Model**: `models.DocumentMetadata`

### `get_document_suggestions` (`tool_constructors.go` — `NewGetDocumentSuggestions`)

- **Endpoint**: `GET /api/documents/{id}/suggestions/`
- **Input**: `id` (required)
- **Output**: AI-generated suggestions (correspondents, document types, storage paths, tags, dates)
- **Model**: `models.DocumentSuggestions`

### `get_next_asn` (`tool_constructors.go` — `NewGetNextASN`)

- **Endpoint**: `GET /api/documents/next_asn/`
- **Input**: None
- **Output**: Next available archive serial number
- **Note**: Response is a bare integer, not a JSON object

### `upload_document` (`upload_document.go` — dedicated file)

- **Endpoint**: `POST /api/documents/post_document/`
- **Content-Type**: `multipart/form-data`
- **Input**: `file_path` (required, must be absolute) + `title`, `correspondent`, `document_type`, `storage_path`, `tags`, `archive_serial_number`, `created` (all optional)
- **Output**: Confirmation with filename, size, and task ID
- **Note**: Dedicated file because it uses `client.PostMultipart`; expects HTTP 200; returns a task ID for async processing

### `download_document` (`download_document.go` — dedicated file)

- **Endpoint**: `GET /api/documents/{id}/download/`
- **Input**: `id`, `save_path` (both required); `original` (optional boolean, default false)
- **Output**: Confirmation with save path, size, and content type
- **Note**: Dedicated file because it streams response body to file; validates `save_path` with `validateFilePath`

### `list_tags` (`tool_constructors.go` — `NewListTags`)

- **Endpoint**: `GET /api/tags/`
- **Input**: `page`, `page_size`, `name` (all optional)
- **Output**: Paginated tag list with inbox flag indicator
- **Model**: `models.PaginatedList[models.Tag]`

### `get_tag` (`tool_constructors.go` — `NewGetTag`)

- **Endpoint**: `GET /api/tags/{id}/`
- **Input**: `id` (required)
- **Output**: Tag details including color, text color, inbox flag, parent, and children
- **Model**: `models.Tag`

### `create_tag` (`create_tag.go` — dedicated file)

- **Endpoint**: `POST /api/tags/`
- **Input**: `name` (required), `color`, `match`, `matching_algorithm`, `is_insensitive`, `is_inbox_tag`, `parent` (all optional)
- **Output**: Created tag details
- **Note**: Dedicated file because `tagSchema` includes fields (color, is_inbox_tag, parent) not in `matchableCreateParams`

### `update_tag` (`tool_constructors.go` — `NewUpdateTag`)

- **Endpoint**: `PATCH /api/tags/{id}/`
- **Input**: `id` (required) + any tag fields
- **Output**: Updated tag details

### `delete_tag` (`tool_constructors.go` — `NewDeleteTag`)

- **Endpoint**: `DELETE /api/tags/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `list_storage_paths` (`tool_constructors.go` — `NewListStoragePaths`)

- **Endpoint**: `GET /api/storage_paths/`
- **Input**: `page`, `page_size`, `name` (all optional)
- **Output**: Paginated storage path list
- **Model**: `models.PaginatedList[models.StoragePath]`

### `get_storage_path` (`tool_constructors.go` — `NewGetStoragePath`)

- **Endpoint**: `GET /api/storage_paths/{id}/`
- **Input**: `id` (required)
- **Output**: Storage path details including path template and matching fields
- **Model**: `models.StoragePath`

### `create_storage_path` (`create_storage_path.go` — dedicated file)

- **Endpoint**: `POST /api/storage_paths/`
- **Input**: `name`, `path` (both required), `match`, `matching_algorithm`, `is_insensitive` (all optional)
- **Output**: Created storage path details
- **Note**: Dedicated file because `create_storage_path` requires both `name` and `path`, which doesn't fit `matchableCreateParams`

### `update_storage_path` (`tool_constructors.go` — `NewUpdateStoragePath`)

- **Endpoint**: `PATCH /api/storage_paths/{id}/`
- **Input**: `id` (required) + any storage path fields
- **Output**: Updated storage path details

### `delete_storage_path` (`tool_constructors.go` — `NewDeleteStoragePath`)

- **Endpoint**: `DELETE /api/storage_paths/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `list_saved_views` (`tool_constructors.go` — `NewListSavedViews`)

- **Endpoint**: `GET /api/saved_views/`
- **Input**: `page`, `page_size` (both optional; no name filter)
- **Output**: Paginated saved view list with dashboard/sidebar flags
- **Model**: `models.PaginatedList[models.SavedView]`

### `get_saved_view` (`tool_constructors.go` — `NewGetSavedView`)

- **Endpoint**: `GET /api/saved_views/{id}/`
- **Input**: `id` (required)
- **Output**: Saved view details including all filter rules with human-readable rule type names
- **Model**: `models.SavedView`

### `create_saved_view` (`create_saved_view.go` — dedicated file)

- **Endpoint**: `POST /api/saved_views/`
- **Input**: `name`, `show_on_dashboard`, `show_in_sidebar`, `filter_rules` (all required); `sort_field`, `sort_reverse`, `page_size`, `display_mode` (all optional)
- **Output**: Created saved view details
- **Note**: Dedicated file because `show_on_dashboard`, `show_in_sidebar`, and `filter_rules` are required boolean/array fields that can't be expressed as `matchableCreateParams`

### `update_saved_view` (`tool_constructors.go` — `NewUpdateSavedView`)

- **Endpoint**: `PATCH /api/saved_views/{id}/`
- **Input**: `id` (required) + any saved view fields
- **Output**: Updated saved view details

### `delete_saved_view` (`tool_constructors.go` — `NewDeleteSavedView`)

- **Endpoint**: `DELETE /api/saved_views/{id}/`
- **Input**: `id` (required)
- **Output**: Confirmation message

### `list_document_notes` (`list_document_notes.go` — dedicated file)

- **Endpoint**: `GET /api/documents/{id}/notes/`
- **Input**: `id` (required), `page`, `page_size` (optional)
- **Output**: Paginated note list for the document
- **Model**: `models.PaginatedList[models.Note]`
- **Note**: Dedicated file because the endpoint is nested under a document ID

### `create_document_note` (`create_document_note.go` — dedicated file)

- **Endpoint**: `POST /api/documents/{id}/notes/`
- **Input**: `id`, `note` (both required)
- **Output**: Confirmation with the updated note list for the document
- **Note**: Dedicated file because the endpoint is nested under a document ID and returns HTTP 200 (not 201) with the full updated notes list

### `delete_document_note` (`delete_document_note.go` — dedicated file)

- **Endpoint**: `DELETE /api/documents/{document_id}/notes/?id={note_id}`
- **Input**: `document_id`, `note_id` (both required)
- **Output**: Confirmation with the updated note list for the document
- **Note**: Dedicated file because the note ID is passed as a query parameter (not path segment), the endpoint is nested under a document ID, and it returns HTTP 200 (not 204) with the remaining notes list

### `get_statistics` (`get_statistics.go` — dedicated file)

- **Endpoint**: `GET /api/statistics/`
- **Input**: None
- **Output**: Document and resource count statistics (formatted as sorted key-value pairs)
- **Note**: Dedicated file because the response is a dynamic JSON object with no fixed schema; unmarshaled as `map[string]interface{}`

### `list_tasks` (`list_tasks.go` — dedicated file)

- **Endpoint**: `GET /api/tasks/`
- **Input**: `status`, `task_name`, `type`, `task_id` (all optional)
- **Output**: Task list formatted as an array summary
- **Model**: `[]models.Task` (bare array, not paginated)
- **Note**: Dedicated file because the endpoint returns a bare JSON array (not paginated) and filtering uses custom query parameters

### `get_task` (`tool_constructors.go` — `NewGetTask`)

- **Endpoint**: `GET /api/tasks/{id}/`
- **Input**: `id` (required)
- **Output**: Task details including UUID, status, type, file name, dates, result, and related document
- **Model**: `models.Task`

### `list_trash` (`tool_constructors.go` — `NewListTrash`)

- **Endpoint**: `GET /api/trash/`
- **Input**: `page`, `page_size` (both optional; no name filter)
- **Output**: Paginated list of soft-deleted documents
- **Model**: `models.PaginatedList[models.Document]`
- **Note**: Uses `paginationOnlySchema()` (no name filter); reuses `formatDocumentList`

## Configuration

**Environment Variables:**
- `PAPERLESS_URL` - Base URL of the Paperless-NGX instance
- `PAPERLESS_TOKEN` - API authentication token

## Response Formatting Guidelines

Tools should return **human-readable formatted strings**, not raw JSON.

## Error Handling

- Always wrap errors with context using `fmt.Errorf("context: %w", err)`
- Check HTTP status codes
- Handle empty results gracefully

## Important Patterns

### Context Propagation
- Always accept and pass `context.Context` through the call chain

### JSON Marshaling
- Use `json.RawMessage` for unknown/dynamic structures

### Resource Cleanup
- Always `defer resp.Body.Close()` after HTTP requests

## Dependencies

Zero external dependencies for production code - uses only Go standard library.

## Version Information

- MCP Protocol Version: `2024-11-05`
- Server Version: `0.1.0`
- Go Version: 1.21+ required

## Resources

- MCP Specification: https://modelcontextprotocol.io/
- Paperless-NGX Documentation: https://docs.paperless-ngx.com/
- Paperless-NGX API: https://docs.paperless-ngx.com/api/
