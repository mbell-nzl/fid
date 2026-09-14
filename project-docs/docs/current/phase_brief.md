# Phase Brief - Phase 23 (Approved)

## Phase Goal

Deliver the Review and Human judgment slice: an independent Reviewer records
coverage and zero or more Findings against an exact Envelope revision and
repository identity; Builder rework preserves review independence; protected or
subjective Obligations receive provenanced Human Outcomes; completion is
re-evaluated from those facts. It also replaces V1 requirement authoring with
references and removes the role-specific verification columns, QA and TL
composites, and role-bound completion behavior this slice owns. Success is a
Human who can tell, from fid output alone, what a Review resolved, what is
merely addressed, what still blocks, and which judgment unblocked it.

## Status

**Approved, not started.** Branch `phase-23` at planning base
`0195897a347dfe1872db0f7be3800c5c3ccdefcf`. DEC-222 is accepted, so the contract
set below is binding authority. Next action: the **Human** activates QA Lead to
author the Heavy-mode QA plan and fresh command and feasibility validation; Tech
Lead performed no handoff. No task may be dispatched until that QA plan exists
and the Phase Runner's readiness check passes.

## Operating Mode

**Heavy.** Four entry criteria met: persistence-contract and migration change
(0013/0015 baseline edits, 0017, 0018); public-interface change (eighteen new
or replaced leaves, seven grant actions, twenty failure codes); an authority
and trust surface (Human Outcome provenance, protected Decisions); and
replacement-bound removal of fail-closed V1 behavior. Matrix:
`../../specs/09_AI_WORKFLOW.md` section 4. Record any promotion here.

## Accepted Contract Set

Accepted under DEC-222 and the only authority for implementation, all under
`../phases/phase-23/`: `contracts.md`,
`public-interface.md`, `schema.md`, `execution-scope.md`,
`v1-convergence-inventory.md`, and `contract-coverage.md` (clause to schema row
to task to assertion).

## Task Index

Full specs and reports: `../phases/phase-23/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Review and judgment schema baseline | Pending | DEC-222 |
| 2 | Typed multi-Obligation Envelope contract | Pending | 1 |
| 3 | Non-mutating review Claim mode and Reviewer Run | Pending | 2 |
| 4 | Immutable Review Records and computed independence | Pending | 3 |
| 5 | Finding dispositions and resolution authority | Pending | 4 |
| 6 | Review evidence, review cycles, and rework exhaustion | Pending | 5 |
| 7 | Review freshness and Human Outcome supersession | Pending | 6 |
| 8 | Generic Concerns and append-only Deviations | Pending | 7 |
| 9 | Decision References and protected Decisions | Pending | 8 |
| 10 | Provenanced Human Outcomes | Pending | 9 |
| 11 | Attention resolution through Human Outcomes | Pending | 10 |
| 12 | Requirement references and a single traceability model | Pending | 11, checkpoint |
| 13 | Phase 23 CLI surface, failure codes, and idempotency | Pending | 12 |
| 14 | Human review and judgment projections | Pending | 13 |
| 15 | Migration 0018 and V1 judgment removal | Pending | 14 |
| 16 | Portable recovery for the review and judgment domain | Pending | 15 |
| 17 | Final capability matrix and availability audit | Pending | 16 |
| 18 | Conformance assembly and wire-surface audit | Pending | 17 |

## Validation Contract

The full contract is `../phases/phase-23/qa/validation-contract.md`: 74
assertions, VAL-001 to VAL-074, all `Draft` until QA Lead refinement. Index only:

| Range | Subject | Owning tasks | Validator |
|---|---|---|---|
| VAL-001..005 | Schema baseline, fence, parity | 1, 15, 16 | QA Lead / Reviewer |
| VAL-006..013 | Obligations, evaluation, review Claims | 2, 3 | QA Lead / Reviewer |
| VAL-014..023 | Review Records, independence, Findings | 4, 5 | QA Lead / Reviewer |
| VAL-024..033 | Review evidence, cycles, exhaustion, freshness | 6, 7 | QA Lead / Reviewer |
| VAL-034..046, 074 | Concerns, Deviations, Decisions, Outcomes | 8, 9, 10 | QA Lead / Reviewer |
| VAL-047..054 | Attention resolution, requirements, links | 11, 12, 17 | QA Lead / Reviewer |
| VAL-055..065 | CLI, Human projections, removal, portability | 13-16 | QA Lead / Reviewer |
| VAL-066..069 | Availability and conformance audit | 17, 18 | QA Lead / Reviewer |
| VAL-070..073 | Comprehension and product judgment | phase-end | Human |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-23/qa/validation-contract.md` | Accepted; rows `Draft` until QA Lead refinement |
| QA plan | `../phases/phase-23/qa/qa-plan.md` | **Missing; blocks task 1.** Human activates QA Lead |
| Mid-phase reviews | `../phases/phase-23/qa/review-M.md` | Windows 1-3, 4-6, 7-9, 10-12, 13-15, 16-18 |
| Independent checkpoint | `../phases/phase-23/qa/checkpoint-after-task-11.md` | Must PASS before task 12 |
| Human fixtures | `../phases/phase-23/qa/human-fixtures.md` | Authored by task 14 |
| Final review / checklist | `../phases/phase-23/qa/` | Pending |

