package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	WriteToolName        = "write"
	WriteToolDescription = `Writes a file to the local filesystem.

Usage:
- This tool will overwrite the existing file if there is one at the provided path.
- If this is an existing file, you MUST use the Read tool first to read the file's contents. This tool will fail if you did not read the file first.
- ALWAYS prefer editing existing files in the codebase. NEVER write new files unless explicitly required.`
)

type WriteInput struct {
	FilePath string `json:"file_path" jsonschema:"required" jsonschema_description:"The absolute path to the file to write (must be absolute, not relative)"`
	Content  string `json:"content" jsonschema:"required" jsonschema_description:"The content to write to the file"`
}

type WriteOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Size    int64  `json:"size"`
}

func NewWriteTool() *ToolDefinition[WriteInput, WriteOutput] {
	return NewToolDefinition(
		WriteToolName,
		WriteToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input WriteInput) (*mcp.CallToolResult, WriteOutput, error) {

			fileExists := false
			if _, err := os.Stat(input.FilePath); err == nil {
				fileExists = true
			}

			err := os.WriteFile(input.FilePath, []byte(input.Content), 0644)
			if err != nil {
				return nil, WriteOutput{}, fmt.Errorf("failed to write file: %w", err)
			}

			info, err := os.Stat(input.FilePath)
			if err != nil {
				return nil, WriteOutput{}, fmt.Errorf("failed to get file info: %w", err)
			}

			var message string
			if fileExists {
				message = fmt.Sprintf("Successfully overwrote file: %s (%d bytes)", input.FilePath, info.Size())
			} else {
				message = fmt.Sprintf("Successfully created file: %s (%d bytes)", input.FilePath, info.Size())
			}

			output := WriteOutput{
				Success: true,
				Message: message,
				Size:    info.Size(),
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: message},
				},
			}, output, nil
		},
	)
}
