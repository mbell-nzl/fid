# Phase 17 - Task 2: Add the Maintenance Grace Configuration Contract - Completion Report

## Summary of what was implemented

- Added `[maintenance].started_without_claim_grace_seconds` to the project-config contract.
- New workspaces serialize the 300-second default; legacy `[project]`-only configs retain compatibility by reading as 300 seconds without a rewrite.
- Explicit zero and positive values round-trip, while negative and malformed values return contextual configuration errors.

## Files created or modified

- `internal/config/config.go`
- `internal/config/config_test.go`
- `cmd/fid/init_test.go`
- `project-docs/docs/phases/phase-17/tasks/task-2-done.md`

## Documentation updated

- No product or specification documentation required; DEC-050 had already synchronized the canonical contract. This completion report records the task evidence.

## Validation Contract Coverage

- VAL-002: PASS — config tests prove legacy/missing defaults, zero and positive serialization/read-back, and negative/malformed rejection; the init test and built `fid init` smoke prove bootstrap writes the default section.

## Tests run

- `go test ./internal/config -run '^TestReadProjectConfigDefaultsMissingMaintenanceGraceToFiveMinutes$' -count=1` - expected RED: compile failure because `ProjectConfig` initially had no maintenance grace field.
- `go test ./internal/config -run '^TestReadProjectConfigDefaultsMissingMaintenanceGraceToFiveMinutes$' -count=1` - PASS after adding the legacy default contract.
- `go test ./internal/config -run '^(TestProjectConfigMaintenanceGraceRoundTrips|TestReadProjectConfigRejectsInvalidMaintenanceGrace)$' -count=1` - expected RED before serialization/parsing, then PASS after implementation.
- `go test ./cmd/fid -run '^TestInitCommandCreatesWorkspaceFiles$' -count=1` - PASS.
- `go test -race ./internal/config ./cmd/fid` - PASS.
- `go build -o /tmp/fid-phase17-task2 ./cmd/fid` - PASS.
- Built CLI smoke: in a disposable workspace, `/tmp/fid-phase17-task2 init` - PASS; emitted `.fid/config.toml` contains `[maintenance]` and `started_without_claim_grace_seconds = 300`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (total statement coverage: 76.2%)
- Conditional commands - N/A (the task specification lists none).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d9c12ecd1b59a26a8239f4c608aca62003f7acdb`

## Commit subject

- `feat(config): add maintenance grace contract`
