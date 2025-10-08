package tools

import (
	"context"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GitShowToolName        = "git_show"
	GitShowToolDescription = `Displays commit or object details using "git show".

Usage:
- Specify a ref (commit, tag, etc.); defaults to HEAD when omitted
- Combine format, name_only, stat, and no_patch to shape the output
- Optionally provide a single path to inspect that file within the selected ref`
)

type GitShowInput struct {
	Ref      string `json:"ref,omitempty" jsonschema_description:"Commit-ish or object to display (defaults to HEAD)."`
	Path     string `json:"path,omitempty" jsonschema_description:"Optional file path to limit the output."`
	Format   string `json:"format,omitempty" jsonschema_description:"Custom pretty format (maps to git show --format)."`
	NameOnly bool   `json:"name_only,omitempty" jsonschema_description:"Show only file names that changed (git show --name-only)."`
	Stat     bool   `json:"stat,omitempty" jsonschema_description:"Include summary statistics (git show --stat)."`
	NoPatch  bool   `json:"no_patch,omitempty" jsonschema_description:"Suppress patch output (git show --no-patch)."`
}

type GitShowOutput struct {
	Content string `json:"content"`
}

func NewGitShowTool(runner *runners.GitRunner) *ToolDefinition[GitShowInput, GitShowOutput] {
	return NewToolDefinition(
		GitShowToolName,
		GitShowToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GitShowInput) (*mcp.CallToolResult, GitShowOutput, error) {
			result, err := runner.Show(ctx, runners.GitShowInput{
				Ref:      input.Ref,
				Path:     input.Path,
				Format:   input.Format,
				NameOnly: input.NameOnly,
				Stat:     input.Stat,
				NoPatch:  input.NoPatch,
			})
			if err != nil {
				return nil, GitShowOutput{}, err
			}

			output := GitShowOutput{Content: result}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, output, nil
		},
	)
}
