# Phase Brief - Phase 21

> Status: **Active / STOPPED After Review 1 / Review-1 Repair Authorized (DEC-163)**
> Branch: `phase-21`
> Product implementation: **AUTHORIZED BY DEC-159 UNDER THE TASK PACK**

## Goal And Mode

Deliver one Heavy one-Issue operational walking skeleton: activated trust,
immutable Gate Registration, one Work Envelope/revision, one gate-evidence
Obligation, generic Actor/Claim/Run authority, fid-executed Receipt Evidence,
guarded terminal-success Issue transition, later Envelope completion, bounded
resume, export create/validate, rebuild parity, and owned V1 removal.

**Operating mode: Heavy.** Authority, public JSON/Event contracts, destructive
schema convergence, local process execution, credentials/redaction, Artifact
safety, completion, and rebuild cross protected boundaries. One numeric phase,
two internal groups, one independent checkpoint; no 21A/21B exists.

## Binding Decisions

- DEC-152 as refined by DEC-158: Phase 20 retrospective controls task size,
  fresh feasibility, Phase Runner two-argument post-commit lint before the next
  dispatch, three-task review cadence, and checkpoint timing.
- DEC-153, DEC-154, DEC-155: Gate D3.2 DP-1..DP-7 stay accepted; exact initial
  slice/schema contract; carry-forward disposition.
- DEC-156: gate evidence guards terminal-success transition before completion;
  both use one shared Obligation evaluator.
- DEC-157: complete wire contract, kebab-case public enums, role-free bootstrap,
  exhaustive V1 disposition, achievable checkpoint, and tasks 19-22 order.
- DEC-158 and DEC-159: QA refinements, two-argument Runner report lint, four
  Human subchecks, three consecutive `<30 s` reference-instance rebuild samples,
  and the Ready plan plus Human approval that authorize dispatch.
- DEC-160: role-free bootstrap remains exact; task 2 owns a nonpersistent,
  seven-name V1 policy vocabulary bridge for inventoried legacy callers only,
  and task 20 deletes it.
- DEC-161: plain task-prefixed schema-v2 identity values are canonical;
  TRG-21-003 governs exact task-2 RUN-RESUME attempt 3 without a verifier edit.
- DEC-162: retained V1 tests are not protected contracts; the task whose
  accepted change invalidates one converges it in that task without weakening.
  Task 2 gains repo-wide test scope; tasks 8, 9, and 20 inherit the rule.
- DEC-163: a QA `REPAIR` review carries an ordered comma-separated owner list in
  one `Owning task` field; every listed owner completes recovery, repair, and
  task completion before the fresh re-review. `CORRECTION` is evidence-only.
  The `workflow-sequence-lint.sh` correction and its drills are authorized.

## Planning Authorities

- Slice semantics/Event fixtures: `../phases/phase-21/contracts.md`; complete
  public interface: `../phases/phase-21/public-interface.md`
