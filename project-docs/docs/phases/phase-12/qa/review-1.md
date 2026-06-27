# Phase 12 - Mid-Phase Review 1 (Tasks 1-2)

Reviewer: QA Reviewer
Date: 2026-06-27
Tasks covered: phase-12/task-1 through phase-12/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-12/tasks/task-1.md`
  - `project-docs/docs/phases/phase-12/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-12/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-12/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-12/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-12/qa/qa-plan.md`
- Current context:
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/docs/current/decisions_active.md`
- Commits or changed files:
  - `e236591de43d994b2b816e4398c1470b4fc4d4c9` (`feat(requirements): add requirement sync`) verified with `git rev-parse --verify`.
  - `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` (`feat(links): add link add command`) verified with `git rev-parse --verify`.
  - Spot-checked changed files under `cmd/fid`, `internal/requirements`, `internal/links`, `internal/policy`, and Phase 12 decisions/reports.

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008

Task 1 report verdict: complete for Review 1. Required commands are recorded, conditional integration is explicitly N/A with a reason, DEC-029 records PI-1, and evidence tracks for VAL-001 partial plus VAL-002..VAL-004.

Task 2 report verdict: complete for Review 1. Required commands are recorded, conditional integration is explicitly N/A with a reason, DEC-030 records PI-2, and evidence tracks for VAL-001 completion plus VAL-005..VAL-008.

