# Phase 1 - Task 6: Minimal `fid issue create` / `fid issue get` + JSON Format Engine

## Objective

Close the tracer bullet: a minimal `fid issue create` writes an `issue.created`
event through the full write path, projects it to `issues`, and
`fid issue get <id> --format json` reads it back. Add the `internal/format/`
output engine with JSON (the default) and the `--format` flag seam. Enforce
`--by <role>` with role-existence validation on writes.

## Requirements

- PRD refs: PRD-OPS-001, PRD-OUT-001, PRD-OUT-004, PRD-ROLE-001, PRD-ROLE-002
- TRD refs: TRD §4.3 (write path, one-write-one-event), TRD-ARCH-001,
  TRD-ARCH-004, TRD-ARCH-012 (format engine seam; JSON only here)
- Design refs: Design Plan CLI §2 (`fid issue create`, `fid issue get`);
  Design Plan §6.1 (issue JSON shape)
- Quality refs: Quality Plan §5.1

## Validation Contract Coverage

- VAL-006: A minimal issue is created via an event and read back as JSON.
- VAL-007: One logical write produces exactly one v1 event (full command path).
- VAL-008: Create issue, delete DB, rebuild → identical issue read-back.
- VAL-016: Writes require `--by <role>` and reject unknown roles.

## Scope

### Files To Create Or Modify

- `cmd/fid/issue.go` (thin `fid issue create` / `fid issue get` adapters)
- `internal/issues/issue.go` (validate + build `issue.created` event; read query)
- `internal/format/format.go` (format dispatch; JSON encoder; `--format` flag)
- `internal/policy/policy.go` (role-existence check) — minimal; full matrix deferred
- corresponding `*_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-6-completion.md`

### Files To Read For Context

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` §2 (issue create/get flags/output)
- `project-docs/specs/05_DESIGN_PLAN.md` §6.1 (issue JSON shape)
- `project-docs/specs/03_TRD_SCHEMA.md` §2.1 (`issues` columns + CHECKs)
- `project-docs/specs/01_PRD.md` PRD-OUT-001/004, PRD-ROLE-001/002
- Task 3/4/5 outputs (write path, projection, init)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- Full v0.1b issue lifecycle: kind hierarchy via `parent_id`, epic metadata,
  workflows/transitions, claims, AC, list filters, `--section`. Create a single
  **standalone** issue (no parent) with the minimal columns that satisfy the
  `issues` CHECK constraints (e.g. a `task` in `task-default`/`todo`).
- The full role×action permission matrix (PRD-ROLE-003) — deferred to v0.1d.
  Phase 1 checks only that `--by` is present and names a seeded role.
- TOON, markdown, and YAML formats (Phase 3). Implement JSON only; `--format`
  with an unimplemented value must error clearly, not silently fall back.

## Acceptance Criteria

- [ ] `fid issue create --kind task --title <t> --by <role>` writes exactly one
  `issue.created` event and projects one `issues` row. (VAL-006, VAL-007)
- [ ] `fid issue get <id>` returns the issue as JSON by default, and
  `--format json` is equivalent; fields match what was created. (VAL-006)
- [ ] After creating an issue, deleting `fid.sqlite`, and running `fid rebuild`,
  `fid issue get <id>` returns identical JSON. (VAL-008)
- [ ] A write with a missing `--by` or an unknown role exits non-zero with a
  clear error; a seeded role is accepted. (VAL-016)
- [ ] `--format` with an unimplemented value (toon/markdown/yaml) exits non-zero
  with a clear "not yet supported" error rather than silently emitting JSON. (no VAL; guards Phase 3 scope)

## TDD Plan

- First failing test: create an issue, then `fid issue get` returns matching
  JSON.
- Additional tests: exactly one `issue.created` line per create; unknown role
  and missing `--by` rejected; valid role accepted; delete-DB + rebuild yields
  identical read-back; unimplemented `--format` value errors.
- Regression risks: minting the issue ID at append time vs build time; emitting
  more than one event; format dispatch silently defaulting; CHECK-constraint
  violations from minimal column set.

## Documentation Updates

- `README.md` - extend the quickstart with `fid issue create` / `fid issue get`
  if user-facing docs need it.
- `project-docs/docs/phases/phase-1/tasks/task-6-completion.md` - completion
  report citing VAL-006/007/008/016 evidence and gate results.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: JSON output is the v0.1 default contract; keep the
  issue JSON shape aligned with Design Plan §6.1 so Phase 3 TOON/markdown/YAML
  can mirror it.
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
user-facing commands). Every required command must be recorded in the completion
report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if create/get/rebuild round-trip tests
  use the `integration` tag with temp workspaces.

## Notes For Developer

- Build the `issue.created` event with resolved values; mint the issue ID at the
  append step (or pre-mint and carry it in the event), but ensure replay applies
  stored values without re-deriving them (TRD-ARCH-002).
- Keep `cmd/fid/issue.go` thin: parse flags, call `internal/issues` and
  `internal/format`, set exit codes (TRD-ARCH-005, VAL-017).
- The `internal/format/` dispatch is the seam Phase 3 extends; design it so
  adding TOON/markdown/YAML later is additive, but implement JSON only now.
- Standalone task issue keeps you inside the minimal Phase 1 schema (no
  `parent_id` hierarchy, no triggers) while still proving the full loop.
