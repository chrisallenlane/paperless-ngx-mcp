# ADR 001: API Endpoint Scope

## Status

Accepted (2026-02-27)

## Context

Paperless-NGX REST API resources typically expose up to six operations:

- `GET /api/{resource}/` - List (paginated)
- `GET /api/{resource}/{id}/` - Get by ID
- `POST /api/{resource}/` - Create
- `PUT /api/{resource}/{id}/` - Full replacement update
- `PATCH /api/{resource}/{id}/` - Partial update
- `DELETE /api/{resource}/{id}/` - Delete

We needed to decide which operations to expose as MCP tools, and how much of
each endpoint's functionality to surface.

## Decisions

### 1. PATCH over PUT (all resources)

We implement `PATCH` for updates. We do not implement `PUT`.

In an MCP tool context, users interact conversationally -- they ask to change
one or two fields at a time ("set the OCR language to eng", "rename that
correspondent to Acme Corp"). PATCH supports this naturally by accepting only
the fields being changed.

PUT requires sending the complete object on every update. This would force the
tool to either:

1. Require the user to specify every field (impractical), or
2. Fetch the current object, merge the user's changes, and send the full
   object (fragile, adds complexity, risks race conditions)

PATCH avoids both problems.

### 2. DELETE: config no, entity resources yes

**Config** (`/api/config/`): DELETE is not implemented. The Paperless-NGX
instance typically has exactly one configuration object (ID 1). Deleting it
would reset all application settings to defaults -- an action that is rarely
intentional, difficult to reverse, and better performed through the admin UI.

**Entity resources** (correspondents, custom fields, etc.): DELETE is
implemented. Unlike the singleton config, entity resources are user-created
objects with a normal lifecycle. Deleting a correspondent or custom field is a
legitimate operation that an AI assistant should support.

### 3. Simple list filtering

List endpoints expose only `page`, `page_size`, and a single `name` filter
(mapped to the API's `name__icontains`). Advanced query parameters
(`name__iexact`, `name__istartswith`, `name__iendswith`, `ordering`, `id__in`)
are not exposed.

This covers the majority of use cases without cluttering the tool schema. In a
conversational context, users typically search by name substring ("find
correspondents matching 'smith'") rather than using precise filter operators.

### 4. No permissions management

Entity resource tools do not expose `owner`, `set_permissions`, or
`full_perms` fields. Permissions management is a niche administrative concern
that is better handled through the Paperless-NGX UI.

## Consequences

- Users cannot do full object replacement via the MCP server. If needed, they
  can PATCH all fields individually or use the Paperless-NGX UI.
- Users cannot delete application configuration via the MCP server. This is
  intentional.
- Users cannot manage object-level permissions via the MCP server.
- Users needing advanced filtering must use the Paperless-NGX UI or API
  directly.
- If future use cases require PUT, advanced filters, or permissions management,
  they can be added as new tools or parameters without changing existing tools.
