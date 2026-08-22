# Phase Brief - Phase 21 (Planning / Awaiting Approval)

> Default execution brief. Product implementation is not authorized.

## Phase Goal

Deliver one complete one-Issue operational walking skeleton: activated trust,
immutable Gate Registration, one Work Envelope and gate-evidence Obligation,
generic Claim/Run authority, fid-executed Receipt evidence, deterministic
completion/resume, export create/validate, rebuild parity, and owned V1
removal. Routine is external kit language and is not fid state.

## Status

**Planning complete; awaiting Human execution approval.** Branch `phase-21`.
Gate D3.2 is accepted by DEC-153, the exact slice is DEC-154, and carry-forward
disposition is DEC-155. Do not dispatch task 1 until the Human accepts this
Heavy plan and the QA Lead supplies a Ready QA plan with fresh feasibility.

## Operating Mode

**Heavy.** This phase changes authority, public JSON/Event contracts, schema,
trusted local execution, redaction/Artifact safety, completion, and resume
projections. Two internal groups remain one numeric phase under ADR-0012.

## Task Index

Full specs are under `../phases/phase-21/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Operational schema and projection | Pending approval | None |
| 2 | Event v2 and durable idempotency | Pending | 1 |
| 3 | Capability discovery matrix | Pending | 1-2 |
| 4 | Activation lifecycle and trust profiles | Pending | 1-3, Review 1 |
| 5 | Activated Gate Registration | Pending | 4 |
| 6 | Envelope, Obligation, and budget values | Pending | 5 |
| 7 | Envelope authority evaluator | Pending | 4-6, Review 2 |
| 8 | Issue Actor/idempotency convergence | Pending | 7 |
| 9 | Retained Actor convergence and policy removal | Pending | 7-8 |
| 10 | Generic mutation Claims | Pending | 6-9, Review 3 |
| 11 | Generic Runs and budget bindings | Pending | 10 |
| 12 | Execution-bound transition foundation | Pending | 7,10-11 |
| CP | Independent Group 1 checkpoint | Required | 1-12, Reviews 1-4 |
| 13 | Portable executor policy | Pending | CP PASS |
| 14 | Live output framing and redaction | Pending | 13 |
| 15 | Trusted Attempt and Receipt execution | Pending | 13-14 |
| 16 | Safe content-addressed Artifacts | Pending | 15, Review 5 |
| 17 | Receipt Evidence and evaluation | Pending | 15-16 |
| 18 | Envelope completion and resume | Pending | 17 |
| 19 | Deterministic export create/validate | Pending | 16-18, Review 6 |
| 20 | Maintenance and full rebuild parity | Pending | 19 |
| 21 | V1 authority/execution removal | Pending | 7-20 |
| 22 | Real-process journey and conformance | Pending | 21, Review 7 |

Reviews occur after tasks 1-3, 4-6, 7-9, 10-12, 13-15, 16-18, and
19-21. Task 22 is followed by final QA and the Human checklist.

## Validation Contract

Full contract: `../phases/phase-21/qa/validation-contract.md` (VAL-001..064,
Draft). Group 1 owns VAL-001..028 plus VAL-062/063; Group 2 owns VAL-029..060
and VAL-064; VAL-061 is phase-end delivery-control evidence. VAL-036 is the
only Human assertion. Every row has exactly one validator.

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-21/qa/validation-contract.md` | Draft |
| QA plan | `../phases/phase-21/qa/qa-plan.md` | Missing; required before task 1 |
| Group 1 checkpoint | `../phases/phase-21/qa/group-1-checkpoint-contract.md` | Contract Ready; review pending |
| Mid-phase reviews | `../phases/phase-21/qa/review-M.md` | Seven windows pending |
| Final review/checklist | `../phases/phase-21/qa/` | Pending |

## Active Constraints

- Registration from finalized activated Gate values precedes Envelope create.
- Actor capabilities/provenance come only from active Activation; caller
  declarations confer no kit authority.
- Only portable standard-library executor controls are accepted; no OS
  containment claim or new dependency exists.
- Attention, Stops, budget consumption, recovery/respawn, import, review,
  Human Outcomes, parallelism, attestation, and telemetry remain deferred.
- `export.portable.v1` stays unavailable; create/validate is not import parity.
- Group 2 requires independent VAL-028 PASS with fresh feasibility evidence.
- Each task runs its full gate and task-local report lint before next dispatch.

## Parity Matrix

| Axis | Live | Rebuild | Export/validate | Evidence |
|---|---|---|---|---|
| State | canonical v2 reads | byte-identical named reads | integrity-bound snapshot | VAL-006/055 |
| Authority | active Activation/grants | same resolved facts | metadata only | VAL-010/018 |
| Evidence | Receipt/Evidence/Artifact | same evaluation/reasons | available material/digests | VAL-043..052 |
| Capability | exact 12-row matrix | identical matrix | snapshot, no import claim | VAL-007/053 |

## Known Risks

- Child-process termination, redaction, and path-race verification may expose
  a platform/dependency need; VAL-028 must surface that before task 13.
- Migration 0014 cannot run until generic authority/Claim/Run replacements
  pass; frozen table scope must not become a second authority.
- The large public-contract surface is controlled by 22 causal tasks and
  three-task independent review windows, not by weakening assertions.

## Success / Exit Criteria

- [ ] Tasks 1-22 complete with required gates and task-local report lint.
- [ ] VAL-028 checkpoint PASS before Group 2.
- [ ] VAL-001..064 have final dispositions; no FAIL/BLOCKED is ignored.
- [ ] Final QA has no unresolved blocking finding and VAL-036 Human check passes.
- [ ] DEC/spec/removal/triage audits and phase sign-off are recorded.

## Next Planning Checkpoint

Human reviews the Heavy mode, DEC-154 contract, validation contract, and task
pack. After approval, QA Lead drafts/refines a Ready Heavy QA plan; Tech Lead
applies any assertion/command refinements before the Phase Runner starts.
