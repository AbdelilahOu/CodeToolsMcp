# Code Tools MCP Server

A Model Context Protocol (MCP) server that provides the same powerful code manipulation tools that Claude Code uses internally. This enables any AI agent to have the same level of code search, reading, editing, and file operations capabilities.

## Features

The server exposes 17 core tools that mirror Claude Code's functionality:

1. **grep** - Fast regex-based code search using ripgrep
2. **glob** - File pattern matching with support for `**/*.ext` patterns
3. **read** - Read files with line numbers and optional range selection
4. **edit** - Perform exact string replacements in files
5. **write** - Create or overwrite files
6. **git_status** - Show working tree changes in porcelain format
7. **git_log** - Query commit history with filtering options
8. **git_diff** - Compare revisions, staged changes, or specific paths
9. **git_show** - Display commit details or object contents
10. **git_branch** - List branches with filters and sorting
11. **list_dir** - Enumerate directory contents
12. **delete** - Delete a single file
13. **remove** - Remove files or directories (supports recursive deletion)
14. **copy** - Copy files or directories
15. **move** - Move or rename files and directories
16. **tree** - Visualise directory structures in ASCII form
17. **run** - Execute shell commands and capture output

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
