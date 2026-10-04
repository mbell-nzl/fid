//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"
)

var p24ArtifactFields = map[string]string{"artifact.verify": "artifact_id"}

func TestP24ArtifactVerifyKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	const artifactName = "artifact.json"
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{artifactName}, p22FinalActions)
	p22Write(t, a.workspace, artifactName, "phase 24 artifact\n")
	p22AcceptedCall(t, a, "p24-artifact-budget", "run.observe-budget", map[string]any{"value": 0})
	gate := p23RunInput(t, a, "gate", "run", a.actor, "p24-artifact-gate", map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "repository": a.repository})
	if gate.ExitCode != 0 || gate.Stderr != "" {
		t.Fatalf("artifact gate = %d stderr=%q", gate.ExitCode, gate.Stderr)
	}
	terminal := p24GateTerminal(t, gate.Stdout)
	response := terminal["response"].(map[string]any)
	artifacts := p21Path(t, response, "result", "receipt", "artifacts").([]any)
	if len(artifacts) != 1 {
		t.Fatalf("gate artifacts = %#v, want one", artifacts)
	}
	artifact := artifacts[0].(map[string]any)
	digest := artifact["digest"].(string)
	if err := os.Remove(filepath.Join(a.workspace, ".fid", "artifacts", digest[len("sha256:"):])); err != nil {
		t.Fatal(err)
	}
	p24StartKeyedCapture(a.workspace)
	verified := p23RunInput(t, a, "artifact", "verify", a.actor, "p24-artifact-verify", map[string]any{"artifact_id": artifact["id"]})
	if verified.ExitCode != 0 || verified.Stderr != "" {
		t.Fatalf("artifact.verify = %d stderr=%q", verified.ExitCode, verified.Stderr)
	}
	entry, found := p24EndKeyedCapture(a.workspace)["artifact.verify"]
	if !found {
		t.Fatal("artifact.verify was not captured")
	}
	p24ExerciseCapturedLeaf(t, a.workspace, a.binary, "artifact.verify", p24ArtifactFields["artifact.verify"], entry)
	p24AssertFamilyConflict(t, a.workspace, a.binary, "artifact.verify", entry, []string{"artifact_id", "--actor"})
	p24AssertFamilyConflict(t, a.workspace, a.binary, "artifact.verify", entry, []string{"--actor"})
	p24AssertRebuildMapParity(t, a.workspace, a.binary)
}
