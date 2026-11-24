package doneeval

import (
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// ReviewWork names the next stream and all profiles one session must review.
type ReviewWork struct {
	Handoff struct {
		ID   string `json:"id"`
		PR   string `json:"pr"`
		Head string `json:"head"`
		Base string `json:"base"`
	} `json:"handoff"`
	Stream           string   `json:"stream"`
	RequiredProfiles []string `json:"required_profiles"`
}

// GatherReviews reads just the frozen hand-off policy and review/approval
// evidence. A recorded fingerprint makes admission independent of local git.
// Live reads also validate the comparison before the first required round,
// without network I/O. Routine work without required reviews needs no diff.
func GatherReviews(root string, conn reviewrecord.Queryer, task string, recorded *fingerprint.Result) (Snapshot, error) {
	s, err := reviewBase(conn, task)
	if err == nil && s.Handoff != nil {
		if recorded == nil && DueReview(s) != nil {
			fp, err := fingerprint.Compute(root, s.Handoff.Base, s.Handoff.Head)
			if err != nil {
				return s, err
			}
			recorded = &fp
		}
		err = gatherReviews(root, conn, &s, recorded)
	}
	return s, err
}

// DueReview uses exactly done's primary and separate gates. Rework is already
// in_progress and cannot create more review work until the next hand-off.
func DueReview(s Snapshot) *ReviewWork {
	if s.Task.State != "in_review" || s.Handoff == nil || s.Policy == nil {
		return nil
	}
	profiles, err := reviewrecord.RequiredProfiles(s.reviewHandoff(), *s.Policy)
	if err != nil {
		return nil
	}
	o := Outcome{Gates: map[Gate]GateResult{}}
	evaluateReviews(&o, s)
	stream := ""
	if o.Gates[PrimaryGate].State != Satisfied {
		stream = "primary"
	} else if o.Gates[SeparateGate].State != Satisfied {
		stream = "separate"
	}
	if stream == "" {
		return nil
	}
	w := &ReviewWork{Stream: stream, RequiredProfiles: profiles}
	w.Handoff.ID, w.Handoff.PR, w.Handoff.Head, w.Handoff.Base = s.Handoff.ID, s.Handoff.PR, s.Handoff.Head, s.Handoff.Base
	return w
}

func ReviewsDone(s Snapshot) bool {
	if s.Task.State != "in_review" || s.Handoff == nil {
		return false
	}
	o := Outcome{Gates: map[Gate]GateResult{}}
	evaluateReviews(&o, s)
	return o.Gates[PrimaryGate].State == Satisfied && o.Gates[SeparateGate].State == Satisfied && o.Gates[MajorsGate].State == Satisfied
}

// AwaitingMerge matches status's approve/merge eligibility. Recorded merges
// awaiting CI must never return to a builder, including after a squash merge.
func AwaitingMerge(conn reviewrecord.Queryer, s Snapshot) (bool, error) {
	if !ReviewsDone(s) {
		return false, nil
	}
	var n int
	if err := conn.QueryRow(`SELECT count(*) FROM merge_observations WHERE task_id=?`, s.Task.ID).Scan(&n); err != nil {
		return false, err
	}
	return n == 0, nil
}
