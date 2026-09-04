# Phase 21 - Mid-Phase Review 4b (tasks 10-12)

Reviewer: qa-reviewer
Date: 2026-08-28
Tasks covered: phase-21/task-10, phase-21/task-11, phase-21/task-12
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-4.md

## Verdict

PASS

No Must-fix finding remains. The bounded test/report repairs preserve accepted
Claim, Run, and execution-transition behavior; no protected authority,
contract, or gate-definition change is needed.

## Scope Reviewed

- Task specs and current completion reports for tasks 10, 11, and 12.
- QA-31 through QA-39; VAL-032 through VAL-041; current phase context;
  DEC-162 through DEC-167; the TRG-21-012 resolution note; and the current
  run ledger.
- Prior Window 4 review and earlier review artifacts as phase context.
- Repair commits `59b079a5a271cc65f7964a471b6383b9c2dd6677`,
  `8c763e9cddee1d4ff0df3379e08b286ef74d5300`, and
  `0ced511b2e6a7e1f252260024176705106e9b793`, plus reports
  `1c4d222b998374c97134e06420fc1a7dcb0cac79`,
  `1f24ea26e0d0ad0b32c33d2dec37ec811b631f27`, and
  `20e6760726edc718fe6d0a79df046835202c9b85`.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-032, VAL-033, VAL-034, VAL-035, VAL-037, VAL-038, VAL-039, VAL-040, VAL-041

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-032 | task 10 | ON TRACK | [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:65) has a qualifying acquire and individually rejects missing/mismatched Issue, missing Envelope, wrong revision, Actor authority, worktree, scope, and TTL with a typed failure plus unchanged Claim-list/journal snapshot. |
| VAL-033 | task 10 | ON TRACK | [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:201) retains conflict, expiry, and wrong-holder no-mutation rows. |
| VAL-034 | task 10 | ON TRACK | [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:242) now rejects stale renew and stale release as `claim-stale`, preserving Claim-list and journal state before retaining exact replay evidence. |
| VAL-035 | task 11 | ON TRACK | [generic_run_test.go](/home/mark/projects/fid/cmd/fid/generic_run_test.go:56) has a qualifying start and rejected inactive Claim, Issue, Envelope/revision, Actor, repository, and injected-budget rows. Each rejected row asserts the canonical typed error plus unchanged Run list and Event count. |
| VAL-037 | task 11 | ON TRACK | [generic_run_test.go](/home/mark/projects/fid/cmd/fid/generic_run_test.go:149) keeps one-way terminal finish/no Event on repeat and proves `run recover`, `run predecessor`, `run continuation`, and `run respawn` are absent. |
| VAL-038 | task 12 | ON TRACK | [task12_execution_transition_test.go](/home/mark/projects/fid/cmd/fid/task12_execution_transition_test.go:114) exercises the qualifying transition and independently-mutated rejection rows for Envelope, Actor, Claim, active Run, Run repository, Issue revision, and workflow edge. The common assertion snapshots Issue status/revision and journal count. |
| VAL-039 | task 12 | ON TRACK | [task12_execution_transition_test.go](/home/mark/projects/fid/cmd/fid/task12_execution_transition_test.go:53) compares complete ordered `workflow.EvaluationReason` JSON to the contract-section-10 golden. [The rejection helper](/home/mark/projects/fid/cmd/fid/task12_execution_transition_test.go:248) asserts the exact CLI envelope and no Issue/Event mutation. |
| VAL-040 | task 12 | ON TRACK | The retained evaluator dual-call test and terminal-transition no-mutation assertion remain recorded in the task report; the focused execution-bound corpus reran fresh below. |
| VAL-041 | task 12 | ON TRACK | The retained direct `TestBuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects` maps whitespace, normalized duplicate, and overflow inputs to no Effects; the exact required Task 12 gate remains recorded PASS and report-lint verified coverage. |

## Gate And Evidence Checks

