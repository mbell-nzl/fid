# Phase Brief - Phase 21

> Status: **Active / task-22 STOPPED / Recovery Authorized (DEC-174)**
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
  walks; task 20's scope was amended to its gates' reach and it gained the
  tagged integration gate. Closed with task 20.
- DEC-174: the §16 journey authors its own workflow (`workflow define
  --applies-to task` then `workflow assign --kind task`, before `issue create`)
  with exactly the edges `unstarted -> started` and `started -> completed`; no
  seeded reference workflow can express the accepted shape. The exact mutation
  prefix is 22 Events and VAL-086 matches. No product, command,
  public-interface, Event-schema, transition, or frozen-workflow change.

## Planning Authorities

- Slice semantics/Event fixtures: `../phases/phase-21/contracts.md`; complete
  public interface: `../phases/phase-21/public-interface.md`
- Live V1 inventory: `../phases/phase-21/v1-convergence-inventory.md`; checkpoint
  `../phases/phase-21/qa/group-1-checkpoint-contract.md`
- Validation: `../phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- QA: `qa/qa-lead-handoff.md`, `qa/qa-plan.md` (Ready at `e8a636e`), audit
  `planning-audit.md`; tasks `../phases/phase-21/tasks/task-<N>-*.md`

## Task Pack

Group 1 (tasks 1-12, VAL-001/003-041), its reviews 1-4, and the independent
checkpoint (VAL-042/043) are complete. Group 2 tasks 13-21 (VAL-036,
044-071, 002, 074-083) and reviews 5-7b are complete. Per-task titles,
dependencies, and VAL ownership are in
`../phases/phase-21/tasks/task-<N>-*.md`.

| Task | Title | Depends on | VAL ownership |
|---:|---|---|---|
| 22 | final interface/journey conformance | 21 | 072-073,084-086 |

Final QA owns task-22 evidence and phase-end VAL-087. VAL-036 is the sole
Human checklist assertion.

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

Mutate once: init; define/assign the journey workflow; Issue; discover;
activate; register; create/activate Envelope; resolve Actor; Claim/Run;
nonterminal transition; execute Gate; persist/link Receipt; evaluate Evidence;
terminal-success transition; finish/release; complete; resume; export/validate.
Exactly 22 Events before rebuild. Delete the projection and rebuild, then
repeat only named reads, evaluation, completion state, resume, export, and
negative checks; never create a second Envelope revision, workflow definition,
or assignment.

## Dispatch Authorization

Tasks 1-21 and reviews 1-7b are complete. Task 22 committed the conformance
corpus (`a99b427`) and the real-process journey guard (`236540c`), then returned
`BLOCKED`; the run is STOPPED after TRG-21-021 / ESC-21-015. DEC-174 resolves
the §16 journey defect without changing any command or product contract.

The Phase Runner appends `RUN-RESUME` retaining TRG-21-021 / ESC-21-015, trigger
`task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain`,
attempt 2, resolution authority DEC-174, continuing task 22, and the exact
`task-22 | in-scope-implementation | task-scope | task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain | 2`
tuple. The bound TASK-RECOVERY uses source/authority `TRG-21-021 / DEC-174`
and no-progress `0/2`.

The Developer continues from `236540c` under the amended
`tasks/task-22-walking-skeleton.md`, keeps the committed corpus and journey
guard, authors the §16 journey workflow, extends the journey through completion,
export, and the post-rebuild repeat, commits durable work before long commands,
and runs the complete byte-identical eight-command block with observed
terminations replacing the eight `not run` rows. After two-argument report lint,
REPAIR RESOLVED, and TASK task-22 COMPLETE, final QA owns task-22 evidence,
VAL-087, and the Human checklist.
