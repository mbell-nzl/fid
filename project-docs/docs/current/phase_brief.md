# Phase Brief - Phase 22 (Planning)

## Phase Goal

Deliver the one-Issue failure-and-recovery journey: interrupt execution,
reconcile repository state, start a linked successor Run, enforce budgets and
typed Stops, dispose Attention, recover activation, stale old proof, import/
export the same state, rebuild, and resume without contradictory facts.

## Status

**Planning; execution unauthorized.** Branch `phase-22`, created from merged
main `89d964a84dbdd297fdf472d6f78b001e2efac702`. Phase 21 is wrapped under
DEC-179 and already merged. Phase 22's proposed contracts and task pack are
written; a Ready QA plan, independent consistency review, TL refinements, and
Human acceptance are pending. Do not dispatch a Developer yet.

## Operating Mode

**Heavy (proposed).** Persistence, public JSON, executor authority, and import
publication change. Full artifacts and fresh QA feasibility are required under
`09_AI_WORKFLOW.md` section 4. Human approves mode with the task list.

## Task Index

All specs are under `../phases/phase-22/tasks/`; reports use `task-N-done.md`.
Tasks are sequential; reviews follow every three tasks.

| Task | Title | Status | Dependency |
|---|---|---|---|
| [1](../phases/phase-22/tasks/task-1.md) | Recovery schema and replay foundation | Pending | Ready QA + approval |
| [2](../phases/phase-22/tasks/task-2.md) | Typed Stops and Attention domain | Pending | 1 |
| [3](../phases/phase-22/tasks/task-3.md) | Lease lifecycle and executor observations | Pending | 2 |
| [4](../phases/phase-22/tasks/task-4.md) | Read-only repository reconciliation | Pending | 3 |
| [5](../phases/phase-22/tasks/task-5.md) | Atomic linked Run recovery | Pending | 4 |
| [6](../phases/phase-22/tasks/task-6.md) | Recovery CLI and early process journey | Pending | 5 |
| [7](../phases/phase-22/tasks/task-7.md) | Durable budget accounting and execution admission | Pending | 6 |
| [8](../phases/phase-22/tasks/task-8.md) | Interrupted Attempts and derived Run closure | Pending | 7 |
| [9](../phases/phase-22/tasks/task-9.md) | Stop, budget, and closure commands | Pending | 8 |
| [10](../phases/phase-22/tasks/task-10.md) | Activation recovery domain | Pending | 9 |
| [11](../phases/phase-22/tasks/task-11.md) | Activation CLI and config boundary | Pending | 10 |
| [12](../phases/phase-22/tasks/task-12.md) | Recovery freshness and actionable resume | Pending | 11 |
| [13](../phases/phase-22/tasks/task-13.md) | Recovery archive format and validation | Pending | 12 + independent checkpoint |
| [14](../phases/phase-22/tasks/task-14.md) | Atomic import into an empty target | Pending | 13 |
| [15](../phases/phase-22/tasks/task-15.md) | Import CLI and three-way recovery parity | Pending | 14 |
| [16](../phases/phase-22/tasks/task-16.md) | Retire escalation state and duplicate recovery truth | Pending | 15 |
| [17](../phases/phase-22/tasks/task-17.md) | Maintenance, capability qualification, and rebuild budget | Pending | 16 |
| [18](../phases/phase-22/tasks/task-18.md) | Final recovery conformance and Human validation fixture | Pending | 17 |

## Validation Contract

`../phases/phase-22/qa/validation-contract.md`: 80 Draft assertions, each
cited by an owning task. QA Lead owns 77 final deterministic validations;
Human owns VAL-076..078. QA Reviewer checks windows without replacing owners.

| Assertions | Scope / tasks | Final validator |
|---|---|---|
| VAL-001..024,079 | Foundation/recovery, tasks 1-6 | QA Lead |
| VAL-025..039 | Budgets/Stops/closure, tasks 7-9 | QA Lead |
| VAL-040..052 | Activation/freshness/resume, tasks 10-12 | QA Lead |
| VAL-054..065 | Portable recovery, tasks 13-15 | QA Lead |
| VAL-066..072 | Removal/maintenance/performance, tasks 16-17 | QA Lead |
| VAL-053,073..075,080 | Final conformance/consistency, task 18 | QA Lead |
| VAL-076..078 | Human recovery/activation/import, task 18 fixture | Human |

