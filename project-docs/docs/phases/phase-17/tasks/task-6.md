# Phase 17 - Task 6: Implement `fid maintenance vacuum`

## Objective

Add the explicit SQLite VACUUM operation as a locked, permission-gated
maintenance command. It may compact only the derived DB; it must preserve all
projected data and journal bytes and report the DB size before and after.

## Validation Contract Coverage

- VAL-015: authorized, locked vacuum with preservation, output, denial, and
  policy-show/enforcement parity.

## Binding Constraints

- Add exactly the documented `maintenance.vacuum` matrix row: `tl` and `human`
  allowed; all other default roles denied. `workflow_runner` is not allowed.
- Require `--by`; reject missing/unknown/denied roles before VACUUM.
- Acquire `.fid/write.lock`, refresh projection under the held-lock pattern,
  then VACUUM the same derived DB. Avoid nested lock acquisition.
- VACUUM is not a domain write: append no journal event and add no activity row.
- Match the documented stderr confirmation, including before/after byte sizes.

## Scope

### Files To Create Or Modify

- `internal/maintenance/vacuum.go` - locked maintenance orchestration.
- `internal/maintenance/vacuum_test.go` - state, lock, and journal preservation.
- `internal/db/vacuum.go` and tests - thin SQLite VACUUM adapter; no policy,
  journal, or CLI knowledge.
- `cmd/fid/maintenance.go` - add `vacuum --by`.
- `cmd/fid/maintenance_vacuum_test.go` - output and role errors.
- `internal/policy/matrix.go` and matrix/policy tests - exact new row and parity.
- `project-docs/specs/02_TRD.md` - record lock/no-journal semantics if absent.
- `project-docs/specs/03_TRD_SCHEMA.md` section 13 - sync the matrix row.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 19 - exact behavior.
- `project-docs/docs/playbook.md` - operator command and safety note.

### Files To Read For Context

- `cmd/fid/rebuild.go`
- `internal/journal/lock.go`, `internal/journal/projection.go`
- `internal/db/db.go`
- `internal/policy/matrix.go`
- `project-docs/specs/02_TRD.md` section 4.3

### Out Of Scope

- Schema migration, `maintenance migrate`, automatic scheduled vacuum.
- Journal compaction or deletion.
- A new database/runtime dependency.

## Acceptance Criteria

- [ ] TL and human can vacuum; every missing/unknown/denied role fails before
  mutation and policy show matches enforcement. (VAL-015)
- [ ] Lock contention proves the command uses `.fid/write.lock` and does not
  deadlock through nested freshness acquisition. (VAL-015)
- [ ] Before/after projected snapshots and journal bytes are identical; stderr
  reports both DB sizes. (VAL-015)
- [ ] No journal/activity event is appended. (VAL-015)

## First Failing Test

`TestMaintenanceVacuumCommandAuthorizesCompactsAndPreservesState` fails because
the maintenance parent has no `vacuum` subcommand or matrix action.

## Documentation Updates

- `project-docs/specs/02_TRD.md` - persistence/lock boundary.
- `project-docs/specs/03_TRD_SCHEMA.md` - matrix row.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - command contract.
- `project-docs/docs/playbook.md` - safe operator use.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate, integration suite, and build are
required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Depends on task-5 for the parent command only; domain/permission work can be
  developed without changing maintenance-check behavior.
- Stop if `modernc.org/sqlite` cannot VACUUM under the established lock/open
  pattern; do not add a dependency or bypass the lock.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-6-done.md`.
