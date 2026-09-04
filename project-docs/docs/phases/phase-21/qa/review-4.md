# Phase 21 - Mid-Phase Review 4 (tasks 10-12)

Reviewer: qa-reviewer
Date: 2026-08-28
Tasks covered: phase-21/task-10, phase-21/task-11, phase-21/task-12
Review type: INITIAL
Prior review: none (review-1/1b, review-2/2b, and review-3/3b inspected as prior-window context)

## Summary

- Review outcome: REPAIR
- Must-fix findings: 3
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-032, VAL-033, VAL-034, VAL-035, VAL-037, VAL-038, VAL-039, VAL-040, VAL-041

The exact task gates and all three two-argument report-lint controls are green, but the committed observable tests do not meet Window 4's required qualifying/non-qualifying matrices. The repairs below are test-and-report scoped within their existing owning tasks and do not require a product, contract, gate, or authority change.

## Scope Reviewed

- Window 4 QA-31 through QA-39 and VAL-032 through VAL-041; task specs and committed completion reports for tasks 10, 11, and 12.
- Current phase context, DEC-162 through DEC-167, the Resolution Note for TRG-21-012, the run ledger, and review-1/1b, review-2/2b, and review-3/3b.
- Task 10 implementation `0b67c8304a84e4e5bcb8df44d79f899e40a16acb` (`feat(claims): add generic mutation claim lifecycle`) and report `a8722b81062002b976927b438c8514c5b7a46e64` (`docs(task): finalize completion report for phase-21 task-10`).
- Task 11 implementation `5ca732ff163ee1d31fc155d3d9ebb3dac1501838` (`feat(workflow): add generic Run lifecycle`) and final report repair `f463e14337867db3deacbd0a3b516fca012e657a` (`docs(task): repair completion report for phase-21 task-11`), whose ledger record confirms the five exact commands and two-argument report-lint PASS.
- Task 12 implementation `16b15d8134e7b29f27cf996486bff997a080fcc3` (`test(integration): converge retained closure journey`) and report repair `92a5e01b79577031fe9d8b0b99661256d76c4e45` (`docs(task): canonicalize task-12 historical gate evidence`). All listed hashes resolve as commits and their report subjects match the reports.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-032 | task 10 | NOT ON TRACK | [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:12) proves one qualifying acquire and stored fields, but contains no rejected incomplete/wrong Issue, Envelope, revision, Actor, worktree, scope, or TTL row with a no-Event assertion. QA-31 requires both states. |
| VAL-033 | task 10 | ON TRACK | [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:72) proves typed active-Claim conflict, wrong-Actor renewal denial, expiry denial, and unchanged journal count. |
| VAL-034 | task 10 | NOT ON TRACK | [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:113) proves advancing renew/release and replay, but it contains no stale-version rejection despite the report's claim. QA-33 requires that non-qualifying lifecycle row. |
| VAL-035 | task 11 | NOT ON TRACK | [generic_run_test.go](/home/mark/projects/fid/cmd/fid/generic_run_test.go:9) verifies one qualifying start and persisted values only. It does not vary active Claim, Issue, Envelope/revision, Actor, repository, or budget binding to prove a failed start has no Event, as QA-34 requires. |
| VAL-037 | task 11 | NOT ON TRACK | [generic_run_test.go](/home/mark/projects/fid/cmd/fid/generic_run_test.go:56) proves explicit terminal finish, repeat no-mutation, and absence of `run recover`/`run respawn`, but does not prove the required predecessor/continuation absence. |
| VAL-038 | task 12 | NOT ON TRACK | [task12_execution_transition_test.go](/home/mark/projects/fid/cmd/fid/task12_execution_transition_test.go:11) proves the all-qualified row; its failures do not provide one independent no-mutation row for every Envelope, Actor, Claim, Run, repository, revision, and edge predicate required by QA-36. In particular envelope-inactive, a usable-Claim/no-active-Run case, and mismatched Run repository are absent. |
| VAL-039 | task 12 | NOT ON TRACK | [task12_execution_transition_test.go](/home/mark/projects/fid/cmd/fid/task12_execution_transition_test.go:58) checks that stderr contains a sequence of codes and checks Event count only. It does not byte-compare the complete reason objects to contracts.md section 10 or prove the Issue revision is unchanged, as QA-37 requires. |
| VAL-040 | task 12 | ON TRACK | [evaluation_test.go](/home/mark/projects/fid/internal/workflow/evaluation_test.go:8) calls the shared evaluator for terminal-success and completion action identities and requires the same ordered absent-Evidence reason. |
| VAL-041 | task 12 | ON TRACK | [revision_test.go](/home/mark/projects/fid/internal/issues/revision_test.go:42) rejects whitespace IDs, normalized duplicates, and revision overflow with nil Effects. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 10 report and required gates | task-10 report compared with task-10 spec | PASS | All five exact commands are recorded PASS with output tokens and observed `exit_code: 0`; historical failed attempts and later byte-identical repairs remain preserved. No additional task-specific QA command is declared. |
| Task 10 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-10-done.md project-docs/docs/phases/phase-21/tasks/task-10-generic-claims.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 10 representative rerun | `go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...` | PASS | Final token: `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 11 report and required gates | task-11 report compared with task-11 spec | PASS | All five exact commands are recorded PASS with output tokens and observed `exit_code: 0`; no additional task-specific QA command is declared. |
| Task 11 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-11-done.md project-docs/docs/phases/phase-21/tasks/task-11-generic-runs.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 11 representative rerun | `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` | PASS | Final token: `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 12 report and required gates | task-12 report compared with amended task-12 spec | PASS | The final Attempt 6 records all six exact commands PASS with output tokens and observed `exit_code: 0`; its historical attempts remain preserved. No additional task-specific QA command is declared. |
| Task 12 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-12-done.md project-docs/docs/phases/phase-21/tasks/task-12-execution-bound-transition.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 12 representative rerun | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 5.808s` (terminated `exit_code: 0`). This is the first green Phase 21 closure-journey evidence required by DEC-167, not continuity evidence. |
| Task 12 affected package rerun | `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` | PASS | Final token: `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Current Window 4 focused corpus | `go test -race ./cmd/fid -run '^(TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL|TestClaimConflictsExpiryAndWrongHolderFailWithoutMutation|TestClaimRenewReleaseVersionAndReplayAreExact|TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|TestRunFinishIsOneWayAndRecoveryRespawnAreAbsent|TestExecutionBound)' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 22.681s` (terminated `exit_code: 0`). The green result establishes the present corpus executes; source inspection above identifies its required missing rows. |

