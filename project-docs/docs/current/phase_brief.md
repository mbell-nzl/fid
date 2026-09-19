# Phase Brief - Phase 23 (Approved)

## Phase Goal

Deliver the Review and Human judgment slice: an independent Reviewer records
coverage and Findings against an exact Envelope revision and repository
identity; Builder rework preserves review independence; protected or subjective
Obligations receive provenanced Human Outcomes; completion is re-evaluated from
those facts. It also replaces V1 requirement authoring with references and
removes the role-specific verification columns, QA and TL composites, and
role-bound completion behavior this slice owns. Success is a Human who can tell,
from fid output alone, what a Review resolved, what is merely addressed, what
still blocks, and which judgment unblocked it.

## Status

**Approved, not started; all amendments approved.** Branch `phase-23` at
planning base `0195897a347dfe1872db0f7be3800c5c3ccdefcf`. The QA Lead's
Heavy-mode plan returned thirteen Must-fix revisions and sixteen improvements;
its 2026-09-19 readiness refresh (`qa-plan.md` section 12) returned R-14, R-15
and I-17..I-26. All are dispositioned and applied; tables and reasoning in
[phase decisions](../phases/phase-23/decisions.md).

The Human approved all five contract amendments on 2026-09-19: DEC-223
link-preserving `evidence_records` rebuild; DEC-224 four named Events raise
`human-outcome-required`, two gaining effect rows; DEC-225 independence
conflicts with the mutating side only; DEC-226 failure-code producers,
precedence, non-committing denials, and the **Failure code** / **Evaluation
reason** split in `GLOSSARY.md`; DEC-227 explicit `-timeout 30m` plus wall-time
budgets. DEC-229 then gave four closed registries owning tasks and corrected
CF-042's cause to a per-run leak with a causal fix in task 1; DEC-230 removes
`decision-not-found`, an added failure code no leaf can return.

Next action: a **fresh QA Lead** confirms the R-14/R-15 edits, then the
Runner's readiness check; no task dispatches before both pass.

## Operating Mode

**Heavy** (`../../specs/09_AI_WORKFLOW.md` section 4). Four entry criteria met:
persistence and migration change; public-interface change (eighteen leaves,
seven grant actions, a new closed failure-code set); an authority and trust
surface; replacement-bound removal of fail-closed V1 behavior.

## Accepted Contract Set

Accepted under DEC-222 and the only implementation authority, all under
`../phases/phase-23/`: `contracts.md`, `public-interface.md`, `schema.md`,
`execution-scope.md`, `v1-convergence-inventory.md`, `contract-coverage.md`.
Amendments are inline with their DEC and approval state.

## Task Index

Full specs and reports: `../phases/phase-23/tasks/`. All 18 are **Pending** and
strictly sequential: each depends on the one before; task 1 also on
DEC-222/223, task 10 on DEC-224, task 12 on the checkpoint after task 11.

1 schema baseline; 2 typed multi-Obligation Envelope; 3 review Claim mode and
Run; 4 Review Records and computed independence; 5 Finding dispositions; 6
review evidence, cycles, exhaustion; 7 freshness and supersession; 8 Concerns
and Deviations; 9 Decision References and protected Decisions; 10 Human
Outcomes; 11 Attention resolution; 12 requirement references and one
traceability model; 13 CLI, codes, idempotency; 14 Human projections; 15
migration 0018 and V1 removal; 16 portable recovery; 17 capability matrix; 18
conformance and wire audit.

## Validation Contract

`../phases/phase-23/qa/validation-contract.md` holds **79** assertions (75
deterministic, 4 Human), all `Ready`. By range: VAL-001..005/075 schema
baseline, fence, parity (tasks 1, 15, 16); VAL-006..013 Obligations,
evaluation, review Claims (2, 3); VAL-014..023 Review Records, independence,
Findings (4, 5); VAL-024..033/077 review evidence, cycles, exhaustion,
freshness (6, 7); VAL-034..046/074 Concerns, Deviations, Decisions, Outcomes
(8, 9, 10); VAL-047..054/078 Attention, Conditions, requirements, links (11,
12, 17); VAL-055..065/076 CLI, Human projections, removal, portability (13-16);
VAL-066..069/079 availability, conformance audit, gate budgets (17, 18);
VAL-070..073 Human comprehension and product judgment at phase end, fixtures
from task 14. Validators are QA Lead or QA Reviewer except those four Human
rows.

