# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 17 **COMPLETE / WRAPPED (2026-07-18)**; **v0.1 closed for
early dogfood (DEC-057)**. **Phase 18 (v0.2a slice 1) is PLANNED — pending Human
approval.** No implementation has begun. v0.2a is sliced into Phases 18–20 (see
the roadmap phase-slicing note and `phase_brief.md`).

Phase 17 outcome:

- All ten tasks (0-9) closed with verified single-commit implementations and the
  required Heavy Mode gates.
- QA Lead final review: **PASS**, no Blocking or Major finding.
- Human checklist: **PASS** (HUMAN-25a..e all PASS, overall PASS).
- VAL-001..025 all resolved. VAL-022 narrowed under DEC-056 (Branch A confirmed);
  VAL-023 unblocked with it. VAL-025 accepted by the Human.
- v0.1 exit criteria demonstrably met: every v0.1-targeted PRD row maps to passing
  evidence (VAL-024), static binary 12.27 MiB (<30 MB), delete-and-rebuild
  per-table parity holds.

## Carry-Forward

- CF-019 CLOSED — markdown formatter split delivered (VAL-001).
- CF-020 (v0.2) — future-`occurred_at` vs `recent_activity_7d` window (DEC-056).
- CF-021 — Go doc comments on `maintenance.VacuumResult`/`Vacuum`; next reopen.
- CF-022 — broader markdown characterization tests; next `internal/format` reopen.
- CF-023 (Medium) — workflow-tooling hardening (report-lint/completion-report
  friction + ESC-NNN uniqueness); flagged for deep retrospective (DEC-058).
- CF-008 (v0.2), CF-017 (v0.2c), CF-018 (next SQL package reopen) unchanged.

## Escalations

None open. All Phase 17 escalations (ESC-004/005/009/010/011) are resolved and
archived in `../phases/phase-17/tl-escalations.md`.

## Planning Checkpoint (2026-07-18)

The trusted-evidence planning checkpoint ran per
`../reviews/2026-07-18-fid-roadmap-update-tl-brief.md`. DEC-059..DEC-066
recorded in `../reviews/2026-07-18-roadmap-update-decisions.md`; specs 00-10
and the roadmap re-sliced (evidence foundation v0.2b before governance/
reports; lean completion v0.2c; work envelope v0.3a). **Pending Human
ratification. No implementation began.**

## Next Action

The Phase 18 QA plan was drafted and came back `BLOCKED BEFORE TASK 1`
(QA-BLK-18-01..03); Phase Runner triaged as TRG-001 → BLOCK and stopped. The
Tech Lead resolved all three blocks on 2026-07-19 (**DEC-067**): the validation
contract is refined to **Ready** (now VAL-001..024), Task 4's dependent-object
inventory is made exact, and Task 5 gains the 100k-rebuild / 10k-read
performance benchmarks. Remaining before task execution:

1. **Human approval** of the Heavy Mode plan (operating mode, contract, five-task
   list).
2. On Phase Runner resume, the QA Lead refreshes the QA plan's IDs/mappings
   (activating the now-active VAL-015..024) and clears its blocked status.

Recommended before execution: the CF-023 workflow-tooling pass (light mode) so
report-lint / escalation-register friction is fixed first. Do not start Phase 18
tasks before Human approval.
