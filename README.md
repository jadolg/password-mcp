# Password Generator MCP Example

This project is a simple example of a password generator tool implemented in Go, using the [mcp-golang](https://github.com/metoro-io/mcp-golang) library. It demonstrates how to register a tool that generates passwords with customizable length, character set, and case.

## Usage

1. **Clone the repository**

   ```bash
   git clone git@github.com:jadolg/password-mcp.git
   cd password-mcp
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Run the tool**

   ```bash
   go run main.go
   ```

The tool will start and wait for requests via the MCP protocol (using stdio transport). You can use an MCP-compatible client to interact with it.

## Features
- Generate passwords of any length
- Choose from different character sets: letters, letters and numbers, or letters, numbers, and symbols
- Select password case: uppercase, lowercase, or mixed
