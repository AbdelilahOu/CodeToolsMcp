package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	RunToolName        = "run"
	RunToolDescription = `Executes shell commands or scripts. Use with caution; output includes stdout, stderr, and exit code.`
)

type RunInput struct {
	Command        string            `json:"command" jsonschema:"required" jsonschema_description:"Executable or script to run."`
	Args           []string          `json:"args,omitempty" jsonschema_description:"Arguments to pass to the command."`
	WorkingDir     string            `json:"working_dir,omitempty" jsonschema_description:"Directory to run the command in."`
	TimeoutSeconds int               `json:"timeout_seconds,omitempty" jsonschema_description:"Timeout in seconds before the command is cancelled."`
	Stdin          string            `json:"stdin,omitempty" jsonschema_description:"Optional standard input data."`
	Env            map[string]string `json:"env,omitempty" jsonschema_description:"Additional environment variables (KEY: VALUE)."`
}

type RunOutput struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

func NewRunTool() *ToolDefinition[RunInput, RunOutput] {
	return NewToolDefinition(
		RunToolName,
		RunToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input RunInput) (*mcp.CallToolResult, RunOutput, error) {
			result, err := runners.RunCommand(ctx, runners.RunCommandInput{
				Command:        input.Command,
				Args:           input.Args,
				WorkingDir:     input.WorkingDir,
				TimeoutSeconds: input.TimeoutSeconds,
				Stdin:          input.Stdin,
				Env:            input.Env,
			})
			if err != nil {
				return nil, RunOutput{}, err
			}

			summary := formatRunSummary(result)

			output := RunOutput{
				Stdout:   result.Stdout,
				Stderr:   result.Stderr,
				ExitCode: result.ExitCode,
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: summary},
				},
			}, output, nil
		},
	)
}

func formatRunSummary(result runners.RunCommandResult) string {
	stdout := result.Stdout
	if stdout == "" {
		stdout = "(empty)"
	}

	stderr := result.Stderr
	if stderr == "" {
		stderr = "(empty)"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Exit code: %d\n", result.ExitCode))
	builder.WriteString("Stdout:\n")
	builder.WriteString(stdout)
	builder.WriteString("\nStderr:\n")
	builder.WriteString(stderr)

	return builder.String()
}
