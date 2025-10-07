# Code Tools MCP Server

A Model Context Protocol (MCP) server that provides the same powerful code manipulation tools that Claude Code uses internally. This enables any AI agent to have the same level of code search, reading, editing, and file operations capabilities.

## Features

The server exposes 5 core tools that mirror Claude Code's functionality:

1. **grep** - Fast regex-based code search using ripgrep
2. **glob** - File pattern matching with support for `**/*.ext` patterns
3. **read** - Read files with line numbers and optional range selection
4. **edit** - Perform exact string replacements in files
5. **write** - Create or overwrite files

## Installation

### Prerequisites

- Go 1.21 or higher
- [ripgrep](https://github.com/BurntSushi/ripgrep) installed and available in PATH

### Building

```bash
make build
```

## Configuration

Create a `config.json` file:

```json
{
  "logging": {
    "level": "INFO",
    "output_file": "code-tools-mcp.log",
    "max_size_mb": 10,
    "console": true
  }
}
```

### Configuration Options

- **logging.level**: Log level (DEBUG, INFO, WARN, ERROR)
- **logging.output_file**: Path to log file
- **logging.max_size_mb**: Maximum log file size in MB before rotation
- **logging.console**: Whether to also log to console

## Usage

### With Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "code-tools": {
      "command": "/path/to/code-tools-mcp",
      "args": ["stdio", "--config", "/path/to/config.json"]
    }
  }
}
```

### Command Line

```bash
# Run with config
./code-tools-mcp stdio --config config.json

# Run without config (uses defaults)
./code-tools-mcp stdio
```

## Tools

### 1. grep

Fast regex-based code search powered by ripgrep.

**Parameters:**
- `pattern` (required): Regular expression pattern to search for
- `path` (optional): Directory or file to search (defaults to current directory)
- `glob` (optional): Glob pattern to filter files (e.g., `"*.js"`, `"**/*.{ts,tsx}"`)
- `type` (optional): File type filter (e.g., `"js"`, `"py"`, `"rust"`, `"go"`)
- `-i` (optional): Case insensitive search
- `-n` (optional): Show line numbers (requires `output_mode: "content"`)
- `-A` (optional): Lines of context after match
- `-B` (optional): Lines of context before match
- `-C` (optional): Lines of context before and after match
- `output_mode` (optional): `"content"` (matching lines), `"files_with_matches"` (default), or `"count"`
- `head_limit` (optional): Limit output to first N results
- `multiline` (optional): Enable multiline mode for cross-line patterns

**Example:**
```json
{
  "pattern": "function.*Error",
  "glob": "**/*.js",
  "output_mode": "content",
  "-n": true,
  "-C": 2
}
```

### 2. glob

Fast file pattern matching sorted by modification time.

**Parameters:**
- `pattern` (required): Glob pattern (e.g., `"**/*.go"`, `"src/**/*.ts"`)
- `path` (optional): Directory to search in (defaults to current directory)

**Example:**
```json
{
  "pattern": "**/*.go",
  "path": "/path/to/project"
}
```

### 3. read

Read file contents with line numbers (cat -n format).

**Parameters:**
- `file_path` (required): Absolute path to file
- `offset` (optional): Line number to start reading from
- `limit` (optional): Number of lines to read (default: 2000)

**Example:**
```json
{
  "file_path": "/path/to/file.go",
  "offset": 100,
  "limit": 50
}
```

### 4. edit

Perform exact string replacements in files.

**Parameters:**
- `file_path` (required): Absolute path to file
- `old_string` (required): Text to replace (must be unique unless using `replace_all`)
- `new_string` (required): Replacement text
- `replace_all` (optional): Replace all occurrences (default: false)

**Example:**
```json
{
  "file_path": "/path/to/file.go",
  "old_string": "func OldName() {\n\treturn nil\n}",
  "new_string": "func NewName() {\n\treturn nil\n}"
}
```

### 5. write

Create new files or overwrite existing ones.

**Parameters:**
- `file_path` (required): Absolute path to file
- `content` (required): File content to write

**Example:**
```json
{
  "file_path": "/path/to/new-file.go",
  "content": "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}"
}
```

## Why This MCP?

Claude Code uses highly optimized tools like ripgrep for code search and specialized file operations for editing. This MCP provides the exact same functionality to any AI agent, enabling:

- **Fast code search**: ripgrep is orders of magnitude faster than grep
- **Precise editing**: Exact string matching prevents accidental modifications
- **Line-numbered reading**: Makes it easy for AI to reference specific lines
- **Pattern matching**: Powerful glob support for finding files
- **Same UX**: Other agents get the same capabilities Claude Code has

## Development

### Project Structure

```
.
├── cmd/
│   ├── main.go          # Entry point
│   └── root.go          # CLI command definitions
├── internal/
│   ├── config/          # Configuration loading
│   ├── logger/          # Logging utilities
│   ├── server/          # MCP server setup
│   └── tools/           # Tool implementations
│       ├── grep.go      # Grep tool
│       ├── glob.go      # Glob tool
│       ├── read.go      # Read tool
│       ├── edit.go      # Edit tool
│       ├── write.go     # Write tool
│       ├── ripgrep.go   # Ripgrep runner
│       ├── tooldef.go   # Tool definition helper
│       └── tools.go     # Tool registration
└── config.json
```
