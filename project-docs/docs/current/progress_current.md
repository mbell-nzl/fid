# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 17 is **APPROVED / AWAITING QA PLAN**. Phase 16 remains wrapped
and no implementation task is active.

Phase 17 kickoff state:

- Heavy Mode approved by the Human on 2026-07-11 (DEC-049).
- DC-A/DC-B/DC-C approved and recorded as DEC-050..052; affected canonical
  specs synchronized.
- VAL-001..025 are Ready in
  `../phases/phase-17/qa/validation-contract.md`.
- Nine ordered one-session task specs are ready under
  `../phases/phase-17/tasks/`.

## Proposed Task Order

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

QA Lead authors `phase-17/qa/qa-plan.md` with fresh command, fixture, artifact,
and human-check feasibility validation. Phase Runner may start Task 1 only after
that plan is Ready.
