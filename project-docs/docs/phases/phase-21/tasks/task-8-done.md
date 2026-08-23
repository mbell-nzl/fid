# Phase 21 - Task 8: Issue Actor And Event Convergence - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Issue create is converted through the durable Event-v2/idempotency path, but Issue update and transition still use the retained V1 writer path and the accepted retained-test convergence is unfinished.
- Recovery attempts: 0
- Next action: Continue task 8 from this checkpoint: complete red/green conversion of update and transition, add the fixed-clock/size corpus, converge affected retained tests, then run every exact task gate.
- Human attention required: no

## Summary of what was implemented

- Converted the Issue create CLI transport to require typed standalone `--actor` attribution and `--idempotency-key`.
- Converted Issue create to emit Event v2 and commit/replay through `idempotency_results`.
- Added red/green CLI corpus coverage for accepted create transport and durable replay identity.

## Files created or modified

- `cmd/fid/issue_create.go`
- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/issue_update.go`
- `cmd/fid/issue_transition.go`
- `cmd/fid/task8_contract_test.go`
- `internal/issues/actor.go`
- `internal/issues/create.go`
- `internal/issues/idempotency_results.go`
- `internal/issues/issue.go`
- `internal/issues/update.go`
- `internal/issues/transition.go`
- `project-docs/docs/phases/phase-21/tasks/task-8-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-8-done.md` - recovery checkpoint and captured red/green evidence.

## Validation Contract Coverage

- VAL-025: BLOCKED - `TestTask8IssueCreateRequiresStandaloneActorAndIdempotencyKey` and `TestTask8IssueCreateReplaysTheOriginalMutation` prove only the create slice; update, transition, and `--by` rejection corpus rows remain unfinished.
- VAL-026: BLOCKED - fixed-clock same-clock ordering test has not yet been added.
- VAL-027: BLOCKED - the creation size inventory and completed Issue corpus have not yet been run.

## Tests run

- `go test -race ./cmd/fid -run '^TestTask8IssueCreateRequiresStandaloneActorAndIdempotencyKey$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.813s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask8IssueCreateReplaysTheOriginalMutation$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.876s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...` - BLOCKED - `not run: task implementation is incomplete` (no observed termination)
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

- `d78bf0edc69f17f9a081918a10cce139edf8de78`

## Commit subject

- `feat(issue): begin actor and idempotency convergence`
