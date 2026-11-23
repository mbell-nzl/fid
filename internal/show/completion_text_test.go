package show

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

var updateCompletionGoldens = flag.Bool("update", false, "rewrite completion goldens")

func completionMock(lane string, flag bool) Result {
	stamp := "2026-11-23T02:47:00Z"
	h := &handoffrecord.Record{View: handoffrecord.View{ID: "HOF-477dfd", PR: "https://github.com/owner/repo/pull/101", Head: strings.Repeat("1", 40), Base: strings.Repeat("0", 40), At: stamp, RequiredChecks: []string{}, ResultIDs: []string{}}, TaskID: "FID-aaaaaa", Lane: lane}
	fp := reviewrecord.Binding{Algorithm: "test", Fingerprint: strings.Repeat("2", 64), Base: h.Base, Head: h.Head}
	c := &Completion{Record: &completionrecord.Record{ID: "CMP-f1ab58", TaskID: h.TaskID, HandoffID: h.ID, MergeID: "MRG-aaaaaa", CIResultIDs: []string{"CIR-aaaaaa"}, ObservedCIResultIDs: []string{"CIR-aaaaaa"}, RoundIDs: []string{"RND-517555"}, Fingerprint: fp, NoOtherModelReview: flag, Actor: "agent:orchestrator", At: stamp}, Handoff: h,
		Merge:     &mergerecord.Merge{ID: "MRG-aaaaaa", TaskID: h.TaskID, HandoffID: h.ID, PR: h.PR, Number: 101, Repository: "owner/repo", BaseBranch: "main", MergeCommit: strings.Repeat("3", 40), HeadSHA: h.Head, MergedAt: "2026-11-23T02:46:00Z"},
		CIResults: []mergerecord.CIResult{{ID: "CIR-aaaaaa", RunID: 45102091114, HeadSHA: strings.Repeat("3", 40), Conclusion: "success", URL: "https://github.com/owner/repo/actions/runs/45102091114"}},
		Rounds:    []reviewrecord.Round{{ID: "RND-517555", Stream: "primary", Verdict: "APPROVE", Reviewer: reviewrecord.Reviewer{Actor: "agent:claude-reviewer", Model: "anthropic/claude-opus-5-5"}}}, Missing: []string{}}
	if lane == "assured" {
		c.ApprovalID, c.RollbackPresent = "PDN-c10c3e", true
		c.Approval = &approvalrecord.Record{ID: c.ApprovalID, Actor: "human:mark", At: "2026-11-23T02:44:00Z"}
		c.RoundIDs = append(c.RoundIDs, "RND-8e5e05")
		c.Rounds = append(c.Rounds, reviewrecord.Round{ID: "RND-8e5e05", Stream: "separate", Verdict: "APPROVE", Reviewer: reviewrecord.Reviewer{Actor: "agent:sol-reviewer", Model: "openai/gpt-6.1-sol"}})
	}
	return Result{IssueView: issues.IssueView{Issue: issues.IssueDetail{IssueSummary: issues.IssueSummary{ID: h.TaskID, Kind: "task", Title: "Completion receipt in fid show", Status: "done", StatusCategory: "completed", Lane: &lane}}}, Completion: c, Now: time.Date(2026, 11, 24, 15, 47, 0, 0, time.UTC)}
}

func TestCompletionGoldens(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"standard", "assured", "assured-flag"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			lane := strings.Split(name, "-")[0]
			r := completionMock(lane, name == "assured-flag")
			for _, ext := range []string{"txt", "json"} {
				got := writeText(t, r, display.Renderer{})
				assertRowsFitTheBox(t, got)
				if ext == "json" {
					raw, err := json.MarshalIndent(r.Fields()["completion"], "", "  ")
					if err != nil {
						t.Fatal(err)
					}
					got = string(raw) + "\n"
				}
				path := filepath.Join("testdata", "completed-"+name+"."+ext)
				if *updateCompletionGoldens {
					if err := os.WriteFile(path, []byte(got), 0o600); err != nil {
						t.Fatal(err)
					}
				}
				want, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if got != string(want) {
					t.Fatalf("%s differs:\n%s\nwant:\n%s", path, got, want)
				}
			}
		})
	}
}
