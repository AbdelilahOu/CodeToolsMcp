package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GrepToolName = "grep"
	GrepToolDescription = `A powerful search tool built on ripgrep

Usage:
- ALWAYS use Grep for search tasks. NEVER invoke 'grep' or 'rg' as a Bash command.
- Supports full regex syntax (e.g., "log.*Error", "function\\s+\\w+")
- Filter files with glob parameter (e.g., "*.js", "**/*.tsx") or type parameter (e.g., "js", "py", "rust")
- Output modes: "content" shows matching lines, "files_with_matches" shows only file paths (default), "count" shows match counts
- Pattern syntax: Uses ripgrep (not grep) - literal braces need escaping (use 'interface\\{\\}' to find 'interface{}' in Go code)
- Multiline matching: By default patterns match within single lines only. For cross-line patterns like 'struct \\{[\\s\\S]*?field', use multiline: true`
)

type GrepInput struct {
	Pattern         string `json:"pattern" jsonschema:"required" jsonschema_description:"The regular expression pattern to search for in file contents"`
	Path            string `json:"path,omitempty" jsonschema_description:"File or directory to search in (rg PATH). Defaults to current working directory."`
	Glob            string `json:"glob,omitempty" jsonschema_description:"Glob pattern to filter files (e.g. \"*.js\", \"*.{ts,tsx}\") - maps to rg --glob"`
	Type            string `json:"type,omitempty" jsonschema_description:"File type to search (rg --type). Common types: js, py, rust, go, java, etc."`
	CaseInsensitive bool   `json:"-i,omitempty" jsonschema_description:"Case insensitive search (rg -i)"`
	LineNumbers     bool   `json:"-n,omitempty" jsonschema_description:"Show line numbers in output (rg -n). Requires output_mode: \"content\", ignored otherwise."`
	ContextAfter    int    `json:"-A,omitempty" jsonschema_description:"Number of lines to show after each match (rg -A). Requires output_mode: \"content\", ignored otherwise."`
	ContextBefore   int    `json:"-B,omitempty" jsonschema_description:"Number of lines to show before each match (rg -B). Requires output_mode: \"content\", ignored otherwise."`
	Context         int    `json:"-C,omitempty" jsonschema_description:"Number of lines to show before and after each match (rg -C). Requires output_mode: \"content\", ignored otherwise."`
	OutputMode      string `json:"output_mode,omitempty" jsonschema_description:"Output mode: \"content\" shows matching lines (supports -A/-B/-C context, -n line numbers, head_limit), \"files_with_matches\" shows file paths (supports head_limit), \"count\" shows match counts (supports head_limit). Defaults to \"files_with_matches\"."`
	HeadLimit       int    `json:"head_limit,omitempty" jsonschema_description:"Limit output to first N lines/entries, equivalent to \"| head -N\". Works across all output modes."`
	Multiline       bool   `json:"multiline,omitempty" jsonschema_description:"Enable multiline mode where . matches newlines and patterns can span lines (rg -U --multiline-dotall). Default: false."`
}

type GrepOutput struct {
	Result string `json:"result"`
}

func NewGrepTool(runner *RipgrepRunner) *ToolDefinition[GrepInput, GrepOutput] {
	return NewToolDefinition(
		GrepToolName,
		GrepToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GrepInput) (*mcp.CallToolResult, GrepOutput, error) {
			if input.OutputMode == "" {
				input.OutputMode = "files_with_matches"
			}

			result, err := runner.Search(ctx, RipgrepSearchInput{
				Pattern:         input.Pattern,
				Path:            input.Path,
				Glob:            input.Glob,
				Type:            input.Type,
				CaseInsensitive: input.CaseInsensitive,
				LineNumbers:     input.LineNumbers,
				ContextAfter:    input.ContextAfter,
				ContextBefore:   input.ContextBefore,
				Context:         input.Context,
				OutputMode:      input.OutputMode,
				HeadLimit:       input.HeadLimit,
				Multiline:       input.Multiline,
			})

			if err != nil {
				return nil, GrepOutput{}, err
			}

			output := GrepOutput{Result: result}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
