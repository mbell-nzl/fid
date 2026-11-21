package doneeval

import (
	"slices"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// Evaluate returns every gate, including non-applicable gates. A hard evidence
// failure wins over waiting; waiting alone is retryable. It performs no I/O.
func Evaluate(s Snapshot) Outcome {
	o := Outcome{Status: Complete, Gates: map[Gate]GateResult{}}
	o.set(TaskGate, Satisfied, EvidenceSatisfied)
	if s.Task.ID == "" || s.Task.Kind != "task" {
		o.set(TaskGate, Missing, TaskUnavailable)
	} else if s.Task.State != "in_review" {
		o.set(TaskGate, RefusedState, TaskNotInReview)
	}
	o.set(HandoffGate, Satisfied, EvidenceSatisfied)
	if s.Handoff == nil {
		o.set(HandoffGate, Missing, HandoffMissing)
	} else if s.Handoff.TaskID != s.Task.ID {
		o.set(HandoffGate, RefusedState, HandoffMismatch)
	}
	for _, gate := range []Gate{LaneGate, PolicyGate, ChecksGate, ArtifactsGate, PrimaryGate, MajorsGate, SeparateGate, ApprovalGate, RollbackGate, PRGate, TreeGate, MergeGate, CIGate} {
		o.set(gate, Missing, PrerequisiteMissing)
	}
	if s.Handoff != nil {
		o.set(LaneGate, Satisfied, EvidenceSatisfied)
		if s.Task.Lane != s.Handoff.Lane {
			o.set(LaneGate, Stale, LaneChanged)
		}
		evaluateLocal(&o, s)
		evaluateReviews(&o, s)
		evaluateRemote(&o, s)
	}
	for _, g := range o.Gates {
		if g.State == WaitingState && o.Status == Complete {
			o.Status = Waiting
		}
		if g.State != Satisfied && g.State != WaitingState {
			o.Status = Refused
		}
	}
	return o
}

func (o *Outcome) set(g Gate, state State, reason Reason, ids ...string) {
	o.Gates[g] = GateResult{state, reason, ids}
}

func evaluateLocal(o *Outcome, s Snapshot) {
	h := s.Handoff
	if s.Policy == nil {
		o.set(PolicyGate, Missing, PolicyMissing)
		return
	}
	p := *s.Policy
	floor, floorErr := policy.ComputeFloor(p, h.Paths())
	if p.RulesVersion >= 3 && (floorErr != nil || floor.Raise(h.Lane) != h.Lane) {
		o.set(PolicyGate, RefusedState, PolicyMismatch)
		return
	}
	required, err := policy.Required(p, h.Lane, s.Task.ID, h.Paths())
	switch {
	case !p.Lanes[h.Lane].ChecksSet:
		o.set(PolicyGate, Missing, ChecksNotConfigured)
		return
	case !p.TestAdequacySet:
		o.set(PolicyGate, Missing, TestPathsNotConfigured)
		return
	case err != nil || !slices.Equal(required, h.RequiredChecks):
		o.set(PolicyGate, RefusedState, PolicyMismatch)
		return
	}
	for _, name := range required {
		definition, exists := p.Checks[name]
		if !exists || h.DefinitionDigests[name] != checkrecord.Digest(checkrecord.DefinitionFor(definition)) {
			o.set(PolicyGate, RefusedState, PolicyMismatch)
			return
		}
	}
	o.set(PolicyGate, Satisfied, EvidenceSatisfied)
	checks := checkread.Evaluate(s.Task.ID, h.Tree, h.RequiredChecks, p.Checks, s.Checks)
	o.set(ChecksGate, Satisfied, EvidenceSatisfied)
	if len(checks.Missing) > 0 {
		o.set(ChecksGate, Missing, ChecksUnsatisfied, checks.Missing...)
	}
	o.set(ArtifactsGate, Satisfied, EvidenceSatisfied)
	switch {
	case !s.ArtifactsVerified:
		o.set(ArtifactsGate, Missing, ArtifactsNotVerified)
	case len(s.ArtifactProblems) > 0:
		ids := []string{}
		for _, p := range s.ArtifactProblems {
			ids = append(ids, p.ArtifactID)
		}
		o.set(ArtifactsGate, Missing, KeptArtifactUnavailable, ids...)
	}
}

func evaluateReviews(o *Outcome, s Snapshot) {
	o.set(MajorsGate, Satisfied, EvidenceSatisfied)
	if len(s.OpenMajors) > 0 {
		ids := []string{}
		for _, f := range s.OpenMajors {
			ids = append(ids, f.ID)
		}
		o.set(MajorsGate, RefusedState, OpenMajor, ids...)
	}
	o.set(PrimaryGate, Satisfied, NotRequired)
	o.set(SeparateGate, Satisfied, NotRequired)
	o.set(ApprovalGate, Satisfied, NotRequired)
	o.set(RollbackGate, Satisfied, NotRequired)
	var profiles []string
	profilesKnown := false
	if s.Policy == nil {
		o.set(PrimaryGate, Missing, PrerequisiteMissing)
	} else {
		var err error
		profiles, err = reviewrecord.RequiredProfiles(*s.Handoff, *s.Policy)
		if err != nil {
			o.set(PrimaryGate, Missing, TestPathsNotConfigured)
		} else {
			profilesKnown = true
			if len(profiles) > 0 {
				o.Gates[PrimaryGate] = reviewGate(s, "primary", profiles)
			}
		}
	}
	if s.Handoff.Lane == "assured" {
		evaluateAssured(o, s, profiles, profilesKnown)
	}
}

func evaluateAssured(o *Outcome, s Snapshot, profiles []string, profilesKnown bool) {
	o.NoOtherModelReview = reviewrecord.NoOtherModelReview(s.Task.ID, s.Rounds)
	if profilesKnown {
		o.Gates[SeparateGate] = reviewGate(s, "separate", profiles)
	} else {
		o.set(SeparateGate, Missing, PrerequisiteMissing)
	}
	if o.Gates[SeparateGate].State == Satisfied {
		separate := o.Gates[SeparateGate].IDs[0]
		session := ""
		for _, r := range s.Rounds {
			if r.ID == separate {
				session = r.Reviewer.SessionID
			}
		}
		for _, r := range s.Rounds {
			if r.Stream == "primary" && r.TaskID == s.Task.ID && r.Reviewer.SessionID == session {
				o.set(SeparateGate, RefusedState, SeparateSessionRequired, separate)
			}
		}
		if session == "" {
			o.set(SeparateGate, RefusedState, SeparateSessionRequired, separate)
		}
	}
	switch s.Approval.State {
	case "current":
		o.set(ApprovalGate, Satisfied, EvidenceSatisfied)
	case "stale":
		o.set(ApprovalGate, Stale, ApprovalStale)
	default:
		o.set(ApprovalGate, Missing, ApprovalMissing)
	}
	if s.GitHub.PR != nil && hasRollback(s.GitHub.PR.Body) {
		o.set(RollbackGate, Satisfied, EvidenceSatisfied)
	} else {
		o.set(RollbackGate, Missing, RollbackMissing)
	}
}

func reviewGate(s Snapshot, stream string, profiles []string) GateResult {
	result := GateResult{State: Missing, Reason: ReviewMissing}
	for _, r := range s.Rounds {
		if r.TaskID != s.Task.ID || r.Stream != stream {
			continue
		}
		if !r.Current {
			if result.Reason == ReviewMissing {
				result = GateResult{State: Stale, Reason: ReviewStale}
			}
			continue
		}
		result = GateResult{State: Satisfied, Reason: EvidenceSatisfied, IDs: []string{r.ID}}
		for _, p := range profiles {
			if !slices.Contains(r.Profiles, p) {
				result.State, result.Reason = Missing, ProfileMissing
			}
		}
	}
	return result
}
