# Phase Brief - Phase 22 (Approved)

## Phase Goal

Deliver the one-Issue failure-and-recovery journey: interrupt execution,
reconcile repository state, start a linked successor Run, enforce budgets and
typed Stops, dispose Attention, recover activation, stale old proof, import/
export the same state, rebuild, and resume without contradictory facts.

## Status

**Executing (DEC-188).** Branch `phase-22` from merged main
`89d964a84dbdd297fdf472d6f78b001e2efac702`. Phase 21 is wrapped (DEC-179) and
merged. Contracts and task pack are refined against QA R-01..14 and all kit
comparison items (DEC-184..187); Human approved the refined pack on 2026-09-05.
See `../phases/phase-22/runner-handoff.md`.

Tasks 1-5 complete; review-1b PASS. Task 6 stopped at TRG-22-008 / ESC-22-003,
resolved under DEC-194: task 4 closed COMPLETE without the `run.checkpoint`
producer or `run.inspect` read API that four accepted records assign to it, so
the owner was never ambiguous and this is undelivered in-scope work. DEC-194
moves execution of four `internal/workflow` files into task 6 with task-4
provenance recorded; task 4 stays COMPLETE and unedited. Task 6 is amended and
ready for a same-task continuation. No open ESC.

## Operating Mode

**Heavy (approved, DEC-188).** Persistence, public JSON, executor authority, and import
publication change. Full artifacts and fresh QA feasibility are required under
`09_AI_WORKFLOW.md` section 4. Mode and ordered task list are approved.

## Task Index

All specs are under `../phases/phase-22/tasks/`; reports use `task-N-done.md`.
Tasks are sequential; reviews follow every three tasks.

| Task | Title | Status | Dependency |
|---|---|---|---|
| [1](../phases/phase-22/tasks/task-1.md) | Recovery schema and replay foundation | Complete | Runner QA readiness refresh |
| [2](../phases/phase-22/tasks/task-2.md) | Typed Stops and Attention domain | Complete | 1 |
| [3](../phases/phase-22/tasks/task-3.md) | Lease lifecycle and executor observations | Complete | 2 |
| [4](../phases/phase-22/tasks/task-4.md) | Read-only repository reconciliation | Complete | 3 |
| [5](../phases/phase-22/tasks/task-5.md) | Atomic linked Run recovery | Complete | 4 |
| [6](../phases/phase-22/tasks/task-6.md) | Recovery CLI and early process journey | Continuation (DEC-194) | 5 |
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

`../phases/phase-22/qa/validation-contract.md`: 80 Ready assertions, each cited
by an owning task. QA Lead owns 77; Human owns VAL-076..078.

| Assertions | Scope / tasks | Final validator |
|---|---|---|
| VAL-001..024,079 | Foundation/recovery, tasks 1-6 | QA Lead |
| VAL-025..039 | Budgets/Stops/closure, tasks 7-9 | QA Lead |
| VAL-040..052 | Activation/freshness/resume, tasks 10-12 | QA Lead |
| VAL-054..065 | Portable recovery, tasks 13-15 | QA Lead |
| VAL-066..072 | Removal/maintenance/performance, tasks 16-17 | QA Lead |
| VAL-053,073..075,080 | Final conformance/consistency, task 18 | QA Lead |
| VAL-076..078 | Human recovery/activation/import, task 18 | Human |

## QA Artifacts

| Artifact | Path under phase-22/ | Status |
|---|---|---|
| QA Lead handoff / TL refinement | qa/qa-lead-handoff.md / qa/tl-refinement.md | Written; accepted DEC-188 |
| QA plan / consistency review | qa/qa-plan.md / qa/consistency-review.md | Ready; refreshed before task 1 |
| Mid-phase reviews | qa/review-1.md through qa/review-6.md | After tasks 3,6,9,12,15,18 |
| Independent checkpoint | qa/checkpoint-after-task-12.md | Required before task 13 |
| Final review / Human checklist | qa/final-review.md / qa/human-checklist.md | Pending |

## Active Constraints

