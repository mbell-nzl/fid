# Phase 21 - Task 5: Activated Gate Registration - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: Phase Runner report-lint and task-6 dispatch control.
- Human attention required: no

## Summary of what was implemented

- Added `fid gate register` and `fid gate resolve` for finalized Activation-owned Gate Definitions.
- Registration resolves the finalized Actor capability, validates the entire portable policy, persists the activated immutable policy through an Event-v2 registration, and returns existing identity for an identical binding.
- Changed bindings retain the Registration ID and create the next immutable version. Runtime policy and argv input remains closed.

## Files created or modified

- `cmd/fid/gate.go`
- `cmd/fid/gate_test.go`
- `cmd/fid/root.go`
- `cmd/fid/contract_surface_test.go`
- `internal/evidence/registration.go`
- `internal/evidence/policy.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-5-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-5-done.md` - task completion evidence and gate results.

## Validation Contract Coverage

- VAL-016: PASS - `TestGateRegisterRequiresFinalizedActivatedDefinition` at `cmd/fid/gate_test.go:15` proves staged registration returns `authority-unavailable`, then finalized Activation resolves and registers the exact activated Definition.
- VAL-017: PASS - `TestGateRegisterRejectsUnsupportedPolicyWithoutPersistence` at `cmd/fid/gate_test.go:71` exercises argv, shell, sandbox, network, write, environment, limits, exit codes, credentials, redaction, and paths; each row asserts `executor-policy-unsupported` and the zero Registration count at `cmd/fid/gate_test.go:101` and `cmd/fid/gate_test.go:105`.
- VAL-018: PASS - `TestGateRegisterReusesBindingAndVersionsChangedBinding` at `cmd/fid/gate_test.go:112` asserts same ID/version and one row for an identical binding, then the same ID/version 2 and changed binding digest for a changed scope; it resolves immutable version 1 at `cmd/fid/gate_test.go:131`.

## Tests run

- `go test -race ./cmd/fid -run '^TestGateRegisterRequiresFinalizedActivatedDefinition$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestGateRegisterRequiresFinalizedActivatedDefinition$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.353s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestGateRegisterRejectsUnsupportedPolicyWithoutPersistence$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestGateRegisterRejectsUnsupportedPolicyWithoutPersistence$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	9.042s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^(TestGateRegisterReusesBindingAndVersionsChangedBinding|TestGateRegisterRejectsRuntimePolicyOverride)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.054s` (terminated `exit_code: 0`)
- `/tmp/fid-task5 --project /tmp/fid-task5-entrypoint.wv7kV4 --json gate resolve REG-f7a68a --version 1` - PASS - `"operation":"gate.resolve"` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	142.193s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 49.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)						72.7%` (terminated `exit_code: 0`)

## Deviations from spec

- The first required focused command failed with `FAIL` (terminated `exit_code: 1`). Root cause: the closed CLI error-code and live-leaf inventories still represented the task-4 contract and excluded the task-5 codes and `gate register`/`gate resolve` leaves. Repair: updated the in-scope exact-contract tests to include the newly specified public errors and commands without changing their assertion strength. The byte-identical later command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	142.193s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5`

## Commit subject

- `feat(gate): register activated definitions`
