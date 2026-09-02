# Phase Brief - Phase 21

> Status: **Active / task-20 STOPPED / Recovery Authorized (DEC-173)**
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
- DEC-168/169: 120s per FEAS probe, 900s the FEAS set, 300s the tagged
  integration command; no lifetime budget on QA child sessions. Checkpoint
  repairs: task 4 owns Activation replay, task 9 the section-7 exit mapping.
- DEC-170/171: an externally interrupted attempt never consumes the no-progress
  budget; commit durable work before long commands; executable in the Runner.
- DEC-172: the unavailable `attention.v1` row is required negative capability
  metadata; task 17 proves the Attention boundary with a named test, not a
  zero-match search.
- DEC-173: a task that owns an absence/inventory control owns every path it
  walks. Task 20's scope is amended (bootstrap, journal, labels, links,
  requirements, db, testdata/journal, integration) with bounded
  subject-preserving dispositions and gains the tagged integration gate. No
  command, expected status, VAL, 0014, frozen-table, or public contract changes.

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
owns VAL-042/043. Group 1 and its checkpoint are complete.

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
scope it owns, and a task's scope must cover every path its controls walk
(DEC-167, DEC-173).

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

Tasks 1-19 and reviews 1-6 are complete. Task 20 committed migration 0014 and
the first V1 source removals (`aef03cf`, report `f189a37`) and returned
`RECOVERY_REQUIRED`; the run is STOPPED after TRG-21-020 / ESC-21-014. DEC-173
resolves the scope-versus-reach contradiction without changing any command.

The Phase Runner appends `RUN-RESUME` retaining TRG-21-020 / ESC-21-014, trigger
`task-20:specification-contradiction:v1-convergence-scope-vs-absence-gates`,
attempt 2, resolution authority DEC-173, continuing task 20, and the exact
`task-20 | in-scope-implementation | task-scope | task-20:specification-contradiction:v1-convergence-scope-vs-absence-gates | 2`
tuple. The bound TASK-RECOVERY uses source/authority `TRG-21-020 / DEC-173`
and no-progress `0/2`.

The Developer continues from `aef03cf` under the amended
`tasks/task-20-v1-removal.md`, commits durable work before long commands, runs
the complete byte-identical amended block (now including the tagged integration
gate), records per-row provenance and `UNCLASSIFIED=0`, and finalizes the
report. After two-argument report lint, REPAIR RESOLVED, and TASK task-20
COMPLETE, task 21 may dispatch. Window 7 review remains after tasks 19-21.
