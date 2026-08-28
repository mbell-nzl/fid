## Workflow outcome
- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented
- Replaced the legacy `fid run record` adapter with generic `run start`, `run finish`, `run get`, and `run list` leaves.
- Added immutable generic Run persistence bound to the active Claim, Envelope revision, Actor, repository observations, and null-counter budget values.
- Made finish terminal-only and fail-closed on repeated attempts; no recover or respawn leaf is registered.
- QA rework attempt 1 adds the complete CLI-observable Run start binding matrix and the prohibited predecessor/continuation absence checks without changing Run behavior or contracts.

## Files created or modified
- `internal/workflow/run.go`
- `internal/workflow/run_read.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `cmd/fid/run.go`
- `cmd/fid/generic_run_test.go`
- `cmd/fid/run_test.go`
- `cmd/fid/contract_surface_test.go`

## Documentation updated
- `project-docs/docs/phases/phase-21/tasks/task-11-done.md`

## Validation Contract Coverage
- VAL-035: `TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets` in `cmd/fid/generic_run_test.go:9` verifies the active Claim, Issue, Envelope/revision, Actor, exact start repository, and immutable null-counter budget binding; it also reads the same Run through `run get` and `run list`. `TestRunStartBindingMatrixRejectsWithoutMutation` at `cmd/fid/generic_run_test.go:56` adds a qualifying CLI row, then independently rejects inactive Claim, mismatched Issue, Envelope, Envelope revision, Actor, repository, and injected budget input. Every rejected row asserts the complete typed `invalid-argument` envelope, unchanged `run list` state, and unchanged journal Event count.
- VAL-037: `TestRunFinishIsOneWayAndRecoveryRespawnAreAbsent` at `cmd/fid/generic_run_test.go:149` verifies explicit terminal finish, rejection without a new Event on repeat, and absence of `run recover`, `run predecessor`, `run continuation`, and `run respawn`; `TestGenericRunUsesNoHistoricalAgentRunsTable` preserves the generic Run Event's legacy `agent_runs` telemetry absence.

## Tests run
- RED - `go test -race ./cmd/fid -run 'TestRun(StartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|FinishIsOneWayAndRecoveryRespawnAreAbsent|GenericRunUsesNoHistoricalAgentRunsTable)$' -count=1` - `FAIL` (terminated `exit_code: 1`). The new generic Run leaves were absent and `--actor` was rejected.
- GREEN - `go test -race ./cmd/fid -run 'TestRun(StartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|FinishIsOneWayAndRecoveryRespawnAreAbsent|GenericRunUsesNoHistoricalAgentRunsTable)$' -count=1` - `ok  github.com/mbell-nzl/fid/cmd/fid  5.541s` (terminated `exit_code: 0`).
- `go test -race ./cmd/fid -run '^TestRunRecord' -count=1` - FAILED ATTEMPT - `FAIL` (terminated `exit_code: 1`). Retained V1 `run record` assertions conflicted with the accepted generic Run replacement.
- Runtime entrypoint: `/tmp/fid-task11 --project /tmp/fid-task11-Sc41KC --json run list` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"run.list","result":{"items":[]}}` (terminated `exit_code: 0`).
- QA rework RED - `go test -race ./cmd/fid/...` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`). The new Run binding-matrix test was deliberately introduced before its matrix assertions.
- QA rework GREEN attempt 1 - `go test -race ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`). The new test referenced an undefined test-helper result type.
- QA rework GREEN attempt 2 - `go test -race ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`). The inactive-Claim setup used an invalid Claim release reason and did not create the intended released Claim.
- QA rework GREEN - `go test -race ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	192.627s` (terminated `exit_code: 0`). This runs the qualifying/rejected Run start matrix, one-way finish, all prohibited continuation-absence checks, and retained generic Run tests.
- `go build -o /tmp/fid-task11-matrix ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).

## Exact quality-gate results
- `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  161.767s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 34.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.7%` (terminated `exit_code: 0`)

- QA rework attempt 1 exact gate:
- `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	(cached)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 33.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)							65.4%` (terminated `exit_code: 0`)

## Deviations from spec
- The first focused full-package attempt failed with `FAIL` (terminated `exit_code: 1`) because `TestLiveLeafCommandsMatchContractInventories` retained `run record`. Root cause: the historical Phase 20 inventory predated the accepted task-11 public interface. Repair: converged the inventory and the retained Run tests to preserve exact leaf-surface strength for `run start|finish|get|list` and legacy absence. Later focused result: `ok  github.com/mbell-nzl/fid/cmd/fid  5.609s` (terminated `exit_code: 0`).
- QA evidence red/green: the first `go test -race ./cmd/fid/...` attempt ended `FAIL` (terminated `exit_code: 1`) because the new Run binding-matrix test was introduced before its matrix assertions. The first complete test attempt then ended `FAIL` (terminated `exit_code: 1`) because `commandResult` was not the local helper type; the next byte-identical rerun ended `FAIL` (terminated `exit_code: 1`) because the released-Claim setup used the non-accepted reason `binding matrix`. Repair: used `commandExecutionResult` and the accepted release reason `completed`, without changing production code, contracts, or assertions. The later byte-identical rerun ended `ok  	github.com/mbell-nzl/fid/cmd/fid	192.627s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead
- None.

## Dependency-budget note
- No dependency changes.

## Commit hash
- Pending commit

## Commit subject
- Pending commit