All fresh commands below ran from `/home/mark/projects/fid`. The completion
reports were compared with their task specs. Their exact required command sets
are complete, with recorded PASS rows and observed terminations; the Window 4
QA-plan evidence requirements add no unrecorded Developer command beyond the
Task 12 tagged integration/unit matrix pair.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Repair commit integrity | `git rev-parse --verify <each reported hash>^{commit}` | PASS | Review-4, all three repair implementations, and all three repair/report commits resolved to the exact reported commits. |
| Repair boundary | `git diff --name-status 11396b45... 20e676072...` | PASS | Only the three bounded test files and their three task reports changed after review-4; no product, task-spec, QA-plan, validation-contract, decision, escalation, or ledger artifact changed in those repairs. |
| Task 10 focused matrix/lifecycle | `go test -race ./cmd/fid -run '^(TestClaimAcquireBindingMatrixRejectsWithoutMutation|TestClaimRenewReleaseVersionAndReplayAreExact)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  9.427s` (terminated `exit_code: 0`). |
| Task 11 focused matrix | `go test -race ./cmd/fid -run '^(TestRunStartBindingMatrixRejectsWithoutMutation|TestRunFinishIsOneWayAndRecoveryRespawnAreAbsent)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  18.269s` (terminated `exit_code: 0`). |
| Task 12 predicate/atomicity | `go test -race ./cmd/fid -run '^TestExecutionBound' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  32.799s` (terminated `exit_code: 0`). |
| Task 12 tagged integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  6.913s` (terminated `exit_code: 0`). It includes `TestV01ClosureJourneyAndRebuildParity`, whose first Phase 21 green is new evidence, not continuity evidence, under DEC-167. |
| Task 10 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-10-done.md project-docs/docs/phases/phase-21/tasks/task-10-generic-claims.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 11 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-11-done.md project-docs/docs/phases/phase-21/tasks/task-11-generic-runs.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 12 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-12-done.md project-docs/docs/phases/phase-21/tasks/task-12-execution-bound-transition.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| DEC-167 convergence preservation | `git show 16b15d8134e7b29f27cf996486bff997a080fcc3 -- internal/dev/start_test.go cmd/fid/integration` | PASS | The Dev diff removes only the obsolete missing-Claim setup probe; the retained integration rows remain typed fail-closed/no-journal assertions with the task-12 report's per-row provenance. TRG-21-012 is not reopened. |

## Per-Task Findings

### Task 10 - Generic Mutation Claims

#### Must-fix

- None.

#### Improvement

- None.

### Task 11 - Generic Runs

#### Must-fix

- None.

#### Improvement

- None.

### Task 12 - Shared Evaluator And Nonterminal Transition

#### Must-fix

- None.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-1: Claim acquire matrix and stale lifecycle evidence | task 10 | RESOLVED | `59b079a5a271cc65f7964a471b6383b9c2dd6677` adds every required rejected acquire binding and stale renew/release no-mutation proof in [generic_claim_test.go](/home/mark/projects/fid/cmd/fid/generic_claim_test.go:65); the focused matrix/lifecycle command and task-10 two-argument report-lint both pass fresh. |
| MF-2: Run binding matrix and prohibited continuation absence | task 11 | RESOLVED | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` adds the qualifying/rejected binding matrix and all four unavailable leaves in [generic_run_test.go](/home/mark/projects/fid/cmd/fid/generic_run_test.go:56); the focused command and task-11 two-argument report-lint both pass fresh. The later report-only commit `1f24ea26e0d0ad0b32c33d2dec37ec811b631f27` is lint-clean. |
| MF-3: Execution predicate matrix, atomicity, and ordered-reason golden | task 12 | RESOLVED | `0ced511b2e6a7e1f252260024176705106e9b793` adds every required predicate row, Issue snapshot, and byte-exact golden in [task12_execution_transition_test.go](/home/mark/projects/fid/cmd/fid/task12_execution_transition_test.go:114); fresh predicate, tagged-integration, and task-12 report-lint controls pass. DEC-167's Dev/tagged rows remain subject-preserving. |

## Blocker Candidates

- None.

## Existing Escalation References

- None. TRG-21-012 and ESC-21-009 are resolved historical context under DEC-167 and do not govern this review outcome.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-4b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
```
