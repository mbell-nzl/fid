# Phase 22 - Task 6: Recovery CLI and early process journey - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: DEC-194 transferred execution of the missing accepted checkpoint and inspection domain work from phase-22/task-4 to this task. The canonical domain producers, read composition, CLI leaves, error mapping, and early recovery process journey are now implemented and verified.
- Recovery attempts: 1
- Next action: none.
- Human attention required: no
- Previous durable recovery history: this report previously recorded `Task outcome: RECOVERY_REQUIRED` after the built-binary journey reached the absent `run inspect` leaf. The implementation/test checkpoint was `c69456fc9bf42f811e4999f1968a48c4a82d3cc5`; the prior finalized recovery report checkpoint was `7f91325ee2c96d50205dca9564806080dc4e3059`.

## Summary of what was implemented

- Added the first red built-binary subprocess journey fixture. It initializes a real Git repository and fid project, activates a recovery-capable actor, creates one Issue and Envelope, acquires and starts a Run, records a terminated observation, and then reaches the absent `run.inspect` leaf.
- Implemented canonical `run.checkpoint`: it enforces the existing write-lock, projection-freshness, idempotency, active-Run, expected-version, activated-harness authority, `run.checkpoint` capability, and Envelope-grant rules; re-observes and reconciles under the lock; rejects a conflicting reconciliation digest; then appends the accepted `run.checkpointed` Event with one checkpoint observation Effect and one runtime-version Effect.
- Implemented read-only `run.inspect`: it composes the Run, reconciliation against the latest accepted checkpoint (or Run start), and latest liveness without appending an Event, taking a mutation write lock, or writing repository state.
- Registered `fid run checkpoint`, `fid run inspect`, and `fid run recover`. Each invokes the canonical workflow operation; recovery errors use the accepted closed JSON envelopes and reason codes.

## Files created or modified

- `internal/workflow/run_checkpoint.go`
- `internal/workflow/run_checkpoint_test.go`
- `internal/workflow/run_inspect.go`
- `internal/workflow/run_inspect_test.go`
- `cmd/fid/run.go`
- `cmd/fid/run_recovery.go`
- `cmd/fid/run_recovery_test.go`
- `cmd/fid/generic_run_test.go`
- `cmd/fid/contract_surface_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-6-done.md`

## Documentation updated

- Updated this durable recovery report with the DEC-194 execution-transfer provenance, preserved red receipts, final gate receipts, and completion outcome.

## Validation Contract Coverage

- VAL-022: PASS - `TestP22EarlyRecoveryJourney` in `cmd/fid/integration/p22_recovery_journey_test.go` runs the built `fid` binary through the recovery journey. The preserved-red-command rerun ends `ok	github.com/mbell-nzl/fid/cmd/fid/integration	1.055s` with `exit_code: 0`; `internal/workflow/run_inspect_test.go` proves inspection is repository-byte-inert.
- VAL-023: PASS - `TestRunCheckpoint` and `TestRunInspect` in `internal/workflow/run_checkpoint_test.go` and `internal/workflow/run_inspect_test.go` cover the accepted checkpoint/inspection domain behavior, while `TestLiveLeafCommandsMatchContractInventories` in `cmd/fid/contract_surface_test.go` retains the three registered recovery leaves.
- VAL-024: PASS - `TestRunRecoveryErrorsUseClosedRecoveryEnvelopes` in `cmd/fid/run_recovery_test.go` checks each recovery-domain error's typed JSON code; `TestReasonCodeDetailSchemaIsExactForEveryCode` in `internal/cli/codes_test.go` checks the added closed detail schemas.
- Preserved prior durable recovery assessment: VAL-022: BLOCKED - `TestP22EarlyRecoveryJourney` in `cmd/fid/integration/p22_recovery_journey_test.go:65` invokes the real built binary and reaches `fid run inspect`; it terminates with `exit_code: 1` because the current binary has no inspect leaf. The remaining checkpoint, expiry, recovery, and successor-inspection steps require a canonical missing predecessor API.
- Preserved prior durable recovery assessment: VAL-023: BLOCKED - no recovery CLI leaf can receive or reject its closed wire input until the missing checkpoint/inspection domain APIs are supplied.
- Preserved prior durable recovery assessment: VAL-024: BLOCKED - no task-6 process error mapping can be correctly implemented until the recovery leaves can invoke the canonical domain operations and their typed errors.

## Tests run

- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable fixture attempt: the predecessor activation contract rejected `phase: 22` with `invalid-argument`; the fixture was corrected to its accepted pre-task-6 `phase: 21` shape.
- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable fixture attempt: activation component digests were unsorted; the fixture now follows the accepted sorted order.
- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable fixture attempt: Envelope creation omitted the retained exact one-revision gate-evidence obligation and registration; the fixture now supplies both.
- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Intended red evidence: `fid --project /tmp/TestP22EarlyRecoveryJourney327508039/001 --json run inspect RUN-c620cb = err:exit status 2 stdout:"" stderr:""`. The focused test itself terminated `exit_code: 1` after the built binary rejected the absent leaf.
- `go test ./internal/workflow -run '^TestRunCheckpoint' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red TDD evidence: the test initially failed because `RunCheckpoint` and `RunCheckpointRequest` were undefined.
- `go test ./internal/workflow -run '^(TestRunCheckpoint|TestRunInspect)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	1.508s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run 'Test.*Run' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Preserved-gate evidence: `TestRunFinishIsOneWayAndRecoveryRespawnAreAbsent` retained `run recover` as an absent leaf.
- `go test ./cmd/fid -run '^TestRunRecoveryErrorsUseClosedRecoveryEnvelopes$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.076s` (terminated `exit_code: 0`)
- `go test ./cmd/fid ./internal/cli ./internal/workflow -run '^(TestRunFinishIsOneWayAndUnsupportedRespawnLeavesAreAbsent|TestReasonCodeDetailSchemaIsExactForEveryCode|TestRunCheckpoint|TestRunInspect)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	1.142s` (terminated `exit_code: 0`)
- Real entrypoint exercise: `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.055s` (terminated `exit_code: 0`). The test builds and invokes the `fid` process.

## Exact quality-gate results

- Preserved prior recovery checkpoint: `gofmt -l cmd internal` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- Preserved prior recovery checkpoint: `golangci-lint run ./...` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- Preserved prior recovery checkpoint: `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- Preserved prior recovery checkpoint: `go tool cover -func=coverage.out` - BLOCKED - not started because its required coverage profile was not produced; no observed termination result.
- Preserved prior recovery checkpoint: `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable first full-gate attempt: `TestLiveLeafCommandsMatchContractInventories` rejected the new recovery leaves because its retained inventory did not include `run checkpoint`, `run inspect`, or `run recover`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 46.0% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `62.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.763s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.055s` (terminated `exit_code: 0`)

## Deviations from spec

- DEC-194 authorizes the only task-scope repair: execution of the missing accepted phase-22/task-4 checkpoint and inspection domain work in task 6. The accepted owner remains phase-22/task-4; no task-4 artifact was edited.
- Required provenance for each transferred file: `internal/workflow/run_checkpoint.go` — accepted owner `phase-22/task-4`; no commit ever created it; DEC-194 transferred execution to task 6; exact command that runs it: `go test -race -coverprofile=coverage.out ./...`.
- Required provenance for each transferred file: `internal/workflow/run_checkpoint_test.go` — accepted owner `phase-22/task-4`; no commit ever created it; DEC-194 transferred execution to task 6; exact command that runs it: `go test -race -coverprofile=coverage.out ./...`.
- Required provenance for each transferred file: `internal/workflow/run_inspect.go` — accepted owner `phase-22/task-4`; no commit ever created it; DEC-194 transferred execution to task 6; exact command that runs it: `go test -race -coverprofile=coverage.out ./...`.
- Required provenance for each transferred file: `internal/workflow/run_inspect_test.go` — accepted owner `phase-22/task-4`; no commit ever created it; DEC-194 transferred execution to task 6; exact command that runs it: `go test -race -coverprofile=coverage.out ./...`.
- Mandatory failed-gate deviation: the first `go test -race -coverprofile=coverage.out ./...` attempt ended `FAIL` with `exit_code: 1`. Root cause: Phase 20/task 15's retained live-leaf inventory omitted the task-6 accepted recovery leaves. Repair: updated the in-scope contract-surface inventory to retain `run checkpoint`, `run inspect`, and `run recover`; the byte-identical later command passed with `exit_code: 0`.
- Mandatory preserved-red-command deviation: the prior `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` attempt ended `FAIL` with `exit_code: 1` after the absent inspect leaf. Root cause: task 4 had not created the canonical producers. Repair: DEC-194 transferred their execution, and the byte-identical rerun passed with `exit_code: 0`.

## Concerns, technical debt, or notes for the Tech Lead

- VAL-013's unchanged tree/checkpoint clause and VAL-022 inspect must be audited against task-6 evidence with task-4 provenance; absence from task 4's report is not a task-6 defect.
- The retained generic-leaf rejection test was introduced by `5ca732ff163ee1d31fc155d3d9ebb3dac1501838` (phase-21/task-11); the live-leaf inventory file was introduced by `ad693228a9e1df57cab1d8f37159b47ecdb08e29` (phase-20/task-15). Task 6's accepted new leaves displaced only their stale absence/inventory assertions; their retained subjects and assertion strength remain intact.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `69941fc1b6262fdd7efaaeecfe72760390dbd511`

## Commit subject

- `feat(recovery): add checkpoint and inspect CLI`
