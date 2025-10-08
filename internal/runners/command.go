package runners

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RunCommandInput struct {
	Command        string
	Args           []string
	WorkingDir     string
	TimeoutSeconds int
	Stdin          string
	Env            map[string]string
}

type RunCommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func RunCommand(ctx context.Context, input RunCommandInput) (RunCommandResult, error) {
	if input.Command == "" {
		return RunCommandResult{}, fmt.Errorf("command is required")
	}

	var cancel context.CancelFunc
	if input.TimeoutSeconds > 0 {
		timeout := time.Duration(input.TimeoutSeconds) * time.Second
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, input.Command, input.Args...)

	if input.WorkingDir != "" {
		cmd.Dir = input.WorkingDir
	}

	if len(input.Env) > 0 {
		env := os.Environ()
		for key, value := range input.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	if input.Stdin != "" {
		cmd.Stdin = strings.NewReader(input.Stdin)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	result := RunCommandResult{
		Stdout:   strings.TrimRight(stdoutBuf.String(), "\n"),
		Stderr:   strings.TrimRight(stderrBuf.String(), "\n"),
		ExitCode: 0,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			return result, nil
		}
		return RunCommandResult{}, fmt.Errorf("failed to run command: %w", err)
	}

	return result, nil
}
