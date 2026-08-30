# Phase 21 - Mid-Phase Review 6b (tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-08-30
Tasks covered: phase-21/task-16, phase-21/task-17, phase-21/task-18
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: `project-docs/docs/phases/phase-21/qa/review-6.md`

## Verdict

REPAIR

## Scope Reviewed

- Task specs: `task-16-artifacts.md`, `task-17-evidence-evaluation.md`, and
  `task-18-completion-resume.md`.
- Completion reports: `task-16-done.md`, `task-17-done.md`, and
  `task-18-done.md`.
- Validation contract: VAL-055 through VAL-068; QA plan Window 6, QA-53
  through QA-66.
- Prior review: `qa/review-6.md`; DEC-172 and closed ESC-21-013 boundary.
- Repair commits and reports: `457366d`, `200b807`, `1a8ecaa`, `bec1d6f`,
  `9847c23`, and `e095eda`; each full reported hash resolves as a commit.
- Current branch: `phase-21`; the only pre-existing worktree modification was
  the runner-owned `project-docs/docs/phases/phase-21/run-ledger.md`.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-013, closed by DEC-172; provenance only
- Covered validation assertions: VAL-055, VAL-056, VAL-057, VAL-058, VAL-059,
  VAL-060, VAL-061, VAL-062, VAL-063, VAL-064, VAL-065, VAL-066, VAL-067,
  VAL-068

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-055 | task-16 | ON TRACK | The repaired focused Artifact suite passed; the safe Artifact test retains byte, digest, metadata, and staging-file assertions. |
| VAL-056 | task-16 | ON TRACK | `cmd/fid/artifact_test.go:74` covers changed-during-copy, directory, missing, non-regular, outside-cwd, oversized, and symlink rows; each checks `artifact-rejected`, no Receipt Artifact link, and an unsatisfied evaluation. |
| VAL-057 | task-16 | ON TRACK | `cmd/fid/artifact_test.go:252` links Evidence before the loss/restore cycle and compares Artifact, Receipt, and Evidence identity snapshots at lines 289-328. |
| VAL-058 | task-17 | ON TRACK | `cmd/fid/evidence_test.go:107` supplies the qualifying-chain mismatch matrix for Obligation, Registration, Run, Issue, and repository target; every rejected row checks zero Event/evidence-record deltas. |
| VAL-059 | task-17 | NOT ON TRACK | The sole dimension test changes validity, required Artifact availability, and tree OID together at `cmd/fid/evidence_test.go:68-81`; no valid/invalid one-dimension pair exists. |
| VAL-060 | task-17 | NOT ON TRACK | The same coupled scenario does not show original-tree fresh versus changed-tree stale independently, and it includes no branch-label control. |
| VAL-061 | task-17 | NOT ON TRACK | The available/missing Artifact comparison is coupled to invalidity and staleness rather than asserting an isolated Artifact reason/status diff. |
| VAL-062 | task-17 | NOT ON TRACK | No repeat or post-rebuild byte-identical ordered-reason golden exists; the one assertion only checks a single in-process three-reason result. |
| VAL-063 | task-17 | ON TRACK | `TestPhase21AttentionBoundary` remains green and requires the exact unavailable `attention.v1` row as the sole production Attention match, with an unavailable-leaf no-mutation probe. |
| VAL-064 | task-18 | ON TRACK | The existing positive guarded-transition and evaluator-identity coverage remains consistent with the task contract. |
| VAL-065 | task-18 | ON TRACK | `cmd/fid/task18_completion_resume_test.go:16-109` covers unsuccessful, missing, invalid, stale, and Artifact-missing states, each with its exact sole reason and unchanged Issue/Event state. |
| VAL-066 | task-18 | ON TRACK | The existing successful completion test retains shared-evaluator and completion-event coverage. |
| VAL-067 | task-18 | ON TRACK | The existing failed-completion test compares Envelope, Issue, Evidence, legacy completion rows, and Event state. |
| VAL-068 | task-18 | ON TRACK | `cmd/fid/task18_completion_resume_test.go:201-270` asserts active and completed durable projections, ordered actions, no recovery/respawn action, and zero Event deltas. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 16 exact command audit | `task-16-done.md` vs. task spec | PASS | All five required commands are recorded PASS with own-output tokens and observed `exit_code: 0`. |
| Task 17 exact command audit | `task-17-done.md` vs. task spec | PASS | All six DEC-172 commands are recorded PASS with observed `exit_code: 0`; the superseded broad search remains a distinct BLOCKED historical attempt. |
| Task 18 exact command audit | `task-18-done.md` vs. task spec | PASS | All six required commands, including tagged integration, are recorded PASS with observed `exit_code: 0`. |
| Task 16 report control | `project-docs/tools/report-lint.sh ...task-16-done.md ...task-16-artifacts.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 17 report control | `project-docs/tools/report-lint.sh ...task-17-done.md ...task-17-evidence-evaluation.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). This mechanical result does not validate the unsupported VAL-059 through VAL-062 PASS claims. |
| Task 18 report control | `project-docs/tools/report-lint.sh ...task-18-done.md ...task-18-completion-resume.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 16 reproduced evidence | `go test -race ./cmd/fid -run '^(TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable|TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 33.443s` (terminated `exit_code: 0`). |
| Task 17 reproduced evidence | `go test -race ./cmd/fid -run '^(TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceLinkRejectsEveryCallerAuthoredProof|TestEvidenceLinkRejectsPersistedReceiptBindingMismatches|TestPhase21AttentionBoundary)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 20.934s` (terminated `exit_code: 0`). Source inspection shows the passing dimension test remains insufficient for QA-57 through QA-60. |
| Task 18 reproduced evidence | `go test -race ./cmd/fid -run '^(TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence|TestTerminalSuccessTransitionRejectsEveryOtherNonqualifyingEvidenceState|TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 33.897s` (terminated `exit_code: 0`). |
| Task 18 tagged integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 7.166s` (terminated `exit_code: 0`). |
| Commit and diff provenance | `git rev-parse --verify <hash>^{commit}`; `git diff --check 22ad81d..HEAD` | PASS | All six inspected reported hashes resolved; `git diff --check` terminated `exit_code: 0` with no output. |

## Per-Task Findings

### Task 16 - Safe Referenced Artifacts

#### Must-fix

- None. MF-16-01 is resolved: the unsafe candidate matrix is observable through
  the gate CLI and the availability lifecycle now includes Evidence identity.

#### Improvement

- None.

### Task 17 - Trusted Evidence Evaluation

#### Must-fix

- MF-17-02 [cmd/fid/evidence_test.go:16,68-81] is still open. The single
  non-qualifying fixture mutates validity, tree OID, and Artifact availability
  before one evaluation, so it cannot prove the one-dimension paired states
  required by QA-57 through QA-59. It has neither the QA-58 branch-label
  control nor QA-60 repeat/post-rebuild byte-identical ordered-reason golden.
  The completion report's VAL-059 through VAL-062 PASS assertions are therefore
  unsupported. Owning task: task-17. Repair boundary: add only the independent
  paired validity, repository-tree, and Artifact fixtures; add a branch-label
  control and repeat/post-rebuild byte-comparison of the full ordered reasons;
  correct the completion report to state only reproduced evidence. Required
  evidence: qualifying and non-qualifying output for all three dimensions,
  exact tree-not-branch control, byte-identical repeat/rebuild golden, exact
  task-17 command block, and a fresh two-argument report lint. Re-review scope:
  inspect the paired fixtures and golden, rerun the focused Evidence/Attention
  command, and rerun the task-17 report control.

#### Improvement

- None. MF-17-01 is resolved: `cmd/fid/evidence_test.go:107-195` covers all
  five persisted binding dimensions and the separate caller-proof test; the
  shared rejection helper at lines 278-285 verifies no Evidence Event or
  `evidence_records` mutation. DEC-172 remains exact: the unavailable
  `attention.v1` row is preserved as negative capability metadata, not a
  proposed Attention surface.

### Task 18 - Guarded Terminal Transition, Completion, And Resume

#### Must-fix

- None. MF-18-01 and MF-18-02 are resolved by the negative transition matrix
  and the active/completed exact resume projections, both reproduced above.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-16-01 | task-16 | RESOLVED | Focused Artifact suite passed with `exit_code: 0`; source verifies every unsafe control/evaluation row and Artifact/Receipt/Evidence lifecycle identity. |
| MF-17-01 | task-17 | RESOLVED | Focused Evidence/Attention suite passed with `exit_code: 0`; persisted mismatch matrix and no Event/record mutation assertions are present. |
| MF-17-02 | task-17 | STILL OPEN | The only dimension test remains coupled at `evidence_test.go:68-81`; branch-label and repeat/rebuild golden coverage are absent. |
| MF-18-01 | task-18 | RESOLVED | Focused terminal negative suite passed with `exit_code: 0`; all five states assert their exact reason and zero Issue/Event mutation. |
| MF-18-02 | task-18 | RESOLVED | Focused resume suite passed with `exit_code: 0`; active and completed projection/action assertions prove zero Event deltas and no recovery/respawn. |

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-17-02 | task-17 | VAL-059, VAL-060, VAL-061, VAL-062 | Add only independent paired validity/freshness/Artifact fixtures, the branch-label control, and repeat/post-rebuild byte-identical ordered-reason golden; update the task-17 report to match. | Paired qualifying/non-qualifying results, branch-label control, repeat/rebuild byte comparison, exact task-17 command block, and report lint. | New task-17 fixture/golden source, focused Evidence/Attention test command, and two-argument report lint. |

## Blocker Candidates

- None. MF-17-02 is bounded, decision-free evidence completion within task-17;
  DEC-172 already closes the Attention authority and requires no product or
  public-interface decision.

## Existing Escalation References

- ESC-21-013, closed by DEC-172; provenance only, not a blocker candidate.

## Final Lines

Review written: project-docs/docs/phases/phase-21/qa/review-6b.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
