# Phase 8 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: QA Reviewer
Date: 2026-06-20
Tasks covered: phase-8/task-1 through phase-8/task-3

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-8/tasks/task-1.md`
  - `project-docs/docs/phases/phase-8/tasks/task-2.md`
  - `project-docs/docs/phases/phase-8/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-8/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-8/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-8/tasks/task-3-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-8/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-8/qa/qa-plan.md`
- Commits or changed files:
  - `cc03c59 feat(issue): add issue completion write path`
  - `d381f4f feat(issue): snapshot gate runs into completions`
  - `d0bbb83 test: add completion gate rebuild parity e2e`
  - `internal/issues/complete.go`
  - `internal/issues/complete_test.go`
  - `cmd/fid/e2e_v01c_complete_test.go`

## Summary

- Blocking findings: 0
- Improvement findings: 3
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | Task report records `TestCompleteCreatesAttemptOneAndPersistsCompletionFields`, `TestIssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields`, and a real CLI transcript with `TC-` read-back for notes, tokens, commits, and files. Representative QA rerun passed: `go test ./internal/issues ./cmd/fid -run 'Test(CompleteCreatesAttemptOneAndPersistsCompletionFields|CompleteSupersedesPriorAttemptAndRetainsHistory|CompleteStoresHumanInterventionStateWithoutActivityLogRows|CompleteRejectsEpicBeforeAppend|IssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields|IssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields|IssueCompleteCommandRejectsEpicBeforeAppend)$'`. Design checked: `05_DESIGN_PLAN.md` command map (`issue complete`), `06_DESIGN_PLAN_CLI.md` `fid issue complete`. |
| VAL-002 | task-1 | ON TRACK | Task report JSON transcript shows attempt 2 as `dev_complete` and attempt 1 retained as `superseded`; focused tests passed in the same representative QA rerun. |
| VAL-003 | task-1 | ON TRACK | Task report and tests cover both flagged and unflagged states. `internal/issues/complete_test.go` asserts flagged note storage and unflagged nil note. |
| VAL-004 | task-1 | ON TRACK | Task report records epic rejection transcript and tests asserting unchanged journal line counts plus zero `task_completions` rows. Focused QA rerun passed. |
| VAL-005 | task-2 | ON TRACK | `internal/issues/complete_test.go` seeds fixture `gate_runs` rows for pass, fail, and non-matching attempt; asserts `gate_results`, `gates_passed=1`, `gates_failed=1`, read-back shape, source-row preservation, and no-gate `[]`/0/0. Representative QA rerun passed: `go test ./internal/issues -run 'TestComplete(SnapshotsMatchingGateRunsAndPreservesSourceRows|LeavesGateSnapshotEmptyWhenNoMatchingGateRunsExist)$'`. Design checked: `05_DESIGN_PLAN.md` command map (`dev gate`), `06_DESIGN_PLAN_CLI.md` `fid dev gate` snapshot note. |
| VAL-006 | task-3 | ON TRACK | `cmd/fid/e2e_v01c_complete_test.go` uses real `fid issue transition` and `fid issue complete` command helpers: reject without completion, assert stderr includes `has_completion_report` and issue ID, assert unchanged journal lines and unchanged issue JSON, then complete and transition successfully. Representative QA rerun passed: `go test -tags integration ./cmd/fid -run TestV01CHasCompletionReportGateAndRebuildParity -count=1`. Design checked: `05_DESIGN_PLAN.md` command map (`issue complete`, `rebuild`), `06_DESIGN_PLAN_CLI.md` `fid issue complete`, `fid issue transition`, and `fid rebuild`. |
| VAL-007 | task-3 | ON TRACK | The same integration test seeds attempt-2 `gate_runs`, records a second real completion, snapshots completions JSON and ordered `gate_runs`, deletes `.fid/fid.sqlite`, runs `fid rebuild`, and asserts exact completion JSON plus `gate_runs` parity. Representative QA rerun passed. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 report shape | Required headings, concerns section, dependency-budget note | PASS | `task-1-done.md` includes all required headings; concerns section is non-empty with `None`; dependency note records no changes. |
| Task 1 required commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; conditional `go build -o /tmp/fid ./cmd/fid` | PASS | All required and triggered conditional commands are reported as PASS. Coverage highlights meet `cmd/fid >=70%` and `internal/issues >=80%`. |
| Task 1 commit | `git rev-parse --verify cc03c59^{commit}` and `git show --no-patch` | PASS | Resolves to `cc03c5960619164b90a6f3e1991510779c3af4ef`; subject and `Refs: phase-8/task-1` match the report. |
| Task 1 QA spot-check | `go test ./internal/issues ./cmd/fid -run 'Test(CompleteCreatesAttemptOneAndPersistsCompletionFields|CompleteSupersedesPriorAttemptAndRetainsHistory|CompleteStoresHumanInterventionStateWithoutActivityLogRows|CompleteRejectsEpicBeforeAppend|IssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields|IssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields|IssueCompleteCommandRejectsEpicBeforeAppend)$'` | PASS | `internal/issues` and `cmd/fid` both passed. |
| Task 2 report shape | Required headings, concerns section, dependency-budget note | PASS | `task-2-done.md` includes all required headings; concerns section is non-empty with `None`; dependency note records no changes. |
| Task 2 required commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; conditional integration/build commands | PASS | Required commands are PASS. Conditional integration/build commands are explicitly NOT RUN with acceptable reasons: no integration tag and no CLI transcript needed. |
| Task 2 commit | `git rev-parse --verify d381f4f^{commit}` and `git show --no-patch` | PASS | Resolves to `d381f4f37d4c41b0889ddbf1f91699936e9f67cf`; subject and `Refs: phase-8/task-2` match the report. |
| Task 2 QA spot-check | `go test ./internal/issues -run 'TestComplete(SnapshotsMatchingGateRunsAndPreservesSourceRows|LeavesGateSnapshotEmptyWhenNoMatchingGateRunsExist)$'` | PASS | Passed from cache; evidence file shows fixture-seeded `gate_runs`, source-row preservation, attempt filtering, and no-gate state. |
| Task 3 report shape | Required headings, concerns section, dependency-budget note | PASS | `task-3-done.md` includes all required headings; concerns section is non-empty with `None`; dependency note records no changes. |
| Task 3 required commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; triggered `go test -tags integration -race ./...`; triggered `go build -o /tmp/fid ./cmd/fid` | PASS | Required commands and triggered conditionals are reported as PASS. Coverage highlights meet `cmd/fid >=70%`, `internal/issues >=80%`, and `internal/workflow >=80%`. |
| Task 3 commit | `git rev-parse --verify d0bbb83^{commit}` and `git show --no-patch` | PASS | Resolves to `d0bbb839614836f442dda6176a5e14d820f61bca`; subject and `Refs: phase-8/task-3` match the report. |
| Task 3 QA spot-check | `go test -tags integration ./cmd/fid -run TestV01CHasCompletionReportGateAndRebuildParity -count=1` | PASS | Passed; test uses real `fid issue complete` and `fid issue transition` command helpers for VAL-006 and fixture-seeded `gate_runs` only for the DEC-C-approved VAL-007 snapshot parity path. |

## Per-Task Findings

### Task 1 - `fid issue complete` core write

#### Blocking

- None.

#### Improvement

- `internal/issues/complete.go:45` `Complete` is a large orchestration function that combines option normalization, policy/projection freshness, lock-scoped validation, ID generation, attempt derivation, gate snapshot lookup, event construction, append, projection, and result assembly. This does not invalidate VAL-001..004, but it is over the project function-size guidance and will get harder to reason about as Task 4 splits the issue package. Fix direction: extract narrow helpers for completion event construction and append/project result assembly during the next structure-focused cleanup.

### Task 2 - Gate snapshot

#### Blocking

- None.

#### Improvement

- `internal/issues/complete_test.go:287` and `cmd/fid/e2e_v01c_complete_test.go:268` now contain separate journal-fixture helpers for seeding `gate_runs`. This is acceptable for the review window and keeps Task 2/3 evidence local, but the duplication can drift when `fid dev gate` lands in v0.1e. Fix direction: when the real gate writer ships, replace or centralize the fixture-only seeding helper behind a shared test helper that preserves the DEC-C audit boundary.

### Task 3 - `has_completion_report` real-data gating + rebuild parity

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-8/tasks/task-3-done.md:46` labels the evidence as an observed command transcript excerpt and uses `<task>` placeholders. The deterministic integration test is sufficient for VAL-006/007, but the task spec asked for a full end-to-end transcript for the completion report. Fix direction: before phase wrap, either paste a literal transcript with concrete IDs or explicitly mark the report's transcript as summarized test evidence so final QA does not treat it as exact CLI output.

## Cross-cutting observations

- Current review source files include uncommitted report/QA-plan state in the worktree (`task-1-done.md`, `task-2-done.md`, `task-3-done.md` modified; `qa-plan.md` untracked at review time). I reviewed the current artifacts as requested and did not revert anything.
- Task 3 preserves the intended fixture boundary: VAL-006 is proven with real commands, while fixture-seeded `gate_runs` remain limited to the DEC-C-approved gate snapshot/parity evidence.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-8/qa/review-1.md
Blocking: 0  Improvement: 3  Escalations: none
VERDICT: CONTINUE
```
