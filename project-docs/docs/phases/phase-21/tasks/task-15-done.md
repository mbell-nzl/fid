# Phase 21 - Task 15: Execution Attempts And Receipts - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added `fid gate run` with closed input, actor capability, and idempotency handling.
- Persisted a durable Attempt before process spawn, then recorded one immutable Receipt after every observed process terminal outcome.
- Added exact failed, spawn-failed, timed-out, and signaled Receipt outcomes with gate-unsuccessful exit class 8; added the frame-only post-stream `fid-error` terminal path.
- Preserved successful idempotent replay as one replayed terminal frame without another spawn, Attempt, or Receipt.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/gate.go`
- `cmd/fid/gate_run_test.go`
- `cmd/fid/generic_run_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/evidence/execution.go`
- `internal/evidence/execution_test.go`
- `internal/evidence/live.go`
- `internal/evidence/live_test.go`
- `internal/journal/idempotency.go`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-15-done.md`

## Validation Contract Coverage

- VAL-051: PASS - `TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame` at `cmd/fid/gate_run_test.go:69` proves the durable `gate.execution-started` Event precedes `gate.receipt-recorded`; the command invokes `/bin/true` at `cmd/fid/gate_run_test.go:17`.
- VAL-052: PASS - `TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame` at `cmd/fid/gate_run_test.go:78` and `cmd/fid/gate_run_test.go:111` proves exactly one successful Attempt/Receipt pair is retained across idempotent replay, while `cmd/fid/gate_run_test.go:98` proves replay emits one terminal frame only.
- VAL-053: PASS - `TestGateRunRecordsExactNonSuccessTerminalOutcomes` at `cmd/fid/gate_run_test.go:122` covers failed, spawn-failed, timed-out, and signaled process results with class-8 terminal response assertions at `cmd/fid/gate_run_test.go:149`; `TestAppendFidErrorClosesStartedStreamWithoutReceipt` at `internal/evidence/live_test.go:139` proves the handled infrastructure-error terminal has no Receipt.
- VAL-054: PASS - `TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame` at `cmd/fid/gate_run_test.go:85` and `cmd/fid/gate_run_test.go:88` proves execution persistence creates no Evidence record; Attempt/Receipt persistence does not satisfy caller proof or Evidence linking.

## Tests run

- `go test -race ./cmd/fid -run '^TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - red test before `gate run` existed.
- `go test -race ./cmd/fid -run '^TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  3.699s` (terminated `exit_code: 0`) - green after implementation.
- `go test -race ./internal/evidence -run '^TestSameGateRepositoryMatchesEqualOptionalValues$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence  1.015s` (terminated `exit_code: 0`).
- `go build -o /tmp/fid-task15 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `/tmp/fid-task15 --json gate run` - observed `invalid-argument` response with `"operation":"gate.run"` for missing required input (terminated `exit_code: 2`).

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  219.353s` (terminated `exit_code: 0`).
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  (cached)  coverage: 32.9% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:                                  (statements)             65.1%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  5.627s` (terminated `exit_code: 0`).

## Deviations from spec

- Required command 1 first failed with `FAIL` (terminated `exit_code: 1`) because `TestLiveLeafCommandsMatchContractInventories` retained a closed command inventory that omitted the accepted `gate run` leaf. The in-scope repair added that leaf and updated the expected count from 50 to 51; the byte-identical later attempt passed with `ok  github.com/mbell-nzl/fid/cmd/fid  219.353s` (terminated `exit_code: 0`).
- No product or public-interface deviation.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