- Live V1 inventory: `../phases/phase-21/v1-convergence-inventory.md`
- Validation: `../phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- Checkpoint: `../phases/phase-21/qa/group-1-checkpoint-contract.md`
- QA: `../phases/phase-21/qa/qa-lead-handoff.md` and `qa/qa-plan.md` (Ready at
  `e8a636e`); audit `../phases/phase-21/planning-audit.md`
- Tasks: `../phases/phase-21/tasks/task-<N>-*.md`

## Group 1 - Authority Foundation

| Task | Title | Depends on | VAL ownership |
|---:|---|---|---|
| 1 | 0013 schema/constraints | approval only | 001 |
| 2 | Event v2, idempotency, role-free bootstrap | 1 | 003-008 |
| 3 | capability discovery | 1-2 | 009-010 |
| 4 | Activation trust | 2-3 | 011-015 |
| 5 | Gate Registration | 4 | 016-018 |
| 6 | Envelope values | 5 | 019-022 |
| 7 | generic authority | 4,6 | 023-024 |
| 8 | Issue Actor/Event convergence | 2,7 | 025-027 |
| 9 | retained attribution/judgment freeze | 7-8 | 028-031 |
| 10 | generic Claims | 6-7 | 032-034 |
| 11 | generic Runs | 10 | 035,037 |
| 12 | shared evaluator/nonterminal transition | 7,10-11 | 038-041 |

QA Reviews follow tasks 1-3, 4-6, 7-9, and 10-12; the independent checkpoint
owns VAL-042/043. Temporary replacement-bound V1 state may remain exactly as
inventoried; PASS proves new generic paths do not consult it. Task 13 needs
checkpoint PASS and fresh Group 2 feasibility.

## Group 2 - Trusted Execution And Convergence

| Task | Title | Depends on | VAL ownership |
|---:|---|---|---|
| 13 | portable executor policy | checkpoint | 044-046 |
| 14 | live framing/redaction/controls | 13 | 036,047-050 |
| 15 | Attempt/Receipt | 13-14 | 051-054 |
| 16 | Artifacts | 15 | 055-057 |
| 17 | Evidence evaluation | 15-16 | 058-063 |
| 18 | terminal guard/completion/resume | 12,17 | 064-068 |
| 19 | export create/validate | 18 | 069-071 |
| 20 | V1 removal and migration 0014 | 1-19 | 002,074-075 |
| 21 | final-schema init/maintenance/rebuild/export parity | 20 | 076-083 |
| 22 | final interface/journey conformance | 21 | 072-073,084-086 |

QA Reviews follow tasks 13-15, 16-18, and 19-21. Final QA owns task 22
evidence and phase-end VAL-087. VAL-036 is the sole Human checklist assertion.

## Dispatch Controls

Before each next dispatch the Phase Runner runs
`project-docs/tools/report-lint.sh <committed-report> <task-spec>` and records
observed termination. Reviews are `qa/review-1.md`..`qa/review-7.md`; checkpoint
`qa/group-1-checkpoint-review.md`; final `qa/final-review.md` and
`qa/human-checklist.md`.

Every external schema-v2 `Originating task`, `Trigger signature`, `Attempt
count`, `Repair category`, and `Repair authority` value is a plain scalar, never
wrapped in Markdown backticks; each trigger begins with its exact originating
task plus `:`. A REVIEW `REPAIR` lists every Must-fix owner in the one `Owning
task` field, in Repair Plan order; each owner is repaired in turn before the
fresh re-review.

## Journey Order

Mutate once: init/Issue; discover; activate; register; create/activate Envelope;
resolve Actor; Claim/Run; nonterminal transition; execute Gate; persist/link
Receipt; evaluate Evidence; terminal-success transition; finish/release;
complete; resume; export/validate. Delete the projection and rebuild, then
repeat only named reads, evaluation, completion state, resume, export, and
negative checks; never create a second Envelope revision.

## Dispatch Authorization

All four planning gates are complete: Heavy shape and Gate D3.2 accepted, plan
Ready at `e8a636e`, no unresolved refinement, Human execution approval under
DEC-159.

Tasks 1-3 are mechanically complete. QA Review 1 returned `REPAIR`: MF-1 to
task-1 (VAL-001 exhaustive 0013 inventory), MF-2 to task-3 (VAL-010 named
deferred-surface probes). The ledger is valid, `STOPPED` at TRG-21-006, and the
Tech Lead must not edit it.

**Precondition.** DEC-163 items 4-5 land first: the authorized
`workflow-sequence-lint.sh` ordered-owner correction plus drills, verified by
`project-docs/tools/drills/run-drills.sh` PASS and a standalone
`VALID state=STOPPED` on the phase-21 ledger, committed on `phase-21`. This is
Human-executed workflow infrastructure, not a Runner or Developer task.

Then the Phase Runner appends RUN-RESUME with governing TRG `TRG-21-006`,
governing ESC `ESC-21-005`, resolution authority `DEC-163`, trigger signature
`task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding`,
attempt count 2, continuing task `task-1`, and this exact next recovery:

`task-1 | qa-rework | task-scope | task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding | 2`

The bound TASK-RECOVERY names `TRG-21-006` and `DEC-163` in both Source and
Repair authority, restarts no-progress at `0/2`, and dispatches only MF-1. After
task 1 completes, MF-2 runs as an ordinary task-3 `qa-rework` recovery under
`task-3:qa-rework:review-1` at attempt 1. Then fresh re-review `review-1b` of
tasks 1-3 with `Re-review of: review-1`. Task 4 is gated on review-1b.
