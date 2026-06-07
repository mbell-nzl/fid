## Summary of what was implemented

- Added package-level unit coverage for `internal/issues`, `internal/policy`, and `internal/journal` to close the Phase 1 QA-11 / COV-001 blocker without changing production behavior.
- `internal/issues` now has direct tests for `Get` not-found behavior, `Create`/`Get` input guards, unsupported-kind rejection, status-category mapping, nullable field projection, helper branches, and a package-local create/get round trip.
- `internal/policy` now has direct tests for missing/unknown/known role validation, workspace and journal-dir guard paths, role-event insert/update/delete and `active=false` handling, and both sides of TRG-002 corrupt-tail tolerance.
- `internal/journal` existing suites were extended to cover replay nil-connection rejection, ordered journal-file filtering, missing-DB rebuild behavior, replay start-offset handling, corrupt-error formatting/unwrapping, projection update/delete and replay-status branches, upgrade-chain validation/error branches, and writer helper branches.

## Files created or modified

- `internal/issues/issue_test.go`
- `internal/policy/policy_test.go`
- `internal/journal/replay_test.go`
- `internal/journal/project_test.go`
- `internal/journal/upgrade_test.go`
- `internal/journal/writer_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-8-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-8-completion.md`

## Validation Contract Coverage

- No new `VAL-NNN` were introduced. This task strengthens the evidence base for already-PASS Phase 1 assertions whose domain logic lives in the covered packages:
  `VAL-006`, `VAL-008`, `VAL-009`, `VAL-010`, `VAL-011`, `VAL-012`, and `VAL-016`.
- `internal/issues/issue_test.go` reinforces `VAL-006` with package-local coverage for:
  `Get` returning `*NotFoundError`, `Create`/`Get` guard errors, unsupported kinds, status-category mapping, populated nullable projection fields, `issueIDExists`, and `countJournalLines`.
- `internal/policy/policy_test.go` reinforces `VAL-016` and explicitly cites `TRG-002` by asserting both current corrupt-tail behaviors:
  a malformed trailing line after at least one valid role event is tolerated for role discovery, while a malformed first line still surfaces an error.
- `internal/journal/replay_test.go` reinforces `VAL-008`, `VAL-010`, and `VAL-012` with package-local coverage for:
  nil-connection rejection, ordered journal-file filtering, missing-database rebuild, replay start-offset catch-up behavior, and `CorruptJournalError` formatting/unwrapping.
- `internal/journal/project_test.go` reinforces `VAL-008`, `VAL-011`, and `VAL-012` with direct coverage for:
  projection input validation, update/delete effect application, replay-status persistence, `PreconditionConflictError` formatting, `ensureSingleRowAffected`, `validatedIdentifier`, and value normalization helpers.
- `internal/journal/upgrade_test.go` reinforces `VAL-009` with direct coverage for:
  invalid upgrade-chain definitions and upgrade-chain error handling for unsupported schema states.
- `internal/journal/writer_test.go` reinforces `VAL-008` and `VAL-012` by covering:
  the public `AppendEvent` wrapper, `parseEventSequence` error branches, `greatestJournalDate`, and `nextSequence`.
- Red evidence for the remediation:
  `go test -coverprofile=/tmp/issues-red.out ./internal/issues` passed functionally but reported only `11.3%` coverage, proving the package was still below COV-001 after the first minimal slice.
- Green evidence for the remediation:
  focused package coverage reached `internal/issues 85.2%`, `internal/policy 85.0%`, and `internal/journal 80.6%`, and the full repository gate preserved `cmd/fid 70.5%`.

## Tests run

- `go test -coverprofile=/tmp/issues-red.out ./internal/issues`
  PASS command; red acceptance state at `11.3%` coverage after only the first minimal `Get` not-found slice.
- `go test -coverprofile=/tmp/issues-green.out ./internal/issues` - PASS (`85.2%`)
- `go test -coverprofile=/tmp/policy.out ./internal/policy` - PASS (`85.0%`)
- `go test -coverprofile=/tmp/journal.out ./internal/journal`
  PASS after iterative red/green expansion; final focused package coverage `80.6%`.
- `gofmt -w internal/issues/issue_test.go internal/policy/policy_test.go internal/journal/replay_test.go internal/journal/project_test.go internal/journal/upgrade_test.go internal/journal/writer_test.go` - PASS
- `go test -coverprofile=/tmp/task8-focused.out ./internal/issues ./internal/policy ./internal/journal` - PASS
- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

Package coverage from the final `coverage.out` profile used by `go tool cover -func`:

- `github.com/mbell-nzl/fid/internal/issues` - `85.2%`
- `github.com/mbell-nzl/fid/internal/policy` - `85.0%`
- `github.com/mbell-nzl/fid/internal/journal` - `80.6%`
- `github.com/mbell-nzl/fid/cmd/fid` - `70.5%`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `Not triggered` `go test -tags integration -race ./...` because all new coverage remained in the default unit test suite and no `integration` build-tag tests were added.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- To be recorded from `git rev-parse --short HEAD` immediately after commit creation in the session close-out.

## Commit subject

- Planned implementation subject: `test(phase-1): remediate task 8 domain coverage`