## QA Artifacts

| Artifact | Path under phase-22/ | Status |
|---|---|---|
| QA Lead handoff | qa/qa-lead-handoff.md | Written |
| QA plan / consistency review | qa/qa-plan.md / qa/consistency-review.md | Pending QA Lead |
| Mid-phase reviews | qa/review-1.md through qa/review-6.md | After tasks 3,6,9,12,15,18 |
| Independent checkpoint | qa/checkpoint-after-task-12.md | Required before task 13 |
| Final review / Human checklist | qa/final-review.md / qa/human-checklist.md | Pending |

## Active Constraints

- Proposed semantic/wire/DDL authority: `../phases/phase-22/contracts.md`,
  `public-interface.md`, and `schema.md`. DEC-180 requires Human acceptance.
- Exact removal and executable ownership: `v1-convergence-inventory.md` and
  `execution-scope.md` in that phase directory; DEC-181.
- Every task runs the ordinary full task gate AND tagged integration gate.
  Cases enter with each wire surface; task 6 is the early real-process journey.
- Terminated/live/unknown liveness is closed. Unknown always blocks. Known-live
  takeover requires an exact current Human force grant; no automatic respawn.
- No recovery Git mutation or budget reset. Adoption is not trusted evidence.
- Claims capability stays false until Phase 24 parallel semantics qualify.
- CF-024 unchanged three-fresh-sample <30 s rebuild rule; CF-034 direct workflow/
  evidence coverage in tasks 2/7; CF-035 consuming-config proposal in task 11.
- CF-033 inventory is recorded; AGENTS maintenance edit remains separately
  scheduled. CF-018 fires only if production sqlquery.go is edited.
- No product code/tests, gates, or agent activation by Tech Lead.
  Runner report lint and schema-v2 ledger rules remain unchanged.

## Parity Matrix

| Axis | Live | Rebuild | Import / fake | Evidence |
|---|---|---|---|---|
| Identity/history | Original entity/Event IDs | Same IDs/order | Same IDs/order | VAL-058,063,073 |
| Recovery | Run/Claim/Stop/Attention facts | Same predicates | Same lineage and unresolved work | VAL-017..021,051/052,063 |
| Evidence/budgets | Observed target and counters | Same derived dimensions | Same proof/counters; local observation explicit | VAL-027..029,048..050,063 |
| Machine/Human | Closed JSON projection | Same JSON | Same semantic projection | VAL-053,075 |
| Archive | Manifest v2 | Canonical equivalent | Re-export equivalent | VAL-054,064 |

## Interfaces / Subsystems In Scope

Existing workflow, claims, evidence, validation, project, journal, db, export,
maintenance and CLI boundaries; no new runtime dependency. Exact ownership is
in the task specs, not inferred from this summary.

## Known Risks

- Limited rebuild headroom; retain reference and per-sample failure semantics.
- Force recovery fences fid authority, not a live process's direct file writes.
- Archive validation/publication and cross-resource agreement need fault evidence.
- Cross-product semantic review is pending against the kit commit pinned in the
  QA handoff; no production adapter qualification is inferred.

## Open Questions

Human acceptance of Heavy mode and DEC-180/182 proposals is pending.
QA must resolve the handoff's feasibility/consistency items before Ready.

## Success / Exit Criteria

All task gates and 80 assertions have evidence-backed outcomes, required reviews
and Human checks pass, owned legacy state is removed, every exception is
explicitly dispositioned, and Tech Lead records spec sync and phase sign-off.
No first-adoption or release authorization is implied.

## Next Planning Checkpoint

QA Lead prepares the plan and independent contract review; Tech Lead refines the
pack and Human approves execution. Phase 23 review/Human judgment follows this
phase's signed-off recovery boundary.
