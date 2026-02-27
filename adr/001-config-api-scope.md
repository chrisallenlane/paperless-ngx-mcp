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

### 5. Unsupported API resources

The following Paperless-NGX API resources are intentionally not exposed through
the MCP server:

**Security-sensitive resources:**

- **users / groups** -- Admin permission management. An MCP tool that can
  delete users or deactivate TOTP would be a security risk. Better handled
  through the admin UI.
- **mail_accounts** -- Contains IMAP/SMTP credentials. Exposing passwords
  through an AI assistant is unacceptable.
- **mail_rules** -- Tightly coupled to mail_accounts. Complex rule
  configuration is better done through the UI.
- **profile** -- Manages auth tokens, TOTP, and social account connections.
  Security-sensitive operations that should not be delegated to an AI.
- **share_links** -- Creates publicly accessible URLs to documents. An MCP
  tool could accidentally expose sensitive documents.
- **oauth / token** -- Authentication flow endpoints. Not applicable to MCP.

**Complexity without conversational value:**

- **workflows / workflow_actions / workflow_triggers** -- Complex multi-step
  automation configuration. Hard to express conversationally, better through
  the UI.
- **saved view `display_fields`** -- Untyped field with no schema definition.
  Skipped for simplicity; other saved view fields are supported.

**Bulk/destructive operations:**

- **bulk_edit_objects** -- Bulk operations on multiple objects. Dangerous in an
  AI context where a misunderstood request could modify many records at once.
- **documents/bulk_edit** -- Same concern as bulk_edit_objects.
- **documents/bulk_download** -- Generates zip archives. Resource-intensive,
  better through the UI.
- **trash restore/empty** -- `POST /api/trash/` can permanently delete
  documents. The read-only `GET /api/trash/` (list) is supported, but write
  operations are too destructive for MCP.

**UI-specific or low-value endpoints:**

- **documents/preview / documents/thumb** -- Return rendered page images. MCP
  clients cannot meaningfully display binary image data.
- **documents/history** -- Audit log. Niche admin use case.
- **documents/email** -- Deprecated endpoint.
- **documents/selection_data** -- Internal UI helper for bulk selection
  widgets. No conversational use case.
- **documents/share_links** (per-document) -- Same security concern as the
  top-level share_links resource.
- **ui_settings** -- UI-specific preferences (theme, sidebar state). No value
  for MCP interactions.
- **remote_version** -- Checks for Paperless-NGX updates. Admin-only concern.
- **logs** -- Server log access. Admin diagnostic that could expose sensitive
  information.
- **processed_mail** -- Admin diagnostic for mail processing. Niche.

**Covered by other tools:**

- **search / search/autocomplete** -- The `list_documents` tool already
  supports a `search` parameter for full-text document search, which covers
  the primary use case. The global `/api/search/` endpoint returns results
  across all resource types (users, mail rules, workflows, etc.), most of
  which are noise in an MCP context. Can be added later if needed.
- **storage_paths/test** -- Testing endpoint for validating path templates.
  Not useful in an MCP context.

### 6. Task monitoring is read-only

The `/api/tasks/` resource exposes four endpoints: list, retrieve, acknowledge,
and run. We implement only list and retrieve (read-only). Acknowledging or
triggering tasks are administrative operations that don't fit conversational
workflows.

## Consequences

- Users cannot do full object replacement via the MCP server. If needed, they
  can PATCH all fields individually or use the Paperless-NGX UI.
- Users cannot delete application configuration via the MCP server. This is
  intentional.
- Users cannot manage object-level permissions via the MCP server.
- Users needing advanced filtering must use the Paperless-NGX UI or API
  directly.
- Users cannot manage users, groups, mail accounts, workflows, or share links
  via the MCP server. These require the Paperless-NGX admin UI.
- Users cannot perform bulk operations (bulk edit, bulk download) via the MCP
  server. Individual operations are available for all supported resources.
- Task monitoring is read-only. Users can check task status but cannot
  acknowledge or trigger tasks.
- Trash listing is read-only. Users can see deleted documents but cannot
  restore or permanently delete them via the MCP server.
- If future use cases require any of the unsupported endpoints, they can be
  added as new tools without changing existing tools.
