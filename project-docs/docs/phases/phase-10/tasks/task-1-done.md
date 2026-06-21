# Phase 10 Task 1 Completion Report

## Summary of what was implemented

- Reconciled CF-003 in `internal/policy` so `KnownRoles` no longer tolerates a
  malformed or un-upgradeable trailing journal line after seeded roles have
  already been discovered.
- `collectRoles` now reports malformed JSON with
  `*journal.CorruptJournalError` using `.fid/journal/<file>:<line>` context, and
  reports upgrade failures with replay-style `upgrade event ... from
  .fid/journal/<file>:<line>` context.
- Rewrote the TRG-002 regression test to assert reported corruption per DEC-D,
  preserved the malformed-first-line case, added a clean multi-file/multi-line
  success case, and added an un-upgradeable trailing-event failure case.

## Files created or modified

- Modified `internal/policy/policy.go`
- Modified `internal/policy/policy_test.go`
- Modified `project-docs/docs/phases/phase-10/decisions.md`
- Created `project-docs/docs/phases/phase-10/tasks/task-1-done.md`

## Documentation updated

- Updated `project-docs/docs/phases/phase-10/decisions.md` with the Task 1
  DEC-D execution note and the intentional `TestKnownRolesTRG002CorruptTailTolerance`
  rewrite rationale.
- `verified-no-change`: `project-docs/specs/02_TRD.md` and
  `project-docs/specs/03_TRD_SCHEMA.md` already matched the intended corruption
  contract; no spec edits were required.

## Validation Contract Coverage

- `VAL-001` PASS.
  Evidence:
  `TestKnownRolesTRG002CorruptTailTolerance/valid_role_before_malformed_trailing_line_reports_corruption`
  failed red with `KnownRoles() error = nil, want corrupt journal error`, then
  passed green once `collectRoles` stopped swallowing the corrupt tail.
- `VAL-001` PASS.
  `TestKnownRolesTRG002CorruptTailTolerance/malformed_first_line_surfaces_an_error`
  now asserts `*journal.CorruptJournalError` with `.fid/journal/2026-06-07.jsonl`
  file context and the expected line number.
- `VAL-001` PASS.
  `TestKnownRolesTRG002CorruptTailTolerance/valid_role_before_un-upgradeable_trailing_event_reports_corruption`
  proves an un-upgradeable trailing event is reported instead of skipped and
  names `.fid/journal/2026-06-07.jsonl:2`.
- `VAL-001` PASS.
  `TestKnownRolesCollectsCleanRolesAcrossFilesAndLines` proves a clean
  multi-file/multi-line role journal still resolves the active role set.
- `VAL-001` PASS.
  Manual CLI evidence with the built binary: after `fid init`, appending a
  malformed second line to the journal and running
  `fid issue create --kind task --title 'Corrupt tail' --by developer` exited
  non-zero with:

  ```text
  error: malformed journal entry at .fid/journal/2026-06-20.jsonl:2
    parse error: unexpected end of JSON input
    read stopped - manual repair required
  ```

## Tests run

- `go test ./internal/policy -run TestKnownRolesTRG002CorruptTailTolerance/valid_role_before_malformed_trailing_line_reports_corruption -count=1`
  RED: failed with `KnownRoles() error = nil, want corrupt journal error`.
- `go test ./internal/policy -run 'TestKnownRolesTRG002CorruptTailTolerance|TestKnownRolesCollectsCleanRolesAcrossFilesAndLines' -count=1`
  GREEN.
- `go test ./internal/policy -count=1`
  PASS.
- Manual CLI evidence using `/tmp/fid` on a temp workspace with an appended
  malformed journal tail:
  `/tmp/fid issue create --kind task --title 'Corrupt tail' --by developer`
  returned the expected corruption error and non-zero exit.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS.
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS.
  Key package coverage from command output: `cmd/fid` 77.8%,
  `internal/issues` 81.8%, `internal/journal` 80.3%, `internal/policy` 88.6%,
  `internal/workflow` 95.5%.
- `go tool cover -func=coverage.out` - PASS.
  Confirmed coverage gate remained above the Phase 10 thresholds, including
  `internal/policy` 88.6% and `cmd/fid` 77.8%.
- `go test -tags integration -race ./...` - PASS.
- `go build -o /tmp/fid ./cmd/fid` - PASS.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `716329cc4cbeda174ce97dcb4e6f3192720dc3b9`

## Commit subject

- `fix(policy): report corrupt tails during role discovery`
