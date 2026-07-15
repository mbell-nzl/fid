# Phase 17 - Task 7: Add Governance Stub-Table CLI Guards - Completion Report

## Summary of what was implemented

- Registered the v0.2 version-gated `decision`, `escalation`, and `triage` command families with exactly the approved ten leaf commands.
- Added a shared typed version-gate error that renders before argument, future-flag, actor, workspace, DB, or journal handling.
- Preserved `fid requirement sync` as the existing shipped v0.1 control surface.

## Files created or modified

- `cmd/fid/governance_stubs.go`
- `cmd/fid/governance_stubs_test.go`
- `cmd/fid/root.go`
- `project-docs/docs/phases/phase-17/tasks/task-7-done.md`

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-016: PASS. The table-driven root CLI test covers the exact 10-command decision/escalation/triage matrix, including missing arguments and representative unknown future flags; each command returns its exact one-line entity/version error, leaves journal bytes and all three governance-table counts unchanged, and a no-workspace probe confirms the guard wins before workspace handling. Root help exposes the guarded families, while the `requirement sync` control writes one journal event successfully.

## Tests run

- `go test ./cmd/fid -run '^TestGovernanceStubCommandsReturnVersionGateBeforeValidation$' -count=1` - expected RED: failed before implementation because `decision` was not registered and no version-gate error was rendered.
- `go test ./cmd/fid -run '^TestGovernanceStubCommandsReturnVersionGateBeforeValidation$' -count=1` - PASS after implementing the minimal guard tree.
- `go test ./cmd/fid -run '^(TestGovernanceStubCommandsReturnVersionGateBeforeValidation|TestGovernanceStubVersionGateWinsBeforeWorkspaceValidation|TestGovernanceStubFamiliesAppearInRootHelp|TestRequirementSyncRemainsShippedAlongsideGovernanceStubGuards)$' -count=1` - PASS.
- `go test ./cmd/fid` - PASS.
- `/tmp/fid triage record --future-flag value` - PASS: exited 1 and wrote `error: commands for triage ship in v0.2` from `/tmp` without a workspace.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (no files listed).
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS.
- `go tool cover -func=coverage.out` - PASS (total statements: 76.4%).
- `go test -tags integration -race ./...` - PASS.
- `go build -o /tmp/fid ./cmd/fid` - PASS.

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
