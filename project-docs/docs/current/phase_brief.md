# Phase Brief - No Active Phase (Phase 18 wrapped)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session.
>
> **Size target:** under 150 lines.

## Status

**No phase is active.** Phase 18 (v0.2a slice 1 of 3) is **COMPLETE / WRAPPED
(2026-07-25)**. Do not start Developer or QA work against this file; the next
work item needs Tech Lead planning and Human approval first.

Phase 18 outcome:

- All five tasks complete and committed, with a terminating, green seven-command
  heavy-mode gate (coverage total `76.8%`; every enforced package threshold met).
- QA Lead final review (rerun, 2026-07-23): **PASS**, no Blocking or Major
  finding, one Improvement finding (preserve CF-023) accepted as written.
- Validation contract: **VAL-001..024 all PASS** (24 PASS / 0 FAIL / 0 BLOCKED /
  0 CARRY-FORWARD) against fresh evidence.
- Human checklist (2026-07-25): **PASS** (HUMAN-01a..d all PASS); Phase 18
  accepted for continued dogfood.
- Wrap decisions: **DEC-073** (sign-off, VAL/TRG audit, CF-024, spec + glossary
  sync) and **DEC-074** (CF-023 promoted to a required pre-Phase-19 pass).
- Archive: `../phases/phase-18/` — `decisions.md`, `progress.md`,
  `tl-escalations.md`, `triage-log.md`, `run-ledger.md`, `tasks/`, `qa/`.

## What Phase 18 Delivered

The v0.1 hard-coded workflow state machines are now **data**: migration `0007`
creates and seeds `workflows`, `workflow_states`, `workflow_transitions`, and
`workflow_transition_requires`; migration `0008` rebuilds `issues` so
`workflow_id` is a real FK and `(workflow_id, status)` a composite FK. Every
consumer (transition engine, `issue get/list/ready`, `dev start/finish/block/
gate`, `maintenance check`) resolves structure from the **projected workflow
tables** via `workflow.Load`; the v0.1 Go `definitions` map survives **only** as
frozen test-parity fixtures, and a production read of it is a defect (VAL-024).
The DEC-066 gate-obligation seam (`required_gates_current`, `has_completion`,
`params_json` gate shape) is present in schema with **zero** behavior. Behavior,
journal bytes, and rendered CLI output are unchanged from v0.1.

Canonical language for the above lives in `../../specs/GLOSSARY.md`:
**Projected workflow tables**, **Workflow registry**, **Frozen workflow
definitions**.

## Next Planning Checkpoint

Ordered. The Tech Lead plans each; the Human approves before execution.

1. **CF-023 workflow-tooling pass — required before Phase 19 (DEC-074).**
   Doc/tooling only, light mode, not a numbered v0.2a phase. Scope: mint
   `ESC-NNN` + triage-log entry on every `BLOCK`; `ESC-NNN`/`TRG-NNN` global
   uniqueness; completion-report/report-lint commit-metadata friction;
   run-ledger chronological append with no duplicates; a wrap-time reopen audit
   for trigger-based carry-forwards; evidence-honesty hardening so a
   non-terminating command cannot be recorded `PASS`.
2. **Deep retrospective — RECOMMENDED before Phase 19 execution (DEC-073).**
   Two recurring themes: the CF-023 workflow-tooling defects (Phase 17 → Phase
   18, five of six triages) and **evidence over-claiming** (three reports
   asserted `PASS` without evidence; each caught only by independent QA rerun).
   Best paired with the CF-023 pass so findings land in one window.
3. **Phase 19 — `fid workflow define/assign/list` + escalation default
   workflow** (v0.2a slice 2 of 3). Planning must also disposition **CF-021**
   (scheduled: doc comments on `maintenance.VacuumResult`/`Vacuum`) and
   **CF-008** (resolution/cancellation reason, paired with PRD-WF-011/012).
4. **Phase 20 — `role_permissions` migration + `fid role create` / `fid policy
   set` + v0.2a closure.**

## Standing Constraints For The Next Phase

- **Phase 19 builds on data, not code.** Custom-workflow definitions become
  journal events; the seeded defaults stay migration DDL (system reference data)
  so v0.1 history is unchanged and FK ordering holds on rebuild.
- **The seam stays behavior-free.** `required_gates_current` / `has_completion`
  get no receipt or gate behavior before v0.2b.
- **Do not reintroduce production reads of the frozen Go definitions map.**
- **Rebuild performance is close to its ceiling (CF-024).** Anything that adds
  per-event projection work must re-measure the 100k-event rebuild against the
  `<30 s` TRD-PERF-003 control.
- `internal/dev` `Finish`/`Block` remain over the 60-line function target
  (CF-017, v0.2e) — do not fold that refactor into a migration phase.

## Carry-Forward Snapshot

Full register in `decisions_active.md`. Open at the end of Phase 18: CF-008,
CF-017, CF-018, CF-020, CF-021, CF-022, CF-023 (now required pre-Phase-19), and
**CF-024** (new: TRD-PERF-003 rebuild headroom, target v0.2b).

## Escalations

None open. Phase 18's ESC-001 is resolved (DEC-072) and archived in
`../phases/phase-18/tl-escalations.md`.
