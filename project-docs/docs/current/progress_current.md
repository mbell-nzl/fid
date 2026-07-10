# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 17 is **PLANNING / AWAITING HUMAN APPROVAL**. Phase 16 remains
wrapped and no implementation task is active.

Tech Lead planning draft created 2026-07-10:

- Proposed Heavy Mode for public CLI/config, persistence maintenance, and v0.1
  closure scope.
- Drafted VAL-001..025 in
  `../phases/phase-17/qa/validation-contract.md`.
- Drafted nine ordered one-session task specs under
  `../phases/phase-17/tasks/`.
- Marked DC-A (maintenance details), DC-B (result/exit semantics), and DC-C
  (exact governance guard set) as Human approval blockers.

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

- CF-019 is proposed for Task 1 because markdown rendering reopens.
- CF-008 remains v0.2.
- CF-017 remains v0.2c.
- CF-018 remains at the next SQL package modification; Phase 17 only exercises
  SQL conformance. A required SQL-package fix stops for CF-018 re-planning.

## Next Action

Human approves or edits Heavy Mode and DC-A/DC-B/DC-C. After approval, Tech
Lead records the decisions and marks the validation contract Ready. QA Lead then
authors `phase-17/qa/qa-plan.md` with fresh feasibility validation before Phase
Runner starts Task 1.
