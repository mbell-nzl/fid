# Phase Brief - Phase 21

> Status: **Active / Task 12 STOPPED / Recovery Authorized (DEC-167)**
> Branch: `phase-21`
> Product implementation: **AUTHORIZED BY DEC-159 UNDER THE TASK PACK**

## Goal And Mode

Deliver one Heavy one-Issue operational walking skeleton: activated trust,
immutable Gate Registration, one Work Envelope/revision, one gate-evidence
Obligation, generic Actor/Claim/Run authority, fid-executed Receipt Evidence,
guarded transition, completion, resume, export, parity, and owned V1 removal.

**Operating mode: Heavy.** Authority, public JSON/Event contracts, destructive
schema convergence, local process execution, credentials/redaction, Artifact
safety, completion, and rebuild cross protected boundaries. One numeric phase,
two internal groups, one independent checkpoint; no 21A/21B exists.

## Binding Decisions

- DEC-152 as refined by DEC-158: Phase 20 retrospective controls task size,
  feasibility, report lint, review cadence, and checkpoint timing.
- DEC-153, DEC-154, DEC-155: Gate D3.2 DP-1..DP-7 stay accepted; exact initial
  slice/schema contract; carry-forward disposition.
- DEC-156: gate evidence guards terminal-success transition before completion;
  both use one shared Obligation evaluator.
- DEC-157: complete wire contract, kebab-case public enums, role-free bootstrap,
  exhaustive V1 disposition, achievable checkpoint, and tasks 19-22 order.
- DEC-158/159: QA refinements, four Human subchecks, three `<30 s` rebuild
  samples, Ready QA, and Human execution approval.
- DEC-160/161: role-free bootstrap stays exact with the task-2 nonpersistent V1
  bridge task 20 deletes; schema-v2 identity values are plain, task-prefixed.
- DEC-162: retained V1 tests are not protected contracts; the invalidating task
  converges them without weakening. Tasks 2, 8, 9, and 20 inherit the rule.
- DEC-163/164: QA `REPAIR` carries one ordered owner list and each owner repairs
  before fresh re-review, `CORRECTION` is evidence-only; ordinary schema-v2
  recovery names its TRG in `Source` and preserves external `Repair authority`.
- DEC-165/166: convergence follows the invalidating change and historical task
  labels never establish current ownership; task 9 owns its incidental reads and
  the Phase 20 request fixture, with V1 behavior byte-identical.
- DEC-167: a required gate must exercise the scope its task owns. Task 12 owns
  its `internal/dev` collateral and the accumulated `cmd/fid/integration/`
  convergence, with per-row provenance. Reaching a later journey step is
  progress; a directive-suppressed attempt does not consume recovery budget.
- DEC-168: 120s bounds each FEAS probe, 900s the FEAS set, 300s the tagged
  integration command; no budget bounds the QA Lead child session's lifetime.
  The checkpoint reruns unchanged and unsplit as attempt 4 under a one-off
  `task-13 | session-continuation | task-scope` wrapper; the Runner waits for the
  QA Lead and dispatches task 13 only on a verified checkpoint PASS.

## Planning Authorities

- Slice semantics/Event fixtures: `../phases/phase-21/contracts.md`; complete
  public interface: `../phases/phase-21/public-interface.md`
- Live V1 inventory: `../phases/phase-21/v1-convergence-inventory.md`; checkpoint
  `../phases/phase-21/qa/group-1-checkpoint-contract.md`
- Validation: `../phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- QA: `qa/qa-lead-handoff.md`, `qa/qa-plan.md` (Ready at `e8a636e`), audit
  `planning-audit.md`; tasks `../phases/phase-21/tasks/task-<N>-*.md`

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
owns VAL-042/043. Temporary replacement-bound V1 state may remain as
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
`qa/group-1-checkpoint-review.md`; final `qa/final-review.md` plus
`qa/human-checklist.md`. A task's required commands must exercise every file
scope it owns (DEC-167).

Every external schema-v2 `Originating task`, `Trigger signature`, `Attempt
count`, `Repair category`, and `Repair authority` value is a plain scalar, never
backtick-wrapped; each trigger begins with its exact originating task plus `:`.
Ordinary recovery names its TRG in `Source` without changing the exact-bound
`Repair authority`. A REVIEW `REPAIR` lists every Must-fix owner in the one
`Owning task` field in Repair Plan order, each repaired before the re-review.

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
DEC-159. Tasks 1-12 and review windows 1-4 are complete, `review-4b` PASS. The
run is STOPPED after TRIAGE TRG-21-014 / ESC-21-010: three fresh QA Lead Group 1
checkpoint sessions were interrupted with no artifact. DEC-168 resolves it - the
cause was orchestration applying the 900s FEAS-set command budget as the whole
child-session lifetime, not an unavailable capability.

The Phase Runner appends `RUN-RESUME` naming that STOP, governing `TRG-21-014`,
governing `ESC-21-010`, trigger signature
`task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted`, attempt 4,
resolution authority `DEC-168`, continuing task `task-13`, and next recovery
`task-13 | session-continuation | task-scope | <that trigger> | 4`. It then
appends one `TASK-RECOVERY task-13` with `Source` and `Repair authority`
`TRG-21-014 / DEC-168`, category `session-continuation`, mode `task-scope`,
attempt 4, and no-progress `0/2`.

That recovery is an orchestration wrapper. Dispatch the QA Lead checkpoint first,
wait for that child to terminate on its own, and never interrupt it at 900s. Every
probe reruns fresh in newly allocated scratch; preserved scratch and any
diagnostic run are inadmissible. Dispatch the task-13 Developer only after the
committed `qa/group-1-checkpoint-review.md` records a verified PASS. DEC-168 items
9-10 fix the REPAIR/TASK COMPLETE closure and every non-PASS branch.
