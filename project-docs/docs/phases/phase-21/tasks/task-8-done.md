# Phase 21 - Task 8: Issue Actor And Event Convergence - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Issue create, update, and transition use typed Actor attribution, Event v2, durable idempotency, and replay. Retained CLI tests still invoke removed `--by` transport, so the exact task gate is red.
- Recovery attempts: 1
- Next action: Migrate retained command tests that create, update, or transition Issues to typed `--actor` and `--idempotency-key`, preserve their assertions, then run every exact task gate.
- Human attention required: no

## Summary of what was implemented

- Converted the Issue create CLI transport to require typed standalone `--actor` attribution and `--idempotency-key`.
- Converted Issue create to emit Event v2 and commit/replay through `idempotency_results`.
- Added red/green CLI corpus coverage for accepted create transport and durable replay identity.
- Converted Issue update and transition to typed Actor/Event-v2/idempotent commit and replay, retaining an in-process V1 test-call compatibility bridge only.
- Added `--by` rejection, fixed-clock aggregate ordering/nonfuture, and creation source-budget coverage.

## Files created or modified

- `cmd/fid/issue_create.go`
- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/issue_update.go`
- `cmd/fid/issue_transition.go`
- `cmd/fid/task8_contract_test.go`
- `internal/issues/actor.go`
- `internal/issues/create.go`
- `internal/issues/create_commit.go`
- `internal/issues/idempotency_results.go`
- `internal/issues/issue.go`
- `internal/issues/update.go`
- `internal/issues/transition.go`
- `internal/issues/transition_rules.go`
- `internal/issues/write_helpers.go`
- `internal/issues/task8_test.go`
- `internal/journal/idempotency.go`
- `project-docs/docs/phases/phase-21/tasks/task-8-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-8-done.md` - recovery checkpoint and captured red/green evidence.

## Validation Contract Coverage

- VAL-025: BLOCKED - update, transition, and `--by` rejection corpus tests pass, but retained CLI test convergence and the exact gate remain unfinished.
- VAL-026: BLOCKED - `TestTask8IssueEventsUseFixedClockAndAggregateOrder` passes, but the exact gate is red.
- VAL-027: BLOCKED - `TestTask8IssueCreationSourceBudget` passes, but the completed Issue corpus and exact gate remain red.

## Tests run

- `go test -race ./cmd/fid -run '^TestTask8IssueCreateRequiresStandaloneActorAndIdempotencyKey$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.813s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask8IssueCreateReplaysTheOriginalMutation$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.876s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - BLOCKED - `not run: task implementation is incomplete` (no observed termination)
- `golangci-lint run ./...` - BLOCKED - `not run: task implementation is incomplete` (no observed termination)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `not run: task implementation is incomplete` (no observed termination)
- `go tool cover -func=coverage.out` - BLOCKED - `not run: coverage profile was not generated because the required suite was not run` (no observed termination)

## Deviations from spec

- The first focused create command failed with `FAIL` (terminated `exit_code: 1`) because `--actor` was not declared. The repair added the typed create transport and Event-v2 result metadata; the byte-identical focused command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	1.813s` (terminated `exit_code: 0`).
- The replay-focused command failed with `FAIL` (terminated `exit_code: 1`) because the second create appended a distinct Event. The repair uses `journal.CommitIdempotent`; the byte-identical focused command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	1.876s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
