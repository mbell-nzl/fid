# Phase 21 - Task 8: Issue Actor And Event Convergence - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: protected-contract-conflict
- Condition: The accepted Phase 21 interface requires Event v2, kebab-case activity values, and attribution-only standalone Actors, while retained tests assert incompatible V1 Event envelopes, snake-case activity values, and role-based standalone authority. Migrating invocations exposes those assertions; changing them violates the assigned requirement to preserve assertions.
- Recovery attempts: 2
- Resolution action: The Tech Lead must authorize either bounded retained-test assertion migration to the accepted Phase 21 contract or a protected change to that contract. Do not reintroduce `--by`, V1 Event fields, snake-case activity values, or role authority without that authority.
- Human attention required: yes

## Summary of what was implemented

- Converted the Issue create CLI transport to require typed standalone `--actor` attribution and `--idempotency-key`.
- Converted Issue create to emit Event v2 and commit/replay through `idempotency_results`.
- Added red/green CLI corpus coverage for accepted create transport and durable replay identity.
- Converted Issue update and transition to typed Actor/Event-v2/idempotent commit and replay, retaining an in-process V1 test-call compatibility bridge only.
- Added `--by` rejection, fixed-clock aggregate ordering/nonfuture, and creation source-budget coverage.
- Migrated retained Issue create/update/transition command transport from `--by` to typed `--actor`; test runners attach unique idempotency keys to those typed invocations.
- Extended the retained in-process bridge to generate unique idempotency keys when the legacy `ActorRole` test-only call shape is used.

## Files created or modified

- `cmd/fid/issue_create.go`
- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/issue_update.go`
- `cmd/fid/issue_transition.go`
- `cmd/fid/task8_contract_test.go`
- `cmd/fid/*_test.go` retained Issue mutation callers
- `cmd/fid/integration/*_test.go` retained Issue mutation callers
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

- `project-docs/docs/phases/phase-21/tasks/task-8-done.md` - protected-conflict checkpoint and captured red/green and gate evidence.

## Validation Contract Coverage

- VAL-025: BLOCKED - `--by` rejection and typed Issue transport are implemented, but retained assertions require removed V1 Event/authority behavior and the exact gate fails.
- VAL-026: BLOCKED - `TestTask8IssueEventsUseFixedClockAndAggregateOrder` passes, but retained activity tests still assert pre-convergence timestamps and snake-case values.
- VAL-027: BLOCKED - `TestTask8IssueCreationSourceBudget` passes, but the exact gate fails on retained incompatible assertions and lint findings.

## Tests run

- `/tmp/fid-task8 init` in `/tmp/fid-task8-entrypoint.bMFurx` - BLOCKED - `error: invalid project initialization argument` (terminated `exit_code: 1`)
- `go test -race ./internal/issues -run '^TestTask8LegacyInProcessIssueMutationsReceiveIdempotencyKeys$'` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race ./internal/issues -run '^TestTask8LegacyInProcessIssueMutationsReceiveIdempotencyKeys$'` - PASS - `ok github.com/mbell-nzl/fid/internal/issues 1.999s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask8IssueCreateRequiresStandaloneActorAndIdempotencyKey$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.813s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask8IssueCreateReplaysTheOriginalMutation$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.876s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - BLOCKED - `4 issues:` (terminated `exit_code: 1`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 70.6%` (terminated `exit_code: 0`)

## Deviations from spec

- The legacy in-process Issue mutation focused command failed with `FAIL` (terminated `exit_code: 1`) because the retained `ActorRole` bridge did not generate an idempotency key. The repair generates a unique bridge-only key; the byte-identical focused command later passed with `ok github.com/mbell-nzl/fid/internal/issues 1.999s` (terminated `exit_code: 0`).
- The first exact focused gate failed with `FAIL` (terminated `exit_code: 1`) because retained CLI Issue mutations still used `--by` and direct retained Issue callers lacked idempotency keys. This attempt migrated only Issue create/update/transition invocations and added bridge-only key generation; its byte-identical rerun remained `FAIL` (terminated `exit_code: 1`) because retained assertions require V1 Event fields, snake-case activities, and standalone role authority. Updating those assertions is prohibited by the assigned preservation constraint.
- The first focused create command failed with `FAIL` (terminated `exit_code: 1`) because `--actor` was not declared. The repair added the typed create transport and Event-v2 result metadata; the byte-identical focused command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	1.813s` (terminated `exit_code: 0`).
- The replay-focused command failed with `FAIL` (terminated `exit_code: 1`) because the second create appended a distinct Event. The repair uses `journal.CommitIdempotent`; the byte-identical focused command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	1.876s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- Protected resolution is required before retained assertions can be aligned with the accepted Phase 21 interface. The lint gate also reports two staticcheck suggestions in `internal/issues/actor.go` and two unused V1 policy helpers in `internal/issues/policy_actions.go`; no repair was made because the task is blocked before an accepted assertion-migration boundary is supplied.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending blocked checkpoint

## Commit subject

- Pending blocked checkpoint
