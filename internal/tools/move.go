package tools

import (
	"context"
	"fmt"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	MoveToolName        = "move"
	MoveToolDescription = `Moves or renames files and directories.`
)

type MoveInput struct {
	Source    string `json:"source" jsonschema:"required" jsonschema_description:"Absolute path to the source file or directory."`
	Target    string `json:"target" jsonschema:"required" jsonschema_description:"Absolute destination path."`
	Overwrite bool   `json:"overwrite,omitempty" jsonschema_description:"Overwrite the destination if it exists (performs rm before move)."`
}

type MoveOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewMoveTool(runner *runners.FileRunner) *ToolDefinition[MoveInput, MoveOutput] {
	return NewToolDefinition(
		MoveToolName,
		MoveToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input MoveInput) (*mcp.CallToolResult, MoveOutput, error) {
			if err := runner.Move(ctx, runners.MoveInput{Source: input.Source, Target: input.Target, Overwrite: input.Overwrite}); err != nil {
				return nil, MoveOutput{}, err
			}

			message := fmt.Sprintf("Moved %s to %s", input.Source, input.Target)
			output := MoveOutput{Success: true, Message: message}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: message},
				},
			}, output, nil
		},
	)
}
