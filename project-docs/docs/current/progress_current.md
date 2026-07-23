# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 17 **COMPLETE / WRAPPED (2026-07-18)**; **v0.1 closed for
early dogfood (DEC-057)**. **Phase 18 (v0.2a slice 1) is IN EXECUTION (Heavy
Mode) — final-review remediation.** All five tasks are implemented and committed;
the QA Lead final review returned **BLOCK** (24 VAL PASS / 0 FAIL, but a
registry-loader deadlock leaves four required gate commands non-terminating),
resolved by **DEC-072**. v0.2a is sliced into Phases 18–20 (see the roadmap
phase-slicing note and `phase_brief.md`).

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

**ESC-001 (Phase 18) — resolved by DEC-072**, pending remediation and QA
re-review before it archives. Raised from the final-review BLOCK: the
`workflow.Load` connection deadlock plus the then-absent QA-30 spec sync. The
spec-sync half is done; the loader fix is the open work. Note the ID collides
with Phase 17's ESC-001 — the known CF-023 defect; left as-is rather than
renumbered mid-phase.

All Phase 17 escalations (ESC-004/005/009/010/011) are resolved and archived in
`../phases/phase-17/tl-escalations.md`.

## Planning Checkpoint (2026-07-18)

The trusted-evidence planning checkpoint ran per
`../reviews/2026-07-18-fid-roadmap-update-tl-brief.md`. DEC-059..DEC-066
recorded in `../reviews/2026-07-18-roadmap-update-decisions.md`; specs 00-10
and the roadmap re-sliced (evidence foundation v0.2b before governance/
reports; lean completion v0.2c; work envelope v0.3a). **Pending Human
ratification. No implementation began.**

## Next Action

Phase 18 is in Heavy-Mode execution. All five tasks are implemented and
committed (Task 5 at `37b3fb5`; DEC-068 VAL-020, DEC-069 VAL-006, DEC-070 gate
completion, and DEC-071 TRD-PERF-003 all remediated and re-reviewed). The QA Lead
final review (`qa/final-review.md`) returned **BLOCK** with **24 VAL PASS / 0
FAIL** but two Blocking findings, triaged as **TRG-006 → ESC-001** and resolved
by **DEC-072**:

- `workflow.Load` issues nested per-workflow queries while the outer `workflows`
  rows still hold the connection, so it deadlocks on a one-connection pool and
  four required gate commands (coverage, unit, integration, benchmark) never
  terminate. The real CLI entrypoint is unaffected (no pool cap in `db.Open`), so
  VAL-010 passes on behavior — but the heavy-mode gate is not green.
- The QA-30 spec sync was absent. **Now applied** (DEC-072).

Remaining to close the phase:

1. Resume the Phase Runner from the ESC-001 / TRG-006 STOP; re-spawn a Developer
   to fix `workflow.Load` (materialize the workflow IDs and release the outer
   rows before any nested query) and add a bounded-timeout single-connection
   regression test, then run the full seven-command phase gate **to completion in
   order**, re-run the focused maintenance evidence and the VAL-024 two-`rg`
   inventory + runtime parity, correct the over-claimed command records in
   `task-5-done.md`, and commit. Do **not** touch the fixtures' `SetMaxOpenConns(1)`
   or the gate commands. A still-hanging command, changed registry ordering, a
   broken parity assertion, or a 100k rebuild at/over `30s` **escalates**.
2. QA Lead re-runs the **final review** against fresh deterministic evidence,
   then the human checklist (4 items, pending) — Task 5 is the last task; no
   further mid-phase review.
3. At wrap: audit all 24 VAL outcomes, confirm the DEC-072 spec-sync verdict
   table, archive ESC-001, and compress current context.
