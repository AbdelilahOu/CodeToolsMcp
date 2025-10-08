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

### 6. git_status

Check the working tree for staged, unstaged, and untracked files.

**Parameters:**
- `short` (optional): When true, request git's `--short` flag. Output always uses porcelain format for easy parsing.

**Example:**
```json
{
  "short": true
}
```

### 7. git_log

Inspect commit history with common git filters.

**Parameters:**
- `oneline` (optional): Return commits in `--oneline` format.
- `limit` (optional): Maximum number of commits to include (maps to `-n`).
- `pattern` (optional): Filter commits whose messages match this regex (uses `--grep`).
- `since` (optional): Only include commits after this git-recognized date expression.
- `until` (optional): Only include commits before this git-recognized date expression.

**Example:**
```json
{
  "oneline": true,
  "limit": 5,
  "since": "1.week",
  "pattern": "bugfix"
}
```

### 8. git_diff

Compare commits, staged changes, or specific paths.

**Parameters:**
- `base` (optional): Left side revision. When omitted, diffs working tree/index.
- `target` (optional): Right side revision when diffing two commits.
- `paths` (optional): Array of file/directory paths to limit the diff.
- `staged` (optional): When true, include staged changes (`--staged`).
- `name_only` (optional): When true, show only filenames (`--name-only`).

**Example:**
```json
{
  "base": "main",
  "target": "feature",
  "name_only": true
}
```

### 9. git_show

Display commit metadata or file contents from a specific ref.

**Parameters:**
- `ref` (optional): Commit-ish or object to show (defaults to `HEAD`).
- `path` (optional): Limit output to a single file.
- `format` (optional): Custom pretty format string (`--format`).
- `name_only` (optional): Only list changed filenames.
- `stat` (optional): Include summary statistics (`--stat`).
- `no_patch` (optional): Suppress patch output.

**Example:**
```json
{
  "ref": "HEAD~1",
  "stat": true,
  "no_patch": true
}
```

### 10. git_branch

List branches with optional filters.

**Parameters:**
- `all` (optional): Show local and remote branches (`--all`).
- `remotes` (optional): Show only remote branches (`--remotes`).
- `contains` (optional): Only branches containing this commit.
- `sort` (optional): Sort key (e.g., `"committerdate"`).

**Example:**
```json
{
  "all": true,
  "sort": "-committerdate"
}
```

### 11. list_dir

List files and directories at a given location.

**Parameters:**
- `path` (required): Absolute path to the directory.
- `recursive` (optional): Walk subdirectories recursively.
- `show_hidden` (optional): Include entries starting with `.`.
- `limit` (optional): Maximum number of entries to return.

**Example:**
```json
{
  "path": "/path/to/project",
  "recursive": false
}
```

### 12. delete

Delete a single file.

**Parameters:**
- `path` (required): Absolute path to the file to delete.

**Example:**
```json
{
  "path": "/path/to/file.tmp"
}
```

### 13. remove

Remove files or directories (with optional recursion).

**Parameters:**
- `path` (required): Absolute path to remove.
- `recursive` (optional): Delete directories and contents (`rm -r`).

**Example:**
```json
{
  "path": "/path/to/build",
  "recursive": true
}
```

### 14. copy

Copy files or directories to a new location.

**Parameters:**
- `source` (required): Absolute path to copy from.
- `target` (required): Absolute path to copy to.
- `overwrite` (optional): Overwrite existing destination.

**Example:**
```json
{
  "source": "/path/to/config.json",
  "target": "/path/to/backup/config.json",
  "overwrite": true
}
```

### 15. move

Move or rename files/directories.

**Parameters:**
- `source` (required): Absolute path to move.
- `target` (required): Absolute destination path.
- `overwrite` (optional): Remove destination before moving.

**Example:**
```json
{
  "source": "/path/to/tmp.log",
  "target": "/path/to/archive/tmp.log"
}
```

### 16. tree

Render an ASCII tree of a directory.

**Parameters:**
- `path` (required): Absolute directory root.
- `depth` (optional): Maximum recursion depth (0 for unlimited).
- `show_hidden` (optional): Include dotfiles.
- `limit` (optional): Maximum number of nodes to display.

**Example:**
```json
{
  "path": "/path/to/project",
  "depth": 2
}
```

### 17. run

Execute shell commands and capture output.

**Parameters:**
- `command` (required): Executable or script to run.
- `args` (optional): Array of arguments.
- `working_dir` (optional): Directory to run the command in.
- `timeout_seconds` (optional): Cancel after this many seconds.
- `stdin` (optional): Data provided to standard input.
- `env` (optional): Additional environment variables.

**Example:**
```json
{
  "command": "bash",
  "args": ["-lc", "echo $MESSAGE"],
  "env": {"MESSAGE": "Hello"}
}
```

## Why This MCP?

Claude Code uses highly optimized tools like ripgrep for code search and specialized file operations for editing. This MCP provides the exact same functionality to any AI agent, enabling:

- **Fast code search**: ripgrep is orders of magnitude faster than grep
- **Precise editing**: Exact string matching prevents accidental modifications
- **Line-numbered reading**: Makes it easy for AI to reference specific lines
- **Pattern matching**: Powerful glob support for finding files
- **Git awareness**: Surface working tree state and commit history without leaving the agent
- **Filesystem control**: List, copy, move, remove, and inspect directories entirely through the MCP interface
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
│   ├── runners/         # External command runners (git, ripgrep, fs, process)
│   │   ├── git.go       # Git command helpers
│   │   ├── ripgrep.go   # Ripgrep command helper
│   │   ├── filesystem.go # Filesystem operations
│   │   └── command.go   # Shell command runner
│   ├── server/          # MCP server setup
│   └── tools/           # Tool implementations
│       ├── grep.go      # Grep tool
│       ├── glob.go      # Glob tool
│       ├── read.go      # Read tool
│       ├── edit.go      # Edit tool
│       ├── write.go     # Write tool
│       ├── gitLog.go   # Git log tool
│       ├── gitStatus.go  # Git status tool
│       ├── gitDiff.go  # Git diff tool
│       ├── gitShow.go  # Git show tool
│       ├── gitBranch.go # Git branch tool
│       ├── listDir.go  # Directory listing tool
│       ├── delete.go    # File delete tool
│       ├── remove.go    # Recursive remove tool
│       ├── copy.go      # Copy tool
│       ├── move.go      # Move tool
│       ├── tree.go      # Directory tree tool
│       ├── run.go       # Command execution tool
│       ├── tooldef.go   # Tool definition helper
│       └── tools.go     # Tool registration
└── config.json
```
