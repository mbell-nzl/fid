# Phase 17 - Task 4: Implement Maintenance Warnings and Info

## Objective

Extend the maintenance domain result with the four v0.1 warning families and
the documented housekeeping counts. Make all time boundaries injectable and
the complete result deterministic so CI and rebuild comparisons are stable.

## Validation Contract Coverage

- VAL-007: stale claims warn at/past TTL and live claims do not.
- VAL-008: started-without-claim honors the configured grace in both states.
- VAL-009: only verifications tied to superseded completions are stale.
- VAL-010: only completions without a completion-linked run are telemetry gaps.
- VAL-011: exact info counts and deterministic result ordering.

## Binding Constraints

- DEC-050/051 are binding; task-2 supplies the config field.
- Measure started-without-claim age from the newest append-ordered
  `status_changed` activity row, falling back to `issues.created_at` only for a
  historical row with no transition activity. A non-status issue update must
  not reset the grace window.
- Exact warning check names are `stale_claims`, `started_without_claim`,
  `stale_verifications`, and `telemetry_gaps`.
- A telemetry run corresponds to a completion only through
  `agent_runs.completion_id`, not merely the same `issue_id`.
- Task-4 classifies findings only. Task-5 owns status promotion and exit codes.

## Scope

### Files To Create Or Modify

- `internal/maintenance/warnings.go` - warning queries/classification.
- `internal/maintenance/info.go` - `issue_count`, `claim_count_live`, and
  `recent_activity_7d`.
- `internal/maintenance/check.go` - compose error/warning/info results and sort.
- `internal/maintenance/warnings_test.go` - boundary and both-state fixtures.
- `internal/maintenance/info_test.go` - exact counts and repeated-run equality.
- `internal/maintenance/test_helpers_test.go` - shared deterministic clock and
  fixture additions.
- `project-docs/specs/03_TRD_SCHEMA.md` - clarify app-layer warning semantics if
  task-2's PRD/config sync is insufficient.

### Files To Read For Context

- `internal/config/config.go` (task-2 result)
- `internal/claims/list.go`
- `internal/issues/issue.go`
- `internal/issues/complete.go`
- `internal/runs/record.go`
- `project-docs/specs/GLOSSARY.md` (started-without-claim, stale verification,
  telemetry gap)
- `scaffold/initial-spec.md` section 9.8.1 (historical detail)

### Out Of Scope

- CLI flags, rendering, exit code promotion, or repair.
- Activity projection gap or scratch journal-vs-DB diff checks (v0.2+).
- Cost backfill for completion rows.

## Acceptance Criteria

- [ ] Stale/live claim boundary states are both covered with an injected time.
  (VAL-007)
- [ ] Live claim, before-grace, at/after-grace, fallback-time, and grace=0
  started states are covered; non-status updates do not reset age. (VAL-008)
- [ ] Superseded/current/NULL completion verification states are classified
  exactly once. (VAL-009)
- [ ] Missing link, issue-only link, and completion-linked run states are all
  covered. (VAL-010)
- [ ] Info counts are exact and two identical audits deeply equal. (VAL-011)

## First Failing Test

`TestCheckWarningsClassifiesQualifyingAndNonQualifyingStates` fails because the
maintenance result has no warning checks.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md` - only if needed to keep the canonical
  app-layer audit rules explicit after task-2.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate is required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Depends on tasks 2 and 3.
- Use SQL for candidate selection and Go for domain classification only where
  that keeps workflow/config semantics single-sourced.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-4-done.md`.
