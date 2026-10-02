package evidence

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestExecuteAppliesAcceptedGatePolicy(t *testing.T) {
	workspace := t.TempDir()
	t.Setenv("FID_CREDENTIAL_GATE_TOKEN", "credential-value")
	t.Setenv("HOME", "inherited-home")
	literal := "literal-value"
	definition := GateDefinition{
		ArtifactBytes:     1024,
		ArtifactPaths:     []string{"coverage.out"},
		Argv:              []string{os.Args[0], "-test.run=^TestExecutorPolicyChild$", "--", "--executor-policy-child", "literal;$(not-a-shell)"},
		CredentialNames:   []string{"GATE_TOKEN"},
		CWD:               ".",
		DefinitionID:      "gate-policy-probe",
		DefinitionVersion: 1,
		Environment: []EnvironmentEntry{
			{Name: "GATE_TOKEN", Source: "credential"},
			{Name: "HOME", Source: "inherit"},
			{Name: "LITERAL", Source: "literal", Value: &literal},
		},
		ExcerptBytes:     128,
		Network:          "allowed",
		OutputBytes:      1024,
		RedactLiterals:   []string{"redaction-canary"},
		Sandbox:          "none",
		Shell:            false,
		SuccessExitCodes: []int{0},
		TimeoutMS:        5000,
		Write:            "allowed",
	}

	result, err := Execute(context.Background(), workspace, definition)
	if err != nil {
		t.Fatalf("Execute() error = %v; stderr = %q", err, result.Stderr)
	}
	var observed executorPolicyChildObservation
	if err := json.Unmarshal(result.Stdout, &observed); err != nil {
		t.Fatalf("decode child observation: %v; stdout = %q", err, result.Stdout)
	}
	resolvedWorkspace, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		t.Fatalf("resolve workspace: %v", err)
	}
	if observed.CWD != resolvedWorkspace {
		t.Fatalf("child cwd = %q, want %q", observed.CWD, resolvedWorkspace)
	}
	if !reflect.DeepEqual(observed.Argv, []string{"literal;$(not-a-shell)"}) {
		t.Fatalf("child argv = %#v, want structured literal argument", observed.Argv)
	}
	wantEnvironment := []string{"GATE_TOKEN=credential-value", "HOME=inherited-home", "LITERAL=literal-value"}
	if !reflect.DeepEqual(observed.Environment, wantEnvironment) {
		t.Fatalf("child environment = %#v, want %#v", observed.Environment, wantEnvironment)
	}
	if result.ExitCode != 0 {
		t.Fatalf("child exit code = %d, want 0", result.ExitCode)
	}
}

func TestExecutorPolicyChild(t *testing.T) {
	t.Parallel()
	marker := "--executor-policy-child"
	index := -1
	for position, arg := range os.Args {
		if arg == marker {
			index = position
			break
		}
	}
	if index < 0 {
		return
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	environment := os.Environ()
	sort.Strings(environment)
	observation := executorPolicyChildObservation{
		Argv:        append([]string(nil), os.Args[index+1:]...),
		CWD:         workingDirectory,
		Environment: environment,
	}
	if err := json.NewEncoder(os.Stdout).Encode(observation); err != nil {
		t.Fatalf("encode child observation: %v", err)
	}
	os.Exit(0)
}

func TestExecuteEnforcesTimeoutLimit(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	definition := GateDefinition{
		ArtifactBytes:     1,
		Argv:              []string{os.Args[0], "-test.run=^TestExecutorTimeoutChild$", "--", "--executor-timeout-child"},
		CWD:               ".",
		DefinitionID:      "gate-policy-timeout",
		DefinitionVersion: 1,
		Network:           "allowed",
		OutputBytes:       1,
		Sandbox:           "none",
		Shell:             false,
		SuccessExitCodes:  []int{0},
		TimeoutMS:         10,
		Write:             "allowed",
	}

	result, err := Execute(context.Background(), workspace, definition)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Execute() error = %v, want timeout; result = %#v", err, result)
	}
}

func TestExecutorTimeoutChild(t *testing.T) {
	t.Parallel()
	if !containsArgument(os.Args, "--executor-timeout-child") {
		return
	}
	time.Sleep(time.Second)
}

type executorPolicyChildObservation struct {
	Argv        []string `json:"argv"`
	CWD         string   `json:"cwd"`
	Environment []string `json:"environment"`
}

func containsArgument(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

func TestExecuteRejectsMissingCredentialBeforeSpawn(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	definition := GateDefinition{
		ArtifactBytes:     1,
		Argv:              []string{os.Args[0], "-test.run=^TestExecutorPolicyChild$", "--", "--executor-policy-child"},
		CredentialNames:   []string{"GATE_TOKEN"},
		CWD:               ".",
		DefinitionID:      "gate-policy-missing-credential",
		DefinitionVersion: 1,
		Environment:       []EnvironmentEntry{{Name: "GATE_TOKEN", Source: "credential"}},
		Network:           "allowed",
		OutputBytes:       1,
		Sandbox:           "none",
		Shell:             false,
		SuccessExitCodes:  []int{0},
		TimeoutMS:         1000,
		Write:             "allowed",
	}

	_, err := Execute(context.Background(), workspace, definition)
	if err == nil || !strings.Contains(err.Error(), "FID_CREDENTIAL_GATE_TOKEN") {
		t.Fatalf("Execute() error = %v, want missing credential source", err)
	}
}
