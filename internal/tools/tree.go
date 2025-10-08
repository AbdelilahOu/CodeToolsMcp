package tools

import (
	"context"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	TreeToolName        = "tree"
	TreeToolDescription = `Visualises a directory structure using an ASCII tree.`
)

type TreeInput struct {
	Path       string `json:"path" jsonschema:"required" jsonschema_description:"Absolute path to the directory root."`
	Depth      int    `json:"depth,omitempty" jsonschema_description:"Limit recursion depth (0 for unlimited)."`
	ShowHidden bool   `json:"show_hidden,omitempty" jsonschema_description:"Include dotfiles in the tree."`
	Limit      int    `json:"limit,omitempty" jsonschema_description:"Maximum number of nodes to display (0 for unlimited)."`
}

type TreeOutput struct {
	Tree string `json:"tree"`
}

func NewTreeTool(runner *runners.FileRunner) *ToolDefinition[TreeInput, TreeOutput] {
	return NewToolDefinition(
		TreeToolName,
		TreeToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input TreeInput) (*mcp.CallToolResult, TreeOutput, error) {
			result, err := runner.Tree(ctx, runners.TreeInput{
				Path:       input.Path,
				Depth:      input.Depth,
				ShowHidden: input.ShowHidden,
				Limit:      input.Limit,
			})
			if err != nil {
				return nil, TreeOutput{}, err
			}

			output := TreeOutput{Tree: result}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
