# Phase 16 - Mid-Phase Review 1 (Tasks 1-2)

Reviewer: qa-reviewer
Date: 2026-07-05
Tasks covered: phase-16/task-1 through phase-16/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-16/tasks/task-1.md`
  - `project-docs/docs/phases/phase-16/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-16/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-16/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-16/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-16/qa/qa-plan.md`
- Design parity checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` command tree entry for `activity list`
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` §8 `fid activity list`
- Commits or changed files:
  - `977974c0ec00d3d9a63f7d7515e633683c9b9ac5` `feat: add activity timeline read command`
  - `931d83a5190e0f4fc74d4d0edf1308e8c49d260d` `feat(activity): backfill issue lifecycle writers`

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | Task 1 report cites shape/order/details tests and all-format dispatch plus built-binary smoke; implementation selects documented fields, deserialises `details`, and orders by `occurred_at ASC, id ASC` in `internal/activity/activity.go`. Reviewer rerun: `go test ./internal/activity -run TestListReturnsChronologicalEntriesWithStructuredDetails` passed as part of the targeted `internal/activity` package run. Current-head CLI rerun is affected by the known DEC-044 fixture issue, not a new Review 1 escalation. |
| VAL-002 | task-1 | ON TRACK | Task 1 report cites domain and CLI `--on` include/exclude coverage. SQL uses parameterised `ref_type = ?` and `ref_id = ?`; `06_DESIGN_PLAN_CLI.md` §8 uses the same `<type>/<id>` contract. |
| VAL-003 | task-1 | ON TRACK | Task 1 report cites exact-boundary and date-only tests. Implementation accepts `YYYY-MM-DD` and RFC3339 and applies `occurred_at >= ?`, satisfying boundary include/exclude semantics. |
| VAL-004 | task-1 | ON TRACK | Task 1 report cites matching/non-matching event-filter tests. Implementation applies parameterised `event = ?` filtering. |
| VAL-005 | task-2 | ON TRACK | Task 2 report records domain/CLI tests and runtime smoke showing `issue_created`, `ac_added`, `ac_verified`, and `task_completed` entries. Reviewer rerun of Task 2 domain/dev evidence passed: `go test ./internal/issues ./internal/dev -run 'TestCreateAppendsIssueCreatedActivityInSameEvent|TestAddAcceptanceCriterionPersistsAndReadsBackInOrdinalOrder|TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent|TestCompleteCreatesAttemptOneAndPersistsCompletionFields|TestCompleteStoresHumanInterventionStateAndProjectsActivity|TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent|TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim'`. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required gate reporting | `task-1-done.md` exact quality-gate section | PASS | Report records `gofmt -l cmd internal`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, integration gate as not required with reason, and `go build -o /tmp/fid ./cmd/fid`. `report-lint.sh` passed. |
| Task 2 required gate reporting | `task-2-done.md` exact quality-gate section | PASS | Report records every task/QA-plan required command, including required `go test -tags integration -race ./...` and build gate. `report-lint.sh` passed. |
| Commit verification | `git rev-parse --verify 977974c^{commit}` and `git rev-parse --verify 931d83a^{commit}` | PASS | Both hashes resolved to full commits listed above. |
| Representative required command | `gofmt -l cmd internal` | PASS | No files listed. |
| Task 1 deterministic rerun | `go test ./internal/activity ./cmd/fid -run 'TestListReturnsChronologicalEntriesWithStructuredDetails|TestActivityListDispatchesFormats|TestActivityListAppliesFilters'` | PARTIAL | `internal/activity` passed. `cmd/fid` failed in `TestActivityListAppliesFilters` with `system clock is behind latest journal date 2026-07-05; refusing to append out of order`; this is the current-head fixture defect already recorded in `phase_brief.md` under DEC-044, so no duplicate ESC is raised in this Review 1. |
| Task 2 deterministic rerun | `go test ./internal/issues ./internal/dev -run 'TestCreateAppendsIssueCreatedActivityInSameEvent|TestAddAcceptanceCriterionPersistsAndReadsBackInOrdinalOrder|TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent|TestCompleteCreatesAttemptOneAndPersistsCompletionFields|TestCompleteStoresHumanInterventionStateAndProjectsActivity|TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent|TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim'` | PASS | Confirms one-event activity rows for create/ac add/ac verify/complete and PI-6 inheritance for `dev finish`/`dev block`. |

## Per-Task Findings

### Task 1 - `fid activity list` projected timeline read

#### Blocking

- None.

#### Improvement

- `cmd/fid/activity_test.go:23` uses fixed July 2026 seeded journal events; at current HEAD the representative CLI rerun fails the journal monotonic-date guard after `fid init`. Fix direction: complete the already-authorized DEC-044 fixture remediation by making the seeded activity events monotonic relative to the initialized workspace instead of relying on fixed dates.

Observation: Task 1 tracks the design surface: `05_DESIGN_PLAN.md` lists `activity list` in the command tree, and `06_DESIGN_PLAN_CLI.md` §8 matches the implemented shape, filters, and formats. The implementation keeps read behavior in `internal/activity` and the CLI adapter in `cmd/fid`, so the new read package boundary is consistent with PI-4.

### Task 2 - CF-011 backfill part 1 - issue create / ac add / ac verify / complete

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-16/tasks/task-2-done.md:40` states rejection-path tests continue to prove zero journal append, but the report does not enumerate the rejection scenario for each of the four VAL-005 writers. Fix direction: in future CF-011 reports, list the per-writer rejection fixtures explicitly so the binary row-vs-zero evidence is self-contained without test spelunking.

Observation: Task 2 satisfies the one-event/zero-on-reject direction for the reviewed window. The code appends one `activity_log` effect into the existing prepared event for create, AC add, AC verify, and completion, and `issues.BuildCompletionEffects` carries `task_completed` into `dev finish`/`dev block` while those composites still append one event.

## Cross-Cutting Observations

- The known `cmd/fid/activity_test.go` date fixture failure is already recorded in the current phase brief and covered by DEC-044; it is not a new Review 1 escalation. It should still be resolved before Task 4 close because it blocks global gate credibility.
- No architecture, data-model, security/privacy, public-interface, or acceptance-criteria ambiguity was found in Tasks 1-2 that requires a new TL-Lite / ESC entry.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-16/qa/review-1.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
