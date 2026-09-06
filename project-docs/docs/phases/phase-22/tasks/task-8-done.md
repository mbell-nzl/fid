# Phase 22 - Task 8: Interrupted Attempts and derived Run closure - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Derived `run.close` and intentional `run.abandon` remain to be implemented and verified.
- Recovery attempts: 0
- Next action: Continue task 8 from this checkpoint; add red/green workflow and journal parity coverage for close/abandon, then run the complete exact gate.
- Human attention required: no

## Summary of what was implemented

- Persisted the gate idempotency claim with the pre-spawn Attempt, and replace that pending result only after a fid-observed Receipt commits. A lost terminal observation now preserves one Attempt without a Receipt and the same key cannot spawn again.
- Added the observed `cancelled` terminal outcome for a caller-cancelled process.
- Repaired the inherited R-14 subprocess fixtures: the signal child blocks after self-signalling, and the timeout child writes a readiness marker then blocks under a five-second real gate deadline.

## Files created or modified

- `internal/evidence/execution.go`
- `internal/evidence/live.go`
- `internal/evidence/live_test.go`
- `cmd/fid/gate_run_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`

## Validation Contract Coverage

- VAL-031: PARTIAL - `TestRunLiveRecordsCancellationAsCancelled` proves cancellation records `cancelled`, with no timeout; `TestGateRunRecordsExactNonSuccessTerminalOutcomes` verifies the retained failed/spawn-failed/timeout/signaled Receipt branches. The full cancel/timeout/signal/spawn/harness matrix remains unfinished.
- VAL-032: PARTIAL - `TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn` proves one durable Attempt, no Receipt, and no same-key second spawn after terminal-observation loss.
- VAL-033: BLOCKED BY INCOMPLETE IMPLEMENTATION - `run.close` is not yet implemented.
- VAL-034: BLOCKED BY INCOMPLETE IMPLEMENTATION - `run.abandon` is not yet implemented.

## Tests run

- `go test -race -run '^TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); red evidence: same-key retry created `2/0` Attempts/Receipts.
- `gofmt -w internal/evidence/execution.go cmd/fid/gate_run_test.go && go test -race -run '^TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.507s` (terminated `exit_code: 0`).
- `go test -race -run '^(TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn|TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame|TestGateRunRecordsExactNonSuccessTerminalOutcomes)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	12.909s` (terminated `exit_code: 0`).
- `gofmt -w internal/evidence/live_test.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`); red evidence: cancellation was incorrectly observed as `timed-out`.
- `gofmt -w internal/evidence/live.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`); the process signal overwrote the newly derived `cancelled` outcome.
- `gofmt -w internal/evidence/live.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	1.031s` (terminated `exit_code: 0`).
- `gofmt -w cmd/fid/gate_run_test.go && go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	17.140s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - not run because task implementation is incomplete.
- `golangci-lint run ./...` - BLOCKED - not run because task implementation is incomplete.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not run because task implementation is incomplete.
- `go tool cover -func=coverage.out` - BLOCKED - not run because task implementation is incomplete.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not run because task implementation is incomplete.
- `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid` - BLOCKED - not run because task implementation is incomplete.
- `go test -bench . -benchmem ./...` - BLOCKED - not run because task implementation is incomplete.
- `go test -bench . -benchmem ./...` - BLOCKED - not run because task implementation is incomplete.

## Deviations from spec

- No contract or protected-definition change. The task is checkpointed before its required derived-closure behavior and final gate.

## Concerns, technical debt, or notes for the Tech Lead

- Retained R-14 fixture provenance: `cmd/fid/gate_run_test.go` signal helper and timeout row were introduced by `b69499116e70d4c88423a30daa2010a4ff221f14` (phase-21/task-15); the timeout extension was last touched by `76ad781727b9e8545e069c4ac0e9f926012e6af9` (phase-21/task-21). DEC-186 and task 8’s cancellation/observed-terminal acceptance invalidate only their fixture mechanisms. The repaired subject remains exact signal and timeout terminal Receipt assertions, executed by `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' ./cmd/fid`.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `bead7efb134f80726022f1e38d7b6749fe550073`

## Commit subject

- `feat(evidence): fence interrupted gate attempts`
