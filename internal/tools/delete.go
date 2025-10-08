package tools

import (
	"context"
	"fmt"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	DeleteToolName        = "delete"
	DeleteToolDescription = `Deletes a single file from the filesystem.`
)

type DeleteInput struct {
	Path string `json:"path" jsonschema:"required" jsonschema_description:"Absolute path to the file to delete."`
}

type DeleteOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewDeleteTool(runner *runners.FileRunner) *ToolDefinition[DeleteInput, DeleteOutput] {
	return NewToolDefinition(
		DeleteToolName,
		DeleteToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input DeleteInput) (*mcp.CallToolResult, DeleteOutput, error) {
			if err := runner.Delete(ctx, runners.DeleteInput{Path: input.Path}); err != nil {
				return nil, DeleteOutput{}, err
			}

			message := fmt.Sprintf("Deleted file %s", input.Path)
			output := DeleteOutput{Success: true, Message: message}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: message},
				},
			}, output, nil
		},
	)
}
