# Phase 21 - Task 6: One-Revision Envelope Values - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: Window 2 QA review.
- Human attention required: no

## Summary of what was implemented

- Added `fid envelope create`, `fid envelope get`, and `fid envelope activate` for exactly one primary Issue, revision 1, no related Issues, and one registration-bound `gate-evidence` Obligation.
- Persisted ordered sources and Authority Grants, immutable null-counter budgets, and the dual terminal-success/completion guarded-action value through Event v2 effects.
- Rejected widened input, including Routine, revisions, cancellation, protected Decision, Stop, Attention, mutable budgets, extra Issues, and unsupported Obligation values before an Event is appended.

## Files created or modified

- `cmd/fid/envelope.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/root.go`
- `cmd/fid/contract_surface_test.go`
- `internal/workflow/envelope.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-6-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-6-done.md`

## Validation Contract Coverage

- VAL-019: PASS - `TestEnvelopeCreateGetActivateRoundTripAndRejectsWidening` in `cmd/fid/envelope_test.go:13` creates revision 1 for one existing Issue and rejects related Issues, a second Obligation, and revision input without journal growth.
- VAL-020: PASS - `cmd/fid/envelope_test.go:89` asserts the persisted `gate-evidence` Obligation has exactly `envelope.complete` and `issue.transition-terminal-success` guards; its table rejects a computed status, another type, and widened guards.
- VAL-021: PASS - `cmd/fid/envelope_test.go:97` verifies the created, fetched, and activated Envelope budgets retain `consumed:null` and `remaining:null`; Run creation and all budget consumption remain unavailable for the later Run task.
- VAL-022: PASS - `cmd/fid/envelope_test.go:42` covers Routine, revision, cancellation, protected Decision, Stop, and Attention input alongside other widening attempts, asserting canonical rejection and unchanged Event count for each.

## Tests run

- Red evidence: `go test -race ./cmd/fid -run TestEnvelopeCreateGetActivateRoundTripAndRejectsWidening -count=1` printed `FAIL` before the Envelope command existed (terminated `exit_code: 1`).
- Green evidence: `go test -race ./cmd/fid -run TestEnvelopeCreateGetActivateRoundTripAndRejectsWidening -count=1` printed `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.699s` (terminated `exit_code: 0`).
- Real entrypoint: `go build -o /tmp/fid ./cmd/fid` completed with `exit 0 (no output)` (terminated `exit_code: 0`); `/tmp/fid --json envelope create` emitted the canonical `invalid-argument` `envelope.create` failure for omitted required flags (terminated `exit_code: 1`).

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t142.905s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 38.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t70.7%` (terminated `exit_code: 0`)

## Deviations from spec

- The first required package gate printed `FAIL` (terminated `exit_code: 1`) because `TestReasonCodeDetailSchemaIsExactForEveryCode` still enumerated 28 codes after task 6 introduced the contract-required `envelope-not-found` code. The test inventory was extended with the exact new code/schema; the byte-identical later command printed `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t141.836s` (terminated `exit_code: 0`).
- The first lint attempt printed `* errcheck: 3` (terminated `exit_code: 1`), and its byte-identical rerun printed `* errcheck: 1` (terminated `exit_code: 1`), because four new row loaders did not discard `rows.Close` errors using the local checked-close convention. Each deferred close now explicitly discards its result; the byte-identical final rerun printed `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
