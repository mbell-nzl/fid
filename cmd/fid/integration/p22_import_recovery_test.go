//go:build integration

package integration

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

// These are declared before the comparison. A common source-repository
// observation makes durable named-read state comparable without rewriting the
// imported Worktree ID. Archive paths, destination roots, and the private
// publication marker are the only environmental exclusions.
var p22ParityDimensions = []string{"resume", "run.inspect", "stop.list", "attention.list", "ordered-events"}
var p22ParityEnvironmentalNormalizations = []string{"archive-path", "destination-root", "publication-metadata"}

func TestRecoveryLiveRebuildImportParity(t *testing.T) {
	t.Parallel()
	adapter := p22NewFIDRecoveryAdapter(t)
	archiveDirectory := t.TempDir()
	archive := filepath.Join(archiveDirectory, "live.zip")
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "export", "create", archive)
	live := p22NamedRecoveryProjection(t, adapter, adapter.workspace, adapter.workspace)
	liveArchive := p22ArchiveEntries(t, archive)

	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "rebuild")
	rebuilt := p22NamedRecoveryProjection(t, adapter, adapter.workspace, adapter.workspace)
	rebuiltArchive := filepath.Join(archiveDirectory, "rebuilt.zip")
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "export", "create", rebuiltArchive)

	target := filepath.Join(t.TempDir(), "destination")
	closureRun(t, filepath.Dir(target), "git", "clone", "-q", adapter.workspace, target)
	imported := p21Run(t, adapter.binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "p22-import-parity")
	if got := p21PathString(t, imported, "result", "import", "project_id"); got != adapter.projectID {
		t.Fatalf("imported project ID = %q, want %q", got, adapter.projectID)
	}
	importedProjection := p22NamedRecoveryProjection(t, adapter, target, adapter.workspace)
	importedArchive := filepath.Join(archiveDirectory, "imported.zip")
	p21Run(t, adapter.binary, target, "--json", "export", "create", importedArchive)

	for _, named := range []struct {
		name  string
		value map[string]any
	}{{"rebuilt", rebuilt}, {"imported", importedProjection}} {
		if !reflect.DeepEqual(live, named.value) {
			t.Fatalf("live and %s named JSON projections differ across dimensions %v: live=%#v %s=%#v", named.name, p22ParityDimensions, live, named.name, named.value)
		}
	}
	for _, named := range []struct {
		name    string
		entries map[string][]byte
	}{{"rebuilt", p22ArchiveEntries(t, rebuiltArchive)}, {"imported", p22ArchiveEntries(t, importedArchive)}} {
		if !reflect.DeepEqual(liveArchive, named.entries) {
			t.Fatalf("live and %s canonical archive entries differ; exclusions=%v", named.name, p22ParityEnvironmentalNormalizations)
		}
	}

	// This is deliberately separate from parity: these are new local facts at
	// the destination, required before a successor can recover imported work.
	destination := *adapter
	destination.workspace = target
	destination.repository = p22CorpusRepository(t, target, destination.projectID)
	destination.p22Observe(t, "terminated")
	destination.p22Expire(t)
	recovery := destination.p22Recover(t)
	if recovery.ExitCode != 0 {
		t.Fatalf("fresh destination recovery = exit %d stdout=%q stderr=%q", recovery.ExitCode, recovery.Stdout, recovery.Stderr)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(recovery.Stdout), &envelope); err != nil {
		t.Fatalf("decode fresh destination recovery: %v", err)
	}
	if successor := p21PathString(t, envelope, "result", "run", "id"); successor == adapter.runID {
		t.Fatalf("destination recovery reused imported Run %q", successor)
	}
	destination.p22Advance(t, "run.recover", recovery.Stdout)
	budget := destination.Execute(t, p22ConformanceCase{Name: "p22-import-successor-budget", Operation: "run.observe-budget", Kind: "accepted"})
	if budget.ExitCode != 0 {
		t.Fatalf("destination successor budget observation = exit %d stdout=%q stderr=%q", budget.ExitCode, budget.Stdout, budget.Stderr)
	}
	gateInput := p21Input(t, target, "p22-import-successor-gate", map[string]any{
		"registration_id":      destination.registrationID,
		"registration_version": destination.registrationVersion,
		"obligation_id":        "OBL-p22-corpus",
		"run_id":               destination.runID,
		"expected_repository":  p24ExpectedRepository(destination.repository),
	})
	gate := p22GateRun(t, destination.binary, target, destination.actor, gateInput, "p22-import-successor-gate")
	if receipt := p21PathString(t, gate, "result", "receipt", "id"); receipt == "" {
		t.Fatal("destination successor Gate produced no Receipt")
	}
}

func p22NamedRecoveryProjection(t *testing.T, adapter *p22FIDRecoveryAdapter, workspace, repository string) map[string]any {
	t.Helper()
	resume := p21Run(t, adapter.binary, workspace, "--project", workspace, "--json", "resume", adapter.issueID, "--repository", repository)
	inspect := p21Run(t, adapter.binary, workspace, "--project", workspace, "--json", "run", "inspect", adapter.runID, "--repository", repository)
	stops := p21Run(t, adapter.binary, workspace, "--project", workspace, "--json", "stop", "list", "--subject", "run:"+adapter.runID, "--limit", "10")
	attention := p21Run(t, adapter.binary, workspace, "--project", workspace, "--json", "attention", "list", "--subject", "run:"+adapter.runID, "--limit", "10")
	return map[string]any{
		"resume":         p21Path(t, resume, "result", "resume"),
		"run.inspect":    p21Path(t, inspect, "result"),
		"stop.list":      p21Path(t, stops, "result"),
		"attention.list": p21Path(t, attention, "result"),
	}
}
