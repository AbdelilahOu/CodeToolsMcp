package runners

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type GitRunner struct{}

func NewGitRunner() *GitRunner {
	return &GitRunner{}
}

type GitStatusInput struct {
	Short bool
}

func (r *GitRunner) Status(ctx context.Context, input GitStatusInput) (string, error) {
	args := []string{"status"}

	if input.Short {
		args = append(args, "--short")
	}

	args = append(args, "--porcelain")

	cmd := exec.CommandContext(ctx, "git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 && stderr.Len() == 0 {
			return "No results available", nil
		}
		if stderr.Len() > 0 {
			return "", fmt.Errorf("git error: %s", stderr.String())
		}
		return "", fmt.Errorf("git command failed: %w", err)
	}

	result := strings.TrimSpace(stdout.String())

	if result == "" {
		return "Working tree clean", nil
	}

	return result, nil
}

type GitLogInput struct {
	Oneline bool
	Limit   int
	Pattern string
	Since   string
	Until   string
}

func (r *GitRunner) Log(ctx context.Context, input GitLogInput) (string, error) {
	args := []string{"log"}

	if input.Oneline {
		args = append(args, "--oneline")
	}

	if input.Limit > 0 {
		args = append(args, "-n", strconv.Itoa(input.Limit))
	}

	if input.Since != "" {
		args = append(args, "--since="+input.Since)
	}

	if input.Until != "" {
		args = append(args, "--until="+input.Until)
	}

	if input.Pattern != "" {
		args = append(args, "--grep="+input.Pattern)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 && stderr.Len() == 0 {
			return "No results available", nil
		}
		if stderr.Len() > 0 {
			return "", fmt.Errorf("git error: %s", stderr.String())
		}
		return "", fmt.Errorf("git command failed: %w", err)
	}

	result := strings.TrimSpace(stdout.String())

	if result == "" {
		return "No matching commits found", nil
	}

	return result, nil
}

type GitDiffInput struct {
	Base     string
	Target   string
	Paths    []string
	Staged   bool
	NameOnly bool
}

func (r *GitRunner) Diff(ctx context.Context, input GitDiffInput) (string, error) {
	args := []string{"diff"}

	if input.Staged {
		args = append(args, "--staged")
	}

	if input.NameOnly {
		args = append(args, "--name-only")
	}

	if input.Base != "" {
		args = append(args, input.Base)
		if input.Target != "" {
			args = append(args, input.Target)
		}
	} else if input.Target != "" {
		args = append(args, input.Target)
	}

	if len(input.Paths) > 0 {
		args = append(args, "--")
		args = append(args, input.Paths...)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("git error: %s", stderr.String())
		}
		return "", fmt.Errorf("git command failed: %w", err)
	}

	result := strings.TrimSpace(stdout.String())

	if result == "" {
		return "No differences found", nil
	}

	return result, nil
}

type GitShowInput struct {
	Ref      string
	Path     string
	Format   string
	NameOnly bool
	Stat     bool
	NoPatch  bool
}

func (r *GitRunner) Show(ctx context.Context, input GitShowInput) (string, error) {
	args := []string{"show"}

	if input.Format != "" {
		args = append(args, "--format="+input.Format)
	}

	if input.NameOnly {
		args = append(args, "--name-only")
	}

	if input.Stat {
		args = append(args, "--stat")
	}

	if input.NoPatch {
		args = append(args, "--no-patch")
	}

	ref := input.Ref
	if ref == "" {
		ref = "HEAD"
	}
	args = append(args, ref)

	if input.Path != "" {
		args = append(args, "--", input.Path)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("git error: %s", stderr.String())
		}
		return "", fmt.Errorf("git command failed: %w", err)
	}

	result := strings.TrimSpace(stdout.String())

	if result == "" {
		return "No content to display", nil
	}

	return result, nil
}

type GitBranchInput struct {
	All      bool
	Remotes  bool
	Contains string
	Sort     string
}

func (r *GitRunner) Branch(ctx context.Context, input GitBranchInput) (string, error) {
	args := []string{"branch", "--list"}

	if input.All {
		args = append(args, "--all")
	}

	if input.Remotes {
		args = append(args, "--remotes")
	}

	if input.Contains != "" {
		args = append(args, "--contains="+input.Contains)
	}

	if input.Sort != "" {
		args = append(args, "--sort="+input.Sort)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("git error: %s", stderr.String())
		}
		return "", fmt.Errorf("git command failed: %w", err)
	}

	result := strings.TrimSpace(stdout.String())

	if result == "" {
		return "No branches found", nil
	}

	return result, nil
}
