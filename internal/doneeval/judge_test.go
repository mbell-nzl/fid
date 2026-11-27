package doneeval

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func TestShadowJudgeNeverChangesReviewGates(t *testing.T) {
	t.Parallel()
	for _, verdict := range []string{"MERGE", "HOLD"} {
		for _, missing := range []string{"none", "primary", "separate"} {
			t.Run(verdict+"/"+missing, func(t *testing.T) {
				t.Parallel()
				s := completeSnapshot("assured", true)
				if missing == "primary" {
					s.Rounds = s.Rounds[1:]
				}
				if missing == "separate" {
					s.Rounds = s.Rounds[:1]
				}
				before := Evaluate(s)
				s.Rounds = append(s.Rounds, reviewrecord.RoundView{Round: reviewrecord.Round{ID: "RND-judge", TaskID: s.Task.ID, HandoffID: s.Handoff.ID, Stream: "judge", Profiles: []string{"merge"}, Verdict: verdict, Reviewer: reviewrecord.Reviewer{SessionID: "judge"}, LaneRaise: &reviewrecord.LaneRaise{Lane: "assured", Reason: "Shadow recommendation"}}, Current: true})
				after := Evaluate(s)
				for _, gate := range []Gate{PrimaryGate, SeparateGate, MajorsGate, LaneGate} {
					if before.Gates[gate].State != after.Gates[gate].State || before.Gates[gate].Reason != after.Gates[gate].Reason {
						t.Fatalf("judge changed %s: %+v / %+v", gate, before.Gates[gate], after.Gates[gate])
					}
				}
			})
		}
	}
}
