package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func TestP24RepositoryHumanViews(t *testing.T) {
	t.Parallel()
	for _, dirty := range []bool{false, true} {
		for _, same := range []bool{false, true} {
			repo := workflow.Repository{ProjectID: "PRJ-one", WorktreeID: "WT-current", HeadOID: strings.Repeat("a", 40), TreeOID: strings.Repeat("b", 40), Dirty: dirty, CaptureSource: "git", DirtyDigest: stringPointer("sha256:private")}
			recorded := repo
			if !same {
				recorded.WorktreeID = "WT-recorded"
			}
			run := workflow.Run{ID: "RUN-one", StartRepository: recorded}
			root := clonePristineWorkspace(t)
			for _, view := range []string{"resume", "run.inspect", "review.record", "review.get", "review.list"} {
				var out bytes.Buffer
				cmd := &cobra.Command{}
				cmd.SetOut(&out)
				cmd.SetContext(context.WithValue(context.Background(), commandWorkingDirectoryKey{}, root))
				cmd.PersistentFlags().String("project", root, "")
				var err error
				switch view {
				case "resume":
					err = writeResumeHuman(cmd, workflow.Resume{Repository: repo, Run: &run}, workflow.RecoveryReadContext{})
				case "run.inspect":
					err = writeRunInspectHuman(cmd, workflow.RunInspection{Run: run, Reconciliation: workflow.Reconciliation{Repository: repo}}, workflow.RecoveryReadContext{})
				default:
					err = writeP23Human(cmd, view, map[string]any{"id": "REV-one", "issue_id": "ISS-one", "envelope_id": "ENV-one", "envelope_revision": 1, "run_id": run.ID, "actor_id": "reviewer", "profile": "default", "verdict": "pass", "independence": "distinct-actor", "repository": repo})
				}
				if err != nil {
					t.Fatal(err)
				}
				state := "clean"
				if dirty {
					state = "dirty"
				}
				text := out.String()
				if !strings.Contains(text, "HEAD aaaaaaaaaaaa; "+state) {
					t.Errorf("%s missing short repository: %s", view, text)
				}
				for _, forbidden := range []string{"worktree_id", "tree_oid", "dirty_digest", "capture_source", "<nil>", "WT-current", "WT-recorded", strings.Repeat("a", 40), strings.Repeat("b", 40), "sha256:private"} {
					if strings.Contains(text, forbidden) {
						t.Errorf("%s leaked %q", view, forbidden)
					}
				}
				if view == "resume" || view == "run.inspect" {
					label := "different worktree"
					if same {
						label = "same worktree"
					}
					if !strings.Contains(text, label) {
						t.Errorf("%s missing %s: %s", view, label, text)
					}
				}
			}
		}
	}
}
