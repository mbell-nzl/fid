# Phase 16 - Task 1: `fid activity list` projected timeline read

## Objective

Ship the human-facing activity timeline read over the existing `activity_log`
projection (PRD-OUT-017). This is the read surface that the CF-011 backfill
(Tasks 2-3) proves against, so it lands first.

- New `internal/activity` read package (PI-4): `List(options)` over `activity_log`
  with `--on <type>/<id>`, `--since <date>`, and `--event <name>` filters,
  returning typed entries in chronological order.
- New `cmd/fid/activity.go`: `fid activity list` cobra command with the four
  flags and all four read formats; read-only, no `--by` (PI-1).

The command reads the `activity_log` rows that transition/update/issue-label/
link-remove/claim writes already project today; Tasks 2-3 widen which writers
populate it.

## Requirements

- PRD refs: PRD-OUT-017
- TRD refs: `03_TRD_SCHEMA.md` §9 (`activity_log` columns); `02_TRD.md` §§2-3
  (package boundaries — add `internal/activity`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §8 `fid activity list`
- Quality refs: `08_QUALITY_PLAN.md` (read freshness, size budgets)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-16/qa/validation-contract.md`.

- VAL-001: `fid activity list` returns the documented entry shape in
  chronological order, all four formats.
- VAL-002: `--on <type>/<id>` includes only the named entity's entries.
- VAL-003: `--since <date>` includes on/after-boundary entries, excludes earlier.
- VAL-004: `--event <name>` includes only matching-event entries.

## Scope

### Files To Create Or Modify

- `internal/activity/activity.go` (new package; PI-4) - `ListOptions`
  (`WorkspaceRoot`, `MigrationsDir`, `On`, `Since`, `Event`), an `Entry` struct
  matching the documented JSON shape (`id`, `journal_event_id`, `ref_type`,
  `ref_id`, `event`, `actor`, `occurred_at`, `details`), and `List(options)`:
  `journal.EnsureProjectionFresh` → `fiddb.Open` → parameterised `SELECT` from
  `activity_log` with optional `ref_type`/`ref_id`, `occurred_at >=`, and
  `event` filters, ordered `occurred_at, id`. Parse stored `details` JSON back
  into a structured value for rendering. Parse `--on` as `<type>/<id>`.
- `internal/activity/activity_test.go` (new) - table-driven tests: full list,
  `--on` include/exclude, `--since` boundary include/exclude, `--event`
  include/exclude, `details` deserialisation, empty result.
- `cmd/fid/activity.go` (new) - `newActivityCommand()` + `list` subcommand;
  `--on`/`--since`/`--event`/`--format` flags; `fidformat.Write`;
  `renderIssueError` on failure (mirror `cmd/fid/stats.go`).
- `cmd/fid/activity_test.go` (new) - CLI tests for each filter + format dispatch.
- `cmd/fid/root.go` - register `newActivityCommand()`.
- `project-docs/specs/02_TRD.md` - §§2-3: add the `internal/activity` read
  package (projected timeline read over `activity_log`). Bump "Last updated".
- `AGENTS.md` - add `internal/activity/` to the Core/domain logic list.

### Files To Read For Context

- `internal/stats/stats.go`, `cmd/fid/stats.go` (read package + cmd + format
  pattern, `EnsureProjectionFresh`, `renderIssueError`)
- `internal/issues/list.go` (read freshness + parameterised filter query pattern)
- `internal/journal/project_effects.go` (`activity_log` columns + how `details`
  is stored as JSON text)
- `internal/format/format.go` (format dispatch, default `json`)
- `06_DESIGN_PLAN_CLI.md` §8 (exact flag table + JSON output shape)

### Out Of Scope

- CF-011 backfill of unlogged writers (Tasks 2-3); this task only reads what is
  already projected.
- `fid changes` entity-grouped digest (v0.2c).
- Any write path or permissions-matrix row (this is a read; PI-1).

## Acceptance Criteria

- [ ] `fid activity list` returns entries with `id`, `journal_event_id`,
  `ref_type`, `ref_id`, `event`, `actor`, `occurred_at`, `details` in
  chronological order, rendering in json/toon/markdown/yaml. (VAL-001)
- [ ] `fid activity list --on issue/<id>` returns only that issue's entries and
  excludes other entities' entries. (VAL-002)
- [ ] `fid activity list --since <date>` includes entries on/after the boundary
  and excludes earlier ones. (VAL-003)
- [ ] `fid activity list --event <name>` includes only entries with that event
  name. (VAL-004)
- [ ] `02_TRD.md` records the `internal/activity` boundary; `AGENTS.md` lists it.
  (doc)

## TDD Plan

- First failing test: `activity.List` over a seeded `activity_log` returns the
  expected entries — fails red (package absent).
- Additional tests: `--on` include/exclude; `--since` boundary (entry exactly on
  the boundary is included, one before is excluded); `--event` include/exclude;
  `details` round-trips to a structured value; empty result returns `[]` not
  null; CLI format dispatch for all four formats.
- Regression risks: keep the SQL parameterised (no string interpolation of
  filter values); preserve chronological ordering deterministically (tie-break
  on `id`); do not regress the read-freshness contract (call
  `EnsureProjectionFresh` before opening); files ≤400 lines, functions ≤60.

## Documentation Updates

- `project-docs/specs/02_TRD.md` - add `internal/activity` read package (§§2-3).
- `AGENTS.md` - list `internal/activity/` under Core/domain logic.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: output shape must match `06_DESIGN_PLAN_CLI.md`
  §8 JSON example across all four formats.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/activity` and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required if any new test in this task
  is integration-tagged (the package set under test includes `cmd/fid`); record
  PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the new
  `activity` command tests.

## Notes For Developer

- PI-1: this is a read command — no `--by`, no matrix row. Mirror the
  `fid stats` read shape exactly (`EnsureProjectionFresh` → open → query →
  `fidformat.Write`).
- `details` is stored as JSON text (`internal/journal/project_effects.go`
  `serializeActivityDetails`); deserialise it back to a structured value so it
  renders as an object, not an escaped string, matching the §8 JSON example.
- `--on` is `<type>/<id>` (e.g. `issue/FID-a1b2c3`); split on the first `/` and
  filter both `ref_type` and `ref_id`.
- Decide `--since` accepts both date-only (`2026-05-27`) and RFC3339; note the
  parsing in the completion report. The boundary semantics (on/after included)
  are what VAL-003 proves.
