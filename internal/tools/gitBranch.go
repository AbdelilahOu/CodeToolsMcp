package tools

import (
	"context"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GitBranchToolName        = "git_branch"
	GitBranchToolDescription = `Lists branches using "git branch".

Usage:
- Toggle flags to list local, remote, or all branches
- Use contains to filter for branches containing a commit and sort to order results
- Output mirrors git branch --list so callers can parse it easily`
)

type GitBranchInput struct {
	All      bool   `json:"all,omitempty" jsonschema_description:"Include both local and remote branches (git branch --all)."`
	Remotes  bool   `json:"remotes,omitempty" jsonschema_description:"Show only remote branches (git branch --remotes)."`
	Contains string `json:"contains,omitempty" jsonschema_description:"Only list branches that contain the specified commit (git branch --contains)."`
	Sort     string `json:"sort,omitempty" jsonschema_description:"Order branches using git's sort keys (git branch --sort)."`
}

type GitBranchOutput struct {
	Branches string `json:"branches"`
}

func NewGitBranchTool(runner *runners.GitRunner) *ToolDefinition[GitBranchInput, GitBranchOutput] {
	return NewToolDefinition(
		GitBranchToolName,
		GitBranchToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GitBranchInput) (*mcp.CallToolResult, GitBranchOutput, error) {
			result, err := runner.Branch(ctx, runners.GitBranchInput{
				All:      input.All,
				Remotes:  input.Remotes,
				Contains: input.Contains,
				Sort:     input.Sort,
			})
			if err != nil {
				return nil, GitBranchOutput{}, err
			}

			output := GitBranchOutput{Branches: result}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
