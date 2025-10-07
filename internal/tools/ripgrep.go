package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type RipgrepRunner struct{}

func NewRipgrepRunner() *RipgrepRunner {
	return &RipgrepRunner{}
}

type RipgrepSearchInput struct {
	Pattern         string
	Path            string
	Glob            string
	Type            string
	CaseInsensitive bool
	LineNumbers     bool
	ContextAfter    int
	ContextBefore   int
	Context         int
	OutputMode      string
	HeadLimit       int
	Multiline       bool
}

func (r *RipgrepRunner) Search(ctx context.Context, input RipgrepSearchInput) (string, error) {
	args := []string{}

	args = append(args, input.Pattern)

	if input.CaseInsensitive {
		args = append(args, "-i")
	}

	if input.Multiline {
		args = append(args, "-U", "--multiline-dotall")
	}

	switch input.OutputMode {
	case "files_with_matches":
		args = append(args, "-l")
	case "count":
		args = append(args, "-c")
	case "content":

		if input.LineNumbers {
			args = append(args, "-n")
		}

		if input.Context > 0 {
			args = append(args, fmt.Sprintf("-C%d", input.Context))
		} else {
			if input.ContextBefore > 0 {
				args = append(args, fmt.Sprintf("-B%d", input.ContextBefore))
			}
			if input.ContextAfter > 0 {
				args = append(args, fmt.Sprintf("-A%d", input.ContextAfter))
			}
		}
	default:
		return "", fmt.Errorf("invalid output_mode: %s", input.OutputMode)
	}

	if input.Type != "" {
		args = append(args, "-t", input.Type)
	}

	if input.Glob != "" {
		args = append(args, "--glob", input.Glob)
	}

	if input.Path != "" {
		args = append(args, input.Path)
	}

	cmd := exec.CommandContext(ctx, "rg", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "No matches found", nil
		}
		if stderr.Len() > 0 {
			return "", fmt.Errorf("ripgrep error: %s", stderr.String())
		}
		return "", fmt.Errorf("ripgrep command failed: %w", err)
	}

	result := stdout.String()

	if input.HeadLimit > 0 && result != "" {
		lines := strings.Split(result, "\n")
		if len(lines) > input.HeadLimit {
			lines = lines[:input.HeadLimit]
		}
		result = strings.Join(lines, "\n")
	}

	return strings.TrimSpace(result), nil
}
