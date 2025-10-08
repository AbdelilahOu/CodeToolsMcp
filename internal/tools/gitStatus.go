package tools

import (
	"context"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GitStatusToolName        = "git_status"
	GitStatusToolDescription = `Shows the repository status using "git status".

Usage:
- Returns the same output as "git status --porcelain" so callers can parse it easily
- Set short: true to keep the compact summary (default behavior); false yields the same porcelain format
- Helpful for checking staged, unstaged, and untracked files before running other tools`
)

type GitStatusInput struct {
	Short bool `json:"short,omitempty" jsonschema_description:"When true include git's --short flag (porcelain output is always used for easy parsing)."`
}

type GitStatusOutput struct {
	Status string `json:"status"`
}

func NewGitStatusTool(runner *runners.GitRunner) *ToolDefinition[GitStatusInput, GitStatusOutput] {
	return NewToolDefinition(
		GitStatusToolName,
		GitStatusToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GitStatusInput) (*mcp.CallToolResult, GitStatusOutput, error) {
			result, err := runner.Status(ctx, runners.GitStatusInput{Short: input.Short})
			if err != nil {
				return nil, GitStatusOutput{}, err
			}

			output := GitStatusOutput{Status: result}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
