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
