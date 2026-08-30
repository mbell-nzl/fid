# Phase Brief - Phase 21

> Status: **Active / task-17 STOPPED / Recovery Authorized (DEC-172)**
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
  feasibility, report lint, review cadence, and checkpoint timing. DEC-153/154/155
  keep Gate D3.2 DP-1..DP-7, the exact slice/schema contract, and carry-forward.
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
- DEC-163/164: QA `REPAIR` carries one ordered owner list, each owner repairs
  before fresh re-review; ordinary recovery names its TRG in `Source`.
- DEC-165/166: convergence follows the invalidating change and historical task
  labels never establish ownership; V1 behavior stays byte-identical.
- DEC-167: a required gate must exercise the scope its task owns; task 12 owns its
  `internal/dev` and `cmd/fid/integration/` convergence with per-row provenance. A
  directive-suppressed attempt does not consume recovery budget.
- DEC-168: 120s bounds each FEAS probe, 900s the FEAS set, 300s the tagged
  integration command; no budget bounds the QA Lead child session's lifetime. The
  Runner waits for the QA Lead and dispatches task 13 only on a verified PASS.
- DEC-169: the checkpoint returned REPAIR. Task 4 owns Activation same-key replay
  ordering across the whole family; task 9 owns the complete section-7
  category-to-exit mapping, not one code. Both gain `cmd/fid/integration/` scope
  and the tagged gate. Order: task 4, task 9, complete checkpoint rerun.
- DEC-170: an externally interrupted attempt never consumes the no-progress
  budget, and no budget bounds a Developer child session's lifetime; only the
  declared per-command budgets apply. Commit durable work before long commands
  and append each result as it terminates. MF-15-02's provenance is satisfied
  forward-only; no commit at or before `09a2a7b` is rewritten.
- DEC-171: parent polling waits never terminate children; typed session routing,
  ledger `INTERRUPTED` evidence, generated prompts, workflow lint, and drill 19
  make DEC-170's liveness accounting executable before attempt 6.
- DEC-172: the unavailable `attention.v1` discovery row is required negative
  capability metadata, not Attention state/effect. Task 17 replaces the
  zero-match search with a boundary test that requires the row as the sole
  production Attention match and rejects every additional surface.

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

All four planning gates are complete under DEC-159. Tasks 1-16 are complete;
task 15's DEC-170/171 repair and independent `review-5b` passed, and task 16
then passed. Task 17's implementation is preserved uncommitted after the
accepted zero-match Attention recipe collided with the required unavailable
capability row. The run is STOPPED after TRG-21-019 / ESC-21-013; DEC-172
resolves the protected contradiction without changing the product interface.

The Phase Runner appends `RUN-RESUME` retaining TRG-21-019 / ESC-21-013, trigger
`task-17:specification-contradiction:task-17-done.md:attention-absence-vs-unavailable-capability`,
attempt 3, resolution authority DEC-172, continuing task 17, and the exact `task-17 | gate-repair | task-scope | task-17:specification-contradiction:task-17-done.md:attention-absence-vs-unavailable-capability | 3` tuple.
The bound TASK-RECOVERY uses source/authority `TRG-21-019 / DEC-172` and
no-progress `0/2`.

The Developer adds `TestPhase21AttentionBoundary`, preserves the old failed
`rg` row as historical evidence, commits the existing task-17 implementation
and test with `Refs: phase-21/task-17` before long commands, and runs the full
revised six-command block. After the report commit, two-argument report lint,
REPAIR RESOLVED, and TASK task-17 COMPLETE, task 18 may dispatch. Window 6
review remains after tasks 16-18.
