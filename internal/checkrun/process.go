// Package checkrun executes policy checks without the kernel record chain.
package checkrun

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	return checkrecord.EnvironmentNames()
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
	output := &boundedOutput{stream: options.Output, failures: make(chan error, 1)}

	if err != nil {
		result.finishOutput(output)
		return result, nil
	}
	result.Toolchain.Digest, err = executableDigest(executable)
	if err != nil {
		result.finishOutput(output)
		return result, nil
	}
	command, launch, err := launchCommand(executable, options.Definition.Run)
	if err != nil {
		return result, err
	}
	defer launch.close()
	command.Dir, command.Env = options.Root, env
	command.Stdout, command.Stderr = output, output
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.WaitDelay = time.Second
	processContext, cancelProcess := context.WithTimeout(ctx, options.Definition.Timeout)
	defer cancelProcess()
	if err := command.Start(); err != nil {
		launch.closeChildEnds()
		result.finishOutput(output)
		return result, nil
	}
	launch.closeChildEnds()
	stop := func() { _ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL) }
	cancelStop := context.AfterFunc(processContext, stop)
	defer cancelStop()
	process, processErr := processIdentity(command.Process.Pid)
	if processErr == nil && options.Started != nil {
		processErr = options.Started(process)
	}
	spawnFailed := false
	if processErr == nil {
		spawnFailed, processErr = launch.release()
	}
	if processErr != nil {
		stop()
	}
	finished := make(chan error, 1)
	go func() { finished <- command.Wait() }()

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
		case processErr = <-output.failures:
			stop()
			waitErr = <-finished
			waiting = false
		case <-processContext.Done():
			result.Outcome = "timed-out"
			if errors.Is(processContext.Err(), context.Canceled) {
				result.Outcome = "cancelled"
			}
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
	select {
	case outputErr := <-output.failures:
		processErr = outputErr
	default:
	}
	if processErr != nil && processContext.Err() == nil {
		return result, fmt.Errorf("observe check process: %w", processErr)
	}
	if spawnFailed {
		result.Outcome = "spawn-failed"
	}
	if processContext.Err() != nil {
		result.Outcome = "timed-out"
		if errors.Is(processContext.Err(), context.Canceled) {
			result.Outcome = "cancelled"
		}
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
	if result.Outcome == "spawn-failed" {
		result.ExitCode = nil
		result.Signal = nil
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
	failures chan error
}

func (output *boundedOutput) Write(data []byte) (int, error) {
	output.mu.Lock()
	defer output.mu.Unlock()
	output.bytes += int64(len(data))
	kept := data[:min(len(data), checkrecord.OutputLimit-output.captured.Len())]
	_, _ = output.captured.Write(kept)
	if output.stream != nil {
		if _, err := output.stream.Write(data); err != nil {
			select {
			case output.failures <- err:
			default:
			}
			return 0, err
		}
	}
	return len(data), nil
}

func (result *ProcessResult) finishOutput(output *boundedOutput) {
	data := output.captured.Bytes()
	result.OutputBytes, result.OutputTruncated = output.bytes, output.bytes > checkrecord.OutputLimit
	result.OutputDigest = checkrecord.BytesDigest(data)
	excerpt := strings.ToValidUTF8(string(data[:min(len(data), checkrecord.ExcerptLimit)]), "�")
	if len(excerpt) > checkrecord.ExcerptLimit {
		excerpt = excerpt[:checkrecord.ExcerptLimit]
		excerpt = strings.ToValidUTF8(excerpt, "")
	}
	result.OutputExcerpt = excerpt
}

func executableDigest(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%x", hash.Sum(nil)), nil
}
