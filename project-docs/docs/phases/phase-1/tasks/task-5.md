# Phase 1 - Task 5: `fid init` Bootstrap Command

## Objective

Implement `fid init`: create the `.fid/` workspace, write canonical
`config.toml`, `pricing.toml`, and `.gitignore`, apply the schema, then run the
full write path once to append and project the `project.initialized` event that
seeds the seven default roles and the initial `projection_meta`. Reject re-init.
This is the first command that exercises the complete write+projection loop.

## Requirements

- PRD refs: PRD-INIT-001, PRD-INIT-002, PRD-INIT-003, PRD-INIT-004,
  PRD-OPS-004
- TRD refs: TRD §4.3 (bootstrap contract, `project.initialized`), TRD CONFIG
  §1-§4, TRD-ARCH-001
- Design refs: Design Plan CLI §1 (`fid init` output and re-init error)
- Quality refs: Quality Plan §5.1

## Validation Contract Coverage

- VAL-001: `fid init` creates the `.fid/` workspace layout and config files.
- VAL-002: `--prefix` sets the prefix; default is uppercased directory name.
- VAL-003: Exactly one `project.initialized` event seeds the seven roles.
- VAL-004: Re-init is rejected with a clear error and no changes.
- VAL-005: Delete `fid.sqlite` + `fid rebuild` reconstructs roles and checkpoint.

## Scope

### Files To Create Or Modify

- `cmd/fid/init.go` (thin `fid init` adapter)
- `internal/config/config.go` (write/read `config.toml`, `pricing.toml`,
  `.gitignore`) — create package
- `internal/issues/` or `internal/bootstrap/` for the `project.initialized`
  event builder (default role profile)
- corresponding `*_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-5-completion.md`

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §4.3 (bootstrap contract + example event)
- `project-docs/specs/04_TRD_CONFIG.md` (full file: layout, config, pricing,
  init behavior)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` §1
- `project-docs/specs/01_PRD.md` PRD-INIT-001..004
- Task 1/3/4 outputs (db/migrate, writer/lock, projection engine)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- Issue create/get and the `--format` engine (Task 6).
- The `[export]` config section (v0.4) — write only `[project]` in `config.toml`.
- Permission-matrix enforcement; `fid init` is the bootstrap exception and takes
  no `--by` flag.

## Acceptance Criteria

- [ ] `fid init` creates `.fid/journal/`, `config.toml`, `pricing.toml`,
  `fid.sqlite`, and `.gitignore` (ignoring `fid.sqlite`); config/pricing parse
  with documented defaults. (VAL-001)
- [ ] `fid init --prefix X` stores prefix `X`; with no flag the prefix is the
  uppercased project directory name. (VAL-002)
- [ ] Init appends exactly one `project.initialized` event and projects it,
  seeding exactly the seven roles (developer, qa_reviewer, qa_lead, tl_lite, tl,
  workflow_runner, human) and the initial `projection_meta` row. (VAL-003)
- [ ] Running `fid init` where `.fid/` already exists exits non-zero with a clear
  `.fid/ already exists` error and leaves the workspace unchanged. (VAL-004)
- [ ] After init, deleting `fid.sqlite` and running `fid rebuild` reconstructs
  the same seven roles and a matching `projection_meta`, without consulting the
  old DB and without mutating the journal. (VAL-005)

## TDD Plan

- First failing test: `fid init` in a temp dir creates the expected file tree
  and `.gitignore` content.
- Additional tests: prefix flag vs default; one `project.initialized` event +
  seven roles projected; re-init rejected with no change; delete-DB + rebuild
  reproduces roles and checkpoint; config/pricing default contents parse.
- Regression risks: bootstrap event written before projection_meta exists
  (special-case the first event); partial workspace left on a failed init;
  default-prefix derivation.

## Documentation Updates

- `README.md` - document `fid init` if the user-facing quickstart needs it.
- `project-docs/docs/phases/phase-1/tasks/task-5-completion.md` - completion
  report citing VAL-001..005 evidence and gate results.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None (only stderr confirmation output here).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
```

The full task gate plus the build is required before task close (this task ships
a user-facing command). Every required command must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if init/rebuild round-trip tests use
  the `integration` tag with temp workspaces.

## Notes For Developer

- `project.initialized` is the only event allowed before valid projection
  metadata exists; project its effects, then write the first `projection_meta`
  row (TRD §4.3 bootstrap).
- `config.toml` and `pricing.toml` are canonical committed files and are NOT
  reconstructed by `fid rebuild`; only the SQLite projection is.
- Keep `cmd/fid/init.go` thin; the workspace/event logic lives in `internal/`
  (TRD-ARCH-005, VAL-017).
- Reuse the writer/lock (Task 3) and projection engine (Task 4); do not
  re-implement the write path.
