# Phase Brief - Phase 23 (Approved)

## Phase Goal

Deliver the Review and Human judgment slice: an independent Reviewer records
coverage and Findings against an exact Envelope revision and repository
identity; Builder rework preserves review independence; protected or subjective
Obligations receive provenanced Human Outcomes; completion is re-evaluated from
those facts. It also replaces V1 requirement authoring with references and
removes the V1 verification columns, composites, and role-bound completion.
Success is a Human who can tell, from fid output alone, what a Review resolved,
what is merely addressed, what still blocks, and which judgment unblocked it.

## Status

**Approved; tasks 1-8 dispatchable now, 9-11 after one QA confirmation pass; no
open decision and no registered escalation.** Branch `phase-23` at planning base
`0195897a347dfe1872db0f7be3800c5c3ccdefcf`.

Seven QA passes returned twenty-seven Must-fix revisions plus five the Tech Lead
found while dispositioning them; **all are applied**. Six Human contract
amendments (DEC-223..227, DEC-235); Tech Lead dispositions DEC-228..234 and
DEC-236. Record: [phase decisions](../phases/phase-23/decisions.md).

**ESC-23-001 is resolved** (DEC-235, Human, Option A) and **QA pass 7
(`qa-plan.md` 18) returned READY EXCEPT TASKS 9, 10, AND 11**, confirming Option
A transcribed correctly in all three specs. **DEC-236 dispositioned all four
Must-fix revisions and three improvements; no task is held and no escalation was
allocated.** Its two standing consequences are the first two Active Constraints
below. R-25 made task 9's byte-identity bar a six-type table; R-27 gave tasks 9
and 10 a `currentStopDetailShapes` proof route that reaches the registry.
VAL-078 is refined and strengthened; the contract stays at 79 rows.

Next action: a QA Lead pass over the DEC-236 batch for tasks 9, 10, and 11, then
the Runner's readiness check. Tasks 1-8 may dispatch in parallel with it.

## Operating Mode

**Heavy** (`../../specs/09_AI_WORKFLOW.md` section 4). Four entry criteria met:
persistence and migration change; public-interface change (eighteen leaves,
seven grant actions, a closed failure-code set); an authority and trust surface;
replacement-bound removal of fail-closed V1 behavior.

## Accepted Contract Set

Accepted under DEC-222 and the only implementation authority, all under
`../phases/phase-23/`: `contracts.md`, `public-interface.md`, `schema.md`,
`execution-scope.md`, `v1-convergence-inventory.md`, `contract-coverage.md`;
amendments inline with their DEC and approval state.

## Task Index

Full specs and reports: `../phases/phase-23/tasks/`. All 18 are **Pending** and
strictly sequential; task 12 also depends on the checkpoint after task 11.

1 schema baseline; 2 typed multi-Obligation Envelope; 3 review Claim mode and
Run; 4 Review Records and independence; 5 Finding dispositions; 6 review
evidence, cycles, exhaustion; 7 freshness and supersession; 8 Concerns and
Deviations; 9 Decision References and protected Decisions; 10 Human Outcomes; 11
Attention resolution; 12 requirement references; 13 CLI, codes, idempotency; 14
Human projections; 15 migration 0018 and V1 removal; 16 portable recovery; 17
capability matrix; 18 conformance audit.

## Validation Contract And QA Artifacts

All under `../phases/phase-23/qa/`. `validation-contract.md` holds **79**
assertions (75 deterministic, 4 Human), all `Ready`; every row names its owning
task, so the contract is its own index. Validators are QA Lead or QA Reviewer
except VAL-070..073, the Human comprehension rows, whose fixtures come from task
14. The QA plan's sections 12-18 are executed and freshly measured; section 18 is
the seventh readiness confirmation, verdict READY EXCEPT TASKS 9, 10, AND 11, and
one pass remains for those three. Mid-phase `review-M.md` for windows 1-3, 4-6,
7-9, 10-12, 13-15, 16-18; `checkpoint-after-task-11.md` must PASS before task 12.

## Active Constraints

