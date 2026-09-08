//go:build integration

package integration

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"modernc.org/sqlite"
)

type p22ActivationCrashInput struct {
	Mode, Operation string
	Options         workflow.ActivationTransitionOptions
	Selection       workflow.ActivationRecoveryRequest
	Completion      workflow.ActivationRecoveryCompleteRequest
}

func TestP22ActivationRecoverySurvivesAppendProjectionInterruption(t *testing.T) {
	for _, direction := range []string{"restore-prior", "complete-forward"} {
		for _, operation := range []string{"recover", "recover-complete"} {
			for _, mode := range []string{"before-append", "after-append"} {
				t.Run(direction+"/"+operation+"/"+mode, func(t *testing.T) {
					p22ActivationFaultCase(t, direction, operation, mode)
				})
			}
		}
	}
}

func p22ActivationFaultCase(t *testing.T, direction, operation, mode string) {
	t.Helper()
	binary := integrationFIDBinary(t)
	workspace, prior, original := p22ActivationUpgradeFaultFixture(t, binary, direction)
	input := p22ActivationCrashInput{Mode: mode, Operation: operation, Options: workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), BootstrapActor: "fq7-bootstrap", IdempotencyKey: "fq7-" + operation}}
	input.Selection = workflow.ActivationRecoveryRequest{ActivationID: original.ID, ExpectedVersion: original.Version, ExpectedStatus: original.Status, Direction: direction, Reason: "recover interrupted installation"}
	if operation == "recover-complete" {
		selected := p22ActivationFaultCLI(t, binary, input.Options, "recover", input.Selection)
		p22ObserveRecoveryResources(t, binary, workspace, original.ID, direction)
		current := p21ActivationStatus(t, workspace, original.ID)
		input.Completion = workflow.ActivationRecoveryCompleteRequest{ActivationID: original.ID, ExpectedVersion: current.Version, SelectedEventID: p21PathString(t, selected, "event", "id")}
		for _, resource := range current.Resources {
			input.Completion.ResourceVersions = append(input.Completion.ResourceVersions, workflow.ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
		}
	}
	before, history := p22FaultProjection(t, workspace), p21JournalBytes(t, workspace)
	files := p22CompleteRepositorySnapshot(t, workspace)
	configBefore := readP21File(t, filepath.Join(workspace, ".fid", "config.toml"))
	p22InterruptActivation(t, input)
	if !reflect.DeepEqual(before, p22FaultProjection(t, workspace)) {
		t.Fatal("interruption partially projected activation transaction")
	}
	afterHistory := p21JournalBytes(t, workspace)
	if mode == "before-append" {
		if string(history) != string(afterHistory) {
			t.Fatal("dry-run failure appended an Event")
		}
	} else {
		if len(strings.Split(strings.TrimSpace(string(afterHistory)), "\n")) != len(strings.Split(strings.TrimSpace(string(history)), "\n"))+1 {
			t.Fatal("post-append fault lacks exactly one durable Event")
		}
	}
	conn := p22FaultDB(t, workspace)
	_, err := conn.Exec(`DROP TRIGGER p22_activation_fault`)
	p22FaultOK(t, err)
	p22FaultOK(t, conn.Close())
	// This actual fresh CLI read must perform catch-up before the retry.
	p21Run(t, binary, workspace, "--json", "activation", "status", original.ID)
	var request any = input.Selection
	if operation == "recover-complete" {
		request = input.Completion
	}
	result := p22ActivationFaultCLI(t, binary, input.Options, operation, request)
	if p22ResultPath(result, "result.replayed") != (mode == "after-append") {
		t.Fatal("wrong interrupted retry disposition")
	}
	if p21PathString(t, result, "result", "activation", "id") != original.ID {
		t.Fatal("recovery changed activation identity")
	}
	if operation == "recover-complete" && p22ResultPath(result, "result.activation.recovery.selected_event_id") != input.Completion.SelectedEventID {
		t.Fatal("completion changed immutable selected Event")
	}
	if mode == "after-append" {
		lines := strings.Split(strings.TrimSpace(string(afterHistory)), "\n")
		var durable map[string]any
		p22FaultOK(t, json.Unmarshal([]byte(lines[len(lines)-1]), &durable))
		if durable["event_id"] != p22ResultPath(result, "event.id") {
			t.Fatal("catch-up lost the appended Event identity")
		}
	}
	after, completedHistory := p22FaultProjection(t, workspace), p21JournalBytes(t, workspace)
	again := p22ActivationFaultCLI(t, binary, input.Options, operation, request)
	if p22ResultPath(again, "result.replayed") != true || !reflect.DeepEqual(p22ResultPath(again, "result.activation"), p22ResultPath(result, "result.activation")) || p22ResultPath(again, "event.id") != p22ResultPath(result, "event.id") {
		t.Fatal("retry lost original activation result/Event")
	}
	if !reflect.DeepEqual(after, p22FaultProjection(t, workspace)) || string(completedHistory) != string(p21JournalBytes(t, workspace)) {
		t.Fatal("exact replay changed durable facts")
	}
	if len(strings.Split(strings.TrimSpace(string(completedHistory)), "\n")) != len(strings.Split(strings.TrimSpace(string(history)), "\n"))+1 {
		t.Fatal("retry created missing/duplicate recovery Event")
	}
	p22AssertActivationFaultFacts(t, workspace, prior, original, operation, direction, before, result)
	p22RejectActivationFaultConflicts(t, binary, input, request)
	if !reflect.DeepEqual(files, p22CompleteRepositorySnapshot(t, workspace)) || string(configBefore) != string(readP21File(t, filepath.Join(workspace, ".fid", "config.toml"))) {
		t.Fatal("activation recovery edited provider/project files")
	}
	// Replay from the durable journal must reconstruct the same full projection.
	p21Run(t, binary, workspace, "--json", "rebuild")
	rebuilt := p22FaultProjection(t, workspace)
	delete(after, "projection_meta")
	delete(rebuilt, "projection_meta")
	if !reflect.DeepEqual(after, rebuilt) {
		t.Fatal("activation recovery rebuild differs from caught-up projection")
	}
}

