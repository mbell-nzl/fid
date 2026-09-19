# Phase Brief - Phase 23 (Approved)

## Phase Goal

Deliver the Review and Human judgment slice: an independent Reviewer records
coverage and Findings against an exact Envelope revision and repository
identity; Builder rework preserves review independence; protected or subjective
Obligations receive provenanced Human Outcomes; completion is re-evaluated from
those facts. It also replaces V1 requirement authoring with references and
removes the role-specific verification columns, QA and TL composites, and
role-bound completion this slice owns. Success is a Human who can tell, from fid
output alone, what a Review resolved, what is merely addressed, what still
blocks, and which judgment unblocked it.

## Status

**Approved and READY; all 18 tasks dispatchable; no open decision or
escalation.** Branch `phase-23` at planning base
`0195897a347dfe1872db0f7be3800c5c3ccdefcf`.
Six QA passes returned twenty-three Must-fix revisions plus five the Tech Lead
found while dispositioning them; **all are applied**. The Human approved five
contract amendments on 2026-09-19 - DEC-223 link-preserving
`evidence_records` rebuild, DEC-224 four named Events raise
`human-outcome-required`, DEC-225 independence over the mutating side only,
DEC-226 failure-code producers, precedence, non-committing denials and the
**Failure code** / **Evaluation reason** split, DEC-227 explicit `-timeout 30m`
plus wall-time budgets - and Tech Lead planning dispositions followed in
DEC-228 to DEC-234. The registry-ownership line runs through them: DEC-231 and
DEC-232 closed it on the strict path and kept the journal `Subject` narrow in
favour of a separate `JudgmentSubject`, DEC-233 widened it to any closed
registry a Phase 23 **value** reaches in any family (rows 7 and 8), and DEC-234
extended it to a registry that reaches a Phase 23 family and declines to
constrain it (row 9). Full record:
[phase decisions](../phases/phase-23/decisions.md).

The sixth pass (`qa-plan.md` 17) returned READY EXCEPT ESC-23-001, confirming
DEC-233 in full with the contract untouched at 79 rows; DEC-234 dispositioned
its four Must-fix revisions and four improvements.

**ESC-23-001 is resolved.** The Human chose Option A on 2026-09-20: the three
new `Condition` fields are `omitempty` and absent on the types they do not apply
to, the six Phase 22 types keep `budget_name` and `capability_id`
present-and-null and stay byte-identical, and exclusivity is enforced
structurally by per-type journal shapes keyed on `type` rather than by the
nullable-binding relation. Recorded as DEC-235 with the whole batch: the
`Condition` row in `public-interface.md` section 2 amended, inventory row 8
restated and row 10 created, tasks 9, 10, and 11 rewritten, and the three held
revisions (R-22, R-23, I-37) applied. VAL-078 keeps its assertion and gains an
evidence-class note only.

Next action: a fresh QA Lead pass over the DEC-235 batch, then the Runner's
readiness check and dispatch from task 1.

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
strictly sequential: each depends on the one before; task 1 also on DEC-222/223,
task 10 on DEC-224, task 12 on the checkpoint after task 11.

1 schema baseline; 2 typed multi-Obligation Envelope; 3 review Claim mode and
Run; 4 Review Records and computed independence; 5 Finding dispositions; 6
review evidence, cycles, exhaustion; 7 freshness and supersession; 8 Concerns
and Deviations; 9 Decision References and protected Decisions; 10 Human
Outcomes; 11 Attention resolution; 12 requirement references and traceability;
13 CLI, codes, idempotency; 14 Human projections; 15 migration 0018 and V1
removal; 16 portable recovery; 17 capability matrix; 18 conformance audit.

## Validation Contract

`../phases/phase-23/qa/validation-contract.md` holds **79** assertions (75
deterministic, 4 Human), all `Ready`. By range: VAL-001..005/075 schema, fence,
parity (1, 15, 16); VAL-006..013 Obligations, evaluation, review Claims (2, 3);
VAL-014..023 Review Records, independence, Findings (4, 5); VAL-024..033/077
review evidence, cycles, exhaustion, freshness (6, 7); VAL-034..046/074
Concerns, Deviations, Decisions, Outcomes (8, 9, 10); VAL-047..054/078
Attention, Conditions, requirements, links (11, 12, 17); VAL-055..065/076 CLI,
Human projections, removal, portability (13-16); VAL-066..069/079 availability,
conformance, gate budgets (17, 18); VAL-070..073 Human comprehension at phase
end, fixtures from task 14. Validators are QA Lead or QA Reviewer except those
four Human rows.

## QA Artifacts

All under `../phases/phase-23/qa/`. QA plan authored, sections 12-17 executed and
freshly measured; section 17 is the sixth readiness confirmation, verdict READY
EXCEPT ESC-23-001, now resolved. One pass remains before dispatch: the DEC-235
batch amends an accepted public value and rewrites three task specs, so it needs
fresh QA Lead eyes.
Mid-phase `review-M.md` for windows 1-3, 4-6, 7-9, 10-12, 13-15, 16-18;
`checkpoint-after-task-11.md` must PASS before task 12; `human-fixtures.md` and
the VAL-070..072 fixtures come from task 14, regenerated at task 18.

