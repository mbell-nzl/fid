# Phase 20 - Task 13: Maintenance inspection contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Migrated `fid maintenance check` to the `fid-cli.v1` Response Envelope with JSON success on stdout and typed `projection-divergence` failures on stderr.
- Preserved identical error, warning, and info findings when `--warn-as-error` promotes a warning result.
- Removed `maintenance check --format`, retained guarded legacy `maintenance vacuum`, and added the 1,000-findings-per-severity read limit without truncation.

## Files created or modified

- `cmd/fid/maintenance.go`
- `cmd/fid/maintenance_test.go`
- `cmd/fid/maintenance_integration_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/integration/e2e_v02a_workflow_authoring_test.go`
- `internal/maintenance/check.go`
- `internal/maintenance/warnings_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-13-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-13-done.md` - task completion evidence.

## Validation Contract Coverage

- VAL-002: PASS - `TestMaintenanceCheckJSONMatchesContract` verifies the exact success/failure envelope partition and the stderr-only `projection-divergence` error at `cmd/fid/maintenance_test.go:29` and `cmd/fid/maintenance_test.go:45`.
- VAL-008: PASS - `TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection` proves Human success is stdout-only at `cmd/fid/maintenance_test.go:155`; the failure branch is stderr-only at `cmd/fid/maintenance_test.go:45`.
- VAL-017: PASS - `TestMaintenanceCheckJSONMatchesContract` asserts `event:null` for the read response at `cmd/fid/maintenance_test.go:30`; the refreshed-projection integration test repeats that assertion at `cmd/fid/maintenance_integration_test.go:64`.
- VAL-027: PASS - `TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo` compares `errors`, `warnings`, and `info` across warn and promoted branches at `cmd/fid/maintenance_test.go:102`; `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat` proves config independence and the 299/300 boundary at `cmd/fid/maintenance_test.go:118` and `cmd/fid/maintenance_test.go:136`; exact closed error and warning check registries are asserted at `internal/maintenance/check_test.go:24` and `internal/maintenance/warnings_test.go:27`.
- VAL-044: PASS - `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat` rejects the removed `--format` request at `cmd/fid/maintenance_test.go:140`; `TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection` proves legacy vacuum rejects `--json` at `cmd/fid/maintenance_test.go:159`; an explicit selected root succeeds at `cmd/fid/maintenance_test.go:165`.
- VAL-045: PASS - error and warning fixtures assert the deterministic `check, entity_id, message` order at `internal/maintenance/check_test.go:58` and `internal/maintenance/warnings_test.go:40`; 1,001 warnings return a typed over-limit error without changing rows at `cmd/fid/maintenance_test.go:183` and `internal/maintenance/warnings_test.go:110`.

## Tests run

- `go test ./cmd/fid -run '^TestMaintenanceCheckJSONMatchesContract$' -count=1` - RED - `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.047s` (terminated `exit_code: 1`). Cause: `maintenance check` had not registered the migrated root JSON transport; repair: registered its capability and added the Response Envelope branch.
- `go test ./cmd/fid -run '^TestMaintenanceCheckJSONMatchesContract$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.056s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestMaintenanceCheck' -count=1` - RED - `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.177s` (terminated `exit_code: 1`). Cause: the refresh integration test still expected default formatter JSON; repair: changed that test to request and assert the canonical JSON envelope explicitly.
- `go test ./cmd/fid -run '^TestMaintenanceCheck' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.201s` (terminated `exit_code: 0`)
- `go test ./internal/maintenance -run '^TestCheckRejectsMoreThanOneThousandWarningsWithoutTruncating$' -count=1` - RED - `FAIL	github.com/mbell-nzl/fid/internal/maintenance	0.231s` (terminated `exit_code: 1`). Cause: the checker returned all 1,001 warnings; repair: added a typed per-severity limit before response construction.
- `go test ./cmd/fid ./internal/maintenance -run '^TestMaintenanceCheck|^TestCheckRejectsMoreThanOneThousandWarningsWithoutTruncating$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/maintenance	0.227s` (terminated `exit_code: 0`)
- `go test ./cmd/fid ./internal/maintenance` - PASS - `ok  	github.com/mbell-nzl/fid/internal/maintenance	0.685s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run 'Maintenance' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.659s` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/tmp.xtkPKxfIiI --json maintenance check` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"maintenance.check","result":{"errors":[],"info":{"claim_count_live":0,"issue_count":0,"recent_activity_7d":0},"status":"pass","warnings":[]},"event":null}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.390s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `161f89ca18256cdc2a0ed3c9da6df094cdb8f427`

## Commit subject

- `feat(maintenance): migrate inspection response contract`