- Accepted semantic/wire/DDL authority: `../phases/phase-22/contracts.md`,
  `public-interface.md`, and `schema.md`. DEC-188 ratifies DEC-180/182/184..187.
- Exact removal and executable ownership: `v1-convergence-inventory.md` and
  `execution-scope.md` in that phase directory; DEC-181.
- Every task runs the ordinary full task gate AND tagged integration gate.
  Cases enter with each wire surface; task 6 is the early real-process journey.
- Terminated/live/unknown liveness is closed. Unknown always blocks. Known-live
  takeover requires an exact current Human force grant; no automatic respawn.
- No recovery Git mutation or budget reset. Adoption is not trusted evidence.
- Claims capability stays false until Phase 24 parallel semantics qualify.
- CF-024 unchanged three-fresh-sample <30 s rebuild rule, early samples at 1/16;
  R-14 synchronized fixture repair at task 8, earlier-failure route in DEC-186.
- CF-034 workflow/evidence coverage in tasks 2/7; CF-035 config in task 11.
- DEC-192: each fenced recovery Event in `internal/journal/project_effects.go`
  is unfenced by its `schema.md` section 5 owner (4,5,7,8,10); task 10 removes
  the empty fence; VAL-003 convergence is scheduled for tasks 5 and 10.
- DEC-193: no test pins an absolute date on a writer-appended Event in a
  live-clock workspace; such a failure is DEC-189 nondeterminism owned by the
  observing task, not convergence.
- DEC-194: task 6 also implements `internal/workflow/run_{checkpoint,inspect}.go`
  and their tests - task 4's undelivered domain scope; record task-4 provenance.
  QA window 2 audits VAL-013's checkpoint clause with it. Read the accepted
  ownership records before calling any owner ambiguous.
- CF-033 recorded; AGENTS edit scheduled; CF-018 fires on sqlquery.go edits.
- Canonical digests: contracts 3.1/3.2; closed 23-action set: interface 2.1.
- Kit dispositions: `kit-contract-mapping.md`; CF-037 retains external adapter
  obligations for Phase 25 qualification.
- No TL product code/tests, gates, or agent activation; Runner lint unchanged.
- In scope: existing workflow, claims, evidence, validation, project, journal,
  db, export, maintenance and CLI boundaries; no new dependency.

## Parity Matrix

| Axis | Live | Rebuild | Import / fake | Evidence |
|---|---|---|---|---|
| Identity/history | Original entity/Event IDs | Same IDs/order | Same IDs/order | VAL-058,063,073 |
| Recovery | Run/Claim/Stop/Attention facts | Same predicates | Same lineage and unresolved work | VAL-017..021,051/052,063 |
| Evidence/budgets | Observed target and counters | Same derived dimensions | Same proof/counters; local observation explicit | VAL-027..029,048..050,063 |
| Machine/Human | Closed JSON projection | Same JSON | Same semantic projection | VAL-053,075 |
| Archive | Manifest v2 | Canonical equivalent | Re-export equivalent | VAL-054,064 |

## Known Risks

- Limited rebuild headroom; retain reference and per-sample failure semantics.
- Force recovery fences fid authority, not a live process's direct file writes.
- Archive validation/publication and cross-resource agreement need fault proof.
- DEC-187 records every pinned-kit difference; no production adapter qualifies.
- A task file scope is not self-enforcing (DEC-194): a named production path
  without a failing acceptance criterion can go unwritten for tasks.

## Open Questions / Success / Exit Criteria

No open questions: all 14 QA refinements are dispositioned; runtime evidence and
the three Human assertion outcomes remain required.
All task gates and 80 assertions have evidence-backed outcomes, required reviews
and Human checks pass, owned legacy state is removed, every exception is
dispositioned, and Tech Lead records spec sync and sign-off. No first-adoption or
release authorization is implied.

## Next Planning Checkpoint

Runner appends RUN-RESUME retaining TRG-22-008 / ESC-22-003 with `Resolution
authority: DEC-194`, then dispatches the bound task-6 same-task continuation at
the next monotonic attempt. Review 2 covers tasks 4-6; the independent task-12
checkpoint must pass before task 13. Human checks and TL wrap-up close it.
