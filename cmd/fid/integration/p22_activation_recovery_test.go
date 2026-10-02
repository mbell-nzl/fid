//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestActivationPrepareResponseLossReplaysDurableResult(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22L")
	input := p21Input(t, workspace, "p22-lost-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22-lost", "phase": 21, "expected_prior_activation_id": nil, "expected_prior_version": nil, "resources": p21ActivationResources()})
	arguments := []string{"--project", workspace, "--json", "activation", "prepare", "--input", input, "--idempotency-key", "p22-lost-prepare"}
	p22RunWithClosedResponsePipe(t, workspace, binary, arguments...)
	replayed := p21Run(t, binary, workspace, arguments...)
	if replayedValue, ok := p21Path(t, replayed, "result", "replayed").(bool); !ok || !replayedValue || p21PathString(t, replayed, "result", "activation", "id") == "" {
		t.Fatalf("prepare retry after response loss = %#v, want durable replay with activation identity", replayed)
	}
}

func p22RunWithClosedResponsePipe(t *testing.T, workspace, binary string, arguments ...string) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, arguments...)
	command.Dir, command.Stdout, command.Stderr = workspace, writer, os.Stderr
	err = command.Run()
	if closeErr := writer.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	if err == nil {
		t.Fatalf("fid %s unexpectedly wrote a response after its reader closed", strings.Join(arguments, " "))
	}
}

// TestActivationRecoveryAtEveryDurableBoundary discards each mutation response
// and opens a fresh fid process to read and recover its durable identity.
func TestActivationRecoveryAtEveryDurableBoundary(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	for _, boundary := range []string{"prepared", "staged", "committed", "finalized"} {
		t.Run(boundary, func(t *testing.T) {
			workspace, activationID := p22ActivationBoundaryFixture(t, binary, boundary)
			sentinelPath := filepath.Join(workspace, "provider-sentinel.txt")
			sentinelBefore, err := os.ReadFile(sentinelPath)
			if err != nil {
				t.Fatal(err)
			}
			status := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "status", activationID)
			if p21PathString(t, status, "result", "activation", "id") != activationID || p21PathString(t, status, "result", "activation", "status") != boundary {
				t.Fatalf("fresh status after lost %s response = %#v, want retained activation %q", boundary, status, activationID)
			}
			if boundary == "finalized" {
				replayed := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "finalize", "--input", filepath.Join(workspace, "activation-finalize.json"), "--idempotency-key", "p22-boundary-finalize")
				if replayedValue, ok := p21Path(t, replayed, "result", "replayed").(bool); !ok || !replayedValue {
					t.Fatalf("finalized response-loss replay = %#v, want original durable result", replayed)
				}
				p22AssertUnchangedSentinel(t, sentinelPath, sentinelBefore)
				return
			}

			activation := p21ActivationStatus(t, workspace, activationID)
			direction := "restore-prior"
			if boundary == "committed" {
				wrong := p21Input(t, workspace, "p22-boundary-wrong-direction", map[string]any{"activation_id": activationID, "expected_version": activation.Version, "expected_status": boundary, "direction": "restore-prior", "reason": "must complete forward"})
				stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "activation", "recover", "--input", wrong, "--idempotency-key", "p22-boundary-wrong-direction")
				if processExitCode(t, err) != 4 || stdout != "" || !strings.Contains(stderr, `"code":"activation-direction-conflict"`) || !strings.Contains(stderr, `"required_direction":"complete-forward"`) {
					t.Fatalf("post-commit restore-prior = exit:%d stdout:%q stderr:%q, want typed complete-forward rejection", processExitCode(t, err), stdout, stderr)
				}
				direction = "complete-forward"
			}
			recoveryInput := p21Input(t, workspace, "p22-boundary-recover", map[string]any{"activation_id": activationID, "expected_version": activation.Version, "expected_status": boundary, "direction": direction, "reason": "resume after lost response"})
			selected := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "recover", "--input", recoveryInput, "--idempotency-key", "p22-boundary-recover")
			selectedID := p21PathString(t, selected, "result", "activation", "id")
			selectedEventID := p21PathString(t, selected, "event", "id")
			if selectedID != activationID || p21PathString(t, selected, "result", "activation", "recovery", "direction") != direction {
				t.Fatalf("recovery after %s = %#v, want retained ID and %s", boundary, selected, direction)
			}

			p22ObserveRecoveryResources(t, binary, workspace, activationID, direction)
			observed := p21ActivationStatus(t, workspace, activationID)
			versions := make([]any, 0, len(observed.Resources))
			for _, resource := range observed.Resources {
				versions = append(versions, map[string]any{"component": resource.Component, "version": resource.Version})
			}
			completeInput := p21Input(t, workspace, "p22-boundary-recover-complete", map[string]any{"activation_id": activationID, "expected_version": observed.Version, "selected_event_id": selectedEventID, "resource_versions": versions})
			completed := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "recover-complete", "--input", completeInput, "--idempotency-key", "p22-boundary-recover-complete")
			wantStatus := "aborted"
			if direction == "complete-forward" {
				wantStatus = "finalized"
			}
			if p21PathString(t, completed, "result", "activation", "id") != activationID || p21PathString(t, completed, "result", "activation", "status") != wantStatus {
				t.Fatalf("complete recovery after %s = %#v, want %s under the same ID", boundary, completed, wantStatus)
			}
			p22AssertUnchangedSentinel(t, sentinelPath, sentinelBefore)
		})
	}
}

