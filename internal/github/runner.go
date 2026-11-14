package github

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	DefaultTimeout     = 30 * time.Second
	DefaultOutputLimit = 8 << 20
)

var ErrOutputLimit = errors.New("gh output exceeded size limit")

// Runner receives gh's arguments (without argv[0]) and the child's environment.
// Nonzero exits and launch failures retain stdout and stderr separately.
type Runner func(context.Context, []string, []string) Result

type Result struct {
	ExitCode       int
	Stdout, Stderr []byte
	Err            error
}

// RunGH bounds the entire child lifetime, including inherited output pipes.
// Both streams share one memory limit. Exceeding it terminates the process;
// neither a partial response nor a successful exit can hide truncation.
func RunGH(ctx context.Context, args, environment []string) Result {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	output := &boundedOutput{limit: DefaultOutputLimit, cancel: cancel}
	command := exec.CommandContext(ctx, "gh", args...)
	command.Env = ghEnvironment(environment)
	command.WaitDelay = 100 * time.Millisecond
	command.Stdout = &outputStream{output: output}
	command.Stderr = &outputStream{output: output, stderr: true}
	err := command.Run()
	result := Result{Stdout: output.stdout.Bytes(), Stderr: output.stderr.Bytes(), Err: err}
	if command.ProcessState != nil {
		result.ExitCode = command.ProcessState.ExitCode()
	}
	if output.exceeded {
		result.Err = ErrOutputLimit
	} else if ctx.Err() != nil {
		result.Err = ctx.Err()
	}
	return result
}

// Keep credentials/config, but disable all documented prompt, paging, terminal,
// debug and update behaviors. Explicit host prevents inherited GH_HOST routing
// a GitHub.com receipt to a different service. No token value is logged.
func ghEnvironment(environment []string) []string {
	overrides := []string{"GH_PROMPT_DISABLED=1", "GH_PAGER=cat", "PAGER=cat", "NO_COLOR=1", "CLICOLOR=0",
		"GH_HOST=github.com", "GH_NO_UPDATE_NOTIFIER=1", "GH_NO_EXTENSION_UPDATE_NOTIFIER=1", "GH_TELEMETRY=0"}
	drop := map[string]bool{"GH_FORCE_TTY": true, "GH_DEBUG": true, "DEBUG": true, "CLICOLOR_FORCE": true, "GH_COLOR_LABELS": true}
	for _, entry := range overrides {
		key, _, _ := strings.Cut(entry, "=")
		drop[key] = true
	}
	result := make([]string, 0, len(environment)+len(overrides))
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		if !drop[key] {
			result = append(result, entry)
		}
	}
	return append(result, overrides...)
}

type boundedOutput struct {
	sync.Mutex
	stdout, stderr bytes.Buffer
	limit          int
	cancel         context.CancelFunc
	exceeded       bool
}

type outputStream struct {
	output *boundedOutput
	stderr bool
}

func (stream *outputStream) Write(data []byte) (int, error) {
	output := stream.output
	output.Lock()
	defer output.Unlock()
	available := output.limit - output.stdout.Len() - output.stderr.Len()
	count := min(len(data), available)
	destination := &output.stdout
	if stream.stderr {
		destination = &output.stderr
	}
	_, _ = destination.Write(data[:count])
	if count < len(data) {
		output.exceeded = true
		output.cancel()
		return count, ErrOutputLimit
	}
	return count, nil
}
