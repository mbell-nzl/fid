# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 19 **IN EXECUTION.** Tasks 1-3 complete (`a509d59`/`3b061e1`,
`10d4d24`, `50533e2`/`697c673`); QA Review 1 done; **Task 4 in progress,
resuming from preserved partial work.** Phase 18 is COMPLETE / WRAPPED
(2026-07-25); v0.1 remains closed for early dogfood (DEC-057); v0.2a slice 1 of
3 is delivered.

Phase 19 (`fid workflow define/assign/list` + escalation default workflow,
v0.2a slice 2 of 3) runs in **heavy mode** with six tasks and a **70-row**
validation contract (`Ready`). Artifacts: `../phases/phase-19/` —
`decisions.md` (DEC-075..DEC-083, DEC-089..DEC-100),
`qa/validation-contract.md`, `qa/qa-plan.md`, `tasks/task-1..6.md`; brief
refreshed in `phase_brief.md`. Branch: `phase-19-workflow-authoring`.

**Task 4 stopped and was cleared to resume (2026-08-01, TRG-19-004 /
ESC-19-004).** The first `capture.sh` run failed on a relative-path bug in the
script Task 4 was itself authoring. Stopping was correct — no commit, six gates
unrun, scratch worktree registered — but the recommended from-scratch restart
was not: the failed run wrote nothing, and DEC-098's baseline guard is
**structural** (detached worktree at the pinned commit + committed fixed
fixture), so working-tree edits cannot reach it. **DEC-099** preserves the
partial work and settles baseline provenance by reproduction against the
committed `SHA256SUMS`. **DEC-100** makes it standing: a task-authored evidence
harness failing while being authored is a **deviation, not a failed gate**.

**Review 1 (Tasks 1-3): BLOCK, resolved by DEC-098.** All 37 in-window
assertions ON TRACK with cited evidence; the sole blocking finding was a Task 4
precondition (ambient `HEAD` vs. the pinned Task 3 implementation commit).

**Earlier stops, all resolved:** Task 1 evidence-token defect (DEC-094/095,
TRG-19-001); Task 2 test-only import cycle (DEC-096/097, TRG-19-002); QA-plan
blockers QA-BLK-01..07 (DEC-089..093). Detail in `../phases/phase-19/`.

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

**None open.** **ESC-19-001** (Task 1 VAL-031 evidence token) was raised
2026-07-26 and resolved 2026-07-29 by DEC-094/DEC-095; it stays in
`tl_escalations_current.md` until the wrap archive. Phase 18's ESC-001 (final-review BLOCK — `workflow.Load`
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

1. **Correct the Task 1 report per DEC-094** (Developer session or human with a
   shell): re-execute the exact VAL-031 inventory, record the observed
   `exit 1 (no output)` on both lines, add the TRG-19-001 correction note,
   commit as its own `docs(task)` commit, rerun `report-lint.sh`.
2. **Phase Runner resumes**: record Task 1 complete, then start Task 2.
3. **Deep retrospective (DEC-073, timing at the Human's discretion)** — an
   input to the agentic-workflow overhaul rather than a driver of a CF-023
   pass; TRG-19-001 is a fresh input for it.
