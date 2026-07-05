# Phase 16 - Mid-Phase Review 2 (Tasks 3-4)

Reviewer: qa-reviewer
Date: 2026-07-05
Tasks covered: phase-16/task-3 through phase-16/task-4

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-16/tasks/task-3.md`
  - `project-docs/docs/phases/phase-16/tasks/task-4.md`
- Completion reports:
  - `project-docs/docs/phases/phase-16/tasks/task-3-done.md`
  - `project-docs/docs/phases/phase-16/tasks/task-4-done.md`
- Prior review:
  - `project-docs/docs/phases/phase-16/qa/review-1.md`
- Validation contract:
  - `project-docs/docs/phases/phase-16/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-16/qa/qa-plan.md`
- Design parity checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` command tree entry for `activity list` and `sql`
  - `project-docs/specs/05_DESIGN_PLAN.md` §4 `SQL escape hatch` interaction pattern
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` §8 `fid activity list`
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` §15 `fid sql`
- Commits or changed files:
  - `fd98a020e3a47027749e2695af5ed346b97e4aaa` `feat(activity): finish task 3 lifecycle coverage`
  - `9cae372284f6c7d3caa861d9a3e3a86e6a07b7cf` `feat: add read-only fid sql command`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-006, VAL-007, VAL-008, VAL-009

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-3 | ON TRACK | Task 3 report records focused domain/CLI coverage for `concern raise`, `concern address`, `issue unclaim`, `link add`, and `dev gate`, including `gate_recorded` under DEC-043 and zero-row rejection paths. The lifecycle e2e expects `gate_recorded`, `claim_released`, `link_added`, `concern_raised`, and `concern_addressed` in the timeline. Reviewer rerun: `go test -tags integration ./cmd/fid -run TestV01EActivityLifecycleCoverageAndRebuildParity` passed. |
| VAL-007 | task-3 | ON TRACK | `cmd/fid/e2e_v01e_activity_test.go` drives create, `dev start`, AC add/verify, `dev gate`, `dev finish`, `link add`, unclaim, and concern raise/address. It snapshots ordered `activity_log`, deletes `fid.sqlite`, runs `fid rebuild --by workflow_runner`, and asserts both ordered `activity_log` and `fid activity list --format json` match after rebuild. Reviewer rerun passed. |
| VAL-008 | task-4 | ON TRACK | Task 4 report records SELECT, EXPLAIN, read PRAGMA, `WITH ... SELECT`, and all output formats. Reviewer rerun: `go test ./cmd/fid ./internal/sqlquery -run 'TestSQLCommandAcceptsReadStatementsAndFormats|TestSQLCommandRejectsMutationsWithDocumentedError|TestClassifyAllowsReadStatements|TestClassifyRejectsMutationsAndMultiStatementInput|TestReadOnlyConnectionRejectsMutationWhenClassifierIsBypassed'` passed. |
| VAL-009 | task-4 | ON TRACK | `cmd/fid/sql_integration_test.go` attempts INSERT/UPDATE/DELETE/DROP and checks unchanged projected row counts; `internal/sqlquery` tests cover multi-statement and `WITH ... DELETE` rejection plus a direct mutation attempt against the read-only connection. The implementation opens SQLite with `mode=ro` and `_pragma=query_only(1)`. Reviewer reruns for the targeted SQL unit/CLI tests and `go test -tags integration ./cmd/fid -run TestSQLCommandRejectsMutationAttemptsWithoutChangingProjectedCounts` passed. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 3 required gate reporting | `task-3-done.md` exact quality-gate section | PASS | Report records `gofmt -l cmd internal`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, required integration gate, and build gate as PASS. `report-lint.sh` passed. |
| Task 4 required gate reporting | `task-4-done.md` exact quality-gate section | PASS | Report records every task/QA-plan required command, including required integration and build gates. `report-lint.sh` passed. |
| Commit verification | `git rev-parse --verify fd98a02^{commit}` and `git rev-parse --verify 9cae372^{commit}` | PASS | Both hashes resolved to the full commits listed in scope. |
| Representative format gate | `gofmt -l cmd internal` | PASS | No files listed. |
| Task 3 deterministic rerun | `go test -tags integration ./cmd/fid -run TestV01EActivityLifecycleCoverageAndRebuildParity` | PASS | Confirms DEC-043 `gate_recorded` coverage and VAL-007 rebuild parity at current HEAD. |
| Task 4 deterministic rerun | `go test -tags integration ./cmd/fid -run TestSQLCommandRejectsMutationAttemptsWithoutChangingProjectedCounts` | PASS | Confirms TRD-SEC-002 mutation rejection and unchanged projected row count for INSERT/UPDATE/DELETE/DROP. |
| Task 4 classifier/backstop rerun | `go test ./cmd/fid ./internal/sqlquery -run 'TestSQLCommandAcceptsReadStatementsAndFormats|TestSQLCommandRejectsMutationsWithDocumentedError|TestClassifyAllowsReadStatements|TestClassifyRejectsMutationsAndMultiStatementInput|TestReadOnlyConnectionRejectsMutationWhenClassifierIsBypassed'` | PASS | Confirms accepted read forms, documented error text, multi-statement/`WITH ... DELETE` rejection, and read-only connection behavior. |
| DEC-044 fixture remediation check | `go test ./cmd/fid -run 'TestActivityList'` | PASS | Review 1's known monotonic-date fixture problem no longer reproduces; Task 4's constrained fixture remediation preserved the Task 1 activity CLI tests. |

## Per-Task Findings

### Task 3 - CF-011 backfill part 2 + timeline-consistency e2e

#### Blocking

- None.

#### Improvement

- None.

Observation: Task 3 is coherent with Review 1 and DEC-043. The e2e explicitly includes `gate_recorded` in the expected full lifecycle timeline and uses ordered rebuild equality, so the previous lifecycle gap is covered without changing the `gate_runs` persistence contract.

### Task 4 - `fid sql` read-only escape hatch

#### Blocking

- None.

#### Improvement

- `internal/sqlquery/sqlquery.go:1` is 595 lines, above the project non-test file target of 400 lines. This is not blocking because behavior and gates are covered, but the scanner/tokenizer and execution/result mapping are separable. Fix direction: when `internal/sqlquery` is next reopened, split classifier/tokenizer/read execution into cohesive files without changing the public package API.

Observation: Task 4 resolves Review 1's DEC-044 fixture issue in the narrow authorized test scope, and the SQL implementation matches the design requirement for friendly classifier rejection plus a read-only SQLite backstop.

## Cross-Cutting Observations

- No new blocking findings, major issues, or escalation triggers were found. The window's required commands are reported PASS, representative reruns passed, and the DEC-043/DEC-044 history is reflected in the task reports.
- No architecture, data-model, public-interface, security/privacy-severity, acceptance-criteria ambiguity, or quality-gate anomaly was found that requires appending an `ESC-NNN`.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-16/qa/review-2.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
