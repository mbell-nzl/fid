# Phase 16 - Task 5 Completion Report

## Summary of what was implemented

- Preserved the existing Task 5 implementation already present in the worktree:
  `fid journal inspect <event-id>` now reads journal events directly, reports
  file and line location, canonical hash parity, actor/command/effects, and the
  PI-5 revert verdict across the `internal/journal` reader and `cmd/fid`
  surface.
- Completed the DEC-046 unblock by remediating the deterministic fixture family
  in `internal/activity/activity_test.go` only. Seeded activity test events now
  derive monotonic journal dates and IDs from the initialized workspace instead
  of hardcoded backdated timestamps, preserving VAL-001..004 and PRD-ERR-014.
- Re-ran the exact Task 5 gate to green and exercised the built
  `fid journal inspect` entrypoint against a fresh workspace using a real
  create -> claim -> transition -> inspect flow.

## Files created or modified

- `cmd/fid/journal.go`
- `cmd/fid/journal_test.go`
- `cmd/fid/root.go`
- `internal/activity/activity_test.go`
- `internal/journal/inspect.go`
- `internal/journal/inspect_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-16/tasks/task-5-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
  - Confirmed the `fid journal inspect` revert-verdict buckets for inverse,
    retraction, and `eligible:false` correction cases.
- `project-docs/docs/phases/phase-16/tasks/task-5-done.md`
  - Replaced the DEC-046 blocker report with this normal completion report.

## Validation Contract Coverage

- `VAL-010` PASS.
  - `internal/journal/inspect_test.go` and `cmd/fid/journal_test.go` cover
    known-event inspection, typed not-found behavior, canonical hash parity with
    `journal.HashEvent`, and all four output formats.
- `VAL-011` PASS.
  - `internal/journal/inspect_test.go` and `cmd/fid/journal_test.go` cover the
    event-type-keyed revert verdicts for inverse, retraction, and
    `eligible:false` correction cases.
- DEC-046 regression preservation PASS.
  - `go test ./internal/activity -run TestListFiltersByOnSinceAndEvent` failed
    before the fixture fix with `system clock is behind latest journal date ...
    refusing to append out of order`, then passed after the helper change. This
    preserved Task 1 VAL-001..004 behavior and PRD-ERR-014 without changing
    production journal ordering.
- Coverage percentages from the required coverage gate:
  - `cmd/fid`: `77.5%`
  - `internal/journal`: `78.9%`

PI-5 event-type -> verdict table:

| Event type | Verdict |
| --- | --- |
| `project.initialized` | `eligible:false` |
| `issue.created` | `eligible:false` |
| `issue.transitioned` | `inverse` |
| `issue.updated` | `inverse` |
| `issue.concern_addressed` | `inverse` |
| `issue.label.added` | `inverse` |
| `issue.label.removed` | `inverse` |
| `link.added` | `inverse` |
| `link.removed` | `inverse` |
| `issue.ac_added` | `retraction` |
| `issue.ac_verified` | `retraction` |
| `issue.concern_raised` | `retraction` |
| `dev.gate.recorded` | `retraction` |
| `label.added` | `eligible:false` |
| `issue.completed` | `eligible:false` |
| `issue.claim.acquired` | `eligible:false` |
| `issue.claim.renewed` | `eligible:false` |
| `issue.claim.released` | `eligible:false` |
| `run.recorded` | `eligible:false` |
| `requirement.synced` | `eligible:false` |
| `dev.started` | `eligible:false` |
| `dev.finished` | `eligible:false` |
| `dev.blocked` | `eligible:false` |
| unknown / future type | safe default `eligible:false` |

## Tests run

- Preserved Task 5 feature evidence from the prior blocked run:
  - `go test ./internal/journal -run 'TestFindEventReturnsStoredEventLocationAndHash|TestFindEventReturnsTypedNotFoundError|TestRevertVerdictClassifiesKnownEventTypes|TestRevertVerdictDefaultsUnknownTypesToSafeNotEligible'`
    - Red: failed with undefined `FindEvent`, `EventNotFoundError`, and
      `RevertVerdict`.
    - Green: passed after the journal inspect reader and verdict classifier were
      implemented.
  - `go test ./cmd/fid -run TestJournalInspect`
    - PASS.
- DEC-046 blocker remediation:
  - `go test ./internal/activity -run TestListFiltersByOnSinceAndEvent`
    - Red: failed with `system clock is behind latest journal date 2026-07-05;
      refusing to append out of order`.
    - Green: passed after the activity fixture derived dates and IDs from the
      workspace journal.
- Required task gate:
  - `gofmt -l cmd internal`
  - `golangci-lint run ./...`
  - `go test -race -coverprofile=coverage.out ./...`
  - `go tool cover -func=coverage.out`
  - `go build -o /tmp/fid ./cmd/fid`
- Real entrypoint exercise:
  - `/tmp/fid init`
  - `/tmp/fid issue create --kind task --title "Smoke inspect event" --by developer`
  - `/tmp/fid issue claim <id> --by developer`
  - `/tmp/fid issue transition <id> --to in_progress --by developer`
  - `/tmp/fid journal inspect EVT-20260705-000004 --format json`
  - Observed result: returned the documented JSON object for
    `issue.transitioned`, including `journal_file`, `line`, canonical `hash`,
    ordered `effects`, and `revert.mode:"inverse"` with restore and
    precondition data.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS
- `go test -tags integration -race ./...` - Not required
  - Reason: Task 5 added no integration-tagged tests, and the DEC-046 resume
    changed only `internal/activity/activity_test.go`.
- `go build -o /tmp/fid ./cmd/fid` - PASS

## Deviations from spec

- DEC-046 added a constrained out-of-scope remediation window to unblock the
  required global gates. The fixture fix touched only
  `internal/activity/activity_test.go`.
- The existing Task 5 feature implementation in `cmd/fid/journal.go`,
  `cmd/fid/journal_test.go`, `cmd/fid/root.go`, `internal/journal/inspect.go`,
  `internal/journal/inspect_test.go`, and
  `project-docs/specs/06_DESIGN_PLAN_CLI.md` was preserved from the earlier
  blocked run per the resume instruction; this session did not redo that
  product work.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `12c5107`

## Commit subject

- `feat: add fid journal inspect command`
