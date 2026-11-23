// Package doneeval is the single source of task completion gates. Evaluate is
// pure over a gathered snapshot; Gather reads local evidence without networking.
package doneeval

import (
	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

type Status string

const (
	Complete Status = "complete"
	Waiting  Status = "waiting"
	Refused  Status = "refused"
)

type State string

const (
	Satisfied    State = "satisfied"
	Missing      State = "missing"
	Stale        State = "stale"
	WaitingState State = "waiting"
	RefusedState State = "refused"
)

type Gate string

const (
	TaskGate      Gate = "task"
	HandoffGate   Gate = "handoff"
	LaneGate      Gate = "lane"
	PolicyGate    Gate = "policy"
	ChecksGate    Gate = "checks"
	ArtifactsGate Gate = "artifacts"
	PrimaryGate   Gate = "primary-review"
	MajorsGate    Gate = "open-majors"
	SeparateGate  Gate = "separate-review"
	ApprovalGate  Gate = "approval"
	RollbackGate  Gate = "rollback"
	PRGate        Gate = "pull-request"
	TreeGate      Gate = "head-tree"
	MergeGate     Gate = "merge-observation"
	CIGate        Gate = "ci"
)

// Reasons are internal domain values. The done command maps them to its public
// refusal codes and human wording; they are not journal or CLI schemas.
type Reason string

const (
	EvidenceSatisfied       Reason = "evidence-satisfied"
	NotRequired             Reason = "not-required"
	PrerequisiteMissing     Reason = "prerequisite-missing"
	TaskUnavailable         Reason = "task-unavailable"
	TaskNotInReview         Reason = "task-not-in-review"
	HandoffMissing          Reason = "handoff-missing"
	HandoffMismatch         Reason = "handoff-mismatch"
	LaneChanged             Reason = "lane-changed"
	PolicyMissing           Reason = "policy-missing"
	PolicyMismatch          Reason = "policy-mismatch"
	ChecksNotConfigured     Reason = "checks-not-configured"
	TestPathsNotConfigured  Reason = "test-paths-not-configured"
	ChecksUnsatisfied       Reason = "checks-unsatisfied"
	KeptArtifactUnavailable Reason = "kept-artifact-unavailable"
	ArtifactsNotVerified    Reason = "artifacts-not-verified"
	ReviewMissing           Reason = "review-missing"
	ReviewStale             Reason = "review-stale"
	ProfileMissing          Reason = "profile-missing"
	SeparateSessionRequired Reason = "separate-session-required"
	OpenMajor               Reason = "open-major"
	ApprovalMissing         Reason = "approval-missing"
	ApprovalStale           Reason = "approval-stale"
	RollbackMissing         Reason = "rollback-missing"
	PRUnavailable           Reason = "pr-unavailable"
	PRMismatch              Reason = "pr-mismatch"
	RepositoryMismatch      Reason = "repository-mismatch"
	BranchMismatch          Reason = "branch-mismatch"
	NotMerged               Reason = "not-merged"
	MergeInvalid            Reason = "merge-invalid"
	HeadTreeMismatch        Reason = "head-tree-mismatch"
	MergeMissing            Reason = "merge-missing"
	MergeMismatch           Reason = "merge-mismatch"
	CIWaiting               Reason = "ci-waiting"
	CINotGreen              Reason = "ci-not-green"
	CIMismatch              Reason = "ci-mismatch"
	CIUnrecorded            Reason = "ci-unrecorded"
)

type GateResult struct {
	State  State
	Reason Reason
	IDs    []string
}
type Outcome struct {
	Status             Status
	Gates              map[Gate]GateResult
	NoOtherModelReview bool
}

// Task retains the revision needed by writers to revalidate this snapshot.
type Task struct {
	ID, Kind, State, Lane string
	Revision              int
}

// Target is trusted local origin identity, supplied by the caller that gathers
// GitHub observations. It is never inferred from the linked PR itself.
type Target struct {
	Repository    github.Repository
	DefaultBranch string
}

// GitHubFacts must come from one fresh, complete adapter observation. CI includes
// pending attempts (which are intentionally not durable). Body and head tree
// have no local record. A nil CI observation can never prove green.
type GitHubFacts struct {
	PR       *github.PullRequest
	HeadTree string
	CI       *github.CI
}

type Snapshot struct {
	Fingerprint       fingerprint.Result
	Task              Task
	Handoff           *handoffrecord.Record
	Policy            *policy.Effective
	Checks            []checkrecord.Record
	ArtifactsVerified bool
	ArtifactProblems  []checkread.ArtifactProblem
	Rounds            []reviewrecord.RoundView // journal order, freshness/other-model from shared reader
	OpenMajors        []reviewrecord.Finding   // task-wide, including stale rounds
	Waiver            waiverrecord.View
	Approval          approvalrecord.View
	Merge             *mergerecord.Merge
	CIResults         []mergerecord.CIResult
	Target            Target
	GitHub            GitHubFacts
}
