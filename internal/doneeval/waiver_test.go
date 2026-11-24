package doneeval

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

func TestWaiverControlsEveryLaneGate(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"standard", "routine"} {
		for _, state := range []string{"missing", "stale", "current"} {
			t.Run(lane+"/"+state, func(t *testing.T) {
				t.Parallel()
				s := completeSnapshot("assured", false)
				s.Policy.Lanes[lane] = config.Lane{ChecksSet: true, Checks: []string{}}
				s.Rounds = s.Rounds[:1]
				s.Approval.State = "missing"
				s.GitHub.PR.Body = ""
				s.Checks = nil
				// The normal Assured floor is waivable; the frozen hand-off stays Assured.
				assured := s.Policy.Lanes["assured"]
				assured.Paths = []string{"notes.md"}
				s.Policy.Lanes["assured"] = assured
				s.Waiver = waiverrecord.View{State: state, Decision: &waiverrecord.Record{Binding: waiverrecord.Binding{Lane: lane}}}
				got := Evaluate(s)
				if state == "current" {
					if got.Status != Complete {
						t.Fatalf("waived %s refused: %+v", lane, got)
					}
					for _, gate := range []Gate{SeparateGate, ApprovalGate, RollbackGate} {
						if got.Gates[gate].Reason != NotRequired {
							t.Fatalf("waived gate %s: %+v", gate, got.Gates[gate])
						}
					}
					if lane == "routine" && got.Gates[PrimaryGate].Reason != NotRequired {
						t.Fatal("routine still requires primary", got)
					}
					if s.Task.Lane != "assured" || s.Handoff.Lane != "assured" {
						t.Fatal("waiver changed recorded lanes")
					}
				} else if got.Status != Refused || got.Gates[SeparateGate].State != Missing || got.Gates[ApprovalGate].Reason != ApprovalMissing || got.Gates[ChecksGate].Reason != ChecksUnsatisfied {
					t.Fatalf("%s waiver relaxed gates: %+v", state, got)
				}
			})
		}
	}
}

func TestWaiverPreservesHardFloor(t *testing.T) {
	t.Parallel()
	for _, path := range []string{".fid.toml", ".github/workflows/ci.yml"} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", false)
			s.Handoff.ChangedPaths = []string{path}
			s.Policy.Lanes["standard"] = config.Lane{ChecksSet: true, Checks: []string{}}
			s.Waiver = waiverrecord.View{State: "current", Decision: &waiverrecord.Record{Binding: waiverrecord.Binding{Lane: "standard"}}}
			if got := Evaluate(s); got.Status != Refused || got.Gates[PolicyGate].Reason != PolicyMismatch {
				t.Fatalf("hard floor waived: %+v", got)
			}
		})
	}
}

func TestWaiverStillRequiresUniversalEvidence(t *testing.T) {
	t.Parallel()
	for _, tc := range []matrixCase{
		{"primary", func(s *Snapshot) { s.Rounds = nil }, PrimaryGate, ReviewMissing, Missing, Refused},
		{"major", func(s *Snapshot) { s.OpenMajors = []reviewrecord.Finding{{ID: "FND-one"}} }, MajorsGate, OpenMajor, RefusedState, Refused},
		{"merge", func(s *Snapshot) { s.Merge = nil }, MergeGate, MergeMissing, Missing, Refused},
		{"CI", func(s *Snapshot) { s.GitHub.CI.Runs = nil }, CIGate, CIWaiting, WaitingState, Waiting},
		{"new lane unset", func(s *Snapshot) { s.Policy.Lanes["standard"] = config.Lane{} }, PolicyGate, ChecksNotConfigured, Missing, Refused},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", false)
			s.Policy.Lanes["standard"] = config.Lane{ChecksSet: true, Checks: []string{}}
			s.Waiver = waiverrecord.View{State: "current", Decision: &waiverrecord.Record{Binding: waiverrecord.Binding{Lane: "standard"}}}
			tc.edit(&s)
			got := Evaluate(s)
			if got.Status != tc.outcome || got.Gates[tc.gate].Reason != tc.reason {
				t.Fatalf("waiver bypassed %s: %+v", tc.gate, got)
			}
		})
	}
}
