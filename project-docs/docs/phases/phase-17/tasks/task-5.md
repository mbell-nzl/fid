# Phase 17 - Task 5: Wire `fid maintenance check` and All Formats

## Objective

Expose the maintenance audit as a read-only CLI command with the approved
status/exit contract and json/toon/markdown/yaml output. Refresh the projection
before the audit, preserve journal bytes, and add focused human-readable
markdown without reopening CF-019.

## Validation Contract Coverage

- VAL-012: healthy/error states and all four output formats.
- VAL-013: warnings-only default versus `--warn-as-error` promotion.
- VAL-014: projection freshness and read-only journal behavior.

## Binding Constraints

- DEC-051 is binding. Default format is JSON; all read formats are required.
- `--warn-as-error` changes only `status` and process exit. Warning rows remain
  warnings, byte/structurally identical across the two runs.
- Render the full result before returning the non-zero outcome. Do not replace
  structured output with a prose error or duplicate it on stderr.
- Read freshness follows existing read commands; this command has no `--by`,
  matrix row, write lock on a fresh projection, or journal event.
- Keep maintenance markdown in `markdown_maintenance.go`; all non-test files
  remain within the 400-line budget after task-1.

## Scope

### Files To Create Or Modify

- `cmd/fid/maintenance.go` - parent command, `check`, flags, output, outcome.
- `cmd/fid/maintenance_test.go` - healthy/error/warning/format CLI behavior.
- `cmd/fid/maintenance_integration_test.go` - stale/missing projection and
  journal-preservation smoke.
- `cmd/fid/root.go` - register the command family.
- `internal/format/markdown.go` - add maintenance dispatch only.
- `internal/format/markdown_maintenance.go` - focused result renderer.
- `internal/format/markdown_test.go` - maintenance markdown cases.
- `project-docs/docs/playbook.md` - concise operator check examples.

### Files To Read For Context

- `cmd/fid/activity.go`, `cmd/fid/sql.go`, `cmd/fid/stats.go`
- `internal/maintenance/` (tasks 3-4)
- `internal/journal/projection.go`
- `internal/format/format.go`
- `project-docs/specs/05_DESIGN_PLAN.md` section 6.5
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 19

### Out Of Scope

- `maintenance vacuum` or `maintenance migrate`.
- Any repair/write command.
- Generic markdown rendering for arbitrary structs.

## Acceptance Criteria

- [ ] Healthy workspace returns zero/`pass`; each seeded structural error is
  listed and returns non-zero/`fail`, in every format. (VAL-012)
- [ ] Warning-only output is zero/`warn` by default and non-zero/`fail` with the
  flag, with equal warning payloads. (VAL-013)
- [ ] Missing/stale projection is refreshed before reading; fresh audit appends
  no event and all cases preserve journal bytes. (VAL-014)
- [ ] Canonical design, CLI, and playbook examples match implementation exactly.
  (VAL-012/013)
- [ ] Every new/modified non-test formatter file remains <=400 lines. (VAL-012)

## First Failing Test

`TestMaintenanceCheckCommandRendersHealthyAndFailingResults` fails because the
root command has no `maintenance` family.

## Documentation Updates

- `project-docs/docs/playbook.md` - operator usage.
- DEC-051 synchronized `05_DESIGN_PLAN.md` and `06_DESIGN_PLAN_CLI.md` at
  kickoff; stop if implementation proves that contract wrong.

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

- Depends on tasks 1, 3, and 4.
- A non-zero audit outcome is expected command behavior, not an encoding error;
  keep the structured result available to CI callers.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-5-done.md`.
