//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestP22HumanGuidance(t *testing.T) {
	t.Run("inspection and imported observations", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapter(t)
		p22AssertGuidanceOutput(t, a, []string{"run", "inspect", a.runID}, "No observation recorded", "unknown")
		p22AcceptedCall(t, a, "guidance-source-observe", "run.observe", nil)
		p22AssertGuidanceOutput(t, a, []string{"run", "inspect", a.runID}, "Changed paths", "Reconciliation", "Liveness", "does not record")
		archive := filepath.Join(t.TempDir(), "history.zip")
		p21Run(t, a.binary, a.workspace, "--json", "export", "create", archive)
		target := filepath.Join(t.TempDir(), "destination")
		closureRun(t, filepath.Dir(target), "git", "clone", "-q", a.workspace, target)
		p21Run(t, a.binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "guidance-import")
		imported := *a
		imported.workspace = target
		before := p21JournalEventCount(t, target)
		p22AssertGuidanceOutput(t, &imported, []string{"resume", a.issueID}, "Imported project", a.projectID, target, "preserved", "Before continuing work here", "run.observe", "run.inspect", "does not grant mutation authority")
		p22AssertGuidanceOutput(t, &imported, []string{"run", "inspect", a.runID}, "Historical liveness", "fresh destination", "Recorded worktree", "Current worktree")
		if after := p21JournalEventCount(t, target); after != before {
			t.Errorf("Human reads appended Events: %d -> %d", before, after)
		}
		p22AcceptedCall(t, &imported, "guidance-destination-observe", "run.observe", nil)
		p22AssertGuidanceOutput(t, &imported, []string{"run", "inspect", a.runID}, "Destination liveness observation", "terminated")
		human, _, err := closureRunCommand(target, a.binary, "--project", target, "run", "inspect", a.runID)
		if err != nil || strings.Contains(human, "Historical liveness") {
			t.Errorf("fresh observation mislabeled historical: %v %s", err, human)
		}
		marker := filepath.Join(target, ".fid", "import.marker.json")
		if err := os.WriteFile(marker, []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
		human, _, err = closureRunCommand(target, a.binary, "--project", target, "resume", a.issueID)
		if err == nil || strings.Contains(human, "Imported project") {
			t.Error("invalid import binding was presented as verified history")
		}
	})
	t.Run("activation directions", func(t *testing.T) {
		binary := integrationFIDBinary(t)
		for _, status := range []string{"prepared", "committed"} {
			t.Run(status, func(t *testing.T) {
				workspace, id := p22ActivationBoundaryFixture(t, binary, status)
				a := &p22FIDRecoveryAdapter{binary: binary, workspace: workspace}
				wants := []string{"Activation", "Resources", "Git commit", "absent"}
				if status == "prepared" {
					wants = append(wants, "restore-prior", "before this activation")
				} else {
					wants = append(wants, "complete-forward", "irreversible commit", "finish installing")
				}
				p22AssertGuidanceOutput(t, a, []string{"activation", "status", id}, wants...)
			})
		}
	})
}

func p22AssertGuidanceOutput(t *testing.T, a *p22FIDRecoveryAdapter, args []string, wants ...string) {
	t.Helper()
	args = append([]string{"--project", a.workspace}, args...)
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
	if err != nil || stderr != "" {
		t.Fatalf("Human command %v: %v %s", args, err, stderr)
	}
	for _, want := range wants {
		if !strings.Contains(stdout, want) {
			t.Errorf("Human command %v missing %q:\n%s", args, want, stdout)
		}
	}
}
