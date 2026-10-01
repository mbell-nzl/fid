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

**Tasks 1-18 complete. The Human checklist FAILED (55/2/6). DEC-268 repairs
A1-A3 in Phase 23 as supplemental tasks 19-21 (DEC-269, Human-approved
2026-10-02, option 1).** Branch `phase-23`.

- DEC-267 resolved the final review's MF-2/BC-A. Task 18 repaired MF-1
  (`b84b17c`, `b37c169`).
- Task 19 (A2): Human-mode Stop-family failures are silent.
- Task 20 (A3): `RequirementRef` gains `requirement_id`; the views are fixed.
- Task 21 (A1): the rework Stop resolves through a `review-rework-extension`
  Human Outcome; the rework budget is excluded from generic admission.

Next action: the Runner resumes using
`qa/resume-after-dec-269.md`. Order: tasks 19, 20, 21, then `review-15`
(window 19-21), then a fresh QA Lead `final-review-1.md` (recorded as
review-16) that regenerates the fixtures and authors `human-checklist-2.md`,
then the Human re-check, then RUN-END. The A4-A8 Human operating experience is
on the wrap-up agenda (DEC-268).

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

Full specs and reports: `../phases/phase-23/tasks/`. Tasks 1-18 are complete;
19-21 are supplemental under DEC-269.

1 schema baseline; 2 typed multi-Obligation Envelope; 3 review Claim mode and
Run; 4 Review Records and independence; 5 Finding dispositions; 6 review
evidence, cycles, exhaustion; 7 freshness and supersession; 8 Concerns and
Deviations; 9 Decision References and protected Decisions; 10 Human Outcomes; 11
Attention resolution; 12 requirement references; 13 CLI, codes, idempotency; 14
Human projections; 15 migration 0018 and V1 removal; 16 portable recovery; 17
capability matrix; 18 conformance audit; 19 Stop-family Human failures; 20
requirement identity; 21 resolvable rework exhaustion.

## Validation Contract And QA Artifacts

All under `../phases/phase-23/qa/`. `validation-contract.md` holds **89**
assertions (85 deterministic, 4 Human; VAL-085..089 under DEC-269), all `Ready`; every row names its owning
task. Validators are QA Lead or QA Reviewer except VAL-070..073, the Human rows,
whose fixtures come from task 14. QA plan section 21 confirms DEC-237/240 READY
in commit `74c83b7`. DEC-242 adds executable inapplicability evidence for Task 5;
window 4-6 reviews it. Mid-phase `review-M.md` for windows 1-3, 4-6,
7-9, 10-12, 13-15, 16-18 (review-14 PASS), 19-21 (`review-15`). The task-11
checkpoint passed at re-review 3. Final: `final-review.md` BLOCK; next
`final-review-1.md`.

## Active Constraints

- DEC-243: Task 6 supplies VAL-025 receipt/review links and mismatch
  classifications; Task 10 owns the real positive Human Outcome link and final
  nine-cell closure. Classifier/seeded acceptance is never a durable-link PASS.
- DEC-242: Task 5 adds optional `inapplicable_findings:[{finding_id,reason}]`;
  omitted/empty preserve Task 4 digests. Existing Finding readback and schema
  suffice. Addressed is nonterminal; terminal dispositions share the Review Event.
- The accepted `Condition` encoding is DEC-235: `omitempty` on `outcome_class`,
  `decision_class`, and `finding_id` and **on no other field**. Tagging
  `budget_name` or `capability_id` changes every stored `semantic_key` and
  `request_digest` in existence - measured - and is the one way to reintroduce
  ESC-23-001 (VAL-082). Exclusivity is structural, in the per-type shapes. No
  task **extends** the `case "Condition"` relation, but task 9's new dispatch
  must still **reach** it, or the Phase 22 nullable binding is silently deleted
  with every gate green (DEC-236 R-24; VAL-081).
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

Wide-reach change: Obligation generalization (task 2) and Evidence widening
(task 6), each isolated. **Registry without an owning task** (DEC-167/173) is
the recurring one: seven re-derivations caught thirteen before dispatch, each
inside the previous fix, and the series has **not** converged - what released
dispatch is that the last additions fail open. Treat a registry rejection in any
task as expected-class and escalate it. Passes 7 and 8 found the family one step
worse: a correct prohibition with no statement of what must still **run**
(R-24), and live Phase 22 behavior with no assertion at all (VAL-081/082).
Expect that shape more than outright omission. Also: late Human-readability
discovery in window 13-15; removal ordering; and test-time headroom (DEC-227).

## Success / Exit Criteria

- [ ] All 21 tasks complete with green quality gates.
- [ ] All 89 `VAL-NNN` outcomes recorded with cited evidence.
- [ ] Final QA review has no unresolved blocking finding.
- [ ] Human checklists accepted, including VAL-070 to VAL-073 (re-check in
  `human-checklist-2.md`).
- [ ] Sign-off recorded and all twelve canonical specs synced.
