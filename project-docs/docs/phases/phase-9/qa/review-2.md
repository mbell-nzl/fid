# Phase 9 - Mid-Phase Review 2 (Task 4)

Reviewer: QA Reviewer
Date: 2026-06-20
Tasks covered: phase-9/task-4

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-9/tasks/task-4.md`
- Completion reports:
  - `project-docs/docs/phases/phase-9/tasks/task-4-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-9/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-9/qa/qa-plan.md`
- Decisions:
  - `project-docs/docs/phases/phase-9/decisions.md` DEC-023, DEC-024
- Escalations:
  - `project-docs/docs/current/tl_escalations_current.md` ESC-006 direction
- Design sections checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` sections 2, 3.1, 3.2, 5.1, 5.2, 5.3
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 6 `fid concern list` / `fid concern address`
- Commits or changed files:
  - `c927de0f81fc43133d45485d789d5d45c2aadd94` `fix: expose concern address read-back deterministically`
  - `internal/issues/concern_list.go`
  - `internal/issues/detail.go`
  - `internal/issues/concern_list_test.go`
  - `internal/issues/issue_test.go`
  - `cmd/fid/concern_test.go`
  - `cmd/fid/e2e_v01c_concern_test.go`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-004, VAL-008, VAL-009

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-004 | task-4 | ON TRACK | `ConcernListItem` exposes `addressed_in_completion_id` and scans the nullable column in `internal/issues/concern_list.go:21`; `IssueConcern` exposes the same field in `internal/issues/detail.go:25`. Tests cover null for open concerns and the addressed `TC-` id after address in `internal/issues/concern_list_test.go:80`, `internal/issues/concern_list_test.go:90`, `internal/issues/issue_test.go:822`, `internal/issues/issue_test.go:934`, and the CLI parity test in `cmd/fid/concern_test.go:205`. Reviewer reran `go test ./internal/issues -run 'TestGetReturnsEnrichedConcernFields|TestGetTreatsAddressedConcernsAsNotOpen|TestGetConcernsOrderSameSecondRowsByAppendOrder|TestListConcernsFiltersByIssueAndStatus|TestListConcernsOrdersSameSecondRowsByAppendOrder' -count=1` with PASS, `go test ./cmd/fid -run 'TestConcernAddressUpdatesConcernAndPrintsConfirmation|TestConcernListFiltersByIssueAndStatus' -count=1` with PASS, and `go test -tags integration ./cmd/fid -run 'TestV01CConcernRebuildParity|TestV01CClosureSmoke' -count=1` with PASS. |
| VAL-008 | task-4 | ON TRACK | Filters remain covered in `internal/issues/concern_list_test.go:5`, with issue/status/combined/empty-match assertions and exact newest-first order at `internal/issues/concern_list_test.go:67`. The public query now orders by `rowid DESC` at `internal/issues/concern_list.go:122`, and same-second append ordering is asserted in `internal/issues/concern_list_test.go:145`. Reviewer reran the internal and CLI targeted tests above with PASS. |
| VAL-009 | task-4 | ON TRACK | Rebuild parity includes `addressed_in_completion_id` in the snapshot struct and query at `cmd/fid/e2e_v01c_concern_test.go:307` and `cmd/fid/e2e_v01c_concern_test.go:329`. Reviewer reran `go test -tags integration ./cmd/fid -run 'TestV01CConcernRebuildParity|TestV01CClosureSmoke' -count=1` with PASS. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Commit verification | `git rev-parse --verify c927de0^{commit}` | PASS | Resolved to `c927de0f81fc43133d45485d789d5d45c2aadd94`. |
| Commit subject | `git show -s --format=%H%n%s%n%b c927de0` | PASS | Subject is `fix: expose concern address read-back deterministically`; body includes `Refs: phase-9/task-4`. |
| Report headings | `task-4-done.md` | PASS | Required sections are present: summary, files, docs, validation coverage, tests, exact quality-gate results, deviations, concerns/technical debt, dependency note, commit hash, and commit subject. |
| Report concerns section | `task-4-done.md` `## Concerns, technical debt, or notes for the Tech Lead` | PASS | Section is non-empty and records `None.` |
| Required command coverage | Task spec exact commands vs `task-4-done.md` | PASS | Report records PASS for all seven required commands: `gofmt -l cmd internal`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, `go test -tags integration -race ./...`, `go test -bench . -benchmem ./...`, and `go build -o /tmp/fid ./cmd/fid`. |
| Coverage thresholds | `task-4-done.md` exact gate results | PASS | Reported `cmd/fid` 77.8%, `internal/issues` 81.8%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%, satisfying DEC-007 thresholds. |
| Public read field, `fid concern list` | `internal/issues/concern_list.go:21` and `internal/issues/concern_list.go:112` | PASS | Field is selected, scanned through `nullStringPointer`, and serialized as `json:"addressed_in_completion_id"`. |
| Public read field, `fid issue get --section concerns` | `internal/issues/detail.go:25` and `internal/issues/detail.go:210` | PASS | Field is selected, scanned through `nullStringPointer`, and serialized as `json:"addressed_in_completion_id"`. |
| Deterministic ordering | `ORDER BY rowid DESC` | PASS | Present on both read queries at `internal/issues/concern_list.go:122` and `internal/issues/detail.go:213`. Same-second order tests exist for list and issue get at `internal/issues/concern_list_test.go:145` and `internal/issues/issue_test.go:939`. |
| Benchmark gate | `go test -bench . -benchmem ./...` | PASS | Reviewer rerun exited 0. Output included `ok github.com/mbell-nzl/fid/cmd/fid 1.812s` and journal benchmarks `BenchmarkRebuildWorkspace1kEvents-6` / `BenchmarkRebuildWorkspace10kEvents-6`; the prior `TestConcernListFiltersByIssueAndStatus` benchmark-run failure is cleared. |
| VAL-004 targeted tests | `go test ./cmd/fid -run 'TestConcernAddressUpdatesConcernAndPrintsConfirmation|TestConcernListFiltersByIssueAndStatus' -count=1` | PASS | Exited 0: `ok github.com/mbell-nzl/fid/cmd/fid 0.127s`. |
| VAL-004/008 internal tests | `go test ./internal/issues -run 'TestGetReturnsEnrichedConcernFields|TestGetTreatsAddressedConcernsAsNotOpen|TestGetConcernsOrderSameSecondRowsByAppendOrder|TestListConcernsFiltersByIssueAndStatus|TestListConcernsOrdersSameSecondRowsByAppendOrder' -count=1` | PASS | Exited 0: `ok github.com/mbell-nzl/fid/internal/issues 0.141s`. |
| VAL-009 integration tests | `go test -tags integration ./cmd/fid -run 'TestV01CConcernRebuildParity|TestV01CClosureSmoke' -count=1` | PASS | Exited 0: `ok github.com/mbell-nzl/fid/cmd/fid 0.161s`. |

## Per-Task Findings

### Task 4 - concern read-back `addressed_in_completion_id` + deterministic ordering

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-9/tasks/task-4-done.md` records the short commit hash `c927de0`; it verified uniquely in git for this review, but full 40-character hashes are more durable for later artifact audits. Fix direction: use the full hash in future completion reports while keeping the subject copied from git output.

## Cross-Cutting Observations

- ESC-006 remains open in `project-docs/docs/current/tl_escalations_current.md` with direction issued and Task 4 as the closeout path. This review found no new blocker and no new escalation trigger; Tech Lead can close ESC-006 during Phase 9 wrap once this review is accepted.
- DEC-023 and DEC-024 are implemented within the intended scope: read-model field exposure only, no schema change, no workflow-edge change, and no activity-log behavior change.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-9/qa/review-2.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
