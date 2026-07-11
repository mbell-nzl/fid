# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 17 is **APPROVED / BLOCKED AT ENTRY**. Phase 16 remains wrapped
and no implementation task is active.

Phase 17 kickoff state:

- Heavy Mode approved by the Human on 2026-07-11 (DEC-049).
- DC-A/DC-B/DC-C approved and recorded as DEC-050..052; affected canonical
  specs synchronized.
- VAL-001..025 are Ready in
  `../phases/phase-17/qa/validation-contract.md`.
- Ten ordered one-session task specs are ready under
  `../phases/phase-17/tasks/`, including DEC-054's tests-only Task 0.
- QA plan feasibility found `internal/journal` at 78.9% and no durable DEC-053
  model-audit evidence. Task 0 addresses coverage; ESC-009 tracks model audit.

## Approved Task Order

0. Restore `internal/journal` COV-001 coverage with tests only.
1. CF-019 formatter split.
2. Maintenance grace config.
3. Structural error checks.
4. Warning/info checks.
5. `maintenance check` CLI/formats.
6. `maintenance vacuum`.
7. Governance stub guards.
8. Read-format conformance.
9. Full v0.1 closure journey.

## Carry-Forward

- CF-019 is scheduled for Task 1 because markdown rendering reopens.
- CF-008 remains v0.2.
- CF-017 remains v0.2c.
- CF-018 remains at the next SQL package modification; Phase 17 only exercises
  SQL conformance. A required SQL-package fix stops for CF-018 re-planning.

## Next Action

Human resolves ESC-009 by supplying the DEC-053 model-audit evidence or
approving the proven-routing fallback. QA Lead then revises `qa-plan.md` to
authorize Task 0 and the entry review. Task 1 starts only after Task 0 passes
and its entry review has no unresolved Blocking finding.
