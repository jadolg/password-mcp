# Password Generator MCP Example

This project is a simple example of a password generator tool implemented in Go, using
the [mcp-golang](https://github.com/metoro-io/mcp-golang) library. It demonstrates how to register a tool that generates
passwords with customizable length, character set, and case.

## Usage

1. **Install**

```bash
go install github.com/jadolg/password-mcp@latest
```

2. **Get the binary path**

```bash
which password-mcp
```

### Stdio (default)

Configure your MCP client using the binary path from the previous step:

```json
{
  "mcpServers": {
    "password-mcp": {
      "command": "/home/myuser/go/bin/password-mcp",
      "args": [],
      "env": {}
    }
  }
}
```

### HTTP

Start the server with the `--http` flag:

```bash
password-mcp --http
```

By default it listens on `:8080`. Use `--addr` to change the address:

```bash
password-mcp --http --addr :9090
```

Then configure your MCP client:

```json
{
  "mcpServers": {
    "password-mcp": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

Adjust the URL if you used a custom `--addr`.

## Features
- Generate passwords of any length
- Choose from different character sets: letters, letters and numbers, or letters, numbers, and symbols
- Select password case: uppercase, lowercase, or mixed