## Per-Task Findings

### Task 10 - Generic Mutation Claims

#### Must-fix

- [generic_claim_test.go:12] QA-31/VAL-032 lacks the required non-qualifying acquire matrix: only one valid acquire is asserted, while missing or mismatched Issue, Envelope, revision, Actor authority, worktree, scope, and TTL have no observable no-mutation proof. [generic_claim_test.go:113] also lacks QA-33/VAL-034's required stale-version rejection even though the completion report claims it. Owning task: phase-21/task-10. Repair boundary: add CLI-observable table rows for every named rejected binding plus stale renew/release version, each preserving Claim/journal state; update task-10 evidence only. Required evidence: focused matrix output with qualifying and every non-qualifying row, all five unchanged exact task-10 gates, and two-argument report-lint PASS. Re-review scope: bounded task-10 tests/report diff; verify complete input/version matrix and no-mutation assertions; rerun the exact task-10 package command and report-lint control.

#### Improvement

- None.

### Task 11 - Generic Runs

#### Must-fix

- [generic_run_test.go:9] QA-34/VAL-035 requires a start-binding matrix for active Claim, Issue, Envelope/revision, Actor, repository, and immutable budget contracts, with any mismatch failing. The committed test asserts only the qualifying start; no rejected row or no-Event assertion exercises those binding failures. [generic_run_test.go:56] proves only recover/respawn absence, leaving QA-35/VAL-037 predecessor/continuation absence unproven. Owning task: phase-21/task-11. Repair boundary: add CLI-observable qualifying/non-qualifying Run-start rows, each with unchanged Event/state proof, and absence checks for every prohibited Run continuation surface; update task-11 evidence only. Required evidence: start/finish/absence matrix output, all five unchanged exact task-11 gates, and two-argument report-lint PASS. Re-review scope: bounded task-11 tests/report diff; verify every binding and unavailable-path row; rerun the exact task-11 package command and report-lint control.

