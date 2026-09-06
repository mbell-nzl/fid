# Phase 22 - Task 8: Interrupted Attempts and derived Run closure - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: VAL-031 and `contracts.md` section 7 require an observed `cancelled` Receipt, while canonical `migrations/0013_p21_operational_kernel.up.sql` permits only `succeeded`, `failed`, `timed-out`, `signaled`, and `spawn-failed` in `execution_receipts.terminal_status`. Task 8 forbids DDL, wire-shape, and contract changes, so no in-scope implementation can persist the required cancelled Receipt.
- Recovery attempts: 3
- Next action: Tech Lead must decide whether to amend the canonical Receipt DDL and task-8 scope for the accepted `cancelled` status, or amend the accepted VAL-031/contracts requirement. Then dispatch the bound same-task continuation.
- Human attention required: yes

## Summary of what was implemented

- Persisted the gate idempotency claim with the pre-spawn Attempt, and replace that pending result only after a fid-observed Receipt commits. A lost terminal observation now preserves one Attempt without a Receipt and the same key cannot spawn again.
- Added the observed `cancelled` terminal outcome for a caller-cancelled process.
- Repaired the inherited R-14 subprocess fixtures: the signal child blocks after self-signalling, and the timeout child writes a readiness marker then blocks under a five-second real gate deadline.
- Added `run.close` and `run.abandon` domain operations. Close accepts only a trusted terminated observation, rejects outstanding Attempts, derives `succeeded` only when the shared Obligation evaluator has no reasons, and otherwise records `failed`. Abandon records a distinct terminal Run, releases its active Claim with `abandoned` provenance, accepts ordinary terminal proof or exact Human force authority over a live observation, and creates no Evidence or successor Run.
- Admitted the two task-8 actions while retaining `run.finish`, removed exactly the `run.closed` and `run.abandoned` DEC-192 availability fence keys, and added live/replay parity coverage for both Events.

## Files created or modified

- `internal/evidence/execution.go`
- `internal/evidence/live.go`
- `internal/evidence/live_test.go`
- `internal/workflow/run_close.go`
- `internal/workflow/run_close_test.go`
- `internal/workflow/activation.go`
- `internal/journal/project_effects.go`
- `internal/journal/checkpoint_projection_test.go`
- `internal/journal/recovery_schema_test.go`
- `cmd/fid/gate_run_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`

## Validation Contract Coverage

- VAL-031: BLOCKED - `TestRunLiveRecordsCancellationAsCancelled` proves the fid-observed `cancelled` process outcome, but `migrations/0013_p21_operational_kernel.up.sql:250-268` rejects that required Receipt status. No task-8-authorized source change can satisfy the cancellation Receipt branch.
- VAL-032: PASS - `TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn` proves one durable Attempt, no Receipt, and no same-key second spawn after terminal-observation loss.
- VAL-033: PARTIAL - `TestCloseRunDerivesTerminalStateFromTrustedObservation` proves trusted terminal observation, shared-evaluator-derived success, versioning, and idempotent replay; the required full gate and integration evidence are blocked by VAL-031's protected contradiction.
- VAL-034: PARTIAL - `TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority` proves ordinary intentional abandonment, Claim release provenance, no Evidence, and unknown-liveness rejection; `TestRunClosureEventsProjectIdenticallyLiveAndReplay` proves the `run.abandoned` Event's live/replay parity.

## Tests run

- `go test -race -run '^TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); red evidence: same-key retry created `2/0` Attempts/Receipts.
- `gofmt -w internal/evidence/execution.go cmd/fid/gate_run_test.go && go test -race -run '^TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.507s` (terminated `exit_code: 0`)
- `go test -race -run '^(TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn|TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame|TestGateRunRecordsExactNonSuccessTerminalOutcomes)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	12.909s` (terminated `exit_code: 0`)
- `gofmt -w internal/evidence/live_test.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`); red evidence: cancellation was incorrectly observed as `timed-out`.
- `gofmt -w internal/evidence/live.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`); the process signal overwrote the newly derived `cancelled` outcome.
- `gofmt -w internal/evidence/live.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	1.031s` (terminated `exit_code: 0`)
- `gofmt -w cmd/fid/gate_run_test.go && go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	17.140s` (terminated `exit_code: 0`)

- `go test ./internal/workflow -run '^TestActivationAdmitsRunCloseAndAbandonActions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red evidence: the accepted action admission rejected `run.abandon` and `run.close`.
- `gofmt -w internal/workflow/activation.go internal/workflow/run_close_test.go && go test ./internal/workflow -run '^TestActivationAdmitsRunCloseAndAbandonActions$' -count=1` - PASS - `ok` (terminated `exit_code: 0`)
- `gofmt -w internal/workflow/run_close_test.go && go test ./internal/workflow -run '^(TestCloseRunDerivesTerminalStateFromTrustedObservation|TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red evidence: `RunCloseRequest`, `CloseRun`, `RunAbandonRequest`, and `AbandonRun` were undefined.
- `go test ./internal/journal -run '^TestRunClosureEventsProjectIdenticallyLiveAndReplay$' - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red evidence: both accepted closure Events were rejected by their DEC-192 availability fences.
- `gofmt -w internal/journal/project_effects.go internal/journal/checkpoint_projection_test.go internal/workflow/run_close.go internal/workflow/run_close_test.go && go test ./internal/journal ./internal/workflow -run '^(TestRunClosureEventsProjectIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestActivationAdmitsRunCloseAndAbandonActions|TestCloseRunDerivesTerminalStateFromTrustedObservation|TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority)$' - FAIL - `FAIL` (terminated `exit_code: 1`)
  - First green attempt preserved: the Claim provenance `run-abandoned` violated the accepted DDL's closed `release_reason` domain.
