# Phase 20 - Task 13: Maintenance inspection contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: Review 5 MF-01's bounded Task 13 evidence repair is implemented and has passed the unchanged exact gate.
- Recovery attempts: 1
- Next action: Fresh independent QA re-review of Review 5 MF-01.
- Human attention required: no

## Summary of what was implemented

- Migrated `fid maintenance check` to the `fid-cli.v1` Response Envelope with JSON success on stdout and typed `projection-divergence` failures on stderr.
- Preserved identical error, warning, and info findings when `--warn-as-error` promotes a warning result.
- Removed `maintenance check --format`, retained guarded legacy `maintenance vacuum`, and added the 1,000-findings-per-severity read limit without truncation.
- Resolved Review 5 MF-01 by replacing the mislabeled 298-second command fixture with exact 299-second non-qualifying and 300-second qualifying cases after an Attempt-config grace override, and by comparing the same fixture's `errors`, `warnings`, and `info` in warn-success and promoted-failure envelopes.

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
- VAL-008: PASS - `TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection` proves Human success is stdout-only at `cmd/fid/maintenance_test.go:171`; the failure branch is stderr-only at `cmd/fid/maintenance_test.go:45`.
- VAL-017: PASS - `TestMaintenanceCheckJSONMatchesContract` asserts `event:null` for the read response at `cmd/fid/maintenance_test.go:30`; the refreshed-projection integration test repeats that assertion at `cmd/fid/maintenance_integration_test.go:64`.
- VAL-027: PASS - `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat` overrides Attempt config at `cmd/fid/maintenance_test.go:118`, seeds exact 299/300-second command fixtures at `cmd/fid/maintenance_test.go:129` and `cmd/fid/maintenance_test.go:130`, requires 299 absent/300 present at `cmd/fid/maintenance_test.go:137`, and compares same-fixture `errors`, `warnings`, and `info` across warn-success/promoted-failure envelopes at `cmd/fid/maintenance_test.go:150`; exact closed error and warning check registries are asserted at `internal/maintenance/check_test.go:24` and `internal/maintenance/warnings_test.go:27`.
- VAL-044: PASS - `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat` rejects the removed `--format` request at `cmd/fid/maintenance_test.go:156`; `TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection` proves legacy vacuum rejects `--json` at `cmd/fid/maintenance_test.go:175`; an explicit selected root succeeds at `cmd/fid/maintenance_test.go:181`.
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
- `go test ./cmd/fid -run '^TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.058s` (terminated `exit_code: 0`)
- The bounded MF-01 repair is test/evidence-only: production already satisfied the exact boundary, so this corrected fixture was green without a production-code change.
- `go test ./cmd/fid -run '^(TestMaintenanceCheckJSONMatchesContract|TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection|TestMaintenanceCheckUsesExplicitProjectRoot|TestMaintenanceCheckRejectsFindingLimitWithoutTruncation)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.295s` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task13-mf01.9Qcsg7 --json init` - FAIL - `{"contract_version":"fid-cli.v1","ok":false,"operation":"project.init","error":{"code":"invalid-argument","category":"invalid-input","retryable":false,"message":"invalid project initialization argument","subject":null,"details":{"allowed":["valid project.init input"],"argument":"prefix","value":"FID-TASK13-MF01.9QCSG7"}}}` (terminated `exit_code: 1`). Cause: the temporary directory's derived prefix contained a dot; repair: retried the real-entrypoint check with a valid generated prefix.
- `/tmp/fid --project /tmp/fidtask13H9DsyG --json init` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-08-15.jsonl","project":{"id":"project-d67f2a9e36d26b9415df5522d536af9c","issue_prefix":"FIDTASK13H9DSYG","root":"/tmp/fidtask13H9DsyG"}},"event":{"id":"EVT-20260815-000001","occurred_at":"2026-08-15T20:15:26.680131112Z"}}` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fidtask13H9DsyG --json maintenance check` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"maintenance.check","result":{"errors":[],"info":{"claim_count_live":0,"issue_count":0,"recent_activity_7d":0},"status":"pass","warnings":[]},"event":null}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.807s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- No specification deviations. Review 5 MF-01 found that the prior report claimed exact 299/300 command evidence while the committed fixture used 298/300. This bounded repair corrects the fixture and report without changing the accepted check, grace, registry, request, result, or quality-gate contract.
- The original Task 13 red/green repairs remain above. The MF-01 fixture correction did not expose a production defect: the exact command test passed after the test-only evidence repair, so no production-code change was made.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `78f2ed7ebaa9078579e42f4e71493f4d0cad114a`

## Commit subject

- `test(maintenance): prove exact grace boundary`
