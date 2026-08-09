# Phase 20 - Task 2: Project selection and config namespace - Completion Report

## Summary of what was implemented

- Added Git-bounded project-root resolution with direct explicit-root selection and `.fid` symlink rejection.
- Replaced the permissive config reader/writer with the exact version-1 project, storage, executor-limit, and maintenance model; bootstrap now writes generated opaque project IDs and `issue_prefix`.
- Registered persistent root `--project` and `--json` flags behind one fail-closed migrated-leaf registry; until a leaf is migrated, explicitly supplied flags reject before its handler.
- Renamed workflow assignment scope to `--project-wide` and detached legacy maintenance warning grace from target config by using the fixed 300-second constant.

## Files created or modified

- `internal/project/resolve.go`
- `internal/project/resolve_test.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/bootstrap/init.go`
- `internal/issues/create.go`
- `cmd/fid/root.go`
- `cmd/fid/root_flags_test.go`
- `cmd/fid/init.go`
- `cmd/fid/init_test.go`
- `cmd/fid/workflow.go`
- `cmd/fid/workflow_assign_test.go`
- `cmd/fid/maintenance.go`
- `cmd/fid/maintenance_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-2-done.md`

## Documentation updated

- Completion report only; the accepted contract already documents the delivered behavior.

## Validation Contract Coverage

- VAL-003: PASS (Task 2 portion) — `TestResolveStopsAtContainingGitRoot`, `TestResolveFindsNearestRootAndRequiresExplicitRootToContainFID`, and `TestResolveRejectsSymlinkedFIDDirectory` in `internal/project/resolve_test.go` prove the Git boundary, nearest discovery/explicit direct-root rule, and symlink rejection.
- VAL-005: PASS (Task 2 portion) — `TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers` in `cmd/fid/root_flags_test.go` proves both explicit root flags reject before legacy and not-yet-migrated handlers; the real entrypoint command below proves the stderr/exit behavior.
- VAL-039: PASS (Task 2 portion) — `TestProjectConfigRoundTripsStrictTargetSchema` in `internal/config/config_test.go` proves opaque `project.id` and `project.issue_prefix` round-trip without the retired key; `TestInitCommandCreatesWorkspaceFiles` in `cmd/fid/init_test.go` proves direct bootstrap emits the strict config.
- VAL-043: PASS (Task 2 config-model portion) — `TestReadProjectConfigRejectsUnknownDuplicateAliasAndInvalidInputs` and `TestReadProjectConfigRejectsArtifactRootSymlinkEscape` in `internal/config/config_test.go` cover unknown/duplicate/alias/version, invalid prefix/limits/path, and symlink escape rejection.
- VAL-027: PASS (Task 2 fixed-policy wiring portion) — `TestMaintenanceCheckUsesFixedLegacyGraceInsteadOfTargetConfig` in `cmd/fid/maintenance_test.go` sets the target Attempt grace to zero yet proves a five-second-old started Issue does not produce the legacy warning.
- VAL-026: PASS (Task 2 flag-collision portion) — `TestWorkflowAssignRejectsWithoutMutation` in `cmd/fid/workflow_assign_test.go` exercises `--project-wide`, including the mutually exclusive kind/project-scope rejection.
- VAL-044: PASS (Task 2 root-dispatch portion) — `TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers` proves flags do not reach handlers before their transport migration. The complete retained-leaf request registry remains assigned to Tasks 3-13 and closes at Task 18.

## Tests run

- `go test ./internal/project -run '^TestResolveStopsAtContainingGitRoot$'` — EXPECTED RED — `FAIL`
- `go test ./internal/config -run '^(TestProjectConfigRoundTripsStrictTargetSchema|TestReadProjectConfigRejectsUnknownDuplicateAliasAndInvalidInputs)$'` — EXPECTED RED — `FAIL`
- `go test ./cmd/fid -run '^TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers$'` — EXPECTED RED — `FAIL`
- `go test ./cmd/fid -run '^(TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers|TestWorkflowAssignRejectsWithoutMutation|TestMaintenanceCheckUsesFixedLegacyGraceInsteadOfTargetConfig|TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename)$'` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.074s`
- `go test ./internal/config ./internal/project` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/project\t(cached)`
- `go test ./...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.582s`
- `/tmp/fid --json journal inspect EVT-missing` — PASS (expected rejection) — `error: --json is not available for journal inspect until its fid-cli.v1 transport is migrated`; observed exit status: 1.
- QA-plan applicability: QA-03 is covered by the project/root-flag tests and real CLI rejection; QA-11 is covered by the workflow collision test; QA-12 is covered by the fixed-grace test; QA-18 is covered by the strict config parser/writer and symlink tests. QA-01 is covered by the exact gates below. Later all-leaf, atomic-init, and maintenance-envelope QA evidence remains with their named owning tasks.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)`
- `golangci-lint run ./...` — PASS — `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` — PASS — `exit 0 (no output)`
- `go tool cover -func=coverage.out` — PASS — `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t76.6%`
- `go test -tags integration -race ./...` — PASS — `exit 0 (no output)`
- `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
