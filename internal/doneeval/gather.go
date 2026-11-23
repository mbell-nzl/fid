package doneeval

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

// Gather reads fid's evidence through a caller-owned projection snapshot. The
// caller refreshes it first and controls its transaction/lock, so completion
// can revalidate after network observations. Target comes from trusted local
// origin metadata. Remote is explicit: this function never contacts GitHub.
// It computes review/approval freshness locally, reads kept bytes, and writes
// no ledger records. Missing evidence stays in the snapshot for Evaluate.
func Gather(root string, conn reviewrecord.Queryer, task string, target Target, remote GitHubFacts) (Snapshot, error) {
	return gather(root, conn, task, target, remote, nil)
}

// Recorded rechecks historical bindings without git, kept files or networking.
// Present-evidence observations are immutable completion attestations.
func Recorded(conn reviewrecord.Queryer, task string, target Target, remote GitHubFacts, fp fingerprint.Result) (Snapshot, error) {
	return gather("", conn, task, target, remote, &fp)
}

func gather(root string, conn reviewrecord.Queryer, task string, target Target, remote GitHubFacts, fp *fingerprint.Result) (Snapshot, error) {
	s, err := reviewBase(conn, task)
	s.Target, s.GitHub = target, remote
	if err != nil || s.Handoff == nil {
		return s, err
	}
	if err := gatherLocal(root, conn, &s, fp); err != nil {
		return s, fmt.Errorf("gather done evidence: %w", err)
	}
	m, err := mergerecord.LatestMerge(conn, task)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return s, fmt.Errorf("gather done merge: %w", err)
	}
	if err == nil {
		s.Merge = &m
		s.CIResults, err = mergerecord.CIResults(conn, task, m.ID)
		if err != nil {
			return s, fmt.Errorf("gather done CI: %w", err)
		}
	}
	return s, nil
}

func reviewBase(conn reviewrecord.Queryer, task string) (Snapshot, error) {
	s := Snapshot{Approval: approvalrecord.View{State: "missing"}}
	var lane sql.NullString
	err := conn.QueryRow(`SELECT id,kind,status,lane,revision FROM issues WHERE id=?`, task).Scan(&s.Task.ID, &s.Task.Kind, &s.Task.State, &lane, &s.Task.Revision)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, nil
		}
		return s, fmt.Errorf("gather done task: %w", err)
	}
	s.Task.Lane = lane.String
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, nil
		}
		return s, fmt.Errorf("gather done hand-off: %w", err)
	}
	s.Handoff = &h
	var raw string
	err = conn.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=?`, h.PolicySnapshot).Scan(&raw)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return s, fmt.Errorf("gather done policy: %w", err)
	}
	if err == nil {
		var p policy.Effective
		if err := json.Unmarshal([]byte(raw), &p); err != nil {
			return s, fmt.Errorf("gather done policy: %w", err)
		}
		s.Policy = &p
	}
	return s, nil
}

func gatherLocal(root string, conn reviewrecord.Queryer, s *Snapshot, recorded *fingerprint.Result) error {
	var err error
	if err := gatherReviews(root, conn, s, recorded); err != nil {
		return err
	}
	s.Checks, err = checkread.Read(conn, s.Task.ID)
	if err != nil {
		return err
	}
	if s.Policy != nil {
		required := s.Handoff.RequiredChecks
		if s.Waiver.State == "current" {
			// Evaluate explains absent or invalid lane configuration as missing
			// evidence. Gathering must not bypass that single gate evaluator.
			required, _ = policy.Required(*s.Policy, s.EffectiveLane(), s.Task.ID, s.Handoff.Paths())
		}
		checks := checkread.Evaluate(s.Task.ID, s.Handoff.Tree, required, s.Policy.Checks, s.Checks)
		if recorded == nil {
			s.ArtifactProblems, err = checkread.VerifyKept(root, checks)
		}
		if err != nil {
			return err
		}
		s.ArtifactsVerified = true
	}
	return nil
}

func gatherReviews(root string, conn reviewrecord.Queryer, s *Snapshot, recorded *fingerprint.Result) error {
	var err error
	s.OpenMajors, err = reviewrecord.OpenMajors(conn, s.Task.ID)
	if err != nil {
		return err
	}
	rounds, err := reviewrecord.Read(conn, s.Task.ID)
	if err != nil {
		return err
	}
	approvals, err := approvalrecord.Read(conn, s.Task.ID)
	if err != nil {
		return err
	}
	waivers, err := waiverrecord.Read(conn, s.Task.ID)
	if err != nil {
		return err
	}
	fp := fingerprint.Result{}
	// Routine work with no review/approval evidence does not need a fingerprint
	// merely to explain missing gates, and may have an empty committed diff.
	if recorded != nil {
		fp = *recorded
	} else if len(rounds) > 0 || len(approvals) > 0 || len(waivers) > 0 {
		fp, err = fingerprint.Compute(root, s.Handoff.Base, s.Handoff.Head)
		if err != nil {
			return err
		}
	}
	s.Fingerprint = fp
	s.Rounds, err = reviewrecord.CurrentViews(conn, s.Task.ID, fp)
	if err != nil {
		return err
	}
	s.Approval, err = approvalrecord.Current(conn, s.Task.ID, fp)
	if err != nil {
		return err
	}
	s.Waiver, err = waiverrecord.Current(conn, s.Task.ID, fp)
	return err
}
