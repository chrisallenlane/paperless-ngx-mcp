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
│   └── tools/                 # MCP tool implementations
│       ├── tool.go            # Tool interface definition
│       ├── helpers.go         # Shared utility functions
│       ├── helpers_test.go    # Helper function tests
│       ├── get_status.go      # System status tool
│       └── get_status_test.go # Status tool tests
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
- `Put(ctx, path, body)` - PUT requests with JSON body
- `Delete(ctx, path)` - DELETE requests

**Testing Support:**
- `HTTPDoer` interface allows mocking HTTP requests
- `NewWithHTTPClient(baseURL, httpClient)` - Test constructor
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

**`doAPIRequest(ctx, client, path)`** - Common HTTP request pattern
**`ParseJSONResponse(body, target)`** - Type-safe JSON parsing

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

### 1. Create the tool file in `internal/tools/`

```go
package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

type MyTool struct {
    client *client.Client
}

func NewMyTool(c *client.Client) *MyTool {
    return &MyTool{client: c}
}

func (t *MyTool) Description() string {
    return "Brief description of what this tool does"
}

func (t *MyTool) InputSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "paramName": map[string]interface{}{
                "type":        "string",
                "description": "Parameter description",
            },
        },
        "required": []string{"paramName"},
    }
}

func (t *MyTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        ParamName string `json:"paramName"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", fmt.Errorf("failed to parse arguments: %w", err)
    }

    if params.ParamName == "" {
        return "", fmt.Errorf("paramName is required")
    }

    body, err := doAPIRequest(ctx, t.client, "/api/endpoint")
    if err != nil {
        return "", fmt.Errorf("API request failed: %w", err)
    }

    var result YourModel
    if err := ParseJSONResponse(body, &result); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    return fmt.Sprintf("Result: %v", result), nil
}
```

### 2. Register in `internal/server/server.go`

Add to `registerTools()`:
```go
s.tools["my_tool"] = tools.NewMyTool(s.client)
```

### 3. Write tests

Create `internal/tools/my_tool_test.go` with:
- Input validation tests
- Description and schema tests

### 4. Rebuild and test

```bash
make check
make build
```

## Code Quality Standards

### Input Validation
Always validate input before making API calls.

### Use Helper Functions
Prefer `doAPIRequest` and `ParseJSONResponse` over duplicating HTTP boilerplate.

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
- One tool per file
- Shared logic in helpers.go
- Type definitions in models.go

## Current Tools

### `get_status` (`internal/tools/get_status.go`)

- **Endpoint**: `GET /api/status/`
- **Input**: None (no parameters)
- **Output**: Human-readable formatted status summary
- **Model**: `models.SystemStatus` — parses version, storage, database, and task statuses

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
