# ADR 001: Config API - PATCH Only, No DELETE

## Status

Accepted (2026-02-27)

## Context

The Paperless-NGX Config API (`/api/config/`) exposes five operations:

- `GET /api/config/` - List configurations
- `GET /api/config/{id}/` - Get a specific configuration
- `PUT /api/config/{id}/` - Full replacement update
- `PATCH /api/config/{id}/` - Partial update
- `DELETE /api/config/{id}/` - Delete a configuration

We needed to decide which of the write operations (PUT, PATCH, DELETE) to
expose as MCP tools.

## Decision

We will implement only `PATCH` for config updates. We will not implement `PUT`
or `DELETE`.

## Rationale

### PATCH over PUT

In an MCP tool context, users interact conversationally -- they ask to change
one or two settings at a time ("set the OCR language to eng", "enable barcode
scanning"). PATCH supports this naturally by accepting only the fields being
changed.

PUT requires sending the complete configuration object on every update. This
would force the tool to either:

1. Require the user to specify every field (impractical), or
2. Fetch the current config, merge the user's changes, and send the full
   object (fragile, adds complexity, risks race conditions)

PATCH avoids both problems. The request body contains only the fields being
modified.

### No DELETE

Deleting application configuration is a destructive operation with no clear
use case from an AI assistant. The Paperless-NGX instance typically has exactly
one configuration object (ID 1). Deleting it would reset all application
settings to defaults -- an action that is:

- Rarely intentional
- Difficult to reverse (all customizations are lost)
- Better performed through the Paperless-NGX admin UI if truly needed

Omitting DELETE reduces the risk of accidental misconfiguration through the
MCP server.

## Consequences

- Users cannot do full config replacement via the MCP server. If needed, they
  can PATCH all fields individually or use the Paperless-NGX UI.
- Users cannot delete configuration via the MCP server. This is intentional.
- If a future use case requires PUT or DELETE, they can be added as separate
  tools without changing the existing `update_config` tool.
