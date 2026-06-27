# Phase 12 - Task 4: `fid link list` And End-To-End Flow

## Objective

Implement `fid link list <id>`: return outgoing and incoming links for an entity
with direction, relation, entity type/id, and resolved entity title, in all four
output formats. Add an end-to-end integration test covering the full Phase 12
flow (requirement sync -> link add -> list -> remove) and proving rebuild
parity.

## Requirements

- PRD refs: PRD-LINK-006
- TRD refs: TRD-ARCH-010; PRD-NFR-008; ADR-0001; `03_TRD_SCHEMA.md` section 10
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid link list`; `05_DESIGN_PLAN.md`
  cross-surface fields
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002, section 5.6
- Decisions: DEC-028 (held-lock projection freshness)
- Depends on: Tasks 1-3 (requirement sync, link add, link remove)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-12/qa/validation-contract.md`.

- VAL-010: `fid link list <id>` returns outgoing + incoming links with
  direction, relation, entity_type, entity_id, entity_title; json/toon/markdown/
  yaml carry the same data.
- VAL-011: end-to-end flow (sync -> issue<->issue link -> issue->requirement
  link -> list both directions -> remove) appends one event per write with no
  nested-lock failure, and `fid rebuild` reproduces identical state.

## Scope

### Files To Create Or Modify

- `internal/links/list.go` (or equivalent) and tests - read-only query returning
  outgoing + incoming links with `entity_title` resolved from issues /
  requirements, deterministic ordering.
- `cmd/fid/link.go` and `cmd/fid/link_test.go` - add the `list` subcommand with
  `--format`.
- `cmd/fid/e2e_v01d_link_test.go` (new) - integration test for the full flow and
  rebuild parity.

### Files To Read For Context

- `internal/claims/list.go` and `cmd/fid/claim.go` for the read-list +
  `EnsureProjectionFresh` + `fidformat.Write` pattern (reads use the unlocked
  freshness helper).
- `internal/issues` enriched read helpers for resolving issue titles.
- `internal/requirements/` (Task 1) for resolving requirement summaries/titles.
- `cmd/fid/e2e_v01c_*_test.go` for the existing end-to-end test style.

### Out Of Scope

- Filtering `link list` by relation or direction (not specified for v0.1).
- Label commands and the full v0.1d lifecycle closure smoke (Phase 13).
- Schema migrations.

## Acceptance Criteria

- [ ] `internal/links` owns the list query; `cmd/fid` contains only flag wiring
  and format dispatch (VAL-010).
- [ ] `fid link list <id>` returns both outgoing and incoming links for the
  entity, each with `direction`, `relation`, `entity_type`, `entity_id`, and a
  resolved `entity_title` (VAL-010).
- [ ] json, toon, markdown, and yaml carry the same links and fields; structures
  parse equivalently to the JSON reference (VAL-010).
- [ ] List ordering is deterministic and stable across `fid rebuild` (VAL-010,
  VAL-011).
- [ ] An end-to-end test runs requirement sync -> issue<->issue link -> issue->
  requirement link -> `link list` (both directions) -> `link remove`, asserts one
  journal event per write with no nested-lock error, and asserts `fid rebuild`
  reproduces identical `links`/`requirements` state (VAL-011).

## TDD Plan

- First failing test: `link list` for an issue with one outgoing and one
  incoming link returns both rows with correct direction and resolved titles.
- Additional tests: requirement endpoint title resolution; all four formats
  carry the same data; deterministic order; empty result for an entity with no
  links; the end-to-end flow + rebuild-parity integration test.
- Regression risks: reads must use the unlocked `EnsureProjectionFresh` helper
  (not the held-lock variant); title resolution must handle both issue and
  requirement endpoints; ordering must be rebuild-stable.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` `fid link list` - reconcile only if implemented output
  differs; otherwise record verified-no-change.
- `project-docs/docs/phases/phase-12/decisions.md` - record any deviation found.
- `project-docs/docs/phases/phase-12/tasks/task-4-done.md` - completion report
  with command results, coverage, and commit hash.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `link list` four-format parity (VAL-010).
- Audit or release-readiness impact: end-to-end + rebuild-parity proof for the
  Phase 12 surface (VAL-011).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
```

The full gate above (including the integration suite) is required before task
close because this task adds the end-to-end / rebuild-parity integration test. A
focused test command may be used during the red/green loop, but it does not
substitute for the full gate at close.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
go test -bench . -benchmem ./...
```

- `go build -o /tmp/fid ./cmd/fid` - required because this task adds a CLI
  subcommand; record PASS.
- `go test -bench . -benchmem ./...` - run only if benchmark/scale fixtures are
  touched; otherwise record N/A with reason.

## Notes For Developer

- Use red/green TDD. Do not add a schema migration unless a spec-backed blocker
  is recorded first.
- `link list` is a read command: use the standard unlocked projection-freshness
  path, not `WithLock`.
- Keep the end-to-end test self-contained in a temp workspace; do not read
  held-out scenario material.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
