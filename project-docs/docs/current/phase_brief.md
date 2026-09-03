# Phase Brief - Phase 21

> Status: **Active / Tasks 1-22 Complete / Ready For Final QA**
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
- DEC-175: a latent production defect no earlier gate could observe belongs to
  the first task whose required gate observes it, with minimum source scope plus
  a regression assertion under a command it already runs. Task 22 owns the
  bounded `internal/workflow/activation.go` Event-type derivation repair
  (`activation.finalized`, not `activation.finalizeted`); task 4 is not
  reopened. No contract, schema, interface, command, corpus, or VAL change.
- DEC-176: the immutable TRG-21-022 attempt-3 TRIAGE makes the schema-v2 STOP
  resume attempt 4. DEC-175 remains the product/scope authority; DEC-176 changes
  only append-only orchestration identity and restarts no-progress at `0/2`.
- DEC-177: Task 22 owns the missing accepted Receipt get/list surface and the
  minimum migration/schema/golden change that makes rejected-Artifact metadata
  durable and canonical. The stale Resume prose is corrected to `fid resume`;
  every accepted wire behavior, 22-Event budget, VAL, and command remains fixed.

## Planning Authorities

- Slice semantics/Event fixtures: `../phases/phase-21/contracts.md`; complete
  public interface: `../phases/phase-21/public-interface.md`
- Live V1 inventory: `../phases/phase-21/v1-convergence-inventory.md`; checkpoint
  `../phases/phase-21/qa/group-1-checkpoint-contract.md`
- Validation: `../phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- QA: `qa/qa-lead-handoff.md`, `qa/qa-plan.md` (Ready at `e8a636e`), audit
  `planning-audit.md`; tasks `../phases/phase-21/tasks/task-<N>-*.md`

## Task Pack

All 22 tasks, reviews 1-7b, and the independent Group 1 checkpoint are
complete. Task 22's final interface/journey conformance report is committed at
`5d8204c`; its product commit is `ac6270a`. Final QA independently verifies
task-22 evidence and phase-end VAL-087. VAL-036 remains the sole Human
checklist assertion.

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

Tasks 1-22 are complete. Task 22 recovery attempt 6 is closed by `REPAIR
task-22 RESOLVED` and `TASK task-22 COMPLETE` at the end of the ledger. Its
eight-command gate and the Phase Runner's authoritative two-argument report
lint passed with observed terminations. The ledger validates as
`ADVANCEABLE`; no escalation or task recovery remains open.

The next authorized action is final QA under `qa/qa-plan.md` section 7.9. The
QA Lead creates `qa/final-review.md`, records an outcome for every VAL, runs
QA-81..93 including the full phase gate in its own session, and presents
`qa/human-checklist.md` for the Human-owned VAL-036 checks. Do not append
`RUN-END` until final QA passes and the Human checklist is complete.
