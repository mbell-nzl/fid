# Phase 17 - Task 6: Implement `fid maintenance vacuum` - Completion Report

## Summary of what was implemented

- Added `fid maintenance vacuum --by <role>` as a locked SQLite VACUUM operation.
- Authorized only `tl` and `human`; missing, unknown, and every other default role are rejected before projection refresh or database mutation.
- Refreshes a missing or stale projection under the held write lock, then compacts only `.fid/fid.sqlite` without adding journal or activity entries.
- Emits the documented stderr confirmation with the measured pre- and post-VACUUM database byte counts.

## Files created or modified

- `cmd/fid/maintenance.go`
- `cmd/fid/maintenance_vacuum_test.go`
- `internal/db/vacuum.go`
- `internal/db/vacuum_test.go`
- `internal/maintenance/vacuum.go`
- `internal/maintenance/vacuum_test.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `internal/policy/authorize_test.go`
- `project-docs/docs/phases/phase-17/tasks/task-6-done.md`

## Documentation updated

- `project-docs/specs/02_TRD.md` — records the authorized, held-lock, no-journal/no-activity VACUUM semantics.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — defines authorization, locking, output stream, byte-size confirmation, and journal-neutral behavior.
- `project-docs/docs/playbook.md` — adds safe operator guidance for compacting the derived database.
- `project-docs/specs/03_TRD_SCHEMA.md` section 13 already contained the exact `maintenance.vacuum` row (`tl` and `human` only), so no synchronization edit was required.

## Validation Contract Coverage

- VAL-015: PASS — focused adapter, maintenance, policy, and CLI tests prove TL/human authorization; missing/unknown/denied role rejection; policy-show parity; write-lock contention without nested-lock deadlock; missing-projection recovery; byte-identical journal; unchanged projected roles/activity; and stderr database-size confirmation.

## Tests run

- `go test ./cmd/fid -run 'TestMaintenanceVacuumCommand(AuthorizesCompactsAndPreservesState|RejectsUnauthorizedRolesBeforeMutation)$' -count=1` - expected RED: failed because `maintenance vacuum` did not yet accept `--by` or enforce the new role contract.
- `go test ./cmd/fid -run 'TestMaintenanceVacuumCommand(AuthorizesCompactsAndPreservesState|RejectsUnauthorizedRolesBeforeMutation)$' -count=1` - PASS after implementation.
- `go test ./internal/db ./internal/maintenance ./internal/policy ./cmd/fid -run '(TestVacuum|TestMaintenanceVacuum|TestMatrixMaintenanceVacuum|TestAuthorizeMaintenanceVacuum)' -count=1` - PASS.
- `/tmp/fid init && /tmp/fid maintenance vacuum --by tl` in a disposable workspace - PASS; observed `vacuum complete — .fid/fid.sqlite: 204800 B → 204800 B` on stderr after initialization.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (total statements: 76.4%)
- `go test -tags integration -race ./...` - PASS
- `go build -o /tmp/fid ./cmd/fid` - PASS
- Conditional commands - N/A (the task specification lists none).

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