func p22InterruptActivation(t *testing.T, input p22ActivationCrashInput) {
	t.Helper()
	encoded, err := json.Marshal(input)
	p22FaultOK(t, err)
	path := filepath.Join(t.TempDir(), "activation-crash.json")
	p22FaultOK(t, os.WriteFile(path, encoded, 0o600))
	cmd := exec.Command(os.Args[0], "-test.run=^TestP22ActivationCrashProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), "FID_P22_ACTIVATION_CRASH="+path, "GOCOVERDIR="+t.TempDir())
	output, err := cmd.CombinedOutput()
	var exited *exec.ExitError
	if !errors.As(err, &exited) || exited.ExitCode() != 23 || len(output) != 0 {
		t.Fatalf("interrupted activation child: %v %q", err, output)
	}
}

func TestP22ActivationCrashProcess(t *testing.T) {
	path := os.Getenv("FID_P22_ACTIVATION_CRASH")
	if path == "" {
		return
	}
	var input p22ActivationCrashInput
	encoded, err := os.ReadFile(path)
	p22FaultOK(t, err)
	p22FaultOK(t, json.Unmarshal(encoded, &input))
	failAt := int32(2)
	if input.Mode == "before-append" {
		failAt = 1
	}
	var calls atomic.Int32
	p22FaultOK(t, sqlite.RegisterScalarFunction("p22_activation_interruption", 0, func(*sqlite.FunctionContext, []driver.Value) (driver.Value, error) {
		if calls.Add(1) == failAt {
			return nil, fmt.Errorf("injected activation interruption")
		}
		return int64(0), nil
	}))
	operation := "INSERT"
	if input.Operation == "recover-complete" {
		operation = "UPDATE OF completed_event_id"
	}
	conn := p22FaultDB(t, input.Options.WorkspaceRoot)
	_, err = conn.Exec(`CREATE TRIGGER p22_activation_fault BEFORE ` + operation + ` ON activation_recoveries BEGIN SELECT p22_activation_interruption(); END`)
	p22FaultOK(t, err)
	p22FaultOK(t, conn.Close())
	if input.Operation == "recover" {
		_, err = workflow.RecoverActivation(input.Options, input.Selection)
	} else {
		_, err = workflow.CompleteActivationRecovery(input.Options, input.Completion)
	}
	if err == nil || !strings.Contains(err.Error(), "injected activation interruption") || calls.Load() != failAt {
		t.Fatalf("wrong activation fault: %d calls, %v", calls.Load(), err)
	}
	os.Exit(23)
}

func p22ActivationFaultCLI(t *testing.T, binary string, options workflow.ActivationTransitionOptions, operation string, request any) map[string]any {
	t.Helper()
	encoded, err := json.Marshal(request)
	p22FaultOK(t, err)
	path := filepath.Join(t.TempDir(), "request.json")
	p22FaultOK(t, os.WriteFile(path, encoded, 0o600))
	return p21Run(t, binary, options.WorkspaceRoot, "--json", "activation", operation, "--actor", options.BootstrapActor, "--input", path, "--idempotency-key", options.IdempotencyKey)
}

func p22ActivationUpgradeFaultFixture(t *testing.T, binary, direction string) (string, workflow.Activation, workflow.Activation) {
	t.Helper()
	workspace, priorID := p22ActivationBoundaryFixture(t, binary, "finalized")
	prior := p21ActivationStatus(t, workspace, priorID)
	resources := make([]any, 0, len(prior.Resources))
	for _, resource := range prior.Resources {
		resources = append(resources, map[string]any{"component": resource.Component, "old_identity": resource.NewIdentity, "new_identity": resource.NewIdentity + ":fq7"})
	}
	prepared := p22ActivationFaultCLI(t, binary, workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, IdempotencyKey: "fq7-prepare"}, "prepare", map[string]any{"backend_identity": "fid", "kit_identity": "fq7-upgrade", "direction": "activate", "phase": 21, "expected_prior_activation_id": prior.ID, "expected_prior_version": prior.Version, "resources": resources})
	id := p21PathString(t, prepared, "result", "activation", "id")
	profile, gates := p22ActivationProfile(), p22ActivationGates()
	profile[0].(map[string]any)["provenance_ref"] = "activation:fq7-new"
	digests := []string{p21Component(t, binary, workspace, id, "authority-profile", "fq7-profile", profile), p21Component(t, binary, workspace, id, "gate-definitions", "fq7-gates", gates)}
	sort.Strings(digests)
	if direction == "complete-forward" {
		p21Activation(t, binary, workspace, "commit", "fq7-commit", map[string]any{"activation_id": id, "expected_status": "staged", "component_digests": digests})
	}
	p22Write(t, workspace, "AGENTS.md", "unchanged project guidance\n")
	p22FaultOK(t, os.MkdirAll(filepath.Join(workspace, ".claude"), 0o700))
	p22Write(t, workspace, ".claude/settings.json", "{\"unchanged\":true}\n")
	return workspace, prior, p21ActivationStatus(t, workspace, id)
}