## Active Constraints

- The accepted `Condition` encoding is DEC-235: `omitempty` on `outcome_class`,
  `decision_class`, and `finding_id` and **on no other field**. Tagging
  `budget_name` or `capability_id` changes every stored `semantic_key` and
  `request_digest` in existence - measured - and is the one way to reintroduce
  ESC-23-001. Exclusivity is structural, in the per-type shapes; the
  `case "Condition"` relation branch is not extended by any task.
- A task owns only the files it names; `execution-scope.md` section 6 names
  every closed registry's owner and two standing rules: a journal Event rule is
  mirrored in `semanticEventRules` in the same task (DEC-229), and a task that
  unfences a family owns that family's entry in **every** closed journal
  registry that family reaches (DEC-231, generalised by DEC-232). DEC-233 widens
  the frame: the obligation covers every closed journal registry a Phase 23
  **value or Effect** reaches, whichever family carries it - nine known,
  inventoried there with their rejection messages - **ten** rows now: 7 and 8
  reached through live Phase 21/22 families, 9 (DEC-234) a registry a Phase 23
  family reaches and declines to constrain, and 10 (DEC-235) the per-type
  `Condition` shape map, the one registry this phase creates rather than finds. Ownership
  is only assignable where it is executable, and a registry absent from that
  inventory is a Developer escalation, not a local fix.
- A **Failure code** names one refused request; an **Evaluation reason**
  describes why one Obligation is unsatisfied. `GLOSSARY.md` defines both; file
  every new identifier in the set it belongs to (DEC-226).
- Every task runs the tagged integration gate from task 1 with an explicit
  `-timeout 30m` and records the `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow` times (DEC-227). Hosted CI does not run that gate (CF-041),
  so a green PR is never phase-level evidence.
- Removal is replacement-bound: task 15 removes nothing whose replacement is not
  already green, and defers nothing because a test broke. Independence is
  computed from Claim and Run history over mutating Claims and their Runs only
  (DEC-225); no role string decides it, and a Human Outcome is never deterministic
  proof.
- Phase 23 adds no denial Event: a Stop-blocked request commits nothing and
  returns the existing Stop and Attention IDs with a null EventRef (DEC-226).
  Envelope revision stays 1, the primary Issue single; parallel Claims,
  multi-Issue Envelopes, `retry-exhausted` remain Phase 24.

## Interfaces / Subsystems In Scope

`internal/review/` (new) owns Review, Finding, Concern, Deviation, Decision
Reference, and Human Outcome rules. `internal/workflow/`, `internal/evidence/`,
`internal/claims/` carry Obligations, protected Decisions, budgets, Stops,
Attention, freshness, review Claim mode, and the closed allowlists in
`activation.go` and `stop.go`. `internal/journal/`, `internal/db/`,
`migrations/` own 0017, 0018, the fence, and the strict-path registries. `internal/links/` is restricted, `internal/requirements/`
removed; `internal/export/`, `internal/maintenance/`, `internal/activity/`,
`cmd/fid/` carry adapters, projections, archive semantics, and convergence.

## Known Risks

Wide-reach change: Obligation generalization (task 2) and Evidence widening
(task 6), each isolated. Registry without an owning task (DEC-167/173) is the
recurring one: six re-derivations caught twelve before dispatch, each inside the
previous fix, and the last two had to widen the frame to find their own. The
sixth still added one (row 9), so the series has **not** converged to zero; what
released dispatch is that the last two additions fail open and block nothing,
not that the search is finished - row 10 was then added by DEC-235. Treat a
registry rejection in any task as expected-class and escalate it. Of the same
family and now closed: a live journal **value** shape that gains a field is not
monotonic over stored history (ESC-23-001, resolved by DEC-235), and `Condition`
was the only one this phase. The residual risk is implementation drift against
that decision - `omitempty` on the wrong field, or a Phase 22 shape row not
transcribed exactly - which task 9's byte-identity criterion exists to catch.
Late Human-readability discovery: the Human-text check in window 13-15. Removal
ordering: task 15 depends on fourteen replacements; a slip amends the task
rather than deferring a removal. Test-time headroom: bounded by DEC-227.

## Success / Exit Criteria

- [ ] All 18 tasks complete with green quality gates.
- [ ] All 79 `VAL-NNN` outcomes recorded with cited evidence.
- [ ] Final QA review has no unresolved blocking finding.
- [ ] Human checklist accepted, including VAL-070 to VAL-073.
- [ ] Sign-off recorded and all twelve canonical specs synced.

## Next Planning Checkpoint

Phase 24: multi-Issue Envelope revision, scope conflict, parallel Claims, partial
progress, generic `retry-exhausted`, `triage_events` and V1 removals.
