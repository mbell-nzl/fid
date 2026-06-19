# Phase 8 - Task 4: CF-012 — behavior-preserving split of `internal/issues/issue.go`

## Objective

Split `internal/issues/issue.go` (677 lines, over the 400-line non-test
tripwire, DEC-007) into cohesive same-package files, each ≤400 non-test lines,
with **no behavior change** and **no public API change**. This is the carry-
forward CF-012 cleanup, mirroring the Phase 6 CF-010 behavior-preserving split of
`transition.go` / `update.go` / `project.go`. Pure refactor: move functions,
keep tests green, delete/weaken nothing.

## Requirements

- PRD refs: N/A (internal refactor)
- TRD refs: N/A
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` file-size budget (≤400 non-test lines);
  DEC-007 size tripwire; CF-012

## Validation Contract Coverage

- VAL-009: `internal/issues/issue.go` is split into cohesive files each ≤400
  non-test lines with no behavior/API change and the full task gate green.

## Scope

### Files To Create Or Modify

- `internal/issues/issue.go` - reduce to a cohesive residual (e.g. shared
  types/errors + small helpers) ≤400 lines.
- New cohesive files in `internal/issues/` (e.g. `create.go` for the creation
  path: `Create`, `validateCreateOptions`, `validateParentConstraint`,
  `defaultWorkflowID`, `defaultIssueOwnerRole`, `issueCreateArgv`; and/or a
  `get.go` / helpers file) - exact cut is the Developer's call, but group by
  cohesion and keep each ≤400 lines.
- No test files need to move; if a test references an unexported symbol, the move
  stays in-package so it keeps compiling.

### Files To Read For Context

- `internal/issues/issue.go` (the file being split).
- `internal/issues/transition.go`, `update.go`, `project.go` (the Phase 6 CF-010
  precedent for how a cohesive split was cut in this package).
- `project-docs/docs/phases/phase-6/decisions.md` (CF-010 closure record).

### Out Of Scope

- Any behavior change, signature change, or new/removed exported symbol.
- Touching `complete.go` or other Task 1–3 files except trivial import hygiene.
- Splitting other already-compliant files.

## Acceptance Criteria

- [ ] `internal/issues/issue.go` and every new split file are ≤400 non-test lines
  (`wc -l`), and the package still builds (VAL-009).
- [ ] No exported symbol is added, removed, or renamed; no function signature
  changes; the diff is function moves + necessary imports only (VAL-009).
- [ ] The full task gate is green and no test was deleted or weakened
  (`go test -race ./...` passes; test count not reduced) (VAL-009).
- [ ] Coverage for `internal/issues` stays ≥ 80% (DEC-007 COV-001).

## TDD Plan

- This is a behavior-preserving refactor: no new failing test is written first.
  The existing `internal/issues` suite is the regression oracle — it must stay
  green before and after the move with no edits beyond mechanical import fixes.
- Regression risks: accidental visibility changes or duplicate declarations;
  run the full package suite after each move.

## Documentation Updates

- `project-docs/docs/phases/phase-8/` completion report - record the resulting
  file list with line counts and confirm no public API change.
- `project-docs/docs/current/decisions_active.md` Carry-Forward Register - CF-012
  is closed at phase wrap by the Tech Lead, not in this task; just record the
  evidence in the completion report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None (no behavior change).
- Audit or release-readiness impact: closes the CF-012 size-tripwire debt
  (spec-sync/decision input for wrap).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `internal/issues` ≥ 80%.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` to confirm the binary still builds after the
  split.

## Notes For Developer

- Do this task **after** Tasks 1–3 land so the new `complete.go` is already in
  place and the split does not churn against in-flight feature work.
- Scope file-size acceptance to the files you move; do not attempt to also split
  other files. The single goal is bringing `issue.go` (and its new siblings)
  under 400 lines without changing behavior.
- If the cleanest cut would change an exported symbol or behavior, stop and raise
  it (triage) — that is no longer behavior-preserving and needs a decision.
