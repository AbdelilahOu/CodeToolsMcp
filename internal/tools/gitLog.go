package tools

import (
	"context"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GitLogToolName        = "git_log"
	GitLogToolDescription = `Reads commit history using "git log" with filters.

Usage:
- Supports typical filters: oneline output, limit (number of commits), grep pattern, since/until date expressions
- Dates accept any value supported by git, e.g. "1.week" or "2024-01-01"
- Combine with other tools to cross-reference file changes or prepare summaries`
)

type GitLogInput struct {
	Oneline bool   `json:"oneline,omitempty" jsonschema_description:"Return commits in git's --oneline format."`
	Limit   int    `json:"limit,omitempty" jsonschema_description:"Maximum number of commits to return (maps to git log -n)."`
	Pattern string `json:"pattern,omitempty" jsonschema_description:"Filter commits whose message matches this regex (git log --grep)."`
	Since   string `json:"since,omitempty" jsonschema_description:"Only show commits more recent than this expression (git log --since)."`
	Until   string `json:"until,omitempty" jsonschema_description:"Only show commits older than this expression (git log --until)."`
}

type GitLogOutput struct {
	Log string `json:"log"`
}

func NewGitLogTool(runner *runners.GitRunner) *ToolDefinition[GitLogInput, GitLogOutput] {
	return NewToolDefinition(
		GitLogToolName,
		GitLogToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GitLogInput) (*mcp.CallToolResult, GitLogOutput, error) {
			result, err := runner.Log(ctx, runners.GitLogInput{
				Oneline: input.Oneline,
				Limit:   input.Limit,
				Pattern: input.Pattern,
				Since:   input.Since,
				Until:   input.Until,
			})
			if err != nil {
				return nil, GitLogOutput{}, err
			}

			output := GitLogOutput{Log: result}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