No blockers were found. No new open ESC is needed.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1, task-2 | ON TRACK | Task 1 report lines 39-42 cover `requirement.sync`; Task 2 report lines 42-47 cover `link.add` and `link.remove`. Spot check passed: `go test ./internal/policy -run 'TestMatrixRequirementSyncRowMatchesSchemaSection13|TestMatrixLinkRowsMatchSchemaSection13'` -> `ok github.com/mbell-nzl/fid/internal/policy`. |
| VAL-002 | task-1 | ON TRACK | Task 1 report lines 43-49 cite create/no-op/update/filter tests. Code spot-check in `cmd/fid/requirement_test.go` lines 52-110 asserts one event, required persisted fields, and requirement insert effects. |
| VAL-003 | task-1 | ON TRACK | Task 1 report lines 50-53 cite all malformed states required by the contract: no parseable IDs, duplicate IDs, unknown `--type`, and missing `--id`, with zero journal/persisted-row deltas. |
| VAL-004 | task-1 | ON TRACK | Task 1 report lines 54-57 cite tl/human allowed and denied/unknown role rejection before append. Spot check passed: `go test ./cmd/fid -run 'TestRequirementSyncCreatesPointersInOneJournalEvent|TestRequirementSyncRejectsMalformedInputBeforeAppend|TestRequirementSyncAuthorizationAndPolicyShowParity'` -> `ok github.com/mbell-nzl/fid/cmd/fid`. |
| VAL-005 | task-2 | ON TRACK | Task 2 report lines 48-54 cite declared + inverse issue link rows, one event, both-direction confirmation, and duplicate rejection. `cmd/fid/link_test.go` lines 39-76 assert one journal line, exactly two link rows, and exactly two `links` insert effects. |
| VAL-006 | task-2 | ON TRACK | Task 2 report lines 55-60 cite issue -> requirement `implements`, `implemented_by` inverse, and missing requirement rejection before append. |
| VAL-007 | task-2 | ON TRACK | Task 2 report lines 61-66 cite incompatible relation/endpoint cases and missing endpoint rejection with no journal growth. |
| VAL-008 | task-2 | ON TRACK | Task 2 report lines 67-74 cite registry equality plus unknown relation and self-link rejection. Spot check passed: `go test ./internal/links -run 'TestRelationInverseMapMatchesDocumentedV01Set|TestAddSupportsInverseDeclaredRelationAndRequirementEndpoints|TestAddRejectsInvalidRequestsBeforeAppend'` -> `ok github.com/mbell-nzl/fid/internal/links`. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 format | `gofmt -l cmd internal` | PASS | Recorded in task-1-done lines 71-71. |
| Task 1 lint | `golangci-lint run ./...` | PASS | Recorded in task-1-done line 72 as `0 issues.` |
| Task 1 race/coverage test | `go test -race -coverprofile=coverage.out ./...` | PASS | Recorded in task-1-done line 73. |
| Task 1 coverage report | `go tool cover -func=coverage.out` | PASS | Recorded in task-1-done lines 74-78 with total 76.5%, `cmd/fid` 78.1%, and domain packages at or above the cited thresholds. |
| Task 1 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required by the task spec for new CLI command registration; recorded in task-1-done line 79. |
| Task 1 conditional integration | `go test -tags integration -race ./...` | N/A | Recorded in task-1-done lines 80-81 with reason: no integration fixtures touched. |
| Task 1 deterministic spot check | `go test ./cmd/fid -run 'TestRequirementSyncCreatesPointersInOneJournalEvent|TestRequirementSyncRejectsMalformedInputBeforeAppend|TestRequirementSyncAuthorizationAndPolicyShowParity'` | PASS | Review-run representative check for VAL-002..VAL-004 and policy-show parity. |
| Task 1 design/parity check | `05_DESIGN_PLAN.md` CLI map; `06_DESIGN_PLAN_CLI.md` `fid requirement sync` | PASS | `05_DESIGN_PLAN.md` only maps the command tree around links; Task 1 detailed behavior is in `06_DESIGN_PLAN_CLI.md` `fid requirement sync`. Task 1 report lines 33-35 says no CLI spec change was needed. |
| Task 2 format | `gofmt -l cmd internal` | PASS | Recorded in task-2-done line 87. |
| Task 2 lint | `golangci-lint run ./...` | PASS | Recorded in task-2-done line 88 as `0 issues.` |
| Task 2 race/coverage test | `go test -race -coverprofile=coverage.out ./...` | PASS | Recorded in task-2-done line 89. |
| Task 2 coverage report | `go tool cover -func=coverage.out` | PASS | Recorded in task-2-done lines 90-94 with total 76.6%, `cmd/fid` 78.2%, `internal/links` 81.6%, and relevant domain packages at or above the cited thresholds. |
| Task 2 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required by the task spec for new CLI command registration; recorded in task-2-done line 95. |
| Task 2 conditional integration | `go test -tags integration -race ./...` | N/A | Recorded in task-2-done lines 96-97 with reason: no integration fixtures touched. |
| Task 2 deterministic spot check | `go test ./internal/links -run 'TestRelationInverseMapMatchesDocumentedV01Set|TestAddSupportsInverseDeclaredRelationAndRequirementEndpoints|TestAddRejectsInvalidRequestsBeforeAppend'` | PASS | Review-run representative check for registry, add success, endpoint compatibility, and before-append rejection. |
| Task 2 policy spot check | `go test ./internal/policy -run 'TestMatrixRequirementSyncRowMatchesSchemaSection13|TestMatrixLinkRowsMatchSchemaSection13'` | PASS | Review-run representative check for the Phase 12 matrix rows. |
| Task 2 design/parity check | `05_DESIGN_PLAN.md` link command map; `06_DESIGN_PLAN_CLI.md` `fid link add`; `03_TRD_SCHEMA.md` sections 10, 13, 14.2-14.3 | PASS | `05_DESIGN_PLAN.md` lines 93-96 map link add/remove/list behavior. Detailed Task 2 behavior is in `06_DESIGN_PLAN_CLI.md` `fid link add` and schema relation/matrix sections. Task 2 report lines 37-38 says no CLI spec change was needed. |
| PI-1 decision record | `project-docs/docs/phases/phase-12/decisions.md` DEC-029 | PASS | Records exact hash input, parser contract, duplicate-ID behavior, and no-migration outcome. |
| PI-2 decision record | `project-docs/docs/phases/phase-12/decisions.md` DEC-030 | PASS | Records closed v0.1 registry, app-layer prechecks, endpoint scope, self-link rejection, and no-migration outcome. |
| Escalation state | `project-docs/docs/current/tl_escalations_current.md` | PASS | Open escalations section is `None`; this review raises none. |

## Per-Task Findings

### Task 1 - Requirement Pointers And `fid requirement sync`

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-12/tasks/task-1-done.md:33` The report says the CLI spec needed no change, but it does not name the specific design-plan sections checked in parity review mode. Fix direction: future parity reports should explicitly list the checked design/spec sections, even when the result is "verified no change."

### Task 2 - Links Domain, Matrix Rows, And `fid link add`

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-12/tasks/task-2-done.md:48` The report proves add emits one event with two `links` inserts, which indirectly supports the QA-plan "no activity-log add" note, and `internal/links/add.go` lines 81-84 has only link insert effects. Fix direction: future reports should state the negative evidence explicitly as "no `activity_log` effect/row on link add" when the QA plan calls it out.

## Cross-Cutting Observations

- The current worktree has modified task done reports for Tasks 1-4 and an untracked Phase 12 QA plan. This review did not alter task reports, specs, decisions, or product code. The requested review artifact is the only file authored by this QA Reviewer.
- Review 1 evidence is sufficient for Tasks 1-2 to proceed. VAL-009..VAL-012 remain outside this review window and are not assessed here.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-12/qa/review-1.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