func p22AssertUnchangedSentinel(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil || string(got) != string(want) {
		t.Fatalf("provider sentinel after activation recovery = %q, %v; want %q", got, err, want)
	}
}

func p22ActivationBoundaryFixture(t *testing.T, binary, boundary string) (string, string) {
	t.Helper()
	return p22ActivationBoundaryFixtureAt(t, binary, t.TempDir(), boundary)
}

func p22ActivationBoundaryFixtureAt(t *testing.T, binary, workspace, boundary string) (string, string) {
	t.Helper()
	p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22A")
	if err := os.WriteFile(filepath.Join(workspace, "provider-sentinel.txt"), []byte("provider bytes remain untouched\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	prepareInput := p21Input(t, workspace, "p22-boundary-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22-activation", "phase": 21, "expected_prior_activation_id": nil, "expected_prior_version": nil, "resources": p21ActivationResources()})
	prepared := p22LostActivationMutation(t, workspace, binary, "prepare", prepareInput, "p22-boundary-prepare")
	activationID := p21PathString(t, prepared, "result", "activation", "id")
	if boundary == "prepared" {
		return workspace, activationID
	}
	activation := p21ActivationStatus(t, workspace, activationID)
	profileInput := p21Input(t, workspace, "p22-boundary-profile", map[string]any{"activation_id": activationID, "component_type": "authority-profile", "component_id": "p22-boundary-profile", "component_version": 1, "expected_prior_digest": nil, "expected_version": activation.Version, "content": p22ActivationProfile()})
	profile := p22LostActivationMutation(t, workspace, binary, "component", profileInput, "p22-boundary-profile")
	profileDigest := p21PathString(t, profile, "result", "component", "digest")
	if boundary == "staged" {
		return workspace, activationID
	}
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p22-boundary-gates", p22ActivationGates())
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	activation = p21ActivationStatus(t, workspace, activationID)
	p21ObserveNewActivationResources(t, binary, workspace, activation)
	activation = p21ActivationStatus(t, workspace, activationID)
	commitInput := p21Input(t, workspace, "p22-boundary-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests, "expected_version": activation.Version})
	p22LostActivationMutation(t, workspace, binary, "commit", commitInput, "p22-boundary-commit")
	if boundary == "committed" {
		return workspace, activationID
	}
	activation = p21ActivationStatus(t, workspace, activationID)
	finalizeInput := p21Input(t, workspace, "activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activation.Version})
	p22LostActivationMutation(t, workspace, binary, "finalize", finalizeInput, "p22-boundary-finalize")
	return workspace, activationID
}

func p22LostActivationMutation(t *testing.T, workspace, binary, operation, input, key string) map[string]any {
	t.Helper()
	arguments := []string{"--project", workspace, "--json", "activation", operation, "--input", input, "--idempotency-key", key}
	p22RunWithClosedResponsePipe(t, workspace, binary, arguments...)
	replayed := p21Run(t, binary, workspace, arguments...)
	if replayedValue, ok := p21Path(t, replayed, "result", "replayed").(bool); !ok || !replayedValue {
		t.Fatalf("%s retry after response loss = %#v, want durable replay", operation, replayed)
	}
	return replayed
}

func p22ObserveRecoveryResources(t *testing.T, binary, workspace, activationID, direction string) {
	t.Helper()
	activation := p21ActivationStatus(t, workspace, activationID)
	state := "old"
	if direction == "complete-forward" {
		state = "new"
	}
	for _, resource := range activation.Resources {
		identity := resource.OldIdentity
		if state == "new" {
			identity = &resource.NewIdentity
		}
		input := p21Input(t, workspace, "p22-boundary-resource-"+resource.Component, map[string]any{"activation_id": activationID, "component": resource.Component, "expected_activation_version": activation.Version, "expected_version": resource.Version, "observed_identity": identity, "state": state, "evidence_ref": "integration:recovery:" + resource.Component, "error_code": nil})
		p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "resource", "--input", input, "--idempotency-key", "p22-boundary-resource-"+resource.Component)
		activation = p21ActivationStatus(t, workspace, activationID)
	}
}