Heavy mode requires fresh command and feasibility validation.

## Active Constraints

- DEC-222 is accepted. Nothing dispatches before the QA plan and readiness check.
- Every task runs the tagged integration gate from task 1. Hosted CI does not
  run it (CF-041), so a green PR is never phase-level evidence.
- Removal is replacement-bound: task 15 removes nothing whose replacement is not
  already green, and defers nothing because a test broke.
- Independence is computed from Claim and Run history; no role string or
  attribution label may decide it.
- A Human Outcome is never stored, evaluated, or rendered as deterministic proof.
- Envelope revision stays 1 and the primary Issue single; parallel Claims,
  multi-Issue Envelopes, and `retry-exhausted` remain Phase 24.

## Interfaces / Subsystems In Scope

- `internal/review/` (new) - Review, Finding, Concern, Deviation, Decision
  Reference, and Human Outcome rules.
- `internal/workflow/`, `internal/evidence/`, `internal/claims/` - Obligations,
  protected Decisions, budgets, Stops, Attention, freshness, review Claim mode.
- `internal/journal/`, `internal/db/`, `migrations/` - 0017, 0018, fence,
  effect rules, parity.
- `internal/links/` restricted; `internal/requirements/` removed.
- `internal/export/`, `internal/maintenance/`, `internal/activity/`, `cmd/fid/`
  - adapters, projections, removal convergence.

## Known Risks

- **Wide-reach change.** Obligation generalization (task 2) and Evidence
  widening (task 6) touch many producers; each is isolated in its own task.
- **Clause without a schema row** (Phase 22 DEC-195). Mitigated by
  `contract-coverage.md`, re-derived independently by QA Lead before approval.
- **Late Human-readability discovery** (Phase 22). Mitigated by a deterministic
  Human-text check in the window 13-15 review against task 14's fixtures.
- **Removal ordering.** Task 15 depends on fourteen replacements; a slip amends
  the task rather than silently deferring a removal.

## Open Questions

Both planning questions are resolved by the Human decisions recorded in
`../phases/phase-23/decisions.md`. None remain open.

## Success / Exit Criteria

- [ ] All 18 tasks complete with green quality gates.
- [ ] All 74 `VAL-NNN` outcomes recorded with cited evidence.
- [ ] Final QA review has no unresolved blocking finding.
- [ ] Human checklist accepted, including VAL-070 to VAL-073.
- [ ] Sign-off recorded and all twelve canonical specs synced.

## Next Planning Checkpoint

Phase 24, parallel and evolving work: multi-Issue Envelope revision, scope
conflict, parallel Claims, partial progress, generic `retry-exhausted`, and the
`triage_events` and V1 Claim and Run removals.
