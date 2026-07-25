# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 19 **PLANNED (2026-07-26), execution gated.** Phase 18 is
COMPLETE / WRAPPED (2026-07-25); v0.1 remains closed for early dogfood
(DEC-057); v0.2a slice 1 of 3 is delivered.

Phase 19 (`fid workflow define/assign/list` + escalation default workflow,
v0.2a slice 2 of 3) is planned in **heavy mode** with six tasks and a 26-row
validation contract. Artifacts: `../phases/phase-19/` — `decisions.md`
(DEC-075..DEC-082), `qa/validation-contract.md`, `tasks/task-1..6.md`; brief
refreshed in `phase_brief.md`.

**Execution needs Human approval of the mode, contract, and task list — no
other precondition.** DEC-083 re-defers CF-023 to the coming agentic-workflow
overhaul and removes the DEC-074 execution gate; two zero-cost mitigations
(evidence-honesty rule in every task spec; phase-scoped `ESC-19-NNN` /
`TRG-19-NNN` with a triage-log row on every `BLOCK`) are carried by the phase
instead. The deep retrospective (DEC-073) is decoupled and becomes an input to
the overhaul.

Phase 18 outcome (full record in `../phases/phase-18/progress.md`): five tasks
closed on a terminating, green seven-command heavy gate (coverage `76.8%`);
final review **PASS** (24 PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD); human
checklist **PASS**; six triages all correct in hindsight; ESC-001 resolved by
DEC-072. Wrap decisions **DEC-073** and **DEC-074**. Delivered: workflow
structure is data (migrations `0007`/`0008`), every consumer resolves through
`workflow.Load`, the frozen Go `definitions` map is test-only, and the DEC-066
gate-obligation seam is present with zero behavior.

## Carry-Forward

Dispositioned at Phase 19 planning (DEC-082):

- **CF-023** — **re-deferred to the workflow overhaul (DEC-083, supersedes
  DEC-074).** No product-quality risk in items (a)-(g); item (h) is carried as
  a rule by every Phase 19 task spec. Stays open in the register.
- **CF-021** — **scheduled: Phase 19 Task 6**, which already touches
  `internal/maintenance` and runs the full Go gate.
- **CF-022** — **partially taken: Phase 19 Task 5** ships a byte-exact golden
  for the new markdown renderer; the backfill for pre-existing renderers stays
  deferred.
- **CF-024** — re-deferred to v0.2b, with a Phase 19 Task 6 **measurement
  obligation**: re-record the 100k-event rebuild (`<30 s`) and 10k-issue read
  (`<500 ms`) against the controls. Do not weaken either target.
- **CF-008** — **re-deferred to v0.2d** (governance vocabulary), where it pairs
  with the `activity_log` ref-type extension. Human confirmation pending.
- **CF-017 (v0.2e), CF-018, CF-020** — re-deferred; triggers not fired.

## Escalations

**None open.** Phase 18's ESC-001 (final-review BLOCK — `workflow.Load`
connection deadlock plus the then-absent QA-30 spec sync) is resolved by DEC-072,
verified by the 2026-07-23 rerun, and archived in
`../phases/phase-18/tl-escalations.md`. All Phase 17 escalations remain archived
in `../phases/phase-17/tl-escalations.md`.

## Planning Checkpoint (2026-07-18)

The trusted-evidence planning checkpoint ran per
`../reviews/2026-07-18-fid-roadmap-update-tl-brief.md`. DEC-059..DEC-066
recorded in `../reviews/2026-07-18-roadmap-update-decisions.md`; specs 00-10 and
the roadmap re-sliced (evidence foundation v0.2b before governance/reports; lean
completion v0.2c; work envelope v0.3a). **Pending Human ratification**; Phase 18
executed within the ratified-by-practice v0.2a slice.

## Next Action

In order:

1. **Human approves the Phase 19 plan** — operating mode (heavy), validation
   contract (VAL-001..026), the six-task list, and the CF-008 re-defer.
2. **QA Lead drafts the Phase 19 QA plan** with fresh command and feasibility
   validation (heavy mode never carries feasibility forward), then Phase
   Runner starts Task 1.
3. **Deep retrospective (DEC-073, timing at the Human's discretion)** — now an
   input to the agentic-workflow overhaul rather than a driver of a CF-023
   pass. Two themes: workflow-tooling reliability (CF-023 recurrence across
   Phases 17-18) and evidence over-claiming (three Phase 18 reports asserted
   `PASS` without evidence, each caught only by independent QA rerun).
