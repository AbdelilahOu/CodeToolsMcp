package tools

import (
	"context"
	"fmt"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	RemoveToolName        = "remove"
	RemoveToolDescription = `Removes files or directories. Set recursive: true to delete directories with contents.`
)

type RemoveInput struct {
	Path      string `json:"path" jsonschema:"required" jsonschema_description:"Absolute path to remove."`
	Recursive bool   `json:"recursive,omitempty" jsonschema_description:"When true, remove directories and their contents (equivalent to rm -r)."`
}

type RemoveOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewRemoveTool(runner *runners.FileRunner) *ToolDefinition[RemoveInput, RemoveOutput] {
	return NewToolDefinition(
		RemoveToolName,
		RemoveToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input RemoveInput) (*mcp.CallToolResult, RemoveOutput, error) {
			if err := runner.Remove(ctx, runners.RemoveInput{Path: input.Path, Recursive: input.Recursive}); err != nil {
				return nil, RemoveOutput{}, err
			}

			action := "Removed"
			if input.Recursive {
				action = "Recursively removed"
			}
			message := fmt.Sprintf("%s %s", action, input.Path)
			output := RemoveOutput{Success: true, Message: message}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: message},
				},
			}, output, nil
		},
	)
}