- `gofmt -w internal/journal/project_effects.go internal/journal/checkpoint_projection_test.go internal/workflow/run_close.go internal/workflow/run_close_test.go && go test ./internal/journal ./internal/workflow -run '^(TestRunClosureEventsProjectIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestActivationAdmitsRunCloseAndAbandonActions|TestCloseRunDerivesTerminalStateFromTrustedObservation|TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority)$' -count=1` - PASS - `ok` (terminated `exit_code: 0`)
- `go test ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' - FAIL - `FAIL` (terminated `exit_code: 1`)
  - DEC-192 made its retained `run.closed` rejection fixture projectable, reaching `run_recoveries.prior_run_id` validation.
- `go test ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/journal 0.104s` (terminated `exit_code: 0`)
- `go test ./internal/workflow ./internal/journal ./internal/evidence ./cmd/fid -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Retained `TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError` expected no idempotency record after task 8's durable pre-spawn Attempt.
- `gofmt -w cmd/fid/gate_run_test.go && go test ./cmd/fid -run '^TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 0.260s` (terminated `exit_code: 0`)
- `go test ./internal/workflow ./internal/journal ./internal/evidence ./cmd/fid -count=1` - PASS - `ok` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-8-done.md project-docs/docs/phases/phase-22/tasks/task-8.md` - FAIL - `FAIL report-lint found 2 failure(s)` (terminated `exit_code: 2`)
  - Root cause: `Recovery attempts` was not a typed integer, and one new `go test` receipt was not an output token. Repair: retain attempt count `3` alone and record `ok` as the exact short terminal token.

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `golangci-lint run ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go tool cover -func=coverage.out` - BLOCKED - not run because its required coverage input is blocked by the protected cancelled-Receipt contradiction.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -bench . -benchmem ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -bench . -benchmem ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.

## Deviations from spec

- No contract, DDL, wire shape, error code, Event semantic, or gate definition was changed. The protected contradiction is recorded rather than silently mapping cancellation to another accepted Receipt status.
- The first closure implementation attempt used Claim `release_reason: "run-abandoned"`; the accepted DDL rejected it. The causal repair uses the existing closed `"abandoned"` value and preserves explicit Claim-release provenance.
- Convergence provenance: `internal/journal/recovery_schema_test.go:94-107` was introduced by `72c5974f6698fa08a4fe1a24099b3a93b3882f1d` (phase-22/task-1) and was moved to its `run.closed` fence subject by task 7's accepted retained-row convergence. DEC-192 and task 8's accepted `run.closed` producer invalidate only that unavailable-Event identity. The subject-preserving repair keeps live/replay rejection coverage by using still-fenced `activation.resource-observed`; executing command: `go test ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1`.
- Convergence provenance: `cmd/fid/gate_run_test.go:180-209` was introduced by `da6970b4ac5562e20c62f6c72ed8568c10d36f08` (phase-21/task-15). Task 8 implementation commit `bead7efb134f80726022f1e38d7b6749fe550073` persistently fences the pre-spawn Attempt and its pending idempotency result, invalidating only the old `0/0` Receipt/idempotency count. The subject-preserving repair retains the frame-only `fid-error` and no-Receipt assertions, and adds the required `1/0/1` Attempt/Receipt/idempotency proof; executing command: `go test ./cmd/fid -run '^TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError$' -count=1`.

## Concerns, technical debt, or notes for the Tech Lead

- Retained R-14 fixture provenance: `cmd/fid/gate_run_test.go` signal helper and timeout row were introduced by `b69499116e70d4c88423a30daa2010a4ff221f14` (phase-21/task-15); the timeout extension was last touched by `76ad781727b9e8545e069c4ac0e9f926012e6af9` (phase-21/task-21). DEC-186 and task 8’s cancellation/observed-terminal acceptance invalidate only their fixture mechanisms. The repaired subject remains exact signal and timeout terminal Receipt assertions, executed by `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' ./cmd/fid`.
- Resolution required: a Tech Lead decision must reconcile the accepted `cancelled` terminal Receipt requirement with the closed canonical Receipt DDL and task-8's explicit DDL prohibition. This task cannot choose a migration/schema change or alter VAL-031/contracts semantics.
- Recovery accounting: attempt 1 made interrupted-attempt and R-14 material progress; attempt 2 was report-only; attempt 3 made closure, action-admission, Event-parity, and retained-fixture convergence progress. The bound no-progress counter remains 0/2.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `b965b24748f57a0f5e42b8dd3054dfe6fd5dbb74`

## Commit subject

- `feat(workflow): derive task-8 Run closure`
