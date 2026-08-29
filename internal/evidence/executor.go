package evidence

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// ExecutionResult contains local process observations only. Receipt persistence
// and public output framing are separate concerns.
type ExecutionResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

// Execute runs one validated GateDefinition without a shell or containment layer.
func Execute(ctx context.Context, workspaceRoot string, definition GateDefinition) (ExecutionResult, error) {
	command, processContext, cancel, err := executorCommand(ctx, workspaceRoot, definition)
	if err != nil {
		return ExecutionResult{}, err
	}
	defer cancel()

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	result := ExecutionResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), ExitCode: -1}
	if command.ProcessState != nil {
		result.ExitCode = command.ProcessState.ExitCode()
	}
	if err != nil {
		if contextErr := processContext.Err(); contextErr != nil {
			return result, fmt.Errorf("execute gate process: %w", contextErr)
		}
		return result, fmt.Errorf("execute gate process: %w", err)
	}
	return result, nil
}

func executorCommand(ctx context.Context, workspaceRoot string, definition GateDefinition) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	if err := ValidateGateDefinition(workspaceRoot, definition); err != nil {
		return nil, nil, nil, err
	}
	workspace, err := filepath.EvalSymlinks(workspaceRoot)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("resolve executor workspace: %w", err)
	}
	cwd, err := filepath.EvalSymlinks(filepath.Join(workspace, definition.CWD))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("resolve executor cwd: %w", err)
	}
	environment, err := executorEnvironment(definition.Environment)
	if err != nil {
		return nil, nil, nil, err
	}
	processContext, cancel := context.WithTimeout(ctx, time.Duration(definition.TimeoutMS)*time.Millisecond)
	command := exec.CommandContext(processContext, definition.Argv[0], definition.Argv[1:]...)
	command.Dir = cwd
	command.Env = environment
	return command, processContext, cancel, nil
}

func executorEnvironment(entries []EnvironmentEntry) ([]string, error) {
	environment := make([]string, 0, len(entries))
	for _, entry := range entries {
		value, err := executorEnvironmentValue(entry)
		if err != nil {
			return nil, err
		}
		environment = append(environment, entry.Name+"="+value)
	}
	return environment, nil
}

func executorEnvironmentValue(entry EnvironmentEntry) (string, error) {
	switch entry.Source {
	case "literal":
		return *entry.Value, nil
	case "inherit":
		return os.Getenv(entry.Name), nil
	case "credential":
		source := "FID_CREDENTIAL_" + entry.Name
		value, found := os.LookupEnv(source)
		if !found {
			return "", fmt.Errorf("executor credential source %s is unavailable", strconv.Quote(source))
		}
		return value, nil
	default:
		return "", fmt.Errorf("executor environment source %q is unsupported", entry.Source)
	}
}
