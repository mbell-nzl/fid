# Phase 21 - Task 15: Execution Attempts And Receipts - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 6
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added `fid gate run` with closed input, actor capability, and idempotency handling.
- Persisted a durable Attempt before process spawn, then recorded one immutable Receipt after every observed process terminal outcome.
- Repaired non-success idempotency so every observed process Receipt is stored as the accepted result and a same-key/digest request replays one terminal frame with the original Receipt binding, without another spawn, Attempt, or Receipt.
- Added negative coverage for Attempt append failure, immutable Receipt mutation, handled `fid-error`, and caller-reported Evidence proof rejection.

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
- `internal/journal/project_effects.go`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-15-done.md`

## Validation Contract Coverage

- VAL-051: PASS - `TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame` at `cmd/fid/gate_run_test.go:21` proves the durable `gate.execution-started` Event precedes `gate.receipt-recorded`; `TestGateRunAttemptAppendFailureDoesNotSpawn` at `cmd/fid/gate_run_test.go:141` injects a rejected Attempt append and proves zero Attempts plus no `/bin/touch` spawn side effect.
- VAL-052: PASS - `TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame` invokes `assertCompleteReceiptBinding` at `cmd/fid/gate_run_test.go:328` for identity, execution, terminal, repository, toolchain, output, and redaction fields; it invokes `assertReceiptMutationRejectedAndImmutable` at `cmd/fid/gate_run_test.go:344`, while `internal/journal/project_effects.go:30` rejects every non-insert Receipt effect.
- VAL-053: PASS - `TestGateRunRecordsExactNonSuccessTerminalOutcomes` at `cmd/fid/gate_run_test.go:210` repeats failed, spawn-failed, timed-out, and signaled same-key/digest requests, asserting one replayed terminal, original Receipt, exit class 8, and one Attempt/Receipt/idempotency result; `TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError` at `cmd/fid/gate_run_test.go:170` proves the CLI-level frame-only `fid-error` exit 7 with null Receipt and no idempotency result. `internal/evidence/execution.go:120` replays the stored Receipt outcome and `internal/evidence/execution.go:339` commits every observed Receipt idempotently.
- VAL-054: PASS - `TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame` at `cmd/fid/gate_run_test.go:21` proves a normal Attempt/Receipt leaves zero Evidence rows, then attempts a caller-authored `fid-executed` Evidence insert and asserts its rejection and continued zero count.

## Tests run

- `go test -race ./cmd/fid -run '^TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - original red test before `gate run` existed.
- `go test -race ./cmd/fid -run '^TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.699s` (terminated `exit_code: 0`)
- Original focused red/green loop is retained above as historical evidence.
- `go test -race ./internal/evidence -run '^TestSameGateRepositoryMatchesEqualOptionalValues$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	1.015s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-task15 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-task15 --json gate run` - observed `invalid-argument` response with `"operation":"gate.run"` for missing required input (terminated `exit_code: 2`).
- `go test -race ./cmd/fid -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	222.988s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-task15-recovery ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-task15-recovery --json gate run` - observed `invalid-argument` response with `"operation":"gate.run"` for missing required input (terminated `exit_code: 2`).

- Historical attempt-5 rows below are unverified because no session-captured command termination was returned; they are retained only as historical context and are not current gate evidence.
- Historical attempt 5: `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` - UNVERIFIED - `ok github.com/mbell-nzl/fid/cmd/fid 234.922s` (termination not session-captured).
- Historical attempt 5: `gofmt -l cmd internal` - UNVERIFIED - `exit 0 (no output)` (termination not session-captured).
- Historical attempt 5: `golangci-lint run ./...` - UNVERIFIED - `0 issues.` (termination not session-captured).
- Historical attempt 5: `go test -race -coverprofile=coverage.out ./...` - UNVERIFIED - `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 32.9% of statements` (termination not session-captured).
- Historical attempt 5: `go tool cover -func=coverage.out` - UNVERIFIED - `total: (statements) 65.2%` (termination not session-captured).
- Recovery attempt 6: `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid	276.663s` (terminated `exit_code: 0`)

- Recovery attempt 6: `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Recovery attempt 6: `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Recovery attempt 6: `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 32.9% of statements` (terminated `exit_code: 0`)
- Recovery attempt 6: `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.2%` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-task15-attempt6 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-task15-attempt6 --json gate run` - observed `invalid-argument` response with `"operation":"gate.run"` for missing required input (terminated `exit_code: 2`).

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	276.663s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)										65.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.627s` (terminated `exit_code: 0`)

## Deviations from spec

- Required command 1 first failed with `FAIL` (terminated `exit_code: 1`) because `TestLiveLeafCommandsMatchContractInventories` retained a closed command inventory that omitted the accepted `gate run` leaf. The in-scope repair added that leaf and updated the expected count from 50 to 51; its later byte-identical attempt passed with `ok  github.com/mbell-nzl/fid/cmd/fid  219.353s` (terminated `exit_code: 0`).
- The first `report-lint.sh` attempt ended `exit_code: 9` with `FAIL report-lint found 9 failure(s)` because PASS rows had prose or punctuation after the observed termination tuple. The later report-only correction passed.
- Review-5 MF-15-01 and MF-15-02 required bounded task-scope rework. The repair changes no accepted contract, public-interface semantics, Event schema, policy, validation contract, migration, decision, or QA artifact.
- The first attempt of `gofmt -l cmd internal` printed `internal/evidence/execution.go` (terminated `exit_code: 0`), which is a format-gate failure despite the tool's zero process status. Formatting only `cmd/fid/gate_run_test.go` and `internal/evidence/execution.go` repaired it; the later byte-identical command was silent and terminated `exit_code: 0`.
- Reviewed report commits `79e3a933cd48568db0866e09255e001497245cca` and `8491cdeed312908e01eb898a3d5795d6000bf746` omit the `Refs: phase-21/task-15` footer required by `AGENTS.md`. Under DEC-170's forward-only disposition they remain unchanged because both are ancestors of the control repair base cited throughout the ledger, triage log, and ESC-21-011/012; every task-15 commit from attempt 6 onward carries that footer.
- No product or public-interface deviation.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `da6970b4ac5562e20c62f6c72ed8568c10d36f08`

## Commit subject

- `fix(evidence): make receipts replayable and immutable`
