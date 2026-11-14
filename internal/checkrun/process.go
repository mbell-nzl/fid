// Package checkrun executes policy checks without the kernel record chain.
package checkrun

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

type Definition = checkrecord.Definition

func EnvironmentNames() []string {
	return []string{"CI", "GOCACHE", "GOMODCACHE", "GOPATH", "HOME", "LANG", "LC_ALL", "PATH", "TMPDIR", "XDG_CACHE_HOME"}
}

func environment() []string {
	result := []string{}
	for _, name := range EnvironmentNames() {
		if value, ok := os.LookupEnv(name); ok {
			result = append(result, name+"="+value)
		}
	}
	return result
}

type ProcessOptions struct {
	Root          string
	Definition    Definition
	Environment   []string
	Version       string
	Output        io.Writer
	Started       func(checkrecord.Process) error
	Renew         func() error
	RenewInterval time.Duration
}

type ProcessResult struct {
	checkrecord.Result
	StopError error
}

func ExecuteProcess(ctx context.Context, options ProcessOptions) (ProcessResult, error) {
	result := ProcessResult{Result: checkrecord.Result{Outcome: "spawn-failed", Artifacts: []checkrecord.Artifact{}, RefusedArtifacts: []checkrecord.RefusedArtifact{}, Reasons: []string{}}}
	if len(options.Definition.Run) == 0 || options.Definition.Timeout <= 0 {
		return result, fmt.Errorf("execute check: invalid definition")
	}
	env := options.Environment
	if env == nil {
		env = environment()
	}
	executable, err := resolveExecutable(options.Root, options.Definition.Run[0], env)
	result.Toolchain = checkrecord.Toolchain{Executable: executable, OS: runtime.GOOS, Arch: runtime.GOARCH, FidVersion: options.Version}
	output := &boundedOutput{stream: options.Output}
	defer func() {}()
	if err != nil {
		result.finishOutput(output)
		return result, nil
	}
	contents, err := os.ReadFile(executable)
	if err != nil {
		result.finishOutput(output)
		return result, nil
	}
	result.Toolchain.Digest = checkrecord.BytesDigest(contents)
	command := exec.Command(executable, options.Definition.Run[1:]...)
	command.Dir, command.Env = options.Root, env
	command.Stdout, command.Stderr = output, output
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.WaitDelay = time.Second
	if err := command.Start(); err != nil {
		result.finishOutput(output)
		return result, nil
	}
	stop := func() { _ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL) }
	process, processErr := processIdentity(command.Process.Pid)
	if processErr == nil && options.Started != nil {
		processErr = options.Started(process)
	}
	if processErr != nil {
		stop()
	}
	finished := make(chan error, 1)
	go func() { finished <- command.Wait() }()
	timer := time.NewTimer(options.Definition.Timeout)
	defer timer.Stop()
	interval := options.RenewInterval
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	waitErr := error(nil)
	result.Outcome = "passed"
	waiting := true
	for waiting {
		select {
		case waitErr = <-finished:
			waiting = false
		case <-timer.C:
			result.Outcome = "timed-out"
			stop()
			waitErr = <-finished
			waiting = false
		case <-ctx.Done():
			result.Outcome = "cancelled"
			stop()
			waitErr = <-finished
			waiting = false
		case <-ticker.C:
			if options.Renew != nil {
				if err := options.Renew(); err != nil {
					result.StopError = err
					result.Outcome = "cancelled"
					stop()
					waitErr = <-finished
					waiting = false
				}
			}
		}
	}
	stop() // A descendant must not outlive its check or retain its pipes.
	if processErr != nil {
		return result, fmt.Errorf("record check process: %w", processErr)
	}
	if command.ProcessState != nil {
		code := command.ProcessState.ExitCode()
		result.ExitCode = &code
		if status, ok := command.ProcessState.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			signal := status.Signal().String()
			result.Signal, result.ExitCode = &signal, nil
			if result.Outcome == "passed" {
				result.Outcome = "signaled"
			}
		}
	}
	if waitErr != nil && result.Outcome == "passed" {
		var exit *exec.ExitError
		if !errors.As(waitErr, &exit) {
			return result, fmt.Errorf("wait for check: %w", waitErr)
		}
		result.Outcome = "failed"
	}
	result.finishOutput(output)
	return result, nil
}

func resolveExecutable(root, name string, env []string) (string, error) {
	paths := []string{name}
	if !strings.ContainsRune(name, filepath.Separator) {
		paths = nil
		for _, entry := range env {
			if value, ok := strings.CutPrefix(entry, "PATH="); ok {
				for _, dir := range filepath.SplitList(value) {
					paths = append(paths, filepath.Join(dir, name))
				}
			}
		}
	}
	for _, candidate := range paths {
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(root, candidate)
		}
		resolved, err := filepath.EvalSymlinks(candidate)
		if err != nil {
			continue
		}
		info, err := os.Stat(resolved)
		if err == nil && info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
			return resolved, nil
		}
	}
	return "", fmt.Errorf("check executable is unavailable")
}

type boundedOutput struct {
	mu       sync.Mutex
	bytes    int64
	captured bytes.Buffer
	stream   io.Writer
}

func (output *boundedOutput) Write(data []byte) (int, error) {
	output.mu.Lock()
	defer output.mu.Unlock()
	output.bytes += int64(len(data))
	kept := data[:min(len(data), checkrecord.OutputLimit-output.captured.Len())]
	_, _ = output.captured.Write(kept)
	if output.stream != nil {
		if _, err := output.stream.Write(kept); err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

func (result *ProcessResult) finishOutput(output *boundedOutput) {
	data := output.captured.Bytes()
	result.OutputBytes, result.OutputTruncated = output.bytes, output.bytes > checkrecord.OutputLimit
	result.OutputDigest = checkrecord.BytesDigest(data)
	result.OutputExcerpt = strings.ToValidUTF8(string(data[:min(len(data), checkrecord.ExcerptLimit)]), "�")
}