- The accepted `Condition` encoding is DEC-235: `omitempty` on `outcome_class`,
  `decision_class`, and `finding_id` and **on no other field**. Tagging
  `budget_name` or `capability_id` changes every stored `semantic_key` and
  `request_digest` in existence - measured - and is the one way to reintroduce
  ESC-23-001. Exclusivity is structural, in the per-type shapes. No task
  **extends** the `case "Condition"` relation, but task 9's new dispatch must
  still **reach** it, or the Phase 22 nullable binding is silently deleted with
  every gate green (DEC-236 R-24).
- The discriminator pairing is `contracts.md` section 11's and is exhaustive:
  `outcome_class`/`human-outcome-recorded`, `decision_class`/`decision-satisfied`,
  `finding_id`/`finding-resolved`, and **none on `review-recorded`** (DEC-236).
  No field serves two types.
- A task owns only the files it names; `execution-scope.md` section 6 names
  every closed registry's owner and two standing rules: a journal Event rule is
  mirrored in `semanticEventRules` in the same task (DEC-229), and a task owns
  the entry in **every** closed journal registry any Phase 23 **value or
  Effect** it carries reaches (DEC-231/232, widened by DEC-233; **ten** rows).
  Ownership is only assignable where it is executable; a registry absent from
  that inventory is a Developer escalation, not a local fix.
- A **Failure code** names one refused request; an **Evaluation reason** says why
  one Obligation is unsatisfied. File every new identifier in its set (DEC-226).
- Every task runs the tagged integration gate from task 1 with an explicit
  `-timeout 30m` and records the `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow` times (DEC-227). Hosted CI does not run it (CF-041), so a
  green PR is never phase-level evidence.
- Removal is replacement-bound: task 15 removes nothing whose replacement is not
  already green, and defers nothing because a test broke. Independence is
  computed from Claim and Run history over mutating Claims and their Runs only
  (DEC-225); no role string decides it, and a Human Outcome is never proof.
- Phase 23 adds no denial Event: a Stop-blocked request commits nothing and
  returns the existing Stop and Attention IDs with a null EventRef (DEC-226).
  Envelope revision stays 1 and the primary Issue single; parallel Claims,
  multi-Issue Envelopes, and `retry-exhausted` remain Phase 24.

## Interfaces / Subsystems In Scope

`internal/review/` (new) owns Review, Finding, Concern, Deviation, Decision
Reference, and Human Outcome rules. `internal/workflow/`, `internal/evidence/`,
`internal/claims/` carry Obligations, protected Decisions, budgets, Stops,
Attention, freshness, review Claim mode, and the closed allowlists in
`activation.go` and `stop.go`. `internal/journal/`, `internal/db/`, `migrations/`
own 0017, 0018, the fence, and the strict-path registries. `internal/links/` is
restricted and `internal/requirements/` removed; `internal/export/`,
`internal/maintenance/`, `internal/activity/`, and `cmd/fid/` carry adapters,
projections, and archive semantics.

## Known Risks

Wide-reach change: Obligation generalization (task 2) and Evidence widening
(task 6), each isolated. **Registry without an owning task** (DEC-167/173) is
the recurring one: seven re-derivations caught thirteen before dispatch, each
inside the previous fix, and the series has **not** converged to zero - what
released dispatch is that the last additions fail open, not that the search is
finished. Treat a registry rejection in any task as expected-class and escalate
it. Pass 7 found the same family one step worse - a correct prohibition with no
statement of what must still **run** (DEC-236 R-24), where the literal reading
deleted an invariant and passed every gate. Expect more of that shape than of
outright omission. Also: late Human-readability discovery in window 13-15;
removal ordering, where a slip amends task 15 rather than deferring a removal;
and test-time headroom, bounded by DEC-227.

## Success / Exit Criteria

- [ ] All 18 tasks complete with green quality gates.
- [ ] All 79 `VAL-NNN` outcomes recorded with cited evidence.
- [ ] Final QA review has no unresolved blocking finding.
- [ ] Human checklist accepted, including VAL-070 to VAL-073.
- [ ] Sign-off recorded and all twelve canonical specs synced.

## Next Planning Checkpoint

Phase 24: multi-Issue Envelope revision, scope conflict, parallel Claims, partial
progress, generic `retry-exhausted`, `triage_events` and V1 removals.
