package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ReadToolName        = "read"
	ReadToolDescription = `Reads a file from the local filesystem. You can access any file directly by using this tool.

Usage:
- The file_path parameter must be an absolute path, not a relative path
- By default, it reads up to 2000 lines starting from the beginning of the file
- You can optionally specify a line offset and limit (especially handy for long files), but it's recommended to read the whole file by not providing these parameters
- Any lines longer than 2000 characters will be truncated
- Results are returned using cat -n format, with line numbers starting at 1
- This tool can only read files, not directories. To read a directory, use glob or bash ls`
)

type ReadInput struct {
	FilePath string `json:"file_path" jsonschema:"required" jsonschema_description:"The absolute path to the file to read"`
	Offset   int    `json:"offset,omitempty" jsonschema_description:"The line number to start reading from. Only provide if the file is too large to read at once"`
	Limit    int    `json:"limit,omitempty" jsonschema_description:"The number of lines to read. Only provide if the file is too large to read at once."`
}

type ReadOutput struct {
	Content    string `json:"content"`
	LineCount  int    `json:"line_count"`
	TotalLines int    `json:"total_lines"`
}

func NewReadTool() *ToolDefinition[ReadInput, ReadOutput] {
	return NewToolDefinition(
		ReadToolName,
		ReadToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input ReadInput) (*mcp.CallToolResult, ReadOutput, error) {
			file, err := os.Open(input.FilePath)
			if err != nil {
				if os.IsNotExist(err) {
					return nil, ReadOutput{}, fmt.Errorf("file does not exist: %s", input.FilePath)
				}
				return nil, ReadOutput{}, fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			totalLines := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				totalLines++
			}
			if err := scanner.Err(); err != nil {
				return nil, ReadOutput{}, fmt.Errorf("failed to count lines: %w", err)
			}

			file.Seek(0, 0)

			offset := input.Offset
			if offset < 0 {
				offset = 0
			}
			limit := input.Limit
			if limit <= 0 {
				limit = 2000
			}

			scanner = bufio.NewScanner(file)
			var lines []string
			lineNum := 1

			for lineNum <= offset && scanner.Scan() {
				lineNum++
			}

			readCount := 0
			for scanner.Scan() && readCount < limit {
				line := scanner.Text()

				if len(line) > 2000 {
					line = line[:2000] + "... (truncated)"
				}

				lines = append(lines, fmt.Sprintf("%6d→%s", lineNum, line))
				lineNum++
				readCount++
			}

			if err := scanner.Err(); err != nil {
				return nil, ReadOutput{}, fmt.Errorf("failed to read file: %w", err)
			}

			content := strings.Join(lines, "\n")

			if len(lines) == 0 && totalLines == 0 {
				content = "(empty file)"
			}

			output := ReadOutput{
				Content:    content,
				LineCount:  len(lines),
				TotalLines: totalLines,
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: content},
				},
			}, output, nil
		},
	)
}
