# MCP Client

A simple Go client for the Model Context Protocol (MCP). This client allows you to connect to an MCP server, initialize a session, list available tools, and call specific tools.

## Installation

```bash
git clone https://github.com/podanypepa/mcpclient.git
cd mcpclient
go build -o mcpclient ./cmd/mcpclient
```

## Usage

By default, the client connects to `http://127.0.0.1:9996/mcp`, initializes a session, and lists available tools.

```bash
./mcpclient
```

### Configuration Flags

- `-url`: MCP server URL (default: `http://127.0.0.1:9996/mcp`)
- `-token`: Optional Bearer token for Authorization header.
- `-timeout`: Per-request timeout (default: `10s`).

### Calling a Tool

To call a specific tool, use the `-tool` flag. You can provide arguments using `-args` (JSON string).

**Example: Calling the 'now' tool**

```bash
./mcpclient -tool now -args '{"format": "2006-01-02T15:04:05Z07:00"}'
```

**Example: Using a custom URL**

```bash
./mcpclient -url http://localhost:8080/mcp -tool mytool
```

## Project Structure

- `cmd/mcpclient`: Main entry point and CLI logic.
- `pkg/mcp`: Reusable MCP client library logic.

## License

See [LICENSE](LICENSE).
