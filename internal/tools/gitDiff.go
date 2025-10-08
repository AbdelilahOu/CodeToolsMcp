package tools

import (
	"context"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GitDiffToolName        = "git_diff"
	GitDiffToolDescription = `Shows differences between revisions using "git diff".

Usage:
- Compare the working tree, index, or specific revisions
- Provide base/target revisions to diff between commits or branches
- Set staged: true to diff staged changes; name_only: true to see only file names
- Supply paths to limit output to particular files or directories`
)

type GitDiffInput struct {
	Base     string   `json:"base,omitempty" jsonschema_description:"Base revision to diff against (left side)."`
	Target   string   `json:"target,omitempty" jsonschema_description:"Target revision to compare (right side)."`
	Paths    []string `json:"paths,omitempty" jsonschema_description:"Optional list of file or directory paths to limit the diff."`
	Staged   bool     `json:"staged,omitempty" jsonschema_description:"Include staged changes by passing --staged."`
	NameOnly bool     `json:"name_only,omitempty" jsonschema_description:"Show only file names that changed (git diff --name-only)."`
}

type GitDiffOutput struct {
	Diff string `json:"diff"`
}

func NewGitDiffTool(runner *runners.GitRunner) *ToolDefinition[GitDiffInput, GitDiffOutput] {
	return NewToolDefinition(
		GitDiffToolName,
		GitDiffToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GitDiffInput) (*mcp.CallToolResult, GitDiffOutput, error) {
			result, err := runner.Diff(ctx, runners.GitDiffInput{
				Base:     input.Base,
				Target:   input.Target,
				Paths:    input.Paths,
				Staged:   input.Staged,
				NameOnly: input.NameOnly,
			})
			if err != nil {
				return nil, GitDiffOutput{}, err
			}

			output := GitDiffOutput{Diff: result}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
