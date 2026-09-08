# Phase 22 - Task 17: Maintenance, capability qualification, and rebuild budget - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added read-only recovery integrity findings for recovery links, budget bindings, activation recoveries, abandoned attempts, stale Claims, and unavailable artifacts. Healthy controls, including a completed complete-forward recovery later superseded by a valid activation, remain finding-free and maintenance does not mutate journal or projection state.
- Qualified the final four completed capability rows exactly, retained the schema-14 v1 archive capability profile, and added a shared fake/fid matrix that compares discovery, export, and the public resume action capability view.
- Added real-binary schema-16 init/rebuild evidence with retired escalation authority absent after both operations.

## Files created or modified

- `internal/maintenance/recovery_integrity.go`
- `internal/maintenance/operational.go`
- `internal/maintenance/check_test.go`
- `internal/workflow/capabilities.go`
- `internal/export/export.go`
- `internal/export/archive_profiles.go`
- `internal/export/archive_v2_test.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/export_test.go`
- `cmd/fid/integration/p22_capability_matrix_test.go`
- `cmd/fid/integration/task22_recovery_schema_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-17-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-17-done.md`

## Validation Contract Coverage

- VAL-069: PASS - `TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation` at `internal/maintenance/check_test.go:150` proves the named broken-link finding and unchanged journal/projection. `TestMaintenanceRecoveryIntegrityMatrixDoesNotMutate` at `internal/maintenance/check_test.go:188` proves clean healthy controls, later supersession at `:204`, all six corrupt classes at `:254`-`:272`, and no mutation.
- VAL-070: PASS - `TestP22QualifiedCapabilityMatrixMatchesFakeDiscoveryAndExport` at `cmd/fid/integration/p22_capability_matrix_test.go:12` compares the shared fake's exact twelve rows to real `fid --json capabilities` at `:15`-`:18`, exported `manifest.json` at `:20`-`:24`, and resume's public `required_capabilities` at `:27`-`:36`. `TestValidateRetainsPinnedV1Profile` in `internal/export/archive_v2_test.go` preserves the schema-14 profile and rejects a tampered legacy capability snapshot without mutation.
- VAL-071: PASS - three separate fresh-process `BenchmarkRebuildWorkspace100kEvents` samples recorded 28,121 ms, 27,904 ms, and 28,157 ms/rebuild, each below the unchanged 30,000 ms CF-024 target. Reference identity was captured by commands 6-11: Go 1.26.2, linux/amd64, Linux 6.6.114.1-microsoft-standard-WSL2 x86_64, 6 logical CPUs, AMD Ryzen 7 5700G with Radeon Graphics, and MemTotal 8133236 kB.
- VAL-072: PASS - tagged real-binary `TestRecoverySchemaFreshBinaryInitAndRebuild` at `cmd/fid/integration/task22_recovery_schema_test.go:14` verifies active schema 16 after init/rebuild and retired `escalations` absence at `:50`-`:58`. `TestFinalSchemaResetsThroughZeroWithoutRetiredAuthority` at `internal/db/escalation_removal_test.go:66` separately verifies migrate-down-through-zero then migrate-up; this is database migration evidence, not a claimed public reset command.

## Tests run

- `go test ./internal/maintenance -run '^TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation$' -count=1` - FAIL - `Check().Errors = []maintenance.Finding(nil), want broken recovery link finding` (terminated `exit_code: 1`); initial intended red before recovery integrity implementation.
- `go test ./internal/maintenance -run '^TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/maintenance	0.139s` (terminated `exit_code: 0`)
- `go test ./internal/maintenance -run '^TestMaintenanceRecoveryIntegrityMatrixDoesNotMutate$' -count=1` - FAIL - `Check().Errors = []maintenance.Finding{...activation-recovery-integrity...}, want later supersession accepted` (terminated `exit_code: 1`); root audit revealed that a later accepted activation supersedes a completed complete-forward recovery.
- `go test ./internal/maintenance -run '^TestMaintenanceRecoveryIntegrityMatrixDoesNotMutate$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/maintenance	0.150s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22QualifiedCapabilityMatrixMatchesFakeDiscoveryAndExport$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.326s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestRecoverySchemaFreshBinaryInitAndRebuild$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.641s` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task17-runtime.aO7RcR --json init --prefix T17` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260908-000001","occurred_at":"2026-09-08T09:19:27.290227801Z"},"ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-09-08.jsonl","project":{"id":"project-aadc24a6569ecd11f319bd7af4b27619","issue_prefix":"T17","root":"/tmp/fid-task17-runtime.aO7RcR"}}}` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task17-runtime.aO7RcR --json capabilities` - PASS - `{"contract_version":"fid-cli.v1",..."operation":"capabilities.get",...}` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task17-runtime.aO7RcR --json maintenance check` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"maintenance.check","result":{"errors":[],"info":{"claim_count_live":0,"issue_count":0,"recent_activity_7d":0},"status":"pass","warnings":[]}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)						63.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	77.237s` (terminated `exit_code: 0`)
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`)
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`)
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`)
- `nproc` - PASS - `6` (terminated `exit_code: 0`)
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`)
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.788s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.642s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.886s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	8.477s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	8.858s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- No contract, threshold, command, reference-instance, or failure-semantics deviation. I reran the full exact gate from the beginning after the in-scope superseded-activation correction and schema-16 evidence addition; all earlier complete gate attempts remain superseded evidence, not substituted results.
- Retained-row convergence provenance: `cmd/fid/capabilities_test.go` originated in `6ff90152fdcdd373c14a434ec736e557144a1617` (phase-21/task-3), displaced by Task 17's final capability qualification, and executes under `go test -race -coverprofile=coverage.out ./...`; `internal/export/archive_profiles.go` and its profile assertion originated in `acb569bb9447a3259991520669167fdc819c1706` (phase-22/task-13), displaced by the same capability qualification, and execute under that command; the retained unavailable-surface assertion originated in `401f609d4d7e24c0ab5ca45d27a16424000ded9f` (phase-21/task-3), displaced by qualified final rows, and executes under that command; `cmd/fid/export_test.go` originated in `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee` (phase-22/task-15), displaced by the final export snapshot, and executes under that command.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