## QA Artifacts

All under `../phases/phase-23/qa/`. QA plan authored; its section 12 readiness
refresh is executed and every result there freshly measured, but the plan is
**NOT READY**, and nothing dispatches, until a fresh QA Lead confirms the
DEC-229 edits and the Runner's readiness check passes. Mid-phase
`review-M.md` for windows 1-3, 4-6, 7-9, 10-12, 13-15, 16-18;
`checkpoint-after-task-11.md` must PASS before task 12; `human-fixtures.md` and
the runnable VAL-070..072 fixtures come from task 14 and are regenerated at
task 18; final review and human checklist pending.

## Active Constraints

- A task owns only the files it names; `execution-scope.md` section 6 names
  every closed registry's owner and the standing rule that a journal Event rule
  change is mirrored in `semanticEventRules` in the same task (DEC-229).
- A **Failure code** names one refused request; an **Evaluation reason**
  describes why one Obligation is unsatisfied. `GLOSSARY.md` defines both; file
  every new identifier in the set it belongs to (DEC-226).
- Every task runs the tagged integration gate from task 1 with an explicit
  `-timeout 30m`, and records the `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow` package times (DEC-227). Hosted CI does not run that gate
  (CF-041), so a green PR is never phase-level evidence.
- Removal is replacement-bound: task 15 removes nothing whose replacement is not
  already green, and defers nothing because a test broke.
- Independence is computed from Claim and Run history, over mutating Claims and
  their Runs only (DEC-225); no role string may decide it. A Human Outcome is
  never stored, evaluated, or rendered as deterministic proof.
- Phase 23 adds no denial Event: a Stop-blocked request commits nothing and
  returns the existing Stop and Attention IDs with a null EventRef (DEC-226).
- Envelope revision stays 1 and the primary Issue single; parallel Claims,
  multi-Issue Envelopes, and `retry-exhausted` remain Phase 24.

## Interfaces / Subsystems In Scope

`internal/review/` (new) owns Review, Finding, Concern, Deviation, Decision
Reference, and Human Outcome rules. `internal/workflow/`, `internal/evidence/`,
`internal/claims/` carry Obligations, protected Decisions, budgets, Stops,
Attention, freshness, and review Claim mode; the closed allowlists are
`workflow/activation.go` (grant actions, tasks 4, 5, 8, 9, 10) and
`workflow/stop.go` (Stop codes and Conditions, tasks 6, 9, 10, 11).
`internal/journal/`, `internal/db/`, `migrations/` own 0017, 0018, the fence,
effect rules, column domains, and parity. `internal/links/` is restricted and
`internal/requirements/` is removed. `internal/export/`,
`internal/maintenance/`, `internal/activity/`, and `cmd/fid/` carry adapters,
projections, the mirrored archive semantics, and removal convergence.

## Known Risks

Wide-reach change: Obligation generalization (task 2) and Evidence widening
(task 6), each isolated in its own task. Clause without a schema row (DEC-195)
and registry without an owning task (DEC-167/173): the coverage matrix and two
independent QA re-derivations caught four before dispatch (R-02, R-12, R-14),
all now owned and assertion-backed. Late Human-readability discovery: the
deterministic Human-text check in window 13-15. Removal ordering: task 15
depends on fourteen replacements; a slip amends the task rather than deferring
a removal. Test-time headroom: bounded by DEC-227.

## Success / Exit Criteria

- [ ] All 18 tasks complete with green quality gates.
- [ ] All 79 `VAL-NNN` outcomes recorded with cited evidence.
- [ ] Final QA review has no unresolved blocking finding.
- [ ] Human checklist accepted, including VAL-070 to VAL-073.
- [ ] Sign-off recorded and all twelve canonical specs synced.

## Next Planning Checkpoint

Phase 24: multi-Issue Envelope revision, scope conflict, parallel Claims,
partial progress, generic `retry-exhausted`, `triage_events` and V1 removals.
