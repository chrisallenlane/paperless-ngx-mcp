# Setup Guide

## Quick Start

1. **Build the server**:
   ```bash
   make build
   # Binary: dist/paperless-ngx-mcp
   ```

2. **Set environment variables**:
   ```bash
   export PAPERLESS_URL="https://paperless.example.com"
   export PAPERLESS_TOKEN="your-api-token"
   ```

3. **Test the server** (optional):
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | ./dist/paperless-ngx-mcp
   ```

## Configuration

### For Claude Code (CLI)

```bash
claude mcp add paperless-ngx /path/to/dist/paperless-ngx-mcp \
  -s user \
  -e PAPERLESS_URL=https://paperless.example.com \
  -e PAPERLESS_TOKEN=your-api-token
```

**Scope options:**
- `-s user` - Available in all projects (recommended)
- `-s local` - Private to current project only
- `-s project` - Save to `.mcp.json` for team sharing

**Verify configuration:**
```bash
claude mcp list
claude mcp get paperless-ngx
```

### For Claude Desktop

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "paperless-ngx": {
      "command": "/path/to/dist/paperless-ngx-mcp",
      "env": {
        "PAPERLESS_URL": "https://paperless.example.com",
        "PAPERLESS_TOKEN": "your-api-token"
      }
    }
  }
}
```

**Restart Claude Desktop** after updating the configuration.

## Troubleshooting

### Server not appearing in Claude Code

```bash
claude mcp list
claude mcp get paperless-ngx

# Try removing and re-adding
claude mcp remove paperless-ngx
claude mcp add paperless-ngx /path/to/dist/paperless-ngx-mcp -s user \
  -e PAPERLESS_URL=https://paperless.example.com \
  -e PAPERLESS_TOKEN=your-api-token
```

### Tools not working

1. Check environment variables are set correctly
2. Verify the binary has execute permissions: `chmod +x dist/paperless-ngx-mcp`
3. Test the server directly with stdin/stdout
4. Check Claude logs for errors

### Binary not found

Use the absolute path to the binary:

```bash
# Good
claude mcp add paperless-ngx /home/user/paperless-ngx-mcp/dist/paperless-ngx-mcp

# Bad (relative path may not work)
claude mcp add paperless-ngx ./dist/paperless-ngx-mcp
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `PAPERLESS_URL` | Yes | Base URL of your Paperless-NGX instance |
| `PAPERLESS_TOKEN` | Yes | API authentication token |

## Security Notes

- Store credentials in environment variables, not in code
- Use `claude mcp add` with env flags rather than hardcoding secrets
- The MCP server runs locally and communicates via stdio (no network exposure)
