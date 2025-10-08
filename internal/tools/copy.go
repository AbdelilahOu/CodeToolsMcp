package tools

import (
	"context"
	"fmt"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	CopyToolName        = "copy"
	CopyToolDescription = `Copies files or directories.`
)

type CopyInput struct {
	Source    string `json:"source" jsonschema:"required" jsonschema_description:"Absolute path to the source file or directory."`
	Target    string `json:"target" jsonschema:"required" jsonschema_description:"Absolute destination path."`
	Overwrite bool   `json:"overwrite,omitempty" jsonschema_description:"Overwrite the destination if it already exists."`
}

type CopyOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewCopyTool(runner *runners.FileRunner) *ToolDefinition[CopyInput, CopyOutput] {
	return NewToolDefinition(
		CopyToolName,
		CopyToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input CopyInput) (*mcp.CallToolResult, CopyOutput, error) {
			if err := runner.Copy(ctx, runners.CopyInput{Source: input.Source, Target: input.Target, Overwrite: input.Overwrite}); err != nil {
				return nil, CopyOutput{}, err
			}

			message := fmt.Sprintf("Copied %s to %s", input.Source, input.Target)
			output := CopyOutput{Success: true, Message: message}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: message},
				},
			}, output, nil
		},
	)
}
