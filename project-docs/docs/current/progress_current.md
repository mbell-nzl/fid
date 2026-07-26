# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 19 **APPROVED (DEC-093, 2026-07-26); QA plan re-verification
is the last precondition to Task 1.** Phase 18 is COMPLETE / WRAPPED
(2026-07-25); v0.1 remains closed for early dogfood (DEC-057); v0.2a slice 1
of 3 is delivered.

Phase 19 (`fid workflow define/assign/list` + escalation default workflow,
v0.2a slice 2 of 3) runs in **heavy mode** with six tasks and a **70-row**
validation contract (`Ready`). Artifacts: `../phases/phase-19/` —
`decisions.md` (DEC-075..DEC-083, DEC-089..DEC-093),
`qa/validation-contract.md`, `qa/qa-plan.md`, `tasks/task-1..6.md`; brief
refreshed in `phase_brief.md`. Branch: `phase-19-workflow-authoring`.

**QA plan blockers dispositioned (2026-07-26).** The QA plan opened BLOCKED on
QA-BLK-01..07. DEC-089 resolves the `internal/workflow` architecture-boundary
contradiction in the task pack's favor and syncs `02_TRD.md` §3; **DEC-090**
keeps PRD-WF-013 (a v0.1 Must) intact by requiring exactly one `cancelled`
state on issue-applicable workflows — the only new implementation obligation;
DEC-091 unbundles 14 assertions and adds VAL-027..070 with corrected sources
and a reproducible v0.1 parity baseline; DEC-092 adds the `workflow` family to
`05_DESIGN_PLAN.md` §3.1; DEC-093 records Human approval and confirms CF-008.
**Next: QA Lead re-verifies the QA plan against the refined contract and moves
it from BLOCKED to Ready.**

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
  with the `activity_log` ref-type extension. **Confirmed by the Human at
  approval (DEC-093).**
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

1. **QA Lead re-verifies the Phase 19 QA plan** against the refined contract
   (VAL-001..070): confirm QA-BLK-01..07 are closed by DEC-089..093, add one
   Section 4 evidence row and one Section 7 checklist step per new assertion,
   re-check the two previously BLOCKED feasibility items (VAL-003 tuple
   reference, VAL-059 baseline), and move the plan to **Ready**.
2. **Phase Runner starts Task 1** once the plan is Ready.
3. **Deep retrospective (DEC-073, timing at the Human's discretion)** — an
   input to the agentic-workflow overhaul rather than a driver of a CF-023
   pass.
