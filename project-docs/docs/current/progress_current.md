# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 17 is **COMPLETE / WRAPPED (2026-07-18)**. No phase is active.
**v0.1 is closed for early dogfood (DEC-057).**

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

## Next Action

**Next planning checkpoint:** slice roadmap **v0.2a - Workflow Engine Migration**
into a new approved phase. Do not start v0.2 work before that phase is planned
and Human-approved. Recommend the CF-023 workflow-tooling deep retrospective
before v0.2a execution so the report-lint / escalation-register friction is
fixed rather than recurring.
