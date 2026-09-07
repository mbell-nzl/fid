# Phase Brief - Phase 22 (Approved)

## Phase Goal

Deliver the one-Issue failure-and-recovery journey: interrupt execution,
reconcile repository state, start a linked successor Run, enforce budgets and
typed Stops, dispose Attention, recover activation, stale old proof,
import/export the same state, rebuild, resume without contradictory facts.

## Status

**Executing (DEC-188).** Branch `phase-22` from merged main
`89d964a84dbdd297fdf472d6f78b001e2efac702`; Phase 21 is wrapped (DEC-179) and
merged. Contracts and task pack are refined against QA R-01..14 and all kit
comparison items (DEC-184..187); Human approved the refined pack on 2026-09-05.
See `../phases/phase-22/runner-handoff.md`. Status is per-task below.

Tasks 1-9 complete; reviews 1b, 2b and 3 PASS. Task 8 closed under DEC-195 with
the exact cancelled-Receipt DDL/filter convergence and all eight gates passed;
its failed history remains in the report. Task 9 completed its final five gates
and actual-operation Stop matrix; independent review 3 passed. Task 10 is
executing under DEC-198, which reconciles prepare versus component staging
while preserving the accepted closed interface and immutable content.
DEC-196 assigns semantic fake strengthening to task 9, Human resume rendering
to task 12 and retained Human fixture generation to task 18. Root orchestrates
under the Human delegation in `delegate-handoff.md`. No open ESC.

## Operating Mode

**Heavy (approved, DEC-188).** Persistence, public JSON, executor authority, and
import publication change; full artifacts and fresh QA feasibility are required
under `09_AI_WORKFLOW.md` section 4. Mode and ordered task list are approved.

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
| [6](../phases/phase-22/tasks/task-6.md) | Recovery CLI and early process journey | Complete (DEC-194) | 5 |
| [7](../phases/phase-22/tasks/task-7.md) | Durable budget accounting and execution admission | Complete | 6 |
| [8](../phases/phase-22/tasks/task-8.md) | Interrupted Attempts and derived Run closure | Complete (DEC-195) | 7 |
| [9](../phases/phase-22/tasks/task-9.md) | Stop, budget, and closure commands | Complete (DEC-196/197) | 8 |
| [10](../phases/phase-22/tasks/task-10.md) | Activation recovery domain | Executing (DEC-198) | 9 |
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
| VAL-053,073..075,080 / VAL-076..078 | Final conformance, task 18 | QA Lead / Human |

## QA Artifacts

| Artifact | Path under phase-22/ | Status |
|---|---|---|
| Handoff / refinement / plan / consistency | qa/qa-lead-handoff.md, qa/tl-refinement.md, qa/qa-plan.md, qa/consistency-review.md | Accepted DEC-188; plan refreshed before task 1 |
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
- Terminated/live/unknown liveness is closed; unknown always blocks. Known-live
  takeover needs an exact current Human force grant; no automatic respawn. No
  recovery Git mutation or budget reset. Adoption is not trusted evidence.
- Claims capability stays false until Phase 24 parallel semantics qualify.
- CF-024 unchanged three-fresh-sample <30 s rebuild rule, early samples at 1/16;
  R-14 synchronized fixture repair at task 8 (DEC-186). CF-034 coverage tasks
  2/7; CF-035 config task 11; CF-033 recorded; CF-018 fires on sqlquery edits.
- DEC-192: each fenced recovery Event in `internal/journal/project_effects.go` is
  unfenced by its `schema.md` section 5 owner (4,5,7,8,10); task 10 removes the
  empty fence; VAL-003 convergence is scheduled for tasks 5 and 10.
- DEC-193: no test pins an absolute date on a writer-appended Event in a
  live-clock workspace; that failure is DEC-189 nondeterminism, not convergence.
- DEC-194: task 6 delivered task 4's `run_{checkpoint,inspect}.go` domain scope
  with task-4 provenance; QA window 2 audits VAL-013's checkpoint clause with it.
- DEC-195: `schema.md` section 1 admits `cancelled` in
  `execution_receipts.terminal_status` plus one composite-CHECK clause. Task 8
  makes those two exact 0013 edits with task-1 provenance, converges the
  `internal/db` golden DDL string, and opens the closed receipt outcome filter.
  VAL-031 stands as written; never substitute an admitted status.
- Canonical digests: contracts 3.1/3.2; closed 23-action set: interface 2.1.
- Kit dispositions: `kit-contract-mapping.md`; CF-037 retains adapter
  obligations for Phase 25.
- No TL product code/tests, gates, or agent activation; Runner lint unchanged.
- In scope: existing workflow, claims, evidence, validation, project, journal,
  db, export, maintenance, CLI boundaries; no new dependency.

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
- Force recovery fences fid authority, not a live process's direct writes.
- Archive validation/publication and cross-resource agreement need fault proof;
  DEC-187 records every pinned-kit difference and no production adapter qualifies.
- Neither a task file scope (DEC-194) nor a contract clause (DEC-195) is
  self-enforcing: a named path with no failing criterion, or a closed-domain
  member with no schema row, fails nothing until code reaches it.

## Open Questions / Success / Exit Criteria

No open questions: all 14 QA refinements are dispositioned; runtime evidence and
the three Human assertion outcomes remain required. Exit criteria: all task gates
and 80 assertions have evidence-backed outcomes, required reviews and Human
checks pass, owned legacy state is removed, every exception is dispositioned, and
Tech Lead records spec sync and sign-off. No release authorization is implied.

## Next Planning Checkpoint

Review 3 passed tasks 7-9, including DEC-196 semantic fake and actual-operation
Stop evidence. Tasks 10-12 lead to review 4 and the independent task-12
checkpoint, which must pass before task 13. Final
deterministic QA and runnable Human fixtures precede tomorrow's Human judgments;
Human outcomes and full phase sign-off remain pending.
