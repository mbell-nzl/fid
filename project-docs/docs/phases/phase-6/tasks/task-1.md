# Phase 6 - Task 1: CF-009 — `fid issue update` one-field-per-line output

## Objective

Resolve carry-forward CF-009 (DEC-016, human HUMAN-013a): the `fid issue update`
confirmation currently renders every changed field as one long comma-separated
line. Render one changed field per line for readability. Behavior of the update
itself (which fields are updatable, journal/activity audit) is unchanged.

## Requirements

- PRD refs: PRD-ISSUE-015, PRD-ISSUE-021 (update surface; unchanged)
- TRD refs: N/A (presentation only)
- Design refs: `06_DESIGN_PLAN_CLI.md` §`fid issue update` (output example)
- Quality refs: `08_QUALITY_PLAN.md` STRUCT-001/002 (size budgets)

## Validation Contract Coverage

- VAL-009: `fid issue update` confirmation renders one changed field per line
  (multi-field update → multiple lines; single-field update → exactly one field
  line).

## Scope

### Files To Create Or Modify

- `cmd/fid/issue.go` — change `issueUpdateSummary` (and the `updated %s:`
  presenter call site) to emit one changed field per line.
- `cmd/fid/issue_update_test.go` — assert per-line rendering for multi-field and
  single-field updates.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — update the `fid issue update`
  output example(s) to the one-field-per-line shape.

### Files To Read For Context

- `cmd/fid/issue.go` (`issueUpdateSummary`, `issueUpdateFieldLabel`,
  `issueUpdateValue`, the update `RunE`)
- `internal/issues/update.go` (`FieldChange` shape)

### Out Of Scope

- Changing which fields are updatable, validation, journal events, or the
  `activity_log` audit.
- The v0.2 `--epic`/`--inbox` reparent example (leave the v0.2 note as-is, but
  keep its rendering consistent with the new per-line shape).

## Acceptance Criteria

- [ ] A multi-field `fid issue update` prints one line per changed field, in a
  stable order (VAL-009).
- [ ] A single-field `fid issue update` prints exactly one field line (VAL-009).
- [ ] The header (`updated <id>:`) remains, and each field line shows
  `<label> <from> -> <to>` (VAL-009).
- [ ] `06_DESIGN_PLAN_CLI.md` update output example matches the new rendering
  (VAL-009).
- [ ] No change to journal events or `activity_log` rows for an update
  (regression: existing update tests still pass).

## TDD Plan

- First failing test: update with two fields asserts two stderr field lines
  (currently one comma-joined line).
- Additional tests: single-field update asserts exactly one field line; header
  line present; field ordering deterministic.
- Regression risks: stderr-format assertions elsewhere (e.g. label add/remove);
  keep those presenters untouched.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid issue update` output example
  to one-field-per-line (and keep the v0.2 promote example consistent).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None (single stderr presenter; JSON/read formats
  unaffected).
- Audit or release-readiness impact: Closes CF-009 / DEC-016.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Hold COV-002
(`cmd/fid` ≥ 70%); do not regress package coverage.

Conditional commands:

- None.

## Notes For Developer

- This is a presentation-only change. Keep `issueUpdateSummary` cohesive; if you
  return a multi-line block, ensure the caller's `writeIssueLine` format string
  does not double-add a trailing newline.
- Preserve the existing changed-field ordering (the order `update.go` produces
  `FieldChange` entries) so output stays deterministic.