#### Improvement

- None.

### Task 12 - Shared Evaluator And Nonterminal Transition

#### Must-fix

- [task12_execution_transition_test.go:11] QA-36/VAL-038 requires one non-mutating failure per Envelope, Actor, Claim, Run, repository, revision, and edge predicate. The test has the qualifying row and several combined failures, but no independent envelope-inactive, usable-Claim/no-active-Run, or mismatched-Run-repository row. [task12_execution_transition_test.go:58] also does not meet QA-37/VAL-039's complete atomicity golden: it uses `strings.Contains` for codes and checks Event count, not the full ordered reason objects against contracts.md section 10 and an unchanged Issue revision. Owning task: phase-21/task-12. Repair boundary: add an execution-bound predicate table with one isolated failing row per required predicate and no Event/revision advance, plus a byte-exact multi-failure reason golden and state snapshot; preserve DEC-167 convergence rows and update task-12 evidence only. Required evidence: affected unit matrix, fresh tagged integration PASS, all six unchanged exact task-12 gates, and two-argument report-lint PASS. Re-review scope: bounded task-12 tests/report diff; inspect the DEC-167 rows remain subject-preserving; rerun the exact task-12 package and tagged-integration commands and report-lint control.

#### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | phase-21/task-10 | VAL-032, VAL-034 | `cmd/fid/generic_claim_test.go` and `task-10-done.md`; add all required failed acquire bindings and stale-version lifecycle rows with no mutation, without changing Claim behavior or contracts. | Focused qualifying/non-qualifying CLI matrix; exact five-command task gate; two-argument report-lint PASS. | Inspect bounded tests/report; verify each listed binding and stale-version rejection; rerun exact task-10 package gate and report lint. |
| MF-2 | phase-21/task-11 | VAL-035, VAL-037 | `cmd/fid/generic_run_test.go` and `task-11-done.md`; add failed Run-start binding rows and complete prohibited continuation absence assertions, without changing Run behavior or contracts. | Focused start/finish/absence matrix with no-Event proof; exact five-command task gate; two-argument report-lint PASS. | Inspect bounded tests/report; verify all bindings and prohibited paths; rerun exact task-11 package gate and report lint. |
| MF-3 | phase-21/task-12 | VAL-038, VAL-039 | `cmd/fid/task12_execution_transition_test.go` and `task-12-done.md`; add isolated predicate failures, revision snapshots, and byte-exact ordered-reason golden, preserving DEC-167 convergence assertions. | Predicate and atomicity matrix; exact six-command task gate including tagged integration; two-argument report-lint PASS. | Inspect bounded tests/report and retained convergence rows; rerun exact task-12 package and tagged-integration gates and report lint. |

## Cross-cutting observations

- The DEC-167 convergence is correctly retained: the committed Task 12 report has per-row provenance for `internal/dev` and the tagged integration rows, and this review's fresh tagged integration command is green. TRG-21-012 is not reopened and no owner routing is inferred.
- All reports pass the canonical two-argument report-lint control. The findings are evidence-strength gaps in otherwise green task gates, not report-format defects or quality-gate failures.

## Blocker Candidates

- None.

## Existing Escalation References

- None. ESC-21-009/TRG-21-012 are resolved historical context under DEC-167 and do not govern these bounded repairs.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-4.md
Must-fix: 3  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
