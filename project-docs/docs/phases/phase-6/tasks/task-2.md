# Phase 6 - Task 2: CF-010 — split oversized non-test files (behavior-preserving)

## Objective

Resolve carry-forward CF-010: three non-test files exceed the
`08_QUALITY_PLAN.md` STRUCT-002 400-line target —
`internal/issues/transition.go` (426), `internal/issues/update.go` (551),
`internal/journal/project.go` (524). Split each into logically cohesive files so
every file is ≤400 non-test lines, **without changing behavior**. Separate event
assembly from append/projection orchestration (per the CF-010 note) rather than
adding any new behavior.

## Requirements

- PRD refs: N/A (internal refactor)
- TRD refs: TRD §4.3 (write-path seam: validate/build vs append/project — the
  split should sharpen this boundary, not move it)
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` STRUCT-001 (≤60-line functions),
  STRUCT-002 (≤400-line files)

## Validation Contract Coverage

- VAL-010: the three named files are each ≤400 non-test lines after a
  behavior-preserving split, and behavior is unchanged (full suite green;
  `fid rebuild` parity holds).

## Scope

### Files To Create Or Modify

- `internal/issues/transition.go` — extract a cohesive unit (e.g. transition
  event assembly / error construction) into a new sibling file.
- `internal/issues/update.go` — extract a cohesive unit (e.g. field-change
  detection / update event assembly) into a new sibling file.
- `internal/journal/project.go` — extract a cohesive unit (e.g. per-effect
  projection handlers) into a new sibling file.
- New sibling files under the same packages as needed (names chosen for
  cohesion, e.g. `transition_event.go`, `update_fields.go`,
  `project_effects.go`).

### Files To Read For Context

- The three files above and their `_test.go` counterparts.
- `internal/journal/writer.go`, `internal/journal/event.go` (effect/event shapes
  referenced by the projection split).

### Out Of Scope

- Any behavior change, signature change to exported functions, schema change, or
  new validation/projection logic.
- Touching files not over the 400-line target.

## Acceptance Criteria

- [ ] `internal/issues/transition.go` ≤400 non-test lines (VAL-010).
- [ ] `internal/issues/update.go` ≤400 non-test lines (VAL-010).
- [ ] `internal/journal/project.go` ≤400 non-test lines (VAL-010).
- [ ] All extracted code stays in its original package; no exported-signature
  changes (VAL-010).
- [ ] Full suite green and `fid rebuild` parity unchanged after the split
  (VAL-010).

## TDD Plan

- This is a behavior-preserving refactor: the existing tests are the safety net.
  Run the full suite before and after; it must stay green with no test edits
  beyond moves.
- If any function is split to satisfy STRUCT-001, add a focused test only if an
  internal seam becomes newly observable; otherwise rely on existing coverage.
- Regression risks: accidental behavior change while moving code; init-order or
  unexported-symbol collisions across new files in the same package.

## Documentation Updates

- None required (no behavior or interface change). Note the split rationale in
  the completion report and cite CF-010 closure.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None — projection output and CLI behavior must be
  byte-identical pre/post split (assert via rebuild parity).
- Audit or release-readiness impact: Closes CF-010; restores STRUCT-002
  compliance for the three files.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Hold COV-001 (`internal/issues`, `internal/journal` ≥ 80%); coverage must not
regress from the split. Record post-split line counts (e.g. `wc -l`) in the
completion report as VAL-010 evidence.

Conditional commands:

- None.

## Notes For Developer

- Do CF-009 (Task 1) first if both are open; Task 1 touches `cmd/fid`, this task
  touches `internal/*`, so they do not conflict.
- Prefer extracting whole functions into a sibling file over reshaping logic.
  The goal is cohesion + size compliance, not redesign.
- Capture a `fid rebuild` parity check (projected state before vs after) as the
  behavior-unchanged evidence for VAL-010.
